import 'dart:async';
import 'dart:io';

import 'package:flutter/material.dart';
import 'package:http/http.dart';

class ListSelector extends StatefulWidget {
  List<String> optionsList = [];
  final String value;
  final String title;
  final Future<List<String>> Function() getOptions;

  ListSelector(
      {super.key,
      required this.title,
      required this.value,
      required this.getOptions});

  @override
  State<StatefulWidget> createState() {
    // ignore: no_logic_in_create_state
    return _ListSelector(editedValue: value, getOptions: getOptions);
  }
}

class _ListSelector extends State<ListSelector> {
  String editedValue;
  List<String> optionsList = [];
  Future<List<String>> Function() getOptions;
  TextEditingController editController = TextEditingController();

  _ListSelector({this.editedValue = "", required this.getOptions});

  @override
  void initState() {
    editController.text = editedValue;
    getOptions().then((result) => setState(() {
          widget.optionsList = result;
          optionsList = result;
        }));
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    final content = Column(children: [
      Expanded(
          child: ListView.builder(
              itemCount: optionsList.length,
              itemBuilder: (context, index) {
                return DecoratedBox(
                  decoration: BoxDecoration(
                      color: (index % 2 == 0) ? Colors.white : Colors.grey[10],
                      border:
                          const Border(bottom: BorderSide(color: Colors.grey))),
                  child: ListTile(
                    title: Text(optionsList[index]),
                    onTap: () => optionOnTap(index),
                    selected: optionsList[index] == editedValue,
                  ),
                );
              })),
      Container(
          padding: const EdgeInsets.fromLTRB(5, 10, 0, 5),
          alignment: Alignment.bottomLeft,
          color: Colors.blue,
          child: const Text("Choisissez dans la liste, ou saisissez ci-dessous",
              textAlign: TextAlign.left,
              style: TextStyle(color: Colors.white))),
      TextField(
        controller: editController,
        onChanged: filterOptions,
        onSubmitted: editOnSubmit,
        autofocus: true,
      ),
    ]);

    return Scaffold(
        appBar: AppBar(title: Text(widget.title)),
        body: Center(child: content));
  }

  void filterOptions(String newValue) {
    setState(() {
      editedValue = newValue;
      optionsList = widget.optionsList
          .where((element) =>
              element.toLowerCase().contains(newValue.toLowerCase()))
          .toList();
    });
  }

  void optionOnTap(int index) {
    String newValue = optionsList[index];
    editController.text = newValue;
    filterOptions(newValue);
    Navigator.pop(context, newValue);
  }

  void editOnSubmit(String newValue) {
    Navigator.pop(context, newValue);
  }
}
