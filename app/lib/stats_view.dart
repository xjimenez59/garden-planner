import 'package:app/api_service.dart';
import 'package:app/recolte_model.dart';
import 'package:app/utils.dart';
import 'package:flutter/material.dart';

class ActionLogStats extends StatefulWidget {
  const ActionLogStats({super.key});

  @override
  State<ActionLogStats> createState() => _ActionLogStatsState();
}

class _ActionLogStatsState extends State<ActionLogStats> {
  late List<Recolte> recoltes = [];

  @override
  void initState() {
    super.initState();
    _getRecoltes();
  }

  @override
  Widget build(BuildContext context) {
    Widget res;

    res = SingleChildScrollView(
        child: Center(child: RecoltesTableau(recoltes: recoltes)));

    return res;
  }

  void _getRecoltes() async {
    recoltes = (await ApiService().getRecoltes())!;
    setState(() {});
    return;
  }
}

class RecoltesTableau extends StatelessWidget {
  final List<Recolte> recoltes;
  const RecoltesTableau({super.key, required this.recoltes});

  @override
  Widget build(BuildContext context) {
    Widget res = DataTable(
        columns: const <DataColumn>[
          DataColumn(
            label: Text('Légume'),
          ),
          DataColumn(
            label: Text('Poids'),
          ),
          DataColumn(
            label: Text('Qté'),
          ),
        ],
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
            cells: <DataCell>[
              DataCell(Text(recoltes[index].legume)),
              DataCell(Text(weightFormat(recoltes[index].poids))),
              DataCell(Text(recoltes[index].qte.toString())),
            ],
/*             selected: selected[index],
            onSelectChanged: (bool? value) {
              setState(() {
                selected[index] = value!;
              });
             },*/
          ),
        ));
    return res;
  }
}
