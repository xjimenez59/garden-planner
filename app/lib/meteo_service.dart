import 'dart:convert';
import 'dart:developer';

import 'package:http/http.dart' as http;

// URL du service météo, injectée via --dart-define=METEO_BASE_URL=...
// Exemple docker-compose : --dart-define=METEO_BASE_URL=https://gardenplanner.jactez.com/meteo
const _meteoBaseUrl = String.fromEnvironment(
  'METEO_BASE_URL',
  defaultValue: 'http://localhost:8082',
);

/// Relevé journalier MétéoFrance (champs utiles pour l'application jardin).
/// Un champ null signifie que la donnée n'est pas disponible pour ce jour.
class Meteo {
  final String poste;
  final DateTime date;

  // Précipitations
  final double? rr; // cumul pluie (mm)
  final double? drr; // durée précipitations (min)

  // Températures
  final double? tm; // moyenne (°C)
  final double? tn; // minimale (°C)
  final double? tx; // maximale (°C)
  final double? dg; // durée gel (min)

  // Vent
  final double? ffm; // vitesse moyenne (m/s)
  final double? fxi; // rafale max instantanée (m/s)
  final double? dxy; // direction vent max (degré)

  // Ensoleillement
  final double? inst; // durée insolation (min)
  final double? qinst; // qualité insolation
  final double? sigma; // fraction d'insolation (%)
  final double? qsigma; // qualité fraction insolation

  // Froid
  final double? nb300; // nb heures < 3 °C
  final double? qnb300; // qualité nb300

  Meteo({
    required this.poste,
    required this.date,
    this.rr,
    this.drr,
    this.tm,
    this.tn,
    this.tx,
    this.dg,
    this.ffm,
    this.fxi,
    this.dxy,
    this.inst,
    this.qinst,
    this.sigma,
    this.qsigma,
    this.nb300,
    this.qnb300,
  });

  factory Meteo.fromJson(Map<String, dynamic> j) {
    double? nd(String key) {
      final v = j[key];
      if (v == null || v == '') return null;
      return double.tryParse(v as String);
    }

    return Meteo(
      poste: j['POSTE'] as String,
      date: DateTime.parse(j['DATE'] as String),
      rr: nd('RR'),
      drr: nd('DRR'),
      tm: nd('TM'),
      tn: nd('TN'),
      tx: nd('TX'),
      dg: nd('DG'),
      ffm: nd('FFM'),
      fxi: nd('FXI'),
      dxy: nd('DXY'),
      inst: nd('INST'),
      qinst: nd('QINST'),
      sigma: nd('SIGMA'),
      qsigma: nd('QSIGMA'),
      nb300: nd('NB300'),
      qnb300: nd('QNB300'),
    );
  }
}

List<Meteo> _meteoListFromJson(String body) =>
    (jsonDecode(body) as List).map((e) => Meteo.fromJson(e)).toList();

class MeteoService {
  /// Retourne les relevés MétéoFrance pour la station [site].
  /// [dateDeb] et [dateFin] sont au format YYYYMMDD.
  /// Par défaut, retourne tous les enregistrements disponibles.
  Future<List<Meteo>> getMeteo(
    String site, {
    String dateDeb = '19000101',
    String dateFin = '21001231',
  }) async {
    try {
      final url = Uri.parse('$_meteoBaseUrl/meteo').replace(queryParameters: {
        'station': site,
        'date_deb': dateDeb,
        'date_fin': dateFin,
      });
      final response = await http.get(url);
      if (response.statusCode == 200) {
        return _meteoListFromJson(response.body);
      }
      log('MeteoService.getMeteo: HTTP ${response.statusCode}');
    } catch (e) {
      log('MeteoService.getMeteo: $e');
    }
    return [];
  }
}
