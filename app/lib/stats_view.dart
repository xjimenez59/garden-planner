import 'package:app/api_service.dart';
import 'package:app/recolte_model.dart';
import 'package:app/stats_recolte_annuelle.dart';
import 'package:app/utils.dart';
import 'package:flutter/material.dart';

class ActionLogStats extends StatefulWidget {
  const ActionLogStats({super.key});

  @override
  State<ActionLogStats> createState() => _ActionLogStatsState();
}

class _ActionLogStatsState extends State<ActionLogStats> {
  late List<Recolte> recoltes = [];
  List<bool> _isOpen = [false, false];

  @override
  void initState() {
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    Widget res;

    res = ListView(
      children: [
        Card(
            child: ListTile(
          title: Text(
            "Récolte annuelle par légume",
          ),
          trailing: Icon(Icons.chevron_right),
          onTap: () {
            Navigator.push(
                context,
                MaterialPageRoute(
                    builder: (context) => StatsRecolteAnnuelle(
                        groupePar: GroupeRecoltes.parLegume)));
          },
        )),
        Card(
            child: ListTile(
          title: Text(
            "Récolte annuelle par Lieu",
          ),
          trailing: Icon(Icons.chevron_right),
          onTap: () {
            Navigator.push(
                context,
                MaterialPageRoute(
                    builder: (context) => StatsRecolteAnnuelle(
                        groupePar: GroupeRecoltes.parLieu)));
          },
        ))
      ],
    );

    return res;
  }
}
