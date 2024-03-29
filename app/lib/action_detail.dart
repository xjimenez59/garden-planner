// ignore_for_file: prefer_const_constructors

import 'dart:convert';

import 'package:app/action_log.dart';
import 'package:app/take_picture.dart';
import 'package:camera/camera.dart';
import 'package:flutter/material.dart';

import 'api_service.dart';
import 'legumes_model.dart';
import 'list_selector.dart';
import 'tags_selector.dart';
import 'utils.dart';

class ActionDetail extends StatefulWidget {
  final ActionLog actionLog;

  const ActionDetail({super.key, required this.actionLog});

  @override
  // ignore: no_logic_in_create_state
  State<StatefulWidget> createState() {
    ActionLog clone =
        ActionLog.fromJson(jsonDecode(jsonEncode(actionLog.toJson())));
    return _ActionDetail(actionLog: clone);
  }
}

class _ActionDetail extends State<ActionDetail> {
  TextEditingController actionInput = TextEditingController();
  TextEditingController dateInput = TextEditingController();
  TextEditingController legumeInput = TextEditingController();
  TextEditingController varieteInput = TextEditingController();
  TextEditingController quantiteInput = TextEditingController();
  TextEditingController poidsInput = TextEditingController();
  TextEditingController notesInput = TextEditingController();
  TextEditingController lieuInput = TextEditingController();

  ActionLog actionLog;

  _ActionDetail({required this.actionLog});

  @override
  void initState() {
    actionInput.text = actionLog.action;
    dateInput.text = dateFormat(actionLog.dateAction);
    legumeInput.text = actionLog.legume;
    varieteInput.text = actionLog.variete;
    quantiteInput.text = (actionLog.qte > 0) ? actionLog.qte.toString() : "";
    poidsInput.text = (actionLog.poids > 0) ? actionLog.poids.toString() : "";
    notesInput.text = actionLog.notes;
    lieuInput.text = actionLog.lieu;

    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    final content = Column(children: [
      TextField(
        //-- Action
        controller: actionInput,
        decoration: const InputDecoration(
            suffixIcon: Icon(Icons.chevron_right),
            labelText: "Action" //label text of field
            ),
        readOnly: true,
        onTap: onActionTap,
      ),
      TextField(
        //-- date
        controller: dateInput,
        decoration: const InputDecoration(
            suffixIcon: Icon(Icons.calendar_today), //icon of text field
            labelText: "date de l'action" //label text of field
            ),
        readOnly: true,
        onTap: onCalendarTap,
      ),
      TextField(
        //-- lieu
        controller: lieuInput,
        decoration: const InputDecoration(
            suffixIcon: Icon(Icons.chevron_right),
            labelText: "Lieu" //label text of field
            ),
        readOnly: true,
        onTap: onLieuTap,
      ),
      TextField(
        //-- légume
        controller: legumeInput,
        decoration: const InputDecoration(
            suffixIcon: Icon(Icons.chevron_right),
            labelText: "Légume" //label text of field
            ),
        readOnly: true,
        onTap: onLegumeTap,
      ),
      TextField(
        //-- variété
        controller: varieteInput,
        decoration: const InputDecoration(
            suffixIcon: Icon(Icons.chevron_right),
            labelText: "Variété" //label text of field
            ),
        readOnly: true,
        onTap: onVarieteTap,
      ),
      TextField(
        //-- Quantité
        controller: quantiteInput,
        keyboardType: TextInputType.number,
        decoration:
            const InputDecoration(labelText: "quantité" //label text of field
                ),

        onChanged: (value) {
          int intVal = int.parse(value);
          if (widget.actionLog.qte != intVal) {
            setState(() {
              actionLog.qte = intVal;
              actionLog.isModified = true;
            });
          }
        },
      ),
      TextField(
        //-- Poids
        controller: poidsInput,
        keyboardType: TextInputType.number,
        decoration: const InputDecoration(
            labelText: "poids (en grammes)" //label text of field
            ),
        onChanged: (value) {
          int intVal = int.parse(value);
          if (widget.actionLog.poids != intVal) {
            setState(() {
              actionLog.poids = intVal;
              actionLog.isModified = true;
            });
          }
        },
      ),
      InkWell(
          onTap: onTagsTap,
          child: Padding(
              padding: EdgeInsets.only(top: 10, right: 10),
              child: Column(children: [
                Align(
                    alignment: Alignment.topLeft,
                    child: Text(
                      "Etiquettes",
                      style:
                          TextStyle(fontSize: 12, color: Colors.grey.shade600),
                    )),
                Row(children: [
                  Expanded(
                      child: Container(
                    padding: const EdgeInsets.only(left: 0, top: 5, right: 10),
                    child: actionLog.tags.isEmpty
                        ? null
                        : Wrap(
                            spacing: 5,
                            runSpacing: 5,
                            children: actionLog.tags.map(
                              (s) {
                                return Chip(
                                  elevation: 0,
                                  shadowColor: Colors.teal,
                                  // pressElevation: 0,
                                  // backgroundColor: Colors.blue[100],
                                  // shape: RoundedRectangleBorder(
                                  //   borderRadius: BorderRadius.circular(7),
                                  // ),
                                  label: Text(s,
                                      style:
                                          TextStyle(color: Colors.blue[900])),
                                );
                              },
                            ).toList(),
                          ),
                  )),
                  Icon(
                    Icons.chevron_right,
                    color: Colors.grey.shade600,
                  )
                ])
              ]))),
      TextField(
        //-- Notes
        controller: notesInput,
        minLines: 3,
        maxLines: 5,
        decoration: InputDecoration(labelText: "Notes" //label text of field
            ),
        onChanged: (value) {
          actionLog.notes = value;
          if (value != widget.actionLog.notes) {
            setState(() {
              actionLog.isModified = true;
            });
          }
        },
      ),
      Row(children: [
        // image count in the list
        // or length of the list
        Text("${actionLog.photos.length} photos"),
        SizedBox(width: 45),
        // icon button to add
        // the image to the list
        ElevatedButton.icon(
          icon: Icon(Icons.add_a_photo),
          label: Text("Ajouter une photo"),
          // when pressed call
          // the add method
          onPressed: onPictureTap,
        )
      ]),
      for (var img in actionLog.photos)
        Dismissible(
            key: Key(img),
            onDismissed: onPhotoDismissed(img),
            background: Container(
                color: Colors.red.shade100, margin: EdgeInsets.only(bottom: 0)),
            child: Padding(
              padding: EdgeInsets.all(12),
              child: Image.network(img),
            ))
    ]);

    String pageTitle =
        "${actionLog.action} du ${actionLog.dateAction.day} ${monthNames[actionLog.dateAction.month]}";
    return WillPopScope(
        onWillPop: onWillPop,
        child: Scaffold(
            appBar: AppBar(
              title: Text(pageTitle),
            ),
            floatingActionButton: actionLog.isModified == false
                ? Container()
                : FloatingActionButton(
                    onPressed: onSaveTap,
                    tooltip: 'Enregistrer',
                    child: const Icon(Icons.save),
                  ),
            body: Align(
                alignment: Alignment.topCenter,
                child: Padding(
                    padding: EdgeInsets.fromLTRB(10, 0, 5, 0),
                    child: SingleChildScrollView(
                      child: content,
                    )))));
  }

  void onSaveTap() async {
    String result = await _postLog(actionLog);
    actionLog.id = result;
    actionLog.isModified = false;
    Navigator.pop(context, actionLog);
  }

  Future<bool> onWillPop() async {
    final oldPhotos = widget.actionLog.photos;
    final newPhotos = actionLog.photos;

    // on suprime les photos si l'utilisateur n'enregistre pas ..
    for (var element
        in newPhotos.where((element) => oldPhotos.contains(element) == false)) {
      await ApiService().deletePicture(element);
    }

    return true;
  }

  void onTagsTap() async {
    final result = await Navigator.push(
        context,
        MaterialPageRoute(
            builder: (context) => TagsSelector(
                  title: "Etiquettes",
                  value: actionLog.tags,
                  getOptions: _getTags,
                )));
    setState(() {
      var diff = (result as List<String>)
          .toSet()
          .difference(widget.actionLog.tags.toSet());
      if (diff.isNotEmpty) {
        actionLog.isModified = true;
      }
      actionLog.tags = result;
    });
  }

  void onActionTap() async {
    final result = await Navigator.push(
        context,
        MaterialPageRoute(
            builder: (context) => ListSelector(
                title: "Qu'avez-vous fait ?",
                value: actionLog.action,
                getOptions: _getActions)));
    setState(() {
      if (widget.actionLog.action != result) {
        actionLog.isModified = true;
      }
      actionLog.action = result;
      actionInput.text = result;
    });
  }

  void onCalendarTap() async {
    DateTime? pickedDate = await showDatePicker(
        context: context,
        initialDate: actionLog.dateAction,
        firstDate: DateTime(1950),
        //DateTime.now() - not to allow to choose before today.
        lastDate: DateTime(2100));

    if (pickedDate != null) {
      if (widget.actionLog.dateAction != pickedDate) {
        actionLog.isModified = true;
      }
      actionLog.dateAction = pickedDate;
      String formattedDate = dateFormat(pickedDate);
      setState(() {
        dateInput.text = formattedDate; //set output date to TextField value.
      });
    } else {}
  }

  void onLieuTap() async {
    final result = await Navigator.push(
        context,
        MaterialPageRoute(
            builder: (context) => ListSelector(
                title: "Emplacement / Lieu",
                value: actionLog.lieu,
                getOptions: _getLieux)));
    setState(() {
      if (widget.actionLog.lieu != result) {
        actionLog.isModified = true;
        actionLog.lieu = result;
        lieuInput.text = result;
      }
    });
  }

  void onPictureTap() async {
    // Ensure that plugin services are initialized
    WidgetsFlutterBinding.ensureInitialized();

    await availableCameras().then((value) async {
      final result = await Navigator.push(
          context,
          MaterialPageRoute(
              builder: (context) => TakePictureScreen(cameras: value)));
      setState(() {
        actionLog.isModified = true;
        actionLog.photos.add(result);
      });
    });
  }

  void onLegumeTap() async {
    final result = await Navigator.push(
        context,
        MaterialPageRoute(
            builder: (context) => ListSelector(
                title: "Choisissez un légume",
                value: actionLog.legume,
                getOptions: _getLegumes)));
    setState(() {
      if (widget.actionLog.legume != result) {
        actionLog.isModified = true;
        actionLog.legume = result;
        legumeInput.text = result;

        actionLog.variete = "";
        varieteInput.text = "";
      }
    });
  }

  void onVarieteTap() async {
    final result = await Navigator.push(
        context,
        MaterialPageRoute(
            builder: (context) => ListSelector(
                title: "Variété de ${actionLog.legume}",
                value: actionLog.variete,
                optionsParam: actionLog,
                getOptions: _getVarietes)));
    setState(() {
      if (widget.actionLog.variete != result) {
        actionLog.isModified = true;
      }
      actionLog.variete = result;
      varieteInput.text = result;
    });
  }

  Future<List<String>> _getActions(dynamic param) async {
    List<String> result = [
      "Semis",
      "Repiquage",
      "Récolte",
      "Semis pleine terre",
      "Plantation",
      "Déparasitage",
      "Photo / Notes",
      "Don"
    ];
    return result;
  }

  // Fonction asynchrone qui sera appelée par ListSelector pour afficher la liste des légumes
  Future<List<String>> _getLegumes(dynamic param) async {
    List<Legume>? listLegumes = (await ApiService().getLegumes());
    List<String> result = [];
    listLegumes?.forEach((element) {
      result.add(element.nom);
    });
    return result;
  }

  Future<List<String>> _getVarietes(dynamic param) async {
    List<String> result = [];
    List<Legume>? listLegumes = (await ApiService().getLegumes());
    if (listLegumes != null) {
      result = listLegumes
          .firstWhere((element) => element.nom == (param as ActionLog).legume)
          .varietes;
    }

    return result;
  }

  Future<List<String>> _getLieux(dynamic param) async {
    List<String>? listLieux = (await ApiService().getLieux());
    List<String> result = [];
    listLieux?.forEach((element) {
      result.add(element);
    });
    return result;
  }

  Future<List<String>> _getTags(dynamic param) async {
    List<String>? result;
    result = (await ApiService().getTags());
    result ??= [];

    return result;
  }

  Future<int> _postLogs(List<ActionLog> logs) async {
    int? result;
    result = (await ApiService().postLogs(logs));
    result ??= 0;

    return result;
  }

  Future<String> _postLog(ActionLog a) async {
    String? result;
    result = (await ApiService().postLog(a));
    if (result == null) {
      result = "";
    } else {}

    return result;
  }

  void Function(DismissDirection direction) onPhotoDismissed(String img) {
    return (DismissDirection direction) async {
      bool deleted = await ApiService().deletePicture(img);
      if (deleted) {
        setState(() {
          actionLog.isModified = true;
          actionLog.photos.remove(img);
        });
      }
    };
  }
}
