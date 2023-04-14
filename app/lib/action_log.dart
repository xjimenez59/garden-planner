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
    required this.parentId,
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
    required this.tags,
  });

  String id;
  String parentId;
  String jardin;
  DateTime dateAction;
  String action;
  String statut;
  String lieu;
  String lot;
  String legume;
  String variete;
  int qte;
  int poids;
  String notes;
  List<String> photos;
  List<String> tags;

  factory ActionLog.fromJson(Map<String, dynamic> json) => ActionLog(
        id: json["_id"],
        parentId: json["_parentId"],
        jardin: json["jardin"],
        dateAction: DateTime.parse(json["dateAction"]),
        action: json["action"],
        statut: json["statut"],
        lieu: json["lieu"],
        lot: json["lot"],
        legume: json["legume"],
        variete: json["variete"],
        qte: json["qte"],
        poids: json["poids"],
        notes: json["notes"],
        photos: List<String>.from(json["photos"].map((x) => x)),
        tags: List<String>.from(json["tags"].map((x) => x)),
      );

  Map<String, dynamic> toJson() => {
        "_id": id,
        "_parentId": parentId,
        "jardin": jardin,
        "dateAction":
            "${dateAction.year.toString().padLeft(4, '0')}-${dateAction.month.toString().padLeft(2, '0')}-${dateAction.day.toString().padLeft(2, '0')}",
        "action": action,
        "statut": statut,
        "lieu": lieu,
        "lot": lot,
        "legume": legume,
        "variete": variete,
        "qte": qte,
        "poids": poids,
        "notes": notes,
        "photos": List<dynamic>.from(photos.map((x) => x)),
        "tags": List<dynamic>.from(tags.map((x) => x)),
      };
}
