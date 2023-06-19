// To parse this JSON data, do
//
//     final actionLog = actionLogFromJson(jsonString);

import 'dart:convert';

List<Recolte> RecolteFromJson(String str) =>
    List<Recolte>.from(json.decode(str).map((x) => Recolte.fromJson(x)));

String RecolteToJson(Recolte data) => json.encode(data.toJson());

class Recolte {
  Recolte({
    required this.legume,
    required this.poids,
    required this.qte,
  });

  String legume;
  int poids;
  int qte;

  factory Recolte.fromJson(Map<String, dynamic> json) => Recolte(
        legume: json["Legume"],
        poids: json["Poids"],
        qte: json["Qte"],
      );

  Map<String, dynamic> toJson() => {
        "Legume": legume,
        "Poids": poids,
        "Qte": qte,
      };
}
