import 'package:app/garden_form.dart';
import 'package:app/garden_model.dart';
import 'package:app/utils.dart';
import 'package:flutter/material.dart';

class GardensView extends StatefulWidget {
  final Garden? activeGarden;
  final List<Garden> gardens;

  const GardensView({super.key, required this.gardens, this.activeGarden});

  @override
  State<StatefulWidget> createState() {
    return _GardensView();
  }
}

class _GardensView extends State<GardensView> {
  List<Garden> gardens = [];
  Garden? activeGarden;

  @override
  void initState() {
    super.initState();

    gardens = widget.gardens;
    activeGarden = widget.activeGarden;

    if (gardens.isNotEmpty && activeGarden == null) {
      activeGarden = gardens.first;
    }
  }

  @override
  Widget build(BuildContext context) {
    Widget result = Scaffold(
        appBar: AppBar(
          title: const Text("Choisissez un jardin"),
        ),
        floatingActionButton: FloatingActionButton(
          onPressed: onNewGardenTap,
          tooltip: 'Ajouter',
          child: const Icon(Icons.add),
        ),
        body: Align(
            alignment: Alignment.topCenter,
            child: GardenListView(
              gardens: gardens,
              selectedGarden: activeGarden,
              onSelectGarden: onSelectGarden,
              onEditGarden: onEditGarden,
            )));

    return result;
  }

  void onSelectGarden(Garden g) {
    setState(() {
      activeGarden = g;
    });
    Navigator.pop(context, g);
  }

  void onNewGardenTap() async {
    Garden? result = await Navigator.push(
        context, MaterialPageRoute(builder: (context) => GardenForm()));
    setState(() {
      if (result != null) {
        gardens!.add(result);
      }
    });
  }

  void onEditGarden(Garden g) async {
    Garden? result = await Navigator.push(context,
        MaterialPageRoute(builder: (context) => GardenForm(garden: g)));
    setState(() {
      if (result != null) {
        int index = gardens.indexOf(g);
        gardens[index].copyFrom(result);
      }
    });
  }
}

class GardenListView extends StatelessWidget {
  final List<Garden> gardens;
  final Garden? selectedGarden;
  final void Function(Garden g) onSelectGarden;
  final void Function(Garden g) onEditGarden;

  const GardenListView(
      {super.key,
      required this.gardens,
      required this.selectedGarden,
      required this.onSelectGarden,
      required this.onEditGarden});

  @override
  Widget build(BuildContext context) {
    return ListView(
      padding: const EdgeInsets.all(8),
      children: gardens
          .map((e) => GardenCard(
                garden: e,
                isActive: e == selectedGarden,
                onSelectGarden: onSelectGarden,
                onEditGarden: onEditGarden,
              ))
          .toList(),
    );
  }
}

class GardenCard extends StatelessWidget {
  final Garden garden;
  final bool isActive;
  final void Function(Garden g) onSelectGarden;
  final void Function(Garden g) onEditGarden;

  const GardenCard(
      {super.key,
      required this.garden,
      this.isActive = false,
      required this.onSelectGarden,
      required this.onEditGarden});

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Card(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: <Widget>[
            ListTile(
              onTap: onGardenTap,
              leading:
                  isActive ? const Icon(Icons.where_to_vote) : const Icon(null),
              title: Text(garden.Nom),
              subtitle: Text(
                  "Fin récoltes : ${monthNames[garden.MoisFinRecolte]} - Fin semis : ${monthNames[garden.MoisFinSemis]}"),
            ),
            Row(
              mainAxisAlignment: MainAxisAlignment.end,
              children: <Widget>[
                TextButton(
                  child: const Text('MODIFIER'),
                  onPressed: () {
                    onEditGarden(garden);
                  },
                ),
                const SizedBox(width: 8),
                TextButton(
                  child: const Text('SUPPRIMER'),
                  onPressed: () {/* ... */},
                ),
                const SizedBox(width: 8),
              ],
            ),
          ],
        ),
      ),
    );
  }

  void onGardenTap() {
    onSelectGarden(garden);
  }
}
