// ignore_for_file: prefer_const_constructors
import 'dart:math';
import 'package:flutter/material.dart';
import 'package:app/meteo_service.dart';

// Mapping WMO weathercode → icône Material
IconData _wmoIcon(int code) {
  if (code == 0) return Icons.wb_sunny;
  if (code <= 3) return Icons.cloud_queue;
  if (code <= 48) return Icons.foggy;
  if (code <= 57) return Icons.grain;
  if (code <= 67) return Icons.water_drop;
  if (code <= 77) return Icons.ac_unit;
  if (code <= 82) return Icons.thunderstorm;
  return Icons.electric_bolt;
}

Color _tempColor(double t) {
  if (t < 10) return Colors.blue.shade600;
  if (t < 20) return Colors.green.shade700;
  return Colors.orange.shade700;
}

// Flèche de direction du vent (0° = Nord, 90° = Est, etc.)
Widget _windArrow(int degrees, double speed) {
  final rad = (degrees - 180) * pi / 180; // pointe dans la direction d'où vient le vent
  return Row(
    mainAxisSize: MainAxisSize.min,
    children: [
      Transform.rotate(
        angle: rad,
        child: Icon(Icons.arrow_upward, size: 12, color: Colors.blueGrey),
      ),
      SizedBox(width: 2),
      Text('${speed.round()}',
          style: TextStyle(fontSize: 10, color: Colors.blueGrey.shade700)),
    ],
  );
}

// Un créneau de 3h dans la frise
Widget _slot(HourlyForecast f, double maxPrecip) {
  final barMaxH = 24.0;
  final barH = maxPrecip > 0 ? (f.precipitation / maxPrecip * barMaxH).clamp(2.0, barMaxH) : 0.0;

  return Container(
    width: 64,
    padding: EdgeInsets.symmetric(horizontal: 4, vertical: 6),
    child: Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        // Heure
        Text(f.heure,
            style: TextStyle(fontSize: 11, fontWeight: FontWeight.bold)),
        SizedBox(height: 4),
        // Icône météo
        Icon(_wmoIcon(f.weatherCode), size: 20, color: Colors.blueGrey.shade600),
        SizedBox(height: 4),
        // Température
        Text('${f.temperature.round()}°',
            style: TextStyle(
                fontSize: 13,
                fontWeight: FontWeight.w600,
                color: _tempColor(f.temperature))),
        SizedBox(height: 4),
        // Barre de précipitations
        Container(
          height: barMaxH,
          alignment: Alignment.bottomCenter,
          child: Container(
            width: 12,
            height: f.precipitation > 0 ? barH : 2,
            decoration: BoxDecoration(
              color: f.precipitation > 0
                  ? Colors.blue.shade300
                  : Colors.grey.shade200,
              borderRadius: BorderRadius.circular(2),
            ),
          ),
        ),
        SizedBox(height: 2),
        // Valeur précip
        Text(
          f.precipitation > 0 ? '${f.precipitation.toStringAsFixed(1)}mm' : '-',
          style: TextStyle(fontSize: 9, color: Colors.blue.shade400),
        ),
        SizedBox(height: 4),
        // Vent
        _windArrow(f.windDir, f.windSpeed),
      ],
    ),
  );
}

class PrevisionsDuJourWidget extends StatefulWidget {
  final List<HourlyForecast> previsions;
  final LuneDay? luneDay;

  const PrevisionsDuJourWidget({
    super.key,
    required this.previsions,
    this.luneDay,
  });

  @override
  State<PrevisionsDuJourWidget> createState() => _PrevisionsDuJourWidgetState();
}

class _PrevisionsDuJourWidgetState extends State<PrevisionsDuJourWidget> {
  bool _expanded = false;

  String _bioemoji() {
    switch (widget.luneDay?.jourBiodynamique) {
      case 'jour_fruit':   return '🍓';
      case 'jour_racine':  return '🥕';
      case 'jour_fleur':   return '🌸';
      case 'jour_feuille': return '🌿';
      default: return '🌙';
    }
  }

  String _biolabel() {
    final raw = widget.luneDay?.jourBiodynamique ?? '';
    if (raw.isEmpty) return '';
    final label = raw.replaceFirst('jour_', '');
    return label[0].toUpperCase() + label.substring(1);
  }

  @override
  Widget build(BuildContext context) {
    if (widget.previsions.isEmpty) return SizedBox.shrink();

    final temps = widget.previsions.map((f) => f.temperature).toList();
    final tmin = temps.reduce(min);
    final tmax = temps.reduce(max);
    final totalPrecip = widget.previsions.fold(0.0, (s, f) => s + f.precipitation);
    final avgWind = widget.previsions.fold(0.0, (s, f) => s + f.windSpeed) /
        widget.previsions.length;
    final maxPrecip = widget.previsions.fold(0.0, (s, f) => s > f.precipitation ? s : f.precipitation);

    // Direction de vent dominante (slot de midi ou milieu de journée)
    final midIdx = widget.previsions.length ~/ 2;
    final dominantWindDir = widget.previsions[midIdx].windDir;

    return GestureDetector(
      onTap: () => setState(() => _expanded = !_expanded),
      child: Container(
        decoration: BoxDecoration(
          color: Color(0xFFF0F4F8),
          border: Border(bottom: BorderSide(color: Color(0x4d9e9e9e))),
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // --- Ligne collapsed ---
            Padding(
              padding: EdgeInsets.symmetric(horizontal: 10, vertical: 6),
              child: Row(
                children: [
                  // Biodynamique
                  Text(_bioemoji(), style: TextStyle(fontSize: 16)),
                  SizedBox(width: 4),
                  Text(_biolabel(),
                      style: TextStyle(
                          fontSize: 12, color: Colors.grey.shade700)),
                  Spacer(),
                  // Températures
                  Icon(Icons.arrow_upward, size: 12, color: Colors.orange),
                  Text('${tmax.round()}°',
                      style: TextStyle(
                          fontSize: 13,
                          fontWeight: FontWeight.w600,
                          color: Colors.orange.shade700)),
                  SizedBox(width: 6),
                  Icon(Icons.arrow_downward, size: 12, color: Colors.blue),
                  Text('${tmin.round()}°',
                      style: TextStyle(
                          fontSize: 13,
                          fontWeight: FontWeight.w600,
                          color: Colors.blue.shade600)),
                  SizedBox(width: 10),
                  // Précipitations
                  Icon(Icons.water_drop, size: 14, color: Colors.blue.shade400),
                  Text(
                    totalPrecip > 0
                        ? ' ${totalPrecip.toStringAsFixed(1)}mm'
                        : ' 0mm',
                    style: TextStyle(fontSize: 12, color: Colors.blue.shade600),
                  ),
                  SizedBox(width: 10),
                  // Vent
                  Transform.rotate(
                    angle: (dominantWindDir - 180) * pi / 180,
                    child: Icon(Icons.arrow_upward,
                        size: 14, color: Colors.blueGrey),
                  ),
                  Text(' ${avgWind.round()}km/h',
                      style: TextStyle(
                          fontSize: 12, color: Colors.blueGrey.shade700)),
                  SizedBox(width: 8),
                  Icon(
                    _expanded
                        ? Icons.keyboard_arrow_up
                        : Icons.keyboard_arrow_down,
                    size: 18,
                    color: Colors.grey,
                  ),
                ],
              ),
            ),

            // --- Frise expanded ---
            if (_expanded)
              SingleChildScrollView(
                scrollDirection: Axis.horizontal,
                padding: EdgeInsets.only(bottom: 6),
                child: Row(
                  children: widget.previsions
                      .map((f) => _slot(f, maxPrecip))
                      .toList(),
                ),
              ),
          ],
        ),
      ),
    );
  }
}
