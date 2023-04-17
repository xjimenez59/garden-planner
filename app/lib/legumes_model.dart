// To parse this JSON data, do
//
//     final actionLog = actionLogFromJson(jsonString);

import 'dart:convert';

List<Legume> LegumeFromJson(String str) =>
    List<Legume>.from(json.decode(str).map((x) => Legume.fromJson(x)));

String LegumeToJson(Legume data) => json.encode(data.toJson());

class Legume {
  Legume(
      {required this.id,
      required this.nom,
      required this.famille,
      required this.notes,
      required this.varietes});

  String id;
  String nom;
  String famille;
  String notes;
  List<String> varietes;

  factory Legume.fromJson(Map<String, dynamic> json) => Legume(
        id: json["_id"],
        nom: json["nom"],
        famille: json["famille"],
        notes: json["notes"],
        varietes: List<String>.from(json["variete"].map((x) => x)),
      );

  Map<String, dynamic> toJson() => {
        "_id": id,
        "nom": nom,
        "famille": famille,
        "notes": notes,
        "variete": List<dynamic>.from(varietes.map((x) => x)),
      };
}
