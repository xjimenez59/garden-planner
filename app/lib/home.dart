///File download from FlutterViz- Drag and drop a tools. For more details visit https://flutterviz.io/

///File download from FlutterViz- Drag and drop a tools. For more details visit https://flutterviz.io/

// ignore_for_file: prefer_const_constructors

import 'dart:js_interop';

import 'package:app/action_log.dart';
import 'package:app/garden_model.dart';
import 'package:app/gardens_view.dart';
import 'package:app/stats_view.dart';
import 'package:flutter/material.dart';
import 'package:app/logs_view.dart';
import 'api_service.dart';
import 'utils.dart';
import 'action_detail.dart';

class MyHomePage extends StatefulWidget {
  const MyHomePage({super.key, required this.title});

  // This widget is the home page of your application. It is stateful, meaning
  // that it has a State object (defined below) that contains fields that affect
  // how it looks.

  // This class is the configuration for the state. It holds the values (in this
  // case the title) provided by the parent (in this case the App widget) and
  // used by the build method of the State. Fields in a Widget subclass are
  // always marked "final".

  final String title;

  @override
  State<MyHomePage> createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> {
  late List<ActionLog> actionLogs = [];
  late List<Garden> jardins = [];
  TextEditingController filterController = TextEditingController();
  int currentPage = 0;
  Garden? selectedGarden;
  @override
  void initState() {
    super.initState();
    _getData();
  }

  @override
  void dispose() {
    filterController.dispose();
    super.dispose();
  }

  void _getData() async {
    jardins = (await ApiService().getGardens())!;
    if (jardins.isNotEmpty) {
      selectedGarden = jardins.first;
      actionLogs = (await ApiService().getLogs(selectedGarden!))!;
    }
    Future.delayed(const Duration(seconds: 1)).then((value) => setState(() {}));
  }

  Future<void> _onRefresh() async {
    _getData();
  }

  void _onFilterChanged(String text) {
    setState(
      () {},
    );
  }

  @override
  Widget build(BuildContext context) {
    List<ActionLog>? filteredActionLogs;
    DateTime lastDate = DateTime(1965);
    if (filterController.text == "" || actionLogs.isEmpty) {
      filteredActionLogs = actionLogs;
    } else {
      filteredActionLogs = actionLogs
          .where((a) => a.toFilterableString().contains(
              filterController.text.withoutDiacriticalMarks.toLowerCase()))
          .toList();
      filteredActionLogs.sort((a, b) => b.dateAction.compareTo(a.dateAction));
    }

    return Scaffold(
      appBar: AppBar(
        // Here we take the value from the MyHomePage object that was created by
        // the App.build method, and use it to set our appbar title.
        title: Text(selectedGarden.isNull
            ? widget.title
            : "Garden Planner - ${selectedGarden!.Nom}"),
      ),
      bottomNavigationBar: NavigationBar(
        destinations: const [
          NavigationDestination(icon: Icon(Icons.home), label: "Home"),
          NavigationDestination(icon: Icon(Icons.add_chart), label: "Récoltes")
        ],
        onDestinationSelected: (int i) {
          setState(() {
            currentPage = i;
          });
        },
        selectedIndex: currentPage,
      ),
      body: RefreshIndicator(
          onRefresh: _onRefresh,
          child: currentPage == 0
              ? Center(
                  // Center is a layout widget. It takes a single child and positions it
                  // in the middle of the parent.
                  child: ListView.builder(
                  itemCount: filteredActionLogs.length + 1,
                  itemBuilder: (context, index) {
                    List<Widget> results = [];

                    if (index == 0) {
                      results.add(TopHomeFilter(
                          jardins: jardins,
                          filterController: filterController,
                          onFilterChanged: _onFilterChanged,
                          onSelectGardenTap: onSelectGardenTap));
                    } else {
                      var a = filteredActionLogs![index - 1];

                      if (!(lastDate.sameDayAs(a.dateAction))) {
                        lastDate = a.dateAction;
                        results.add(DaySeparator(date: a.dateAction, icon: ""));
                      }
                      bool showDivider = (index == filteredActionLogs.length) ||
                          (filteredActionLogs[index]
                              .dateAction
                              .sameDayAs(a.dateAction));
                      results.add(Dismissible(
                          key: Key(a.id),
                          onDismissed: onTileDismissed(index - 1),
                          background: Container(
                              color: Colors.red.shade100,
                              margin: EdgeInsets.only(bottom: 0)),
                          child: InkWell(
                            onTap: onTileTap(a),
                            child: ActionListTile(
                                actionLog: a, showDivider: showDivider),
                          )));
                    }

                    if (results.length > 1) {
                      return Column(children: results);
                    }
                    return results.first;
                  },
                ))
              : currentPage == 1
                  ? ActionLogStats()
                  : Container()),
      floatingActionButton: (currentPage == 0 && !selectedGarden.isNull)
          ? FloatingActionButton(
              onPressed: onNewActionLogTap,
              tooltip: 'Ajouter une action',
              child: const Icon(Icons.add),
            )
          : null,
    );
  }

  void onNewActionLogTap() async {
    DateTime today = DateTime.now();
    today = DateTime(today.year, today.month, today.day);
    ActionLog a = ActionLog(dateAction: today, jardinId: selectedGarden!.ID);
    ActionLog result = await Navigator.push(context,
        MaterialPageRoute(builder: (context) => ActionDetail(actionLog: a)));
    setState(() {
      //-- result est une copie de a. On doit donc "recharger" a avec les valeurs modifiées.
      a.updateFrom(result);
      actionLogs.add(a);
      actionLogs.sort((a, b) => b.dateAction.compareTo(a.dateAction));
    });
  }

  void Function() onTileTap(ActionLog a) {
    return () async {
      ActionLog? result = await Navigator.push(context,
          MaterialPageRoute(builder: (context) => ActionDetail(actionLog: a)));
      if (result != null) {
        setState(() {
          //-- result est une copie de a. On doit donc "recharger" a avec les valeurs modifiées.
          a.updateFrom(result);
        });
      }
    };
  }

  void onSelectGardenTap() async {
    Garden? result = await Navigator.push(
        context,
        MaterialPageRoute(
            builder: (context) => GardensView(
                  gardens: jardins,
                  activeGarden: selectedGarden,
                )));
    if (result != null) {
      selectedGarden = result;

      actionLogs = (await ApiService().getLogs(selectedGarden!))!;
    }
    setState(() {});
  }

  void Function(DismissDirection direction) onTileDismissed(int index) {
    return (DismissDirection direction) async {
      bool deleted = await ApiService().deleteLog(actionLogs[index].id);
      if (deleted) {
        setState(() {
          actionLogs.removeAt(index);
        });
      }
    };
  }
}
