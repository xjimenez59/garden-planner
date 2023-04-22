// ignore_for_file: prefer_const_constructors

import 'package:app/action_log.dart';
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
  State<StatefulWidget> createState() {
    // ignore: no_logic_in_create_state
    return _ActionDetail(actionLog: actionLog);
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

  ActionLog actionLog;

  _ActionDetail({required this.actionLog});

  @override
  void initState() {
    actionInput.text = actionLog.action;
    dateInput.text = dateFormat(actionLog.dateAction);
    legumeInput.text = actionLog.legume;
    varieteInput.text = actionLog.variete;
    quantiteInput.text = actionLog.qte.toString();
    poidsInput.text = actionLog.poids.toString();
    notesInput.text = actionLog.notes;

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
        onSubmitted: (value) {
          actionLog.qte = int.parse(value);
        },
      ),
      TextField(
        //-- Poids
        controller: poidsInput,
        keyboardType: TextInputType.number,
        decoration: const InputDecoration(
            labelText: "poids (en grammes)" //label text of field
            ),
        onSubmitted: (value) {
          actionLog.poids = int.parse(value);
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
                    child: actionLog.tags.length < 1
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
        decoration:
            const InputDecoration(labelText: "Notes" //label text of field
                ),
        onSubmitted: (value) {
          actionLog.notes = value;
        },
      ),
    ]);

    String pageTitle =
        "${actionLog.action} du ${actionLog.dateAction.day} ${monthNames[actionLog.dateAction.month]}";
    return Scaffold(
        appBar: AppBar(title: Text(pageTitle)),
        body: Align(
            alignment: Alignment.topCenter,
            child: Padding(
                padding: EdgeInsets.fromLTRB(10, 0, 5, 0),
                child: SingleChildScrollView(
                  child: content,
                ))));
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
      actionLog.action = result;
      actionInput.text = result;
    });
  }

  void onCalendarTap() async {
    DateTime? pickedDate = await showDatePicker(
        context: context,
        initialDate: DateTime.now(),
        firstDate: DateTime(1950),
        //DateTime.now() - not to allow to choose before today.
        lastDate: DateTime(2100));

    if (pickedDate != null) {
      String formattedDate = dateFormat(pickedDate);
      setState(() {
        dateInput.text = formattedDate; //set output date to TextField value.
      });
    } else {}
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
      actionLog.legume = result;
      legumeInput.text = result;

      actionLog.variete = "";
      varieteInput.text = "";
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
      actionLog.variete = result;
      varieteInput.text = result;
    });
  }

  Future<List<String>> _getActions(dynamic param) async {
    List<String> result = [
      "Semis",
      "Semis pleine terre",
      "Repiquage",
      "Plantation",
      "Déparasitage",
      "Récolte",
      "Photo / Notes"
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

  Future<List<String>> _getTags(dynamic param) async {
    List<String> result = ["ceci", "cela", "et voilà"];

    return result;
  }
}