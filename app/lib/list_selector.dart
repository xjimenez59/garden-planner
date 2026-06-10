import 'dart:async';

import 'package:app/utils.dart';
import 'package:flutter/material.dart';

class ListSelector extends StatefulWidget {
  List<String> optionsList = [];
  final String value;
  final String title;
  final Future<List<String>> Function(dynamic param) getOptions;
  // Options secondaires (catalogue de référence) — affichées en gris, dédupliquées
  final Future<List<String>> Function(dynamic param)? getSecondaryOptions;
  final dynamic optionsParam;

  ListSelector(
      {super.key,
      required this.title,
      required this.value,
      required this.getOptions,
      this.getSecondaryOptions,
      this.optionsParam});

  @override
  State<StatefulWidget> createState() {
    // ignore: no_logic_in_create_state
    return _ListSelector(editedValue: value, getOptions: getOptions);
  }
}

class _ListSelector extends State<ListSelector> {
  String editedValue;
  List<String> _allOptions = [];   // liste complète (primaires + secondaires)
  List<String> optionsList = [];   // liste filtrée affichée
  Set<String> _primaryValues = {}; // pour la coloration
  Future<List<String>> Function(dynamic param) getOptions;
  TextEditingController editController = TextEditingController();

  _ListSelector({this.editedValue = "", required this.getOptions});

  @override
  void initState() {
    editController.text = editedValue;
    _loadOptions();
    super.initState();
  }

  Future<void> _loadOptions() async {
    final primary = await getOptions(widget.optionsParam);
    final primarySet = primary.map((s) => s.toLowerCase()).toSet();

    List<String> secondary = [];
    if (widget.getSecondaryOptions != null) {
      final raw = await widget.getSecondaryOptions!(widget.optionsParam);
      // Dédupliquer insensiblement à la casse
      secondary = raw
          .where((s) => !primarySet.contains(s.toLowerCase()))
          .toList();
    }

    if (!mounted) return;
    setState(() {
      _primaryValues = primary.toSet();
      _allOptions = [...primary, ...secondary];
      widget.optionsList = _allOptions;
      optionsList = _allOptions;
    });
  }

  @override
  Widget build(BuildContext context) {
    final content = Column(children: [
      Expanded(
          child: optionsList.isEmpty
              ? const Center(child: CircularProgressIndicator())
              : ListView.builder(
                  itemCount: optionsList.length,
                  itemBuilder: (context, index) {
                    final item = optionsList[index];
                    final isPrimary = _primaryValues.contains(item);
                    return DecoratedBox(
                      decoration: BoxDecoration(
                          color: (index % 2 == 0)
                              ? Colors.white
                              : Colors.grey[10],
                          border: const Border(
                              bottom: BorderSide(color: Colors.grey))),
                      child: ListTile(
                        title: Text(
                          item,
                          style: TextStyle(
                            color: isPrimary
                                ? Colors.black87
                                : Colors.grey.shade500,
                          ),
                        ),
                        onTap: () => optionOnTap(index),
                        selected: item == editedValue,
                      ),
                    );
                  })),
      Container(
          padding: const EdgeInsets.fromLTRB(5, 10, 0, 5),
          alignment: Alignment.bottomLeft,
          color: Colors.blue,
          child: const Text("Choisissez dans la liste, ou saisissez ci-dessous",
              textAlign: TextAlign.left,
              style: TextStyle(color: Colors.white))),
      TextField(
        controller: editController,
        onChanged: filterOptions,
        onSubmitted: editOnSubmit,
        autofocus: true,
      ),
    ]);

    return Scaffold(
        appBar: AppBar(title: Text(widget.title)),
        body: Center(child: content));
  }

  void filterOptions(String newValue) {
    setState(() {
      editedValue = newValue;
      if (newValue == "" || newValue == widget.value) {
        optionsList = _allOptions;
      } else {
        optionsList = _allOptions
            .where((element) => element.withoutDiacriticalMarks
                .toLowerCase()
                .contains(newValue.withoutDiacriticalMarks.toLowerCase()))
            .toList();
      }
    });
  }

  void optionOnTap(int index) {
    String newValue = optionsList[index];
    editController.text = newValue;
    filterOptions(newValue);
    Navigator.pop(context, newValue);
  }

  void editOnSubmit(String newValue) {
    Navigator.pop(context, newValue);
  }
}
