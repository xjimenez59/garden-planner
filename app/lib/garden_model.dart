import 'dart:convert';

class Garden {
  String ID = "";
  String Nom;
  String Notes = "";
  int MoisFinRecolte = 3;
  int MoisFinSemis = 10;
  String Localisation;
  int Surface;

  Garden(
      {this.ID = "",
      required this.Nom,
      this.Notes = "",
      this.MoisFinRecolte = 3,
      this.MoisFinSemis = 10,
      this.Localisation = "",
      this.Surface = 0});

  factory Garden.fromJson(Map<String, dynamic> json) => Garden(
      ID: json["_id"],
      Nom: json["nom"],
      Notes: json["notes"],
      MoisFinRecolte: json["moisFinRecolte"],
      MoisFinSemis: json["moisFinSemis"],
      Localisation: json["localisation"],
      Surface: json["surface"]);

  Map<String, dynamic> toJson() => {
        "_id": ID,
        "nom": Nom,
        "notes": Notes,
        "moisFinRecolte": MoisFinRecolte,
        "moisFinSemis": MoisFinSemis,
        "localisation": Localisation,
        "surface": Surface,
      };
}

List<Garden> GardenFromJson(String str) =>
    List<Garden>.from(json.decode(str).map((x) => Garden.fromJson(x)));

String GardenToJson(Garden data) => json.encode(data.toJson());
