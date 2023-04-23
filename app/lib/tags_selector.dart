import 'dart:async';
import 'dart:io';

import 'package:flutter/material.dart';
import 'package:http/http.dart';

class TagsSelector extends StatefulWidget {
  List<String> optionsList = [];
  final List<String> value;
  final String title;
  final Future<List<String>> Function(dynamic param) getOptions;
  final dynamic optionsParam;

  TagsSelector(
      {super.key,
      required this.title,
      required this.value,
      required this.getOptions,
      this.optionsParam});

  @override
  State<StatefulWidget> createState() {
    // ignore: no_logic_in_create_state
    return _TagsSelector(
        selectedTags: List.from(value), getOptions: getOptions);
  }
}

class _TagsSelector extends State<TagsSelector> {
  List<String> optionsList = [];
  List<String> selectedTags = ["ceci", "cela", "et tout ça"];
  Future<List<String>> Function(dynamic param) getOptions;
  TextEditingController editController = TextEditingController();
  TextEditingController tagsController = TextEditingController();
  late FocusNode editFocusNode;

  _TagsSelector({this.selectedTags = const [], required this.getOptions});

  @override
  void initState() {
    editController.text = "";
    editFocusNode = FocusNode();
    getOptions(widget.optionsParam).then((result) => setState(() {
          widget.optionsList = result;
          optionsList = result;
          filterOptions();
        }));
    super.initState();
  }

  @override
  void dispose() {
    // Clean up the focus node when the Form is disposed.
    editFocusNode.dispose();
    editController.dispose();
    tagsController.dispose();
    super.dispose();
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
                    selected: selectedTags.contains(optionsList[index]),
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
        onChanged: editOnChanged,
        onSubmitted: editOnSubmit,
        autofocus: true,
        focusNode: editFocusNode,
      ),
      Padding(
        padding: const EdgeInsets.only(left: 10, right: 10),
        child: Wrap(
          spacing: 5,
          runSpacing: 5,
          children: selectedTags.map(
            (s) {
              return Chip(
                elevation: 0,
                shadowColor: Colors.teal,
                // pressElevation: 0,
                // backgroundColor: Colors.blue[100],
                // shape: RoundedRectangleBorder(
                //   borderRadius: BorderRadius.circular(7),
                // ),
                label: Text(s, style: TextStyle(color: Colors.blue[900])),
                onDeleted: () {
                  setState(
                    () {
                      selectedTags.remove(s);
                      filterOptions();
                    },
                  );
                },
              );
            },
          ).toList(),
        ),
      ),
    ]);

    return Scaffold(
        appBar: AppBar(title: Text(widget.title)),
        body: Center(child: content));
  }

  void editOnChanged(String text) {
    filterOptions();
  }

  void filterOptions() {
    setState(() {
      optionsList = widget.optionsList
          .where((element) => selectedTags.contains(element) == false)
          .toList();
      if (editController.text == "" || editController.text == widget.value) {
        //--- rien de plus
      } else {
        optionsList = optionsList
            .where((element) => element
                .toLowerCase()
                .contains(editController.text.toLowerCase()))
            .toList();
      }
    });
  }

  void optionOnTap(int index) {
    String newValue = optionsList[index];
    editController.text = "";
    addTag(newValue);
  }

  void editOnSubmit(String text) {
    if (text == "") {
      Navigator.pop(context, selectedTags);
    } else {
      editController.clear();
      addTag(text);
      editFocusNode.requestFocus();
    }
  }

  void addTag(String value) {
    setState(() {
      selectedTags.add(value);
      filterOptions();
    });
  }
}
