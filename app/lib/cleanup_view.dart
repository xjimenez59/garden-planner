// ignore_for_file: prefer_const_constructors

import 'package:app/api_service.dart';
import 'package:app/cleanup_model.dart';
import 'package:app/garden_model.dart';
import 'package:app/utils.dart';
import 'package:flutter/material.dart';

// ─── Similarité par trigrammes ───────────────────────────────────────────────

String _norm(String s) =>
    s.toLowerCase().withoutDiacriticalMarks.replaceAll(RegExp(r'[^a-z]'), '');

String _firstLetter(String value) {
  final n = _norm(value);
  return n.isEmpty ? '' : n[0].toUpperCase();
}

Set<String> _trigrams(String normalized) {
  if (normalized.length < 3) return {normalized};
  final result = <String>{};
  for (var i = 0; i <= normalized.length - 3; i++) {
    result.add(normalized.substring(i, i + 3));
  }
  return result;
}

double _similarity(Set<String> ta, Set<String> tb) {
  if (ta.isEmpty && tb.isEmpty) return 1.0;
  final intersection = ta.intersection(tb).length;
  final union = ta.length + tb.length - intersection;
  return union == 0 ? 1.0 : intersection / union;
}

const _kThreshold = 0.30;
const _kItemHeight = 48.0;
const _kItemHeightWithSugg = 64.0;
const _kAlphaBarWidth = 24.0;
const _kMinItemsForAlphaBar = 20;

class _RefEntry {
  final String name;
  final Set<String> tgrams;
  _RefEntry(this.name, this.tgrams);
}

class _Suggestion {
  final String name;
  final double score;
  _Suggestion(this.name, this.score);
}

// ─── Tabs ─────────────────────────────────────────────────────────────────────

const _tabFields = ['legumes', 'varietes', 'lieux', 'tags'];
const _tabLabels = ['Légumes', 'Variétés', 'Lieux', 'Tags'];

class CleanupView extends StatefulWidget {
  final Garden garden;
  const CleanupView({super.key, required this.garden});

  @override
  State<CleanupView> createState() => _CleanupViewState();
}

class _CleanupViewState extends State<CleanupView>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;
  late List<ScrollController> _scrollControllers;

  final List<List<CleanupItem>> _data = [[], [], [], []];
  final List<bool> _loading = [false, false, false, false];
  // Suggestions précomputées par tab × item
  final List<List<List<_Suggestion>>> _sugg = [[], [], [], []];

  List<_RefEntry> _refLegumes = [];
  List<_RefEntry> _refVarietes = []; // toutes variétés, dédupliquées

  @override
  void initState() {
    super.initState();
    _scrollControllers = List.generate(4, (_) => ScrollController());
    _tabController = TabController(length: 4, vsync: this);
    _tabController.addListener(() {
      if (!_tabController.indexIsChanging) _loadTab(_tabController.index);
    });
    _loadReferences();
    for (var i = 0; i < 4; i++) { _loadTab(i); }
  }

  @override
  void dispose() {
    _tabController.dispose();
    for (final c in _scrollControllers) { c.dispose(); }
    super.dispose();
  }

  // ── Référentiel ─────────────────────────────────────────────────────────────

  Future<void> _loadReferences() async {
    final refs = await ApiService().getLegumesReference();
    if (!mounted) return;
    final seen = <String>{};
    setState(() {
      _refLegumes = refs
          .map((r) => _RefEntry(r.legume, _trigrams(_norm(r.legume))))
          .toList();
      _refVarietes = refs
          .expand((r) => r.varietes)
          .where(seen.add)
          .map((v) => _RefEntry(v, _trigrams(_norm(v))))
          .toList();
      // Recalculer les suggestions pour les tabs déjà chargés
      for (var i = 0; i < 4; i++) {
        if (_data[i].isNotEmpty) _rebuildSugg(i);
      }
    });
  }

  void _rebuildSugg(int ti) {
    final refs = ti == 0
        ? _refLegumes
        : ti == 1
            ? _refVarietes
            : <_RefEntry>[];
    _sugg[ti] = _data[ti].map((item) {
      if (refs.isEmpty) return <_Suggestion>[];
      final ut = _trigrams(_norm(item.value));
      return refs
          .map((r) => _Suggestion(r.name, _similarity(ut, r.tgrams)))
          .where((s) => s.score >= _kThreshold && s.name != item.value)
          .toList()
        ..sort((a, b) => b.score.compareTo(a.score));
    }).toList();
  }

  // ── Données ─────────────────────────────────────────────────────────────────

  Future<void> _loadTab(int ti) async {
    if (!mounted) return;
    setState(() => _loading[ti] = true);
    var items =
        await ApiService().getCleanupList(widget.garden.ID, _tabFields[ti]);
    if (!mounted) return;
    // Tri normalisé côté Flutter (gère les accents mieux que SQLite par défaut)
    items.sort((a, b) => _norm(a.value).compareTo(_norm(b.value)));
    setState(() {
      _data[ti] = items;
      _loading[ti] = false;
      _rebuildSugg(ti);
    });
  }

  // ── Navigation alphabétique ─────────────────────────────────────────────────

  void _scrollToLetter(int ti, String letter) {
    double offset = 0;
    for (var i = 0; i < _data[ti].length; i++) {
      if (_firstLetter(_data[ti][i].value) == letter) {
        _scrollControllers[ti].animateTo(
          offset,
          duration: Duration(milliseconds: 180),
          curve: Curves.easeOut,
        );
        return;
      }
      final hasSugg = i < _sugg[ti].length && _sugg[ti][i].isNotEmpty;
      offset += (hasSugg ? _kItemHeightWithSugg : _kItemHeight) + 1; // +1 divider
    }
  }

  // ── Actions ─────────────────────────────────────────────────────────────────

  void _showActions(int ti, CleanupItem item, List<_Suggestion> suggestions) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      builder: (_) => _ActionSheet(
        item: item,
        field: _tabFields[ti],
        gardenId: widget.garden.ID,
        allItems: _data[ti],
        suggestions: suggestions,
        onDone: () => _loadTab(ti),
      ),
    );
  }

  // ── Build ────────────────────────────────────────────────────────────────────

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('Nettoyage — ${widget.garden.Nom}'),
        bottom: TabBar(
          controller: _tabController,
          tabs: _tabLabels.map((l) => Tab(text: l)).toList(),
        ),
      ),
      body: TabBarView(
        controller: _tabController,
        children: List.generate(4, (ti) {
          if (_loading[ti]) {
            return Center(child: CircularProgressIndicator());
          }
          if (_data[ti].isEmpty) {
            return Center(
              child: Text('Aucune valeur',
                  style: TextStyle(color: Colors.grey.shade500)),
            );
          }
          final showAlpha = _data[ti].length >= _kMinItemsForAlphaBar;
          return Stack(
            children: [
              ListView.separated(
                controller: _scrollControllers[ti],
                padding: showAlpha
                    ? EdgeInsets.only(right: _kAlphaBarWidth)
                    : EdgeInsets.zero,
                itemCount: _data[ti].length,
                separatorBuilder: (_, __) =>
                    Divider(height: 1, color: Colors.grey.shade200),
                itemBuilder: (context, j) {
                  final item = _data[ti][j];
                  final suggestions =
                      j < _sugg[ti].length ? _sugg[ti][j] : <_Suggestion>[];
                  final best = suggestions.isNotEmpty ? suggestions.first : null;
                  return ListTile(
                    dense: true,
                    title: Text(item.value),
                    subtitle: best != null
                        ? Text(
                            '≈ ${best.name}  (${(best.score * 100).round()}%)',
                            style: TextStyle(
                                fontSize: 11, color: Colors.amber.shade700),
                          )
                        : null,
                    trailing: Row(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        if (best != null)
                          Padding(
                            padding: EdgeInsets.only(right: 4),
                            child: Icon(Icons.auto_fix_high,
                                size: 16, color: Colors.amber.shade600),
                          ),
                        _CountBadge(item.countLabel),
                        SizedBox(width: 4),
                        Icon(Icons.chevron_right, color: Colors.grey.shade400),
                      ],
                    ),
                    onTap: () => _showActions(ti, item, suggestions),
                  );
                },
              ),
              if (showAlpha)
                Positioned(
                  right: 0,
                  top: 0,
                  bottom: 0,
                  child: _AlphaBar(
                    items: _data[ti],
                    onLetterTap: (letter) => _scrollToLetter(ti, letter),
                  ),
                ),
            ],
          );
        }),
      ),
    );
  }
}

// ─── Barre alphabétique ───────────────────────────────────────────────────────

class _AlphaBar extends StatelessWidget {
  final List<CleanupItem> items;
  final void Function(String letter) onLetterTap;

  const _AlphaBar({required this.items, required this.onLetterTap});

  @override
  Widget build(BuildContext context) {
    final letters = items
        .map((i) => _firstLetter(i.value))
        .where((l) => l.isNotEmpty)
        .toSet()
        .toList()
      ..sort();

    return Container(
      width: _kAlphaBarWidth,
      color: Colors.white.withValues(alpha: 0.88),
      child: SingleChildScrollView(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: letters
              .map((letter) => GestureDetector(
                    onTap: () => onLetterTap(letter),
                    child: SizedBox(
                      width: _kAlphaBarWidth,
                      height: 22,
                      child: Center(
                        child: Text(
                          letter,
                          style: TextStyle(
                            fontSize: 11,
                            fontWeight: FontWeight.w600,
                            color: Colors.blue.shade700,
                          ),
                        ),
                      ),
                    ),
                  ))
              .toList(),
        ),
      ),
    );
  }
}

// ─── Badge compteur ───────────────────────────────────────────────────────────

class _CountBadge extends StatelessWidget {
  final String label;
  const _CountBadge(this.label);

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: EdgeInsets.symmetric(horizontal: 6, vertical: 2),
      decoration: BoxDecoration(
        color: Colors.grey.shade200,
        borderRadius: BorderRadius.circular(10),
      ),
      child: Text(label,
          style: TextStyle(fontSize: 11, color: Colors.grey.shade700)),
    );
  }
}

// ─── Bottom sheet d'actions ───────────────────────────────────────────────────

class _ActionSheet extends StatefulWidget {
  final CleanupItem item;
  final String field;
  final String gardenId;
  final List<CleanupItem> allItems;
  final List<_Suggestion> suggestions;
  final VoidCallback onDone;

  const _ActionSheet({
    required this.item,
    required this.field,
    required this.gardenId,
    required this.allItems,
    required this.suggestions,
    required this.onDone,
  });

  @override
  State<_ActionSheet> createState() => _ActionSheetState();
}

class _ActionSheetState extends State<_ActionSheet> {
  bool _renaming = false;
  late TextEditingController _renameCtrl;
  bool _busy = false;

  @override
  void initState() {
    super.initState();
    _renameCtrl = TextEditingController(text: widget.item.value);
  }

  @override
  void dispose() {
    _renameCtrl.dispose();
    super.dispose();
  }

  Future<void> _doRename(String newValue) async {
    if (newValue.trim().isEmpty || newValue == widget.item.value) {
      Navigator.pop(context);
      return;
    }
    setState(() => _busy = true);
    await ApiService().renameCleanupValue(
        widget.gardenId, widget.field, widget.item.value, newValue.trim());
    if (mounted) Navigator.pop(context);
    widget.onDone();
  }

  Future<void> _doDelete({String action = 'clear', String with_ = ''}) async {
    setState(() => _busy = true);
    await ApiService().deleteCleanupValue(
        widget.gardenId, widget.field, widget.item.value,
        action: action, with_: with_);
    if (mounted) Navigator.pop(context);
    widget.onDone();
  }

  void _pickMergeTarget({required bool forDelete}) async {
    final others =
        widget.allItems.where((i) => i.value != widget.item.value).toList();
    final chosen = await showModalBottomSheet<CleanupItem>(
      context: context,
      isScrollControlled: true,
      builder: (_) => _ValuePicker(items: others),
    );
    if (chosen == null) return;
    if (forDelete) {
      await _doDelete(action: 'replace', with_: chosen.value);
    } else {
      await _doRename(chosen.value);
    }
  }

  void _confirmDelete() {
    if (widget.item.count == 0) {
      _doDelete();
      return;
    }
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: Text('Supprimer "${widget.item.value}"'),
        content: Text(
            '${widget.item.countLabel} log${widget.item.count > 1 ? 's utilisent' : ' utilise'} cette valeur.'),
        actions: [
          TextButton(
              onPressed: () => Navigator.pop(ctx), child: Text('Annuler')),
          TextButton(
              onPressed: () {
                Navigator.pop(ctx);
                _doDelete();
              },
              child: Text('Vider le champ')),
          TextButton(
              onPressed: () {
                Navigator.pop(ctx);
                _pickMergeTarget(forDelete: true);
              },
              child: Text('Remplacer par...')),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final hasSuggestions = widget.suggestions.isNotEmpty && !_renaming;
    return Padding(
      padding:
          EdgeInsets.only(bottom: MediaQuery.of(context).viewInsets.bottom),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Padding(
            padding: EdgeInsets.fromLTRB(16, 16, 16, 4),
            child: Row(children: [
              Expanded(
                child: Text(widget.item.value,
                    style:
                        TextStyle(fontSize: 16, fontWeight: FontWeight.bold)),
              ),
              _CountBadge(widget.item.countLabel),
            ]),
          ),
          Divider(),
          if (hasSuggestions) ...[
            Padding(
              padding: EdgeInsets.fromLTRB(16, 4, 16, 0),
              child: Text('Suggestions',
                  style: TextStyle(
                      fontSize: 11,
                      color: Colors.grey.shade500,
                      letterSpacing: 0.5)),
            ),
            for (final s in widget.suggestions.take(3))
              ListTile(
                dense: true,
                leading: Icon(Icons.auto_fix_high,
                    size: 20, color: Colors.amber.shade600),
                title: Text('→ "${s.name}"'),
                trailing: Text('${(s.score * 100).round()}%',
                    style:
                        TextStyle(fontSize: 12, color: Colors.grey.shade500)),
                onTap: _busy ? null : () => _doRename(s.name),
              ),
            Divider(),
          ],
          if (_renaming) ...[
            Padding(
              padding: EdgeInsets.symmetric(horizontal: 16),
              child: TextField(
                controller: _renameCtrl,
                autofocus: true,
                decoration: InputDecoration(
                  labelText: 'Nouveau nom',
                  border: OutlineInputBorder(),
                ),
                onSubmitted: _doRename,
              ),
            ),
            Padding(
              padding: EdgeInsets.symmetric(horizontal: 16, vertical: 8),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.end,
                children: [
                  TextButton(
                      onPressed: () => setState(() => _renaming = false),
                      child: Text('Annuler')),
                  SizedBox(width: 8),
                  ElevatedButton(
                      onPressed:
                          _busy ? null : () => _doRename(_renameCtrl.text),
                      child: Text('Valider')),
                ],
              ),
            ),
          ] else ...[
            ListTile(
              leading: Icon(Icons.edit_outlined),
              title: Text('Renommer'),
              onTap: () => setState(() => _renaming = true),
            ),
            ListTile(
              leading: Icon(Icons.merge_outlined),
              title: Text('Fusionner avec...'),
              onTap: () => _pickMergeTarget(forDelete: false),
            ),
            ListTile(
              leading: Icon(Icons.delete_outline, color: Colors.red.shade400),
              title: Text('Supprimer',
                  style: TextStyle(color: Colors.red.shade400)),
              onTap: _confirmDelete,
            ),
          ],
          SizedBox(height: 8),
        ],
      ),
    );
  }
}

// ─── Picker de valeur existante ───────────────────────────────────────────────

class _ValuePicker extends StatefulWidget {
  final List<CleanupItem> items;
  const _ValuePicker({required this.items});

  @override
  State<_ValuePicker> createState() => _ValuePickerState();
}

class _ValuePickerState extends State<_ValuePicker> {
  String _filter = '';

  @override
  Widget build(BuildContext context) {
    final filtered = widget.items
        .where((i) =>
            _filter.isEmpty ||
            i.value.toLowerCase().contains(_filter.toLowerCase()))
        .toList();

    return DraggableScrollableSheet(
      expand: false,
      initialChildSize: 0.5,
      maxChildSize: 0.9,
      builder: (_, ctrl) => Column(
        children: [
          Padding(
            padding: EdgeInsets.fromLTRB(16, 12, 16, 4),
            child: TextField(
              autofocus: true,
              decoration: InputDecoration(
                hintText: 'Rechercher...',
                prefixIcon: Icon(Icons.search, size: 20),
                isDense: true,
                contentPadding: EdgeInsets.symmetric(vertical: 12),
                border: OutlineInputBorder(),
              ),
              onChanged: (v) => setState(() => _filter = v),
            ),
          ),
          Expanded(
            child: ListView.separated(
              controller: ctrl,
              itemCount: filtered.length,
              separatorBuilder: (_, __) =>
                  Divider(height: 1, color: Colors.grey.shade200),
              itemBuilder: (_, i) => ListTile(
                dense: true,
                title: Text(filtered[i].value),
                trailing: _CountBadge(filtered[i].countLabel),
                onTap: () => Navigator.pop(context, filtered[i]),
              ),
            ),
          ),
        ],
      ),
    );
  }
}
