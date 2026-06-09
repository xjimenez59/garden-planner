//File download from FlutterViz- Drag and drop a tools. For more details visit https://flutterviz.io/

//File download from FlutterViz- Drag and drop a tools. For more details visit https://flutterviz.io/

// ignore_for_file: prefer_const_constructors

import 'package:app/action_log.dart';
import 'package:app/cleanup_view.dart';
import 'package:app/garden_model.dart';
import 'package:app/gardens_view.dart';
import 'package:app/meteo_service.dart';
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
  List<Meteo> meteoData = [];
  List<LuneDay> luneData = [];
  List<HourlyForecast> previsions = [];
  TextEditingController filterController = TextEditingController();
  final ScrollController _scrollController = ScrollController();
  bool _showToTopButton = false;
  bool _hasMore = true;
  bool _loadingMore = false;
  String? _beforeCursor;
  // Mode recherche serveur
  List<ActionLog> _searchResults = [];
  bool _searchHasMore = false;
  bool _searchLoading = false;
  String? _searchBeforeCursor;
  int currentPage = 0;
  Garden? selectedGarden;

  @override
  void initState() {
    super.initState();
    _scrollController.addListener(_onScroll);
    _getData();
  }

  @override
  void dispose() {
    filterController.dispose();
    _scrollController.dispose();
    super.dispose();
  }

  void _onScroll() {
    final pos = _scrollController.position;
    final showBtn = pos.pixels > 400;
    if (showBtn != _showToTopButton) setState(() => _showToTopButton = showBtn);
    if (pos.pixels <= pos.maxScrollExtent * 0.66) return;
    if (filterController.text.isNotEmpty) {
      if (_searchHasMore && !_searchLoading) _loadSearchPage(filterController.text);
    } else {
      if (_hasMore && !_loadingMore) _loadMoreLogs();
    }
  }

  Future<void> _loadMoreLogs() async {
    if (selectedGarden == null || !_hasMore || _loadingMore) return;
    setState(() => _loadingMore = true);
    final page = await ApiService().getLogs(selectedGarden!, before: _beforeCursor);
    if (page != null && page.logs.isNotEmpty) {
      final prevOldest = _beforeCursor;
      setState(() {
        actionLogs.addAll(page.logs);
        _hasMore = page.hasMore;
        _beforeCursor = page.oldestDate;
        _loadingMore = false;
      });
      // Charger meteo/lune uniquement sur la nouvelle tranche
      if (page.oldestDate != null && prevOldest != null) {
        await _loadMeteoRange(page.oldestDate!, prevOldest);
        setState(() {});
      }
    } else {
      setState(() {
        _hasMore = false;
        _loadingMore = false;
      });
    }
  }

  void _getData() async {
    jardins = (await ApiService().getGardens())!;
    if (jardins.isNotEmpty) {
      selectedGarden = jardins.first;
      await _resetAndLoadLogs(selectedGarden!);
    }
    setState(() {});
    if (selectedGarden != null) {
      await _loadPrevisions(selectedGarden!);
      setState(() {});
    }
  }

  Future<void> _resetAndLoadLogs(Garden garden) async {
    actionLogs = [];
    meteoData = [];
    luneData = [];
    _hasMore = true;
    _beforeCursor = null;

    final page = await ApiService().getLogs(garden);
    if (page != null) {
      actionLogs = page.logs;
      _hasMore = page.hasMore;
      _beforeCursor = page.oldestDate;
      if (page.oldestDate != null) {
        final today = DateTime.now();
        final dateFin = '${today.year}-${today.month.toString().padLeft(2, '0')}-${today.day.toString().padLeft(2, '0')}';
        await _loadMeteoRange(page.oldestDate!, dateFin);
      }
    }
  }

  String _fmtYMDcompact(String ymd) => ymd.replaceAll('-', '');

  Future<void> _loadMeteoRange(String dateDeb, String dateFin) async {
    final newLune = await MeteoService().getLuneRange(dateDeb, dateFin);
    final existingDates = luneData.map((l) => l.date).toSet();
    luneData.addAll(newLune.where((l) => !existingDates.contains(l.date)));

    if (selectedGarden == null || selectedGarden!.MeteofSite.isEmpty) return;
    final newMeteo = await MeteoService().getMeteo(
      selectedGarden!.MeteofSite,
      dateDeb: _fmtYMDcompact(dateDeb),
      dateFin: _fmtYMDcompact(dateFin),
    );
    final existingMeteo = meteoData.map((m) => m.date).toSet();
    meteoData.addAll(newMeteo.where((m) => !existingMeteo.contains(m.date)));
  }

  Future<void> _loadPrevisions(Garden garden) async {
    if (garden.MeteofSite.isEmpty) return;
    previsions = await MeteoService().getPrevisions(garden.MeteofSite);
    // Charger le lune pour toute la période de prévision (aujourd'hui + 6 j)
    final today = DateTime.now();
    String fmt(DateTime d) =>
        '${d.year}-${d.month.toString().padLeft(2, '0')}-${d.day.toString().padLeft(2, '0')}';
    final newLune = await MeteoService()
        .getLuneRange(fmt(today), fmt(today.add(Duration(days: 6))));
    final existingDates = luneData.map((l) => l.date).toSet();
    luneData.addAll(newLune.where((l) => !existingDates.contains(l.date)));
  }

Future<void> _onRefresh() async {
    _getData();
  }

  void _onFilterChanged(String text) {
    if (text.isEmpty) {
      setState(() {
        _searchResults = [];
        _searchHasMore = false;
        _searchBeforeCursor = null;
      });
    } else {
      _searchResults = [];
      _searchHasMore = false;
      _searchBeforeCursor = null;
      _searchLoading = false;
      _loadSearchPage(text);
    }
  }

  Future<void> _loadSearchPage(String query) async {
    if (_searchLoading || selectedGarden == null) return;
    setState(() => _searchLoading = true);
    final page = await ApiService().getLogs(
      selectedGarden!,
      before: _searchBeforeCursor,
      limit: 30,
      search: query.withoutDiacriticalMarks,
    );
    if (!mounted || filterController.text != query) {
      if (mounted) setState(() => _searchLoading = false);
      return;
    }
    if (page != null) {
      final newDates = page.logs
        .map((l) => '${l.dateAction.year}-${l.dateAction.month.toString().padLeft(2, '0')}-${l.dateAction.day.toString().padLeft(2, '0')}')
        .toSet().toList(); // dédupliqué : plusieurs logs peuvent avoir la même date
      setState(() {
        _searchResults.addAll(page.logs);
        _searchHasMore = page.hasMore;
        _searchBeforeCursor = page.oldestDate;
        _searchLoading = false;
      });
      await _loadMeteoForDates(newDates);
      if (mounted) setState(() {});
    } else {
      setState(() => _searchLoading = false);
    }
  }

  // Charge meteo+lune uniquement pour les dates non encore couvertes en mémoire.
  Future<void> _loadMeteoForDates(List<String> dates) async {
    final coveredMeteo = meteoData.map((m) =>
        '${m.date.year}${m.date.month.toString().padLeft(2,'0')}${m.date.day.toString().padLeft(2,'0')}').toSet();
    final coveredLune = luneData.map((l) =>
        '${l.date.year}-${l.date.month.toString().padLeft(2,'0')}-${l.date.day.toString().padLeft(2,'0')}').toSet();

    // dates reçues sont en YYYY-MM-DD
    final missingMeteo = dates.where((d) => !coveredMeteo.contains(d.replaceAll('-', ''))).toList();
    final missingLune  = dates.where((d) => !coveredLune.contains(d)).toList();

    if (selectedGarden != null && selectedGarden!.MeteofSite.isNotEmpty && missingMeteo.isNotEmpty) {
      final compact = missingMeteo.map((d) => d.replaceAll('-', '')).toList();
      final newMeteo = await MeteoService().getMeteoForDates(selectedGarden!.MeteofSite, compact);
      meteoData.addAll(newMeteo.where((m) => !coveredMeteo.contains(
          '${m.date.year}${m.date.month.toString().padLeft(2,'0')}${m.date.day.toString().padLeft(2,'0')}')));
    }

    if (missingLune.isNotEmpty) {
      final newLune = await MeteoService().getLuneForDates(missingLune);
      luneData.addAll(newLune.where((l) => !coveredLune.contains(
          '${l.date.year}-${l.date.month.toString().padLeft(2,'0')}-${l.date.day.toString().padLeft(2,'0')}')));
    }
  }

  @override
  Widget build(BuildContext context) {
    final bool isSearching = filterController.text.isNotEmpty;
    final List<ActionLog> displayedLogs = isSearching ? _searchResults : actionLogs;
    final bool isLoadingDisplay = isSearching ? _searchLoading : _loadingMore;

    // Insérer un en-tête "Aujourd'hui" si aucun log n'existe pour la date du jour
    final todayMidnight = DateTime.now();
    final todayDate = DateTime(todayMidnight.year, todayMidnight.month, todayMidnight.day);
    final bool needsTodayHeader = !isSearching &&
        (displayedLogs.isEmpty || !displayedLogs[0].dateAction.sameDayAs(todayDate));
    final int offset = needsTodayHeader ? 1 : 0;

    return Scaffold(
      appBar: AppBar(
        // Here we take the value from the MyHomePage object that was created by
        // the App.build method, and use it to set our appbar title.
        title: Text(selectedGarden == null
            ? widget.title
            : "Garden Planner - ${selectedGarden!.Nom}"),
      ),
      bottomNavigationBar: NavigationBar(
        destinations: const [
          NavigationDestination(icon: Icon(Icons.home), label: "Home"),
          NavigationDestination(icon: Icon(Icons.add_chart), label: "Récoltes"),
          NavigationDestination(icon: Icon(Icons.tune), label: "Données"),
        ],
        onDestinationSelected: (int i) {
          setState(() {
            currentPage = i;
          });
        },
        selectedIndex: currentPage,
      ),
      body: currentPage == 0
              ? Column(
                  children: [
                    TopHomeFilter(
                        jardins: jardins,
                        filterController: filterController,
                        onFilterChanged: _onFilterChanged,
                        onSelectGardenTap: onSelectGardenTap),
                    Expanded(
                      child: Stack(
                        children: [
                      RefreshIndicator(
                        onRefresh: _onRefresh,
                        child: ListView.builder(
                          controller: _scrollController,
                          itemCount: displayedLogs.length + 1 + offset,
                          itemBuilder: (context, index) {
                            // En-tête "Aujourd'hui" synthétique (aucun log ce jour)
                            if (needsTodayHeader && index == 0) {
                              return DaySeparator(
                                date: todayDate,
                                luneDay: luneData
                                    .where((l) => l.date == todayDate)
                                    .firstOrNull,
                                previsions: previsions,
                                luneData: luneData,
                              );
                            }
                            final logIndex = index - offset;
                            // Indicateur de chargement en bas de liste
                            if (logIndex == displayedLogs.length) {
                              if (isLoadingDisplay) {
                                return const Padding(
                                  padding: EdgeInsets.all(16),
                                  child: Center(child: CircularProgressIndicator()),
                                );
                              }
                              return const SizedBox.shrink();
                            }

                            List<Widget> results = [];
                            var a = displayedLogs[logIndex];

                            final prev = logIndex > 0 ? displayedLogs[logIndex - 1] : null;
                            final isNewDay = prev == null || !prev.dateAction.sameDayAs(a.dateAction);
                            if (isNewDay) {
                              final isNewYear = prev == null
                                  ? a.dateAction.year != DateTime.now().year
                                  : prev.dateAction.year != a.dateAction.year;
                              if (isNewYear) {
                                results.add(YearSeparator(year: a.dateAction.year));
                              }
                              Meteo? meteo = meteoData
                                  .where((m) => m.date == a.dateAction)
                                  .firstOrNull;
                              LuneDay? luneDay = luneData
                                  .where ((l) => l.date == a.dateAction)
                                  .firstOrNull;
                              final isToday = a.dateAction.sameDayAs(DateTime.now());
                              results.add(DaySeparator(
                                  date: a.dateAction,
                                  meteo: meteo,
                                  luneDay: luneDay,
                                  previsions: isToday ? previsions : null,
                                  luneData: isToday ? luneData : null));
                            }
                            bool showDivider =
                                logIndex < displayedLogs.length - 1 &&
                                    displayedLogs[logIndex + 1]
                                        .dateAction
                                        .sameDayAs(a.dateAction);
                            results.add(Dismissible(
                                key: Key(a.id),
                                onDismissed: onTileDismissed(logIndex),
                                background: Container(
                                    color: Colors.white,
                                    margin: EdgeInsets.only(bottom: 0)),
                                child: InkWell(
                                  onTap: onTileTap(a),
                                  child: ActionListTile(
                                      actionLog: a, showDivider: showDivider),
                                )));

                            if (results.length > 1) {
                              return Column(children: results);
                            }
                            return results.first;
                          },
                        ),
                      ),
                        ],
                      ),
                    ),
                  ],
                )
              : currentPage == 1
                  ? ActionLogStats()
                  : selectedGarden != null
                      ? CleanupView(garden: selectedGarden!)
                      : Center(child: Text('Sélectionnez un jardin')),
      floatingActionButton: (currentPage == 0 && selectedGarden != null)
          ? Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                if (_showToTopButton) ...[
                  FloatingActionButton(
                    heroTag: 'toTop',
                    tooltip: "Aujourd'hui",
                    onPressed: () => _scrollController.animateTo(
                      0,
                      duration: Duration(milliseconds: 400),
                      curve: Curves.easeOut,
                    ),
                    child: Icon(Icons.vertical_align_top),
                  ),
                  SizedBox(height: 8),
                ],
                FloatingActionButton(
                  heroTag: 'add',
                  onPressed: onNewActionLogTap,
                  tooltip: 'Ajouter une action',
                  child: const Icon(Icons.add),
                ),
              ],
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
      await _resetAndLoadLogs(selectedGarden!);
      await _loadPrevisions(selectedGarden!);
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
