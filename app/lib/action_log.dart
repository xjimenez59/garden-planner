// To parse this JSON data, do
//
//     final actionLog = actionLogFromJson(jsonString);

import 'dart:convert';

//ActionLog actionLogFromJson(String str) => ActionLog.fromJson(json.decode(str));

List<ActionLog> actionLogFromJson(String str) =>
    List<ActionLog>.from(json.decode(str).map((x) => ActionLog.fromJson(x)));

String actionLogToJson(ActionLog data) => json.encode(data.toJson());

class ActionLog {
  ActionLog({
    required this.id,
    required this.jardin,
    required this.dateAction,
    required this.action,
    required this.statut,
    required this.lieu,
    required this.lot,
    required this.legume,
    required this.variete,
    required this.qte,
    required this.poids,
    required this.notes,
    required this.photos,
  });

  String id;
  String jardin;
  DateTime dateAction;
  String action;
  String statut;
  String lieu;
  String lot;
  String legume;
  String variete;
  double qte;
  double poids;
  String notes;
  List<String> photos;

  factory ActionLog.fromJson(Map<String, dynamic> json) => ActionLog(
        id: json["id"],
        jardin: json["jardin"],
        dateAction: DateTime.parse(json["dateAction"]),
        action: json["action"],
        statut: json["statut"],
        lieu: json["lieu"],
        lot: json["Lot"],
        legume: json["legume"],
        variete: json["variete"],
        qte: json["qte"]?.toDouble(),
        poids: json["poids"]?.toDouble(),
        notes: json["notes"],
        photos: List<String>.from(json["photos"].map((x) => x)),
      );

  Map<String, dynamic> toJson() => {
        "id": id,
        "jardin": jardin,
        "dateAction":
            "${dateAction.year.toString().padLeft(4, '0')}-${dateAction.month.toString().padLeft(2, '0')}-${dateAction.day.toString().padLeft(2, '0')}",
        "action": action,
        "statut": statut,
        "lieu": lieu,
        "Lot": lot,
        "legume": legume,
        "variete": variete,
        "qte": qte,
        "poids": poids,
        "notes": notes,
        "photos": List<dynamic>.from(photos.map((x) => x)),
      };
}
