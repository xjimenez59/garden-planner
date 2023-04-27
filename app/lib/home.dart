///File download from FlutterViz- Drag and drop a tools. For more details visit https://flutterviz.io/

///File download from FlutterViz- Drag and drop a tools. For more details visit https://flutterviz.io/

// ignore_for_file: prefer_const_constructors

import 'package:app/action_log.dart';
import 'package:flutter/material.dart';

import 'api_service.dart';
import 'utils.dart';
import 'action_detail.dart';

class DaySeparator extends StatelessWidget {
  final DateTime date;
  final String icon;

  const DaySeparator({super.key, required this.date, required this.icon});

  @override
  Widget build(BuildContext context) {
    String strDate = dateFormat(date);
    int weeknum = weekNum(date);

    return Container(
      margin: EdgeInsets.fromLTRB(0, 10, 0, 0),
      padding: EdgeInsets.all(0),
      width: MediaQuery.of(context).size.width,
      //height: 40,
      decoration: BoxDecoration(
        color: Color(0x12000000),
        shape: BoxShape.rectangle,
        borderRadius: BorderRadius.zero,
        border: Border.all(color: Color(0x4d9e9e9e), width: 1),
      ),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.start,
        crossAxisAlignment: CrossAxisAlignment.center,
        mainAxisSize: MainAxisSize.max,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.start,
            crossAxisAlignment: CrossAxisAlignment.center,
            mainAxisSize: MainAxisSize.max,
            children: [
              Expanded(
                flex: 1,
                child: Padding(
                  padding: EdgeInsets.fromLTRB(5, 0, 0, 0),
                  child: Text(
                    strDate,
                    textAlign: TextAlign.start,
                    overflow: TextOverflow.clip,
                    style: TextStyle(
                      fontWeight: FontWeight.w400,
                      fontStyle: FontStyle.normal,
                      fontSize: 14,
                      color: Color(0xff000000),
                    ),
                  ),
                ),
              ),
              Padding(
                padding: EdgeInsets.all(2),
                child: Chip(
                  labelPadding: EdgeInsets.all(0),
                  label: Text("Sem $weeknum"),
                  labelStyle: TextStyle(
                    fontSize: 12,
                    fontWeight: FontWeight.w400,
                    fontStyle: FontStyle.normal,
                    color: Color(0x87000000),
                  ),
                  backgroundColor: Color(0x003a57e8),
                  elevation: 0,
                  shadowColor: Color(0xff808080),
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(15.0),
                  ),
                ),
              ),
              IconButton(
                icon: Icon(Icons.wb_sunny),
                onPressed: () {},
                color: Color(0xff212435),
                iconSize: 18,
              ),
            ],
          ),
        ],
      ),
    );
  }
}

class TopHomeFilter extends StatelessWidget {
  final TextEditingController filterController;
  final void Function(String text) onFilterChanged;

  const TopHomeFilter(
      {required this.filterController, required this.onFilterChanged});

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: EdgeInsets.all(0),
      padding: EdgeInsets.all(0),
      width: MediaQuery.of(context).size.width,
      height: 100,
      decoration: BoxDecoration(
        color: Color(0x1f000000),
        shape: BoxShape.rectangle,
        borderRadius: BorderRadius.zero,
        border: Border.all(color: Color(0x189e9e9e), width: 1),
      ),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.start,
        crossAxisAlignment: CrossAxisAlignment.center,
        mainAxisSize: MainAxisSize.max,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.start,
            crossAxisAlignment: CrossAxisAlignment.center,
            mainAxisSize: MainAxisSize.max,
            children: [
              Padding(
                padding: EdgeInsets.fromLTRB(5, 0, 0, 0),
                child: Icon(
                  Icons.location_on,
                  color: Color(0xff212435),
                  size: 18,
                ),
              ),
              Expanded(
                flex: 1,
                child: Container(
                    width: 130,
                    height: 50,
                    padding: EdgeInsets.symmetric(vertical: 4, horizontal: 8),
                    decoration: BoxDecoration(
                      color: Color(0x00ffffff),
                      borderRadius: BorderRadius.circular(0),
                    ),
                    child: DropdownButtonHideUnderline(
                      child: DropdownButton(
                        value: "Potager Jactez",
                        items: ["Potager Jactez", "Jardin partagé Tropark"]
                            .map<DropdownMenuItem<String>>((String value) {
                          return DropdownMenuItem<String>(
                            value: value,
                            child: Text(value),
                          );
                        }).toList(),
                        style: TextStyle(
                          color: Color(0xff000000),
                          fontSize: 16,
                          fontWeight: FontWeight.w400,
                          fontStyle: FontStyle.normal,
                        ),
                        onChanged: (value) {},
                        elevation: 8,
                        isExpanded: true,
                      ),
                    )),
              ),
            ],
          ),
          Row(
            mainAxisAlignment: MainAxisAlignment.start,
            crossAxisAlignment: CrossAxisAlignment.center,
            mainAxisSize: MainAxisSize.max,
            children: [
              Expanded(
                flex: 1,
                child: TextField(
                  controller: filterController,
                  onSubmitted: onFilterChanged,
                  textAlign: TextAlign.start,
                  style: TextStyle(
                    fontWeight: FontWeight.w400,
                    fontStyle: FontStyle.normal,
                    fontSize: 14,
                    color: Color(0xff000000),
                  ),
                  decoration: InputDecoration(
                    disabledBorder: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(20.0),
                      borderSide:
                          BorderSide(color: Color(0x39000000), width: 1),
                    ),
                    focusedBorder: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(20.0),
                      borderSide:
                          BorderSide(color: Color(0x39000000), width: 1),
                    ),
                    enabledBorder: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(20.0),
                      borderSide:
                          BorderSide(color: Color(0x39000000), width: 1),
                    ),
                    hintText: "Rechercher...",
                    hintStyle: TextStyle(
                      fontWeight: FontWeight.w400,
                      fontStyle: FontStyle.normal,
                      fontSize: 14,
                      color: Color(0xff000000),
                    ),
                    filled: true,
                    fillColor: Color(0xfff2f2f3),
                    isDense: true,
                    contentPadding:
                        EdgeInsets.symmetric(vertical: 8, horizontal: 12),
                    prefixIcon:
                        Icon(Icons.search, color: Color(0xff212435), size: 24),
                  ),
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }
}

class ActionListTile extends StatelessWidget {
  final ActionLog actionLog;
  final bool showDivider;

  const ActionListTile(
      {super.key, required this.actionLog, this.showDivider = true});

  @override
  Widget build(BuildContext context) {
    var lignes = [
      Align(
        alignment: Alignment.centerLeft,
        child: Text(
          "${actionLog.variete}${actionLog.lieu == "" ? "" : ' / ${actionLog.lieu}'}",
          textAlign: TextAlign.left,
          overflow: TextOverflow.clip,
          style: TextStyle(
            fontWeight: FontWeight.w400,
            fontStyle: FontStyle.normal,
            fontSize: 12,
            color: Color(0xff000000),
          ),
        ),
      )
    ];
    if (actionLog.notes != "") {
      lignes.add(Align(
        alignment: Alignment.centerLeft,
        child: Text(
          actionLog.notes,
          textAlign: TextAlign.start,
          overflow: TextOverflow.clip,
          style: TextStyle(
            fontWeight: FontWeight.w400,
            fontStyle: FontStyle.normal,
            fontSize: 12,
            color: Color(0xff000000),
          ),
        ),
      ));
    }

    Widget? tagLine;
    if (actionLog.tags.isNotEmpty) {
      List<Widget> chips = [];
      chips = actionLog.tags.map(
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
          );
        },
      ).toList();

      tagLine = Padding(
          padding: const EdgeInsets.only(left: 10, right: 10),
          child: Align(
              alignment: Alignment.topRight,
              child: Wrap(spacing: 5, runSpacing: 5, children: chips)));
    }

    var tile = Column(
      mainAxisAlignment: MainAxisAlignment.start,
      crossAxisAlignment: CrossAxisAlignment.center,
      mainAxisSize: MainAxisSize.max,
      children: [
        Padding(
          padding: EdgeInsets.fromLTRB(10, 0, 0, 0),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.start,
            crossAxisAlignment: CrossAxisAlignment.center,
            mainAxisSize: MainAxisSize.max,
            children: [
              Expanded(
                flex: 1,
                child: Text(
                  "${actionLog.action} - ${actionLog.legume} ${actionLog.qte > 0 ? '(${actionLog.qte})' : actionLog.poids > 0 ? '(${actionLog.poids}g)' : ''}",
                  textAlign: TextAlign.start,
                  overflow: TextOverflow.clip,
                  style: TextStyle(
                    fontWeight: FontWeight.w400,
                    fontStyle: FontStyle.normal,
                    fontSize: 14,
                    color: Color(0xff000000),
                  ),
                ),
              ),
              IconButton(
                icon: Icon(Icons.more_vert),
                onPressed: () {},
                color: Color(0xff212435),
                iconSize: 22,
              ),
            ],
          ),
        ),
        Padding(
          padding: EdgeInsets.fromLTRB(10, 0, 0, 0),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.start,
            crossAxisAlignment: CrossAxisAlignment.center,
            mainAxisSize: MainAxisSize.max,
            children: [
              Expanded(
                flex: 1,
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.start,
                  crossAxisAlignment: CrossAxisAlignment.center,
                  mainAxisSize: MainAxisSize.max,
                  children: lignes,
                ),
              ),
            ],
          ),
        ),
      ],
    );

    if (tagLine != null) {
      tile.children.add(tagLine);
    }

    if (true == showDivider) {
      tile.children.add(Divider(
        color: Color(0xff808080),
        height: 1,
      ));
    }

    return tile;
  }
}

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
  TextEditingController filterController = TextEditingController();
  int currentPage = 0;
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
    actionLogs = (await ApiService().getLogs())!;
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
          .where((a) => a
              .toJson()
              .toString()
              .toLowerCase()
              .contains(filterController.text.toLowerCase()))
          .toList();
      filteredActionLogs.sort((a, b) => b.dateAction.compareTo(a.dateAction));
    }

    return Scaffold(
      appBar: AppBar(
        // Here we take the value from the MyHomePage object that was created by
        // the App.build method, and use it to set our appbar title.
        title: Text(widget.title),
      ),
      bottomNavigationBar: NavigationBar(
        destinations: const [
          NavigationDestination(icon: Icon(Icons.home), label: "Home"),
          NavigationDestination(icon: Icon(Icons.add_chart), label: "Whatever")
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
          child: Center(
            // Center is a layout widget. It takes a single child and positions it
            // in the middle of the parent.
            child: ListView.builder(
              itemCount: filteredActionLogs.length + 1,
              itemBuilder: (context, index) {
                List<Widget> results = [];

                if (index == 0) {
                  results.add(TopHomeFilter(
                      filterController: filterController,
                      onFilterChanged: _onFilterChanged));
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
            ),
          )),
      floatingActionButton: FloatingActionButton(
        onPressed: () async {
          DateTime today = DateTime.now();
          today = DateTime(today.year, today.month, today.day);
          ActionLog a = ActionLog(dateAction: today);
          ActionLog result = await Navigator.push(
              context,
              MaterialPageRoute(
                  builder: (context) => ActionDetail(actionLog: a)));
          setState(() {
            //-- result est une copie de a. On doit donc "recharger" a avec les valeurs modifiées.
            a.updateFrom(result);
            actionLogs.add(a);
            actionLogs.sort((a, b) => b.dateAction.compareTo(a.dateAction));
          });
        },
        tooltip: 'Ajouter une action',
        child: const Icon(Icons.add),
      ),
    );
  }

  void Function() onTileTap(ActionLog a) {
    return () async {
      ActionLog result = await Navigator.push(context,
          MaterialPageRoute(builder: (context) => ActionDetail(actionLog: a)));
      setState(() {
        //-- result est une copie de a. On doit donc "recharger" a avec les valeurs modifiées.
        a.updateFrom(result);
      });
    };
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
