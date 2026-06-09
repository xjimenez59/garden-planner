// ignore_for_file: prefer_const_constructors
import 'dart:math';
import 'package:flutter/material.dart';
import 'package:app/meteo_service.dart';

// Mapping WMO weathercode → icône Material (public pour logs_view)
IconData wmoIconData(int code) {
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

Widget _windArrow(int degrees, double speed) {
  final rad = (degrees - 180) * pi / 180;
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

const double _slotW = 64.0;

const _joursSemaine = ['Lun', 'Mar', 'Mer', 'Jeu', 'Ven', 'Sam', 'Dim'];

String _labelJour(String dateStr) {
  final d = DateTime.parse(dateStr);
  return '${_joursSemaine[d.weekday - 1]}. ${d.day}';
}

List<MapEntry<String, List<HourlyForecast>>> _groupByDay(
    List<HourlyForecast> forecasts) {
  final map = <String, List<HourlyForecast>>{};
  for (final f in forecasts) {
    final key = f.time.length >= 10 ? f.time.substring(0, 10) : f.time;
    map.putIfAbsent(key, () => []).add(f);
  }
  return map.entries.toList();
}

String _luneEmoji(LuneDay? lune) {
  switch (lune?.jourBiodynamique) {
    case 'jour_fruit':   return '🍓';
    case 'jour_racine':  return '🥕';
    case 'jour_fleur':   return '🌸';
    case 'jour_feuille': return '🌿';
    default: return '🌙';
  }
}

Widget _slotHeader(HourlyForecast f) {
  return SizedBox(
    width: _slotW,
    child: Padding(
      padding: EdgeInsets.symmetric(horizontal: 4).copyWith(top: 4),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Text(f.heure,
              style: TextStyle(fontSize: 11, fontWeight: FontWeight.bold)),
          SizedBox(height: 3),
          Text('${f.temperature.round()}°',
              style: TextStyle(
                  fontSize: 13,
                  fontWeight: FontWeight.w600,
                  color: _tempColor(f.temperature))),
          SizedBox(height: 3),
          Icon(wmoIconData(f.weatherCode),
              size: 20, color: Colors.blueGrey.shade600),
          SizedBox(height: 3),
        ],
      ),
    ),
  );
}

Widget _slotFooter(HourlyForecast f) {
  return SizedBox(
    width: _slotW,
    child: Padding(
      padding: EdgeInsets.symmetric(horizontal: 4).copyWith(bottom: 6),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          SizedBox(height: 4),
          _windArrow(f.windDir, f.windSpeed),
        ],
      ),
    ),
  );
}

class _PrecipAreaPainter extends CustomPainter {
  final List<double> values;

  _PrecipAreaPainter({required this.values});

  @override
  void paint(Canvas canvas, Size size) {
    final n = values.length;
    if (n == 0) return;

    final maxVal = values.reduce(max);

    const topMargin = 6.0;
    double yFor(double v) {
      if (maxVal == 0) return size.height;
      return size.height - (v / maxVal * (size.height - topMargin));
    }

    double xFor(int i) => _slotW * i + _slotW / 2;

    final fillPath = Path();
    final strokePath = Path();

    fillPath.moveTo(xFor(0), size.height);
    fillPath.lineTo(xFor(0), yFor(values[0]));
    strokePath.moveTo(xFor(0), yFor(values[0]));

    for (int i = 0; i < n - 1; i++) {
      final x0 = xFor(i);
      final y0 = yFor(values[i]);
      final x1 = xFor(i + 1);
      final y1 = yFor(values[i + 1]);
      final cp1x = x0 + (x1 - x0) * 0.4;
      final cp2x = x1 - (x1 - x0) * 0.4;
      fillPath.cubicTo(cp1x, y0, cp2x, y1, x1, y1);
      strokePath.cubicTo(cp1x, y0, cp2x, y1, x1, y1);
    }

    fillPath.lineTo(xFor(n - 1), size.height);
    fillPath.close();

    canvas.drawPath(
      fillPath,
      Paint()
        ..color = Colors.blue.shade200.withValues(alpha: 0.55)
        ..style = PaintingStyle.fill,
    );
    canvas.drawPath(
      strokePath,
      Paint()
        ..color = Colors.blue.shade400
        ..style = PaintingStyle.stroke
        ..strokeWidth = 1.5
        ..strokeCap = StrokeCap.round
        ..strokeJoin = StrokeJoin.round,
    );

    if (maxVal > 0) {
      final peakIdx = values.indexWhere((v) => v == maxVal);
      final px = xFor(peakIdx);
      final py = yFor(maxVal);
      final label = maxVal >= 1
          ? '${maxVal.round()}mm'
          : '${maxVal.toStringAsFixed(1)}mm';
      final tp = TextPainter(
        text: TextSpan(
          text: label,
          style: TextStyle(
              fontSize: 9,
              color: Colors.blue.shade700,
              fontWeight: FontWeight.w600),
        ),
        textDirection: TextDirection.ltr,
      )..layout();
      tp.paint(
        canvas,
        Offset(
          (px - tp.width / 2).clamp(2, size.width - tp.width - 2),
          (py - tp.height - 2).clamp(0, size.height - tp.height),
        ),
      );
    }
  }

  @override
  bool shouldRepaint(_PrecipAreaPainter old) => old.values != values;
}

Widget _buildDaySection(
    String dateStr, List<HourlyForecast> slots, bool isToday, LuneDay? lune) {
  final dayWidth = _slotW * slots.length;
  final label = isToday
      ? 'Aujourd\'hui ${DateTime.parse(dateStr).day}'
      : _labelJour(dateStr);

  return Container(
    decoration: BoxDecoration(
      color: isToday ? Color(0xFFE8F0F8) : null,
      border: Border(right: BorderSide(color: Colors.grey.shade300, width: 1)),
    ),
    child: Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Container(
          width: dayWidth,
          padding: EdgeInsets.symmetric(horizontal: 6, vertical: 3),
          decoration: BoxDecoration(
            border: Border(bottom: BorderSide(color: Colors.grey.shade200)),
          ),
          child: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              Text(_luneEmoji(lune), style: TextStyle(fontSize: 12)),
              SizedBox(width: 4),
              Text(
                label,
                style: TextStyle(
                  fontSize: 11,
                  fontWeight: isToday ? FontWeight.bold : FontWeight.w500,
                  color: isToday ? Colors.blue.shade700 : Colors.grey.shade600,
                ),
              ),
            ],
          ),
        ),
        Row(children: slots.map(_slotHeader).toList()),
        SizedBox(
          height: 40,
          width: dayWidth,
          child: CustomPaint(
            painter: _PrecipAreaPainter(
                values: slots.map((f) => f.precipitation).toList()),
          ),
        ),
        Row(children: slots.map(_slotFooter).toList()),
      ],
    ),
  );
}

/// Frise scrollable des prévisions multi-jours (vue dépliée uniquement).
/// Utilisée dans DaySeparator pour le séparateur "Aujourd'hui".
class PrevisionsFriseWidget extends StatelessWidget {
  final List<HourlyForecast> previsions;
  final List<LuneDay> luneData;

  const PrevisionsFriseWidget({
    super.key,
    required this.previsions,
    this.luneData = const [],
  });

  LuneDay? _luneForDate(String dateStr) {
    final d = DateTime.parse(dateStr);
    return luneData
        .where((l) =>
            l.date.year == d.year &&
            l.date.month == d.month &&
            l.date.day == d.day)
        .firstOrNull;
  }

  @override
  Widget build(BuildContext context) {
    final days = _groupByDay(previsions);
    final todayStr = DateTime.now().toIso8601String().substring(0, 10);
    return SingleChildScrollView(
      scrollDirection: Axis.horizontal,
      padding: EdgeInsets.only(bottom: 6),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: days
            .map((e) => _buildDaySection(
                  e.key,
                  e.value,
                  e.key == todayStr,
                  _luneForDate(e.key),
                ))
            .toList(),
      ),
    );
  }
}
