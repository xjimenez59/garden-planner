// ignore_for_file: prefer_const_constructors

import 'package:app/action_log.dart';
import 'package:app/garden_model.dart';
import 'package:app/utils.dart';
import 'package:flutter/material.dart';

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
  final List<Garden> jardins;
  final void Function() onSelectGardenTap;

  const TopHomeFilter(
      {required this.jardins,
      required this.filterController,
      required this.onFilterChanged,
      required this.onSelectGardenTap});

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
              Expanded(
                flex: 1,
                child: TextField(
                  controller: filterController,
                  onChanged: onFilterChanged,
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
              IconButton(
                  onPressed: onSelectGardenTap, icon: Icon(Icons.location_pin))
            ],
          ),
        ],
      ),
    );
  }
}

class HorizontalImageListview extends StatelessWidget {
  final List<String> imgUrlList;

  const HorizontalImageListview({super.key, required this.imgUrlList});
  @override
  Widget build(BuildContext context) {
    if (imgUrlList.isEmpty) {
      return Container();
    }
    Widget result = Container(
        height: 200,
        child: ListView.builder(
            scrollDirection: Axis.horizontal,
            itemCount: imgUrlList.length,
            itemBuilder: (context, index) {
              return Container(
                  padding: EdgeInsets.all(10),
                  height: 150,
                  child: Image.network(imgUrlList[index]));
            }));
    return result;
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
                  "${actionLog.action} - ${actionLog.legume} ${actionLog.poids > 0 ? '(${actionLog.poids}g)' : actionLog.qte > 0 ? '(${actionLog.qte})' : ''}",
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
              /* IconButton(
                icon: Icon(Icons.more_vert),
                onPressed: () {},
                color: Color(0xff212435),
                iconSize: 22,
              ), */
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

    if (actionLog.photos.isNotEmpty) {
      tile.children.add(HorizontalImageListview(imgUrlList: actionLog.photos));
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
