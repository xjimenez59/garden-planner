import 'package:app/api_service.dart';
import 'package:app/garden_model.dart';
import 'package:flutter/material.dart';

class GardenForm extends StatefulWidget {
  final Garden? garden;
  const GardenForm({super.key, this.garden});

  @override
  State<StatefulWidget> createState() {
    return _GardenForm();
  }
}

class _GardenForm extends State<GardenForm> {
  Garden garden = Garden(Nom: "");

  TextEditingController nomInput = TextEditingController();
  TextEditingController notesInput = TextEditingController();
  TextEditingController moisFinRecolteInput = TextEditingController();
  TextEditingController moisFinSemis = TextEditingController();
  TextEditingController localisationInput = TextEditingController();
  TextEditingController surfaceInput = TextEditingController();

  bool isModified = false;

  @override
  void initState() {
    super.initState();
    if (widget.garden != null) {
      garden = Garden.fromJson(widget.garden!.toJson()); // copie
    }

    nomInput.text = garden.Nom;
    notesInput.text = garden.Notes;
    moisFinRecolteInput.text = garden.MoisFinRecolte.toString();
    moisFinSemis.text = garden.MoisFinSemis.toString();
    localisationInput.text = garden.Localisation;
    surfaceInput.text = garden.Surface.toString();
  }

  @override
  Widget build(BuildContext context) {
    Widget content = Padding(
        padding: EdgeInsets.all(10),
        child: Column(
          children: [
            TextField(
              controller: nomInput,
              decoration: InputDecoration(labelText: "Nom du jardin"),
              onChanged: (value) {
                garden.Nom = value;
                refresh();
              },
            ),
            TextField(
              controller: localisationInput,
              decoration:
                  InputDecoration(labelText: "adresse , code postal ou ville"),
              onChanged: (value) {
                garden.Localisation = value;
                refresh();
              },
            ),
            TextField(
              controller: surfaceInput,
              keyboardType: TextInputType.number,
              decoration: InputDecoration(labelText: "Surface cultivée"),
              onChanged: (value) {
                garden.Surface = (value == "") ? 0 : int.parse(value);
                refresh();
              },
            ),
            TextField(
              controller: moisFinRecolteInput,
              keyboardType: TextInputType.number,
              decoration:
                  InputDecoration(labelText: "Mois de fin des récoltes"),
              onChanged: (value) {
                garden.MoisFinRecolte = (value == "") ? 0 : int.parse(value);
                refresh();
              },
            ),
            TextField(
              controller: moisFinSemis,
              keyboardType: TextInputType.number,
              decoration: InputDecoration(labelText: "Mois de début des semis"),
              onChanged: (value) {
                garden.MoisFinSemis = (value == "") ? 0 : int.parse(value);
                refresh();
              },
            ),
            TextField(
              controller: notesInput,
              minLines: 3,
              maxLines: 6,
              decoration: InputDecoration(labelText: "Notes"),
              onChanged: (value) {
                garden.Notes = value;
                refresh();
              },
            ),
          ],
        ));

    return Scaffold(
      appBar: AppBar(
          title: Text(
        garden.Nom == "" ? "Nouveau jardin" : garden.Nom,
      )),
      body: content,
      floatingActionButton: isModified
          ? FloatingActionButton(
              onPressed: onSaveTap,
              tooltip: 'Enregistrer',
              child: const Icon(Icons.save),
            )
          : Container(),
    );
  }

  void refresh() {
    setState(() {
      Garden g = (widget.garden == null) ? Garden(Nom: "") : widget.garden!;
      isModified = (garden.Nom != g.Nom) ||
          (garden.Notes != g.Notes) ||
          (garden.MoisFinRecolte != g.MoisFinRecolte) ||
          (garden.MoisFinSemis != g.MoisFinSemis) ||
          (garden.Localisation != g.Localisation) ||
          (garden.Surface != g.Surface);
    });
  }

  void onSaveTap() async {
    String result = await _postGarden(garden);
    garden.ID = result;
    Navigator.pop(context, garden);
  }

  Future<String> _postGarden(Garden g) async {
    String? result;
    result = (await ApiService().postGarden(g));
    if (result == null) {
      result = "";
    } else {}

    return result;
  }
}
