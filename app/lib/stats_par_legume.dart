import 'package:app/api_service.dart';
import 'package:app/recolte_model.dart';
import 'package:app/utils.dart';
import 'package:flutter/material.dart';

class StatsParLegume extends StatefulWidget {
  const StatsParLegume({super.key});

  @override
  State<StatsParLegume> createState() => _StatsParLegumeState();
}

class _StatsParLegumeState extends State<StatsParLegume> {
  late List<Recolte> recoltes = [];

  @override
  void initState() {
    super.initState();
    _getRecoltes();
  }

  @override
  Widget build(BuildContext context) {
    Widget res;

    res = Scaffold(
        appBar: AppBar(title: const Text("Récoltes par Légume")),
        body: SingleChildScrollView(
            child: Center(child: RecoltesTableau(recoltes: recoltes))));

    return res;
  }

  void _getRecoltes() async {
    recoltes = (await ApiService().getRecoltes())!;
    recoltes.sort((a, b) => a.legume.withoutDiacriticalMarks
        .toLowerCase()
        .compareTo(b.legume.withoutDiacriticalMarks.toLowerCase()));
    recoltes.removeWhere((element) => element.legume.isEmpty);
    setState(() {});

    return;
  }
}

class RecoltesTableau extends StatelessWidget {
  final List<Recolte> recoltes;
  const RecoltesTableau({super.key, required this.recoltes});

  @override
  Widget build(BuildContext context) {
    Map<int, int> totaux = totalRecoltes(recoltes);

    List<DataColumn> columns = [
      const DataColumn(label: Text('Légume')),
    ];
    List<int> annees = totaux.keys.toList();
    annees.sort((a, b) => b.compareTo(a));
    for (var a in annees) {
      columns.add(DataColumn(
          label: Column(
              children: [Text(a.toString()), Text(weightFormat(totaux[a]!))])));
    }

    Widget res = DataTable(
      columns: columns,
      rows: List<DataRow>.generate(
        recoltes.length,
        (int index) => DataRow(
          color: MaterialStateProperty.resolveWith<Color?>(
              (Set<MaterialState> states) {
            // All rows will have the same selected color.
            if (states.contains(MaterialState.selected)) {
              return Theme.of(context).colorScheme.primary.withOpacity(0.08);
            }
            // Even rows will have a grey color.
            if (index.isEven) {
              return Colors.grey.withOpacity(0.3);
            }
            return null; // Use default value for other states and odd rows.
          }),
          cells: [
            DataCell(Text(recoltes[index].legume)),
            for (int a in annees)
              DataCell(Text(weightFormat(recoltes[index]
                  .annees
                  .firstWhere((e) => e.annee == a,
                      orElse: () => RecolteAnnee(annee: a, poids: 0, qte: 0))
                  .poids))),
          ],
/*             selected: selected[index],
            onSelectChanged: (bool? value) {
              setState(() {
                selected[index] = value!;
              });
             },*/
        ),
      ),
    );
    return res;
  }
}
