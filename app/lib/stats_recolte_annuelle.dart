import 'package:app/api_service.dart';
import 'package:app/recolte_model.dart';
import 'package:app/utils.dart';
import 'package:flutter/material.dart';
import 'package:flutter_expandable_table/flutter_expandable_table.dart';

enum GroupeRecoltes { parLegume, parLieu }

class StatsRecolteAnnuelle extends StatefulWidget {
  GroupeRecoltes groupePar;
  StatsRecolteAnnuelle({super.key, this.groupePar = GroupeRecoltes.parLegume});

  @override
  State<StatsRecolteAnnuelle> createState() => _StatsRecolteAnnuelleState();
}

class _StatsRecolteAnnuelleState extends State<StatsRecolteAnnuelle> {
  late List<Recolte> recoltes = [];

  @override
  void initState() {
    super.initState();
    _getRecolteAnnuelle();
  }

  @override
  Widget build(BuildContext context) {
    Widget res;

    res = Scaffold(
        appBar: AppBar(
            title: Text(
                "Récolte Annuelle ${widget.groupePar == GroupeRecoltes.parLegume ? "par Légume" : "par Lieu"}")),
        body: recoltes.isEmpty
            ? const Center(child: Text("Chargement en cours"))
            : Container(
                child: RecolteAnnuelleDataTable(
                    recoltes: recoltes, groupePar: widget.groupePar)));

    return res;
  }

  void _getRecolteAnnuelle() async {
    recoltes = (await ApiService().getRecolteAnnuelle())!;
    setState(() {});
    return;
  }
}

class RecolteAnnuelleDataTable extends StatefulWidget {
  final List<Recolte> recoltes;
  GroupeRecoltes groupePar;
  RecolteAnnuelleDataTable({
    super.key,
    required this.recoltes,
    this.groupePar = GroupeRecoltes.parLegume,
  });

  @override
  State<RecolteAnnuelleDataTable> createState() => _RecolteAnnuelleDataTable();
}

class _RecolteAnnuelleDataTable extends State<RecolteAnnuelleDataTable> {
  Map<int, int> totauxGlobaux = {};
  Map<String, Map<int, int>> totGroupedRow = {};
  List<Recolte> recoltes = []; // version filtree classee
  List<int> annees = [];

  late ExpandableTableController controller;

  @override
  void initState() {
    super.initState();
    recoltes = widget.recoltes;
    totauxGlobaux = totalRecoltes(recoltes);
    annees = totauxGlobaux.keys.toList();
    annees.sort((a, b) => b.compareTo(a));
    _sortRecolteAnnuelle();
    controller = _buildController();
  }

  ExpandableTableController _buildController() {
    ExpandableTableController res = ExpandableTableController(
      firstHeaderCell: ExpandableTableCell(
          child: Container(
              color: Colors.grey.shade300,
              padding: EdgeInsets.only(left: 25, right: 5),
              margin: EdgeInsets.only(bottom: 10),
              child: Text(""))),
      headers: _buildHeaders(),
      rows: buildRows(annees, widget.groupePar),
      headerHeight: 60,
      firstColumnWidth: 150,
      defaultsRowHeight: 40,
    );

    return res;
  }

  @override
  Widget build(BuildContext context) {
    Widget res = ExpandableTable(controller: controller);
    return res;
  }

  void _sortRecolteAnnuelle() {
    List<Recolte> sortedRecolte = recoltes.sublist(0);
    recoltes = sortedRecolte;
    recoltes.sort((a, b) {
      String strSortA = (widget.groupePar == GroupeRecoltes.parLieu)
          ? "${a.lieu}-${a.legume}"
          : "${a.legume}-${a.lieu}";
      String strSortB = (widget.groupePar == GroupeRecoltes.parLieu)
          ? "${b.lieu}-${b.legume}"
          : "${b.legume}-${b.lieu}";

      return strSortA.withoutDiacriticalMarks
          .toLowerCase()
          .compareTo(strSortB.withoutDiacriticalMarks.toLowerCase());
    });
    recoltes.removeWhere((element) => element.legume.isEmpty);
    _calculeTotauxGroupe(widget.groupePar);
  }

  void _calculeTotauxGroupe(GroupeRecoltes groupePar) {
    for (var r in recoltes) {
      String strGroupe =
          (groupePar == GroupeRecoltes.parLegume) ? r.legume : r.lieu;
      Map<int, int> curGroupedRow =
          totGroupedRow.putIfAbsent(strGroupe, () => {});
      for (var a in r.annees) {
        curGroupedRow.update(a.annee, (value) => value + a.poids,
            ifAbsent: () => a.poids);
      }
    }
  }

  List<ExpandableTableHeader> _buildHeaders() {
    List<ExpandableTableHeader> headers = [];
    for (var a in annees) {
      headers.add(ExpandableTableHeader(
          cell: ExpandableTableCell(
              child: Container(
                  color: Colors.grey.shade300,
                  padding: EdgeInsets.only(left: 5, right: 5),
                  margin: EdgeInsets.only(bottom: 10),
                  child: Column(
                    children: [
                      Text(a.toString()),
                      Text(weightFormat(totauxGlobaux[a]!))
                    ],
                    crossAxisAlignment: CrossAxisAlignment.center,
                    mainAxisAlignment: MainAxisAlignment.center,
                  )))));
    }
    return headers;
  }

  List<ExpandableTableRow> buildRows(
      List<int> annees, GroupeRecoltes groupePar) {
    List<ExpandableTableRow> rows = [];
    totGroupedRow.forEach((key, curTot) {
      //-- on commence par ajouter une ligne de "regroupement"
      ExpandableTableRow row;
      List<ExpandableTableCell> cells = [];
      for (int a in annees) {
        String strCell = curTot.containsKey(a) ? weightFormat(curTot[a]!) : "";
        cells.add(ExpandableTableCell(
            child: Container(
                alignment: Alignment.center,
                decoration: BoxDecoration(
                    border: Border(
                        bottom: BorderSide(
                  color: Colors.grey.shade300,
                  width: 1,
                ))),
                child: Text(strCell))));
      }
      row = ExpandableTableRow(
        firstCell: _buildFirstRowCell(key == "" ? "(non spécifié)" : key),
        cells: cells,
        //childrenExpanded: false
      );

      //-- puis toutes les lignes de ce regroupement
      List<ExpandableTableRow> subRows = [];
      for (var r in recoltes.where((e) =>
          (groupePar == GroupeRecoltes.parLegume ? e.legume : e.lieu) == key)) {
        cells = [];
        for (int a in annees) {
          RecolteAnnee? recolteAnnee =
              r.annees.where((e) => e.annee == a).firstOrNull;
          String strCell =
              recolteAnnee == null ? "" : weightFormat(recolteAnnee.poids);
          cells.add(ExpandableTableCell(
              child: Container(
                  decoration: BoxDecoration(
                      border: Border(
                          bottom: BorderSide(
                    color: Colors.grey.shade300,
                    width: 1,
                  ))),
                  alignment: Alignment.center,
                  child: Text(strCell))));
        }
        String strCell =
            (groupePar == GroupeRecoltes.parLieu) ? r.legume : r.lieu;
        if (strCell.isEmpty) {
          strCell = "(non précisé)";
        }

        subRows.add(ExpandableTableRow(
            firstCell: ExpandableTableCell(
                child: Container(
                    margin: EdgeInsets.only(left: 40),
                    alignment: Alignment.centerLeft,
                    decoration: BoxDecoration(
                        border: Border(
                            bottom: BorderSide(
                      color: Colors.grey.shade300,
                      width: 1,
                    ))),
                    child: Text(strCell))),
            cells: cells));
      }
      row.children = subRows;
      rows.add(row);
    });
    return rows;
  }

  ExpandableTableCell _buildFirstRowCell(String content) {
    return ExpandableTableCell(
      builder: (context, details) => Padding(
        padding: const EdgeInsets.only(left: 16.0),
        child: Container(
            child: Row(
              children: [
                SizedBox(
                  width: 24 * details.row!.address.length.toDouble(),
                  child: details.row?.children != null
                      ? Align(
                          alignment: Alignment.centerRight,
                          child: AnimatedRotation(
                            duration: const Duration(milliseconds: 250),
                            turns: details.row?.childrenExpanded == true
                                ? 0.25
                                : 0,
                            child: const Icon(
                              Icons.keyboard_arrow_right,
                              color: Colors.black,
                            ),
                          ),
                        )
                      : null,
                ),
                Text(content),
              ],
            ),
            decoration: BoxDecoration(
                border: Border(
                    bottom: BorderSide(
              color: Colors.grey.shade300,
              width: 1,
            )))),
      ),
    );
  }
}
