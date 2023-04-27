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
    this.id = "",
    this.parentId = "",
    this.jardin = "",
    required this.dateAction,
    this.action = "",
    this.statut = "",
    this.lieu = "",
    this.legume = "",
    this.variete = "",
    this.qte = 0,
    this.poids = 0,
    this.notes = "",
    this.photos = const [],
    this.tags = const [],
  });

  bool isModified = false;
  String id;
  String parentId;
  String jardin;
  DateTime dateAction;
  String action;
  String statut;
  String lieu;
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
        "legume": legume,
        "variete": variete,
        "qte": qte,
        "poids": poids,
        "notes": notes,
        "photos": List<dynamic>.from(photos.map((x) => x)),
        "tags": List<dynamic>.from(tags.map((x) => x)),
      };

  void updateFrom(ActionLog a) {
    isModified = a.isModified;
    id = a.id;
    parentId = a.parentId;
    jardin = a.jardin;
    dateAction = a.dateAction;
    action = a.action;
    statut = a.statut;
    lieu = a.lieu;
    legume = a.legume;
    variete = a.variete;
    qte = a.qte;
    poids = a.poids;
    notes = a.notes;
    photos = a.photos;
    tags = a.tags;
  }
}
