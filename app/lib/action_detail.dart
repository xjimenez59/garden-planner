import 'package:app/action_log.dart';
import 'package:flutter/material.dart';

import 'api_service.dart';
import 'legumes_model.dart';
import 'list_selector.dart';
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
  TextEditingController dateInput = TextEditingController();
  TextEditingController legumeInput = TextEditingController();
  ActionLog actionLog;

  _ActionDetail({required this.actionLog});

  @override
  void initState() {
    dateInput.text = dateFormat(actionLog.dateAction);
    legumeInput.text = actionLog.legume;
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    final content = Column(children: [
      TextField(
        controller: dateInput,
        decoration: const InputDecoration(
            icon: Icon(Icons.calendar_today), //icon of text field
            labelText: "date de l'action" //label text of field
            ),
        readOnly: true,
        onTap: onCalendarTap,
      ),
      TextField(
        controller: legumeInput,
        decoration: const InputDecoration(
            icon: Icon(Icons.space_bar),
            suffixIcon: Icon(Icons.chevron_right),
            labelText: "Légume" //label text of field
            ),
        readOnly: true,
        onTap: onLegumeTap,
      ),
      Text(actionLog.dateAction.toString()),
      Text(actionLog.action),
      Text(actionLog.legume),
    ]);

    return Scaffold(
        appBar: AppBar(title: Text(actionLog.id)),
        body: Center(child: content));
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
                title: "Choisissez un légume :",
                value: actionLog.legume,
                getOptions: _getLegumes)));
    setState(() {
      actionLog.legume = result;
      legumeInput.text = result;
    });
  }

  Future<List<String>> _getLegumes() async {
    List<Legume>? listLegumes = (await ApiService().getLegumes());
    List<String> result = [];
    listLegumes?.forEach((element) {
      result.add(element.nom);
    });
    return result;
  }
}
