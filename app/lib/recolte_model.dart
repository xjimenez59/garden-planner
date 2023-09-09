import 'dart:convert';

class RecolteAnnee {
  int annee;
  int poids;
  int qte;

  RecolteAnnee({required this.annee, required this.poids, required this.qte});

  factory RecolteAnnee.fromJson(Map<String, dynamic> json) => RecolteAnnee(
      annee: json["Annee"], poids: json["Poids"], qte: json["Qte"]);
}

class Recolte {
  String legume;
  List<RecolteAnnee> annees;

  Recolte({required this.legume, required this.annees});

  factory Recolte.fromJson(Map<String, dynamic> json) {
    var listeDynamic = json["Annees"] as List;
    List<RecolteAnnee> annees =
        listeDynamic.map((a) => RecolteAnnee.fromJson(a)).toList();
    return Recolte(legume: json["Legume"], annees: annees);
  }
}

List<Recolte> RecolteFromJson(String str) =>
    List<Recolte>.from(json.decode(str).map((x) => Recolte.fromJson(x)));

Map<int, int> totalRecoltes(List<Recolte> recolteJardin) {
  Map<int, int> result = {};

  for (var r in recolteJardin) {
    for (var a in r.annees) {
      result.update(a.annee, (value) => value + a.poids,
          ifAbsent: () => a.poids);
    }
  }
  return result;
}
