// ignore_for_file: prefer_const_constructors

import 'dart:math';

import 'package:app/action_log.dart';
import 'package:app/garden_model.dart';
import 'package:app/meteo_service.dart';
import 'package:app/previsions_view.dart';
import 'package:app/utils.dart';
import 'package:flutter/material.dart';
import 'package:material_symbols_icons/symbols.dart';

IconData _meteoIcon(Meteo m) {
  final soleil = ((m.sigma ?? 0.0) >= 0.8) || ((m.inst ?? 0.0) >= 6 * 60.0);
  final pluie = ((m.drr ?? 0.0) > 60.0) || ((m.rr ?? 0.0) > 2.0);
  final nuages = pluie || ((m.sigma ?? 1.0) < 0.8) || ((m.inst ?? 12 * 60.0) < 6 * 60.0);
  if (soleil && !nuages) return Icons.wb_sunny;
  if (soleil && pluie) return Symbols.weather_mix;
  if (soleil && nuages && !pluie) return Icons.cloud_queue;
  if (!soleil && pluie) return Icons.water_drop;
  if (!soleil && nuages && !pluie) return Icons.cloud;
  return Icons.remove;
}

class DaySeparator extends StatefulWidget {
  final DateTime date;
  final Meteo? meteo;
  final LuneDay? luneDay;
  // Non-null uniquement pour aujourd'hui → active la frise dépliable
  final List<HourlyForecast>? previsions;
  final List<LuneDay>? luneData;

  const DaySeparator({
    super.key,
    required this.date,
    this.meteo,
    this.luneDay,
    this.previsions,
    this.luneData,
  });

  @override
  State<DaySeparator> createState() => _DaySeparatorState();
}

class _DaySeparatorState extends State<DaySeparator> {
  bool _expanded = false;

  static const _jours = ['Lun', 'Mar', 'Mer', 'Jeu', 'Ven', 'Sam', 'Dim'];
  static const _mois  = ['janv.', 'févr.', 'mars', 'avr.', 'mai', 'juin',
                          'juil.', 'août', 'sept.', 'oct.', 'nov.', 'déc.'];

  bool get _isToday =>
      widget.previsions != null && widget.previsions!.isNotEmpty;

  String _bioemoji() {
    switch (widget.luneDay?.jourBiodynamique) {
      case 'jour_fruit':   return '🍓';
      case 'jour_racine':  return '🥕';
      case 'jour_fleur':   return '🌸';
      case 'jour_feuille': return '🌿';
      default: return '';
    }
  }

  // Résumé météo à droite : depuis Meteo (historique) ou depuis previsions (aujourd'hui)
  Widget _meteoSummary() {
    final m = widget.meteo;

    if (m != null) {
      // Données MétéoFrance historiques
      return Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(_meteoIcon(m), size: 16, color: Colors.blueGrey.shade600),
          SizedBox(width: 4),
          if (m.tx != null || m.tm != null) ...[
            Icon(Icons.arrow_upward, size: 10, color: Colors.orange),
            Text('${(m.tx ?? m.tm)!.round()}°',
                style: TextStyle(
                    fontSize: 12,
                    fontWeight: FontWeight.w600,
                    color: Colors.orange.shade700)),
          ],
          if (m.tn != null || m.tm != null) ...[
            SizedBox(width: 3),
            Icon(Icons.arrow_downward, size: 10, color: Colors.blue),
            Text('${(m.tn ?? m.tm)!.round()}°',
                style: TextStyle(
                    fontSize: 12,
                    fontWeight: FontWeight.w600,
                    color: Colors.blue.shade600)),
          ],
          if ((m.rr ?? 0) > 0) ...[
            SizedBox(width: 6),
            Icon(Icons.water_drop, size: 12, color: Colors.blue.shade400),
            Text(
              m.rr! < 1 ? '${m.rr!.toStringAsFixed(1)}mm' : '${m.rr!.round()}mm',
              style: TextStyle(fontSize: 11, color: Colors.blue.shade600),
            ),
          ],
          if ((m.ffm ?? 0) > 0) ...[
            SizedBox(width: 6),
            Transform.rotate(
              angle: ((m.dxy ?? 180) - 180) * pi / 180,
              child: Icon(Icons.arrow_upward, size: 12, color: Colors.blueGrey),
            ),
            Text('${m.ffm!.round()}',
                style: TextStyle(fontSize: 11, color: Colors.blueGrey.shade700)),
          ],
        ],
      );
    }

    // Prévisions Open-Meteo (aujourd'hui, données MF pas encore disponibles)
    final slots = widget.previsions;
    if (slots == null || slots.isEmpty) return SizedBox.shrink();

    final todayStr = DateTime.now().toIso8601String().substring(0, 10);
    final todaySlots =
        slots.where((f) => f.time.startsWith(todayStr)).toList();
    final src = todaySlots.isNotEmpty ? todaySlots : slots;

    final tmax = src.map((f) => f.temperature).reduce(max);
    final tmin = src.map((f) => f.temperature).reduce(min);
    final totalPrecip = src.fold(0.0, (s, f) => s + f.precipitation);
    final avgWind = src.fold(0.0, (s, f) => s + f.windSpeed) / src.length;
    final midSlot = src[src.length ~/ 2];

    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Icon(wmoIconData(midSlot.weatherCode),
            size: 16, color: Colors.blueGrey.shade600),
        SizedBox(width: 4),
        Icon(Icons.arrow_upward, size: 10, color: Colors.orange),
        Text('${tmax.round()}°',
            style: TextStyle(
                fontSize: 12,
                fontWeight: FontWeight.w600,
                color: Colors.orange.shade700)),
        SizedBox(width: 3),
        Icon(Icons.arrow_downward, size: 10, color: Colors.blue),
        Text('${tmin.round()}°',
            style: TextStyle(
                fontSize: 12,
                fontWeight: FontWeight.w600,
                color: Colors.blue.shade600)),
        if (totalPrecip > 0) ...[
          SizedBox(width: 6),
          Icon(Icons.water_drop, size: 12, color: Colors.blue.shade400),
          Text(
            totalPrecip < 1
                ? '${totalPrecip.toStringAsFixed(1)}mm'
                : '${totalPrecip.round()}mm',
            style: TextStyle(fontSize: 11, color: Colors.blue.shade600),
          ),
        ],
        if (avgWind > 0) ...[
          SizedBox(width: 6),
          Transform.rotate(
            angle: (midSlot.windDir - 180) * pi / 180,
            child: Icon(Icons.arrow_upward, size: 12, color: Colors.blueGrey),
          ),
          Text('${avgWind.round()}',
              style:
                  TextStyle(fontSize: 11, color: Colors.blueGrey.shade700)),
        ],
      ],
    );
  }

  @override
  Widget build(BuildContext context) {
    final wnum = weekNum(widget.date);
    final dayLabel = _isToday
        ? 'Aujourd\'hui'
        : '${_jours[widget.date.weekday - 1]}. ${widget.date.day} ${_mois[widget.date.month - 1]}';
    final emoji = _bioemoji();

    return GestureDetector(
      onTap: _isToday ? () => setState(() => _expanded = !_expanded) : null,
      child: Container(
        margin: EdgeInsets.fromLTRB(0, _isToday ? 0 : 10, 0, 0),
        decoration: BoxDecoration(
          color: _isToday ? Color(0xFFE8F0F8) : Color(0x12000000),
          border: Border.all(color: Color(0x4d9e9e9e), width: 1),
        ),
        width: MediaQuery.of(context).size.width,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Padding(
              padding: EdgeInsets.symmetric(horizontal: 8, vertical: 5),
              child: Row(
                children: [
                  // Gauche : label jour + semaine + biodynamique
                  Text(dayLabel,
                      style: TextStyle(
                          fontSize: 13,
                          fontWeight: _isToday ? FontWeight.bold : FontWeight.w600,
                          color: _isToday
                              ? Colors.blue.shade700
                              : Colors.black87)),
                  SizedBox(width: 6),
                  Text('S$wnum',
                      style: TextStyle(
                          fontSize: 11, color: Colors.grey.shade600)),
                  if (emoji.isNotEmpty) ...[
                    SizedBox(width: 5),
                    Text(emoji, style: TextStyle(fontSize: 13)),
                  ],
                  Spacer(),
                  // Droite : synthèse météo
                  _meteoSummary(),
                  if (_isToday) ...[
                    SizedBox(width: 6),
                    Icon(
                      _expanded
                          ? Icons.keyboard_arrow_up
                          : Icons.keyboard_arrow_down,
                      size: 18,
                      color: Colors.grey,
                    ),
                  ],
                ],
              ),
            ),
            // Frise dépliable (aujourd'hui uniquement) — AnimatedSize évite le flash d'overflow
            if (_isToday)
              ClipRect(
                child: AnimatedSize(
                  duration: Duration(milliseconds: 280),
                  curve: Curves.easeInOut,
                  child: _expanded
                      ? PrevisionsFriseWidget(
                          previsions: widget.previsions!,
                          luneData: widget.luneData ?? [],
                        )
                      : SizedBox.shrink(),
                ),
              ),
          ],
        ),
      ),
    );
  }
}

class YearSeparator extends StatelessWidget {
  final int year;
  const YearSeparator({super.key, required this.year});

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: EdgeInsets.fromLTRB(0, 16, 0, 0),
      padding: EdgeInsets.symmetric(vertical: 6),
      decoration: BoxDecoration(
        border: Border.symmetric(
          horizontal: BorderSide(color: Colors.brown.shade200, width: 1),
        ),
        color: Colors.brown.shade50,
      ),
      width: double.infinity,
      child: Center(
        child: Text(
          '$year',
          style: TextStyle(
            fontSize: 13,
            fontWeight: FontWeight.bold,
            color: Colors.brown.shade600,
            letterSpacing: 2,
          ),
        ),
      ),
    );
  }
}

class TopHomeFilter extends StatelessWidget {
  final TextEditingController filterController;
  final void Function(String text) onFilterChanged;
  final List<Garden> jardins;
  final void Function() onSelectGardenTap;

  const TopHomeFilter(
      {required this.jardins,
      required this.filterController,
      required this.onFilterChanged,
      required this.onSelectGardenTap});

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: EdgeInsets.all(0),
      padding: EdgeInsets.all(0),
      width: MediaQuery.of(context).size.width,
      decoration: BoxDecoration(
        color: Color(0x1f000000),
        shape: BoxShape.rectangle,
        borderRadius: BorderRadius.zero,
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.start,
        crossAxisAlignment: CrossAxisAlignment.center,
        mainAxisSize: MainAxisSize.max,
        children: [
          Expanded(
            child: TextField(
              controller: filterController,
              onChanged: onFilterChanged,
              decoration: InputDecoration(
                isDense: true,
                prefixIcon: Icon(Icons.search, size: 20),
                hintText: 'Rechercher...',
                border: InputBorder.none,
                suffixIcon: filterController.text.isNotEmpty
                    ? IconButton(
                        onPressed: () {
                          filterController.clear();
                          onFilterChanged('');
                        },
                        icon: Icon(Icons.close, size: 20),
                      )
                    : null,
              ),
            ),
          ),
          IconButton(
            onPressed: onSelectGardenTap,
            tooltip: 'Changer de jardin',
            icon: Icon(Icons.location_pin),
          ),
        ],
      ),
    );
  }
}

class HorizontalImageListview extends StatelessWidget {
  final List<String> imgUrlList;

  const HorizontalImageListview({super.key, required this.imgUrlList});

  @override
  Widget build(BuildContext context) {
    if (imgUrlList.isEmpty) return Container();
    return Container(
        height: 200,
        child: ListView.builder(
            scrollDirection: Axis.horizontal,
            itemCount: imgUrlList.length,
            itemBuilder: (context, index) {
              return Container(
                  padding: EdgeInsets.all(10),
                  height: 150,
                  child: Image.network(imgUrlList[index]));
            }));
  }
}

class ActionListTile extends StatelessWidget {
  final ActionLog actionLog;
  final bool showDivider;

  const ActionListTile(
      {super.key, required this.actionLog, this.showDivider = true});

  @override
  Widget build(BuildContext context) {
    var lignes = [
      Align(
        alignment: Alignment.centerLeft,
        child: Text(
          "${actionLog.variete}${actionLog.lieu == "" ? "" : ' / ${actionLog.lieu}'}",
          textAlign: TextAlign.left,
          overflow: TextOverflow.clip,
          style: TextStyle(
            fontWeight: FontWeight.w400,
            fontStyle: FontStyle.normal,
            fontSize: 12,
            color: Color(0xff000000),
          ),
        ),
      )
    ];

    if (actionLog.notes != "") {
      lignes.add(Align(
        alignment: Alignment.centerLeft,
        child: Text(
          actionLog.notes,
          textAlign: TextAlign.start,
          overflow: TextOverflow.clip,
          style: TextStyle(
            fontWeight: FontWeight.w400,
            fontStyle: FontStyle.normal,
            fontSize: 12,
            color: Color(0xff000000),
          ),
        ),
      ));
    }

    Widget? tagLine;
    if (actionLog.tags.isNotEmpty) {
      final chips = actionLog.tags.map((s) {
        return Chip(
          elevation: 0,
          shadowColor: Colors.teal,
          label: Text(s, style: TextStyle(color: Colors.blue[900])),
        );
      }).toList();
      tagLine = Padding(
          padding: const EdgeInsets.only(left: 10, right: 10),
          child: Align(
              alignment: Alignment.topRight,
              child: Wrap(spacing: 5, runSpacing: 5, children: chips)));
    }

    var tile = Column(
      mainAxisAlignment: MainAxisAlignment.start,
      crossAxisAlignment: CrossAxisAlignment.center,
      mainAxisSize: MainAxisSize.max,
      children: [
        Padding(
          padding: EdgeInsets.fromLTRB(10, 0, 0, 0),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.start,
            crossAxisAlignment: CrossAxisAlignment.center,
            mainAxisSize: MainAxisSize.max,
            children: [
              Expanded(
                flex: 1,
                child: Text(
                  "${actionLog.action} - ${actionLog.legume} "
                  "${actionLog.poids > 0 ? '(${actionLog.poids}g)' : actionLog.qte > 0 ? '(${actionLog.qte})' : ''}",
                  textAlign: TextAlign.start,
                  overflow: TextOverflow.clip,
                  style: TextStyle(
                    fontWeight: FontWeight.w400,
                    fontStyle: FontStyle.normal,
                    fontSize: 14,
                    color: Color(0xff000000),
                  ),
                ),
              ),
            ],
          ),
        ),
        Padding(
          padding: EdgeInsets.fromLTRB(10, 0, 0, 0),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.start,
            crossAxisAlignment: CrossAxisAlignment.center,
            mainAxisSize: MainAxisSize.max,
            children: [
              Expanded(
                flex: 1,
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.start,
                  crossAxisAlignment: CrossAxisAlignment.center,
                  mainAxisSize: MainAxisSize.max,
                  children: lignes,
                ),
              ),
            ],
          ),
        ),
      ],
    );

    if (tagLine != null) tile.children.add(tagLine);
    if (actionLog.photos.isNotEmpty) {
      tile.children.add(HorizontalImageListview(imgUrlList: actionLog.photos));
    }
    if (showDivider) {
      tile.children.add(Divider(color: Color(0xff808080), height: 1));
    }

    return tile;
  }
}
