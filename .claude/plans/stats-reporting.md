# Plan — Statistiques et reporting du jardin

## Contexte

Les données d'`action_log` (date, action, légume, variété, lieu, qte, poids) permettent des analyses agronomiques précieuses que l'app n'exploite pas encore au-delà des tableaux de récoltes. L'objectif est d'enrichir la vue `stats_view.dart` existante avec trois nouvelles vues et un export Excel.

**Données disponibles** : `action_log` contient `date_action`, `action` (Semis, Repiquage, Plantation, Récolte…), `legume`, `variete`, `lieu`, `qte` (quantité), `poids` (grammes). Les vues existantes (`stats_recolte_annuelle.dart`) utilisent déjà `flutter_expandable_table`.

**Décisions** :
- Tri des légumes : **alphabétique** pour l'instant (pas de groupement par famille — champ inexistant en DB)
- Granularité timeline : **mois** (12 colonnes/an) — suffisant visuellement, évite la surcharge
- **Multi-années** : toutes les années sont affichées côte à côte (12 colonnes par an), scroll horizontal libre — même logique que les vues récoltes existantes. Pas de sélecteur d'année, pas de sous-lignes N/N-1. En portrait on voit ~1,5 an, en paysage 2-3 ans. Paramètre API `?years=4` (4 dernières années par défaut).
- Taille des bulles : proportionnelle à `sqrt(qte || poids || 1)`, plafonnée entre r=4 et r=18 px
- Export : fichier **Excel `.xlsx`** généré côté Go via `excelize`, téléchargé directement depuis l'app (ouverture URL dans onglet sur Flutter Web, ou `share_plus` sur mobile)
- Pas d'intégration Google Sheets dans cette itération

---

## Nouvelles routes API

| Méthode | Route | Description |
|---------|-------|-------------|
| `GET` | `/garden/:gardenId/stats/timeline?years=N` | Événements pour la frise (légume × année × mois × action), N dernières années (défaut 4) |
| `GET` | `/garden/:gardenId/stats/monthly?year=YYYY` | Nombre d'actions par mois et type d'action |
| `GET` | `/garden/:gardenId/stats/durations` | Délai moyen semis→récolte par légume/variété |
| `GET` | `/garden/:gardenId/export/excel?year=YYYY` | Export Excel (.xlsx) de tous les action_log du jardin |

Paramètre `year` : si omis → année courante. Toutes les routes sont scopées au `gardenId`.

---

## 1. Backend Go

### Fichiers à créer / modifier

| Fichier | Rôle |
|---------|------|
| `api/models/stats.go` | Requêtes SQL pour les 3 endpoints stats |
| `api/controllers/stats.go` | Handlers Gin + handler export Excel |
| `api/main.go` | Ajout des 4 routes |

### Modèles Go (`api/models/stats.go`)

```go
// Timeline
type TimelineEvent struct {
    Legume string `json:"legume"`
    Annee  int    `json:"annee"`  // N ou N-1
    Mois   int    `json:"mois"`   // 1-12
    Action string `json:"action"`
    Qte    int    `json:"qte"`
    Poids  int    `json:"poids"`
}

// Monthly distribution
type MonthlyCount struct {
    Mois   int    `json:"mois"`
    Action string `json:"action"`
    Count  int    `json:"count"`
}

// Durations
type DurationStat struct {
    Legume  string `json:"legume"`
    Variete string `json:"variete"`
    AvgDays int    `json:"avg_days"`
    Samples int    `json:"samples"`
}
```

### Requêtes SQL clés

```sql
-- Timeline : N dernières années (param years=4 → 4 ans)
SELECT legume,
       CAST(strftime('%Y', date_action) AS INTEGER) AS annee,
       CAST(strftime('%m', date_action) AS INTEGER) AS mois,
       action,
       SUM(qte)   AS qte,
       SUM(poids) AS poids
FROM action_log
WHERE jardin_id = ?
  AND CAST(strftime('%Y', date_action) AS INTEGER) >= CAST(strftime('%Y','now') AS INTEGER) - ? + 1
  AND legume != ''
GROUP BY legume, annee, mois, action
ORDER BY legume, annee, mois, action

-- Répartition mensuelle (toutes actions)
SELECT CAST(strftime('%m', date_action) AS INTEGER) AS mois,
       action,
       COUNT(*) AS count
FROM action_log
WHERE jardin_id = ?
  AND strftime('%Y', date_action) = ?
GROUP BY mois, action
ORDER BY mois, action

-- Durées semis→récolte (auto-jointure sur même légume, même jardin, même année)
SELECT s.legume, s.variete,
       AVG(JULIANDAY(r.date_action) - JULIANDAY(s.date_action)) AS avg_days,
       COUNT(*) AS samples
FROM action_log s
JOIN action_log r
  ON r.jardin_id = s.jardin_id
 AND r.legume    = s.legume
 AND r.action    = 'Récolte'
 AND r.date_action > s.date_action
 AND strftime('%Y', r.date_action) = strftime('%Y', s.date_action)
WHERE s.jardin_id = ?
  AND s.action IN ('Semis', 'Plantation')
GROUP BY s.legume, s.variete
HAVING avg_days > 0
ORDER BY s.legume, s.variete

-- Export CSV : toutes les colonnes de action_log pour un jardin/année
SELECT id, date_action, action, legume, variete, lieu, qte, poids, notes, tags
FROM action_log
WHERE jardin_id = ?
  AND (? = '' OR strftime('%Y', date_action) = ?)
ORDER BY date_action
```

### Export Excel (handler `ExportExcel`)

Dépendance Go : `github.com/xuri/excelize/v2`

- Feuille "Actions" : toutes les colonnes de `action_log` (id, date, action, legume, variete, lieu, qte, poids, notes, tags)
- Ligne d'en-tête en gras, filtre auto activé → prêt pour tableau croisé dynamique Excel
- Réponse avec headers :
```
Content-Type: application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
Content-Disposition: attachment; filename="jardin-YYYY.xlsx"
```

---

## 2. Flutter

### Dépendances à ajouter (`app/pubspec.yaml`)

```yaml
fl_chart: ^0.70.0   # bar charts (monthly distribution)
```

Pas de lib pour la frise — rendu personnalisé avec `CustomPainter`.

### Navigation — deux changements

**1. Tab "Récoltes" → "Bilan"** (`app/lib/home.dart`)
- Renommer le label et changer l'icône : `Icons.bar_chart` (ou `Icons.analytics`)
- Aucun changement structurel — juste le `NavigationDestination`

**2. `stats_view.dart` — ajout de 3 entrées**
La vue existante est déjà une liste de tiles de navigation ; ajouter à la suite :
- "Frise temporelle"
- "Activité par mois"
- "Délais semis → récolte"

**3. Export dans "Données"** (`app/lib/cleanup_view.dart`)
Ajouter une icône `Icons.file_download` dans l'AppBar de `CleanupView`. Tap → `showModalBottomSheet` :
```
┌────────────────────────────┐
│  Exporter les données      │
│  Année : ←  2025  →        │
│  [Télécharger Excel]       │
└────────────────────────────┘
```
Le bouton ouvre `exportExcelUrl(gardenId, year)` dans un nouvel onglet (`url_launcher` ou `dart:html` selon la plateforme). L'export n'apparaît pas dans la vue Bilan.

### A. Frise temporelle (`app/lib/stats_timeline.dart`)

Widget principal : `StatsTimeline(gardenId)`.

**Layout** :
```
┌──────────┬──────────── 2023 ────────────┬──────────── 2024 ────────────┬── 2025 ──...
│ Légume   │ J  F  M  A  M  J  J  A  S  O  N  D │ J  F  M ...│ J  F  M ...
╠══════════╪══════════════════════════════╪══════════════════════════════╪══════
│ Tomate   │         ●     ▲        ■  ■  │    ●  ▲        ■  ■         │  ●  ▲  ...
╠══════════╪══════════════════════════════╪══════════════════════════════╪══════
│ Carotte  │  ●                 ▶  ▶      │ ●              ▶            │  ●  ...
└──────────┴──────────────────────────────┴──────────────────────────────┴──────
```
- **Colonne gauche fixe** (noms légumes, une ligne par légume) — scroll vertical
- **Groupes d'années** en header de section, chacun avec 12 colonnes mois — scroll horizontal libre
- Une ligne par légume (pas de sous-lignes) — scroll vertical
- Séparateur vertical entre groupes d'années pour aider la lecture
- **Légende** en bas fixe : couleurs par action
- Pas de sélecteur d'année — chargement des N dernières années (défaut 4)

**Couleurs par action** :
| Action | Couleur |
|--------|---------|
| Semis | `Colors.green` |
| Repiquage | `Colors.lightBlue` |
| Plantation | `Colors.orange` |
| Récolte | `Colors.red.shade700` |
| Autres | `Colors.grey` |

**Taille des bulles** :
```dart
double _radius(int qte, int poids) {
  final v = (qte > 0 ? qte : 0) + (poids > 0 ? poids / 100 : 0);
  return v > 0 ? (4.0 + sqrt(v) * 1.5).clamp(4.0, 20.0) : 4.0;
}
```

**Rendu** : `CustomPainter` par cellule (légume × mois) — dessine N cercles côte à côte si plusieurs actions dans le même mois, centrés dans la cellule. Tooltip au tap (via `GestureDetector` + `Overlay`).

**Scroll synchronisé** : `LinkedScrollControllerGroup` (ou deux `ScrollController` avec `addListener`) pour synchroniser header et contenu sur l'axe horizontal.

### B. Activité par mois (`app/lib/stats_monthly.dart`)

`fl_chart` `BarChart` groupé.
- X : 12 mois (labels "Jan"…"Déc")
- Y : nombre d'actions
- Une `BarGroup` par mois, avec un `BarChartRodData` par type d'action présent
- Couleurs identiques à la frise temporelle
- Légende en dessous
- Sélecteur d'année identique

### C. Délais semis→récolte (`app/lib/stats_durations.dart`)

`ListView` triée par `avg_days` croissant :
```
┌──────────────────────────────────────────┐
│ Radis            19 jours  (3 mesures)   │
│ Laitue           52 jours  (5 mesures)   │
│ Haricot vert     68 jours  (2 mesures)   │
│ Tomate          112 jours  (7 mesures)   │
└──────────────────────────────────────────┘
```
Variété affichée en sous-titre gris si non vide.

### Modèles Dart (`app/lib/stats_model.dart`) — nouveau fichier

```dart
class TimelineEvent { String legume, action; int annee, mois, qte, poids; }
// Années disponibles déduites des données reçues — pas de param côté Flutter
// Flutter regroupe par (legume, annee, mois) pour le rendu
class MonthlyCount  { int mois, count; String action; }
class DurationStat  { String legume, variete; int avgDays, samples; }
```

### Méthodes `ApiService` à ajouter (`app/lib/api_service.dart`)

```dart
Future<List<TimelineEvent>> getTimeline(String gardenId, int year)
Future<List<MonthlyCount>>  getMonthlyStats(String gardenId, int year)
Future<List<DurationStat>>  getDurationStats(String gardenId)
String                      exportExcelUrl(String gardenId, int? year) // url directe → téléchargement .xlsx
```

---

## Vérification

1. Tab bas renommé "Bilan" avec icône `bar_chart` ; les deux vues récoltes existantes sont toujours accessibles
2. `GET /garden/:id/stats/timeline?years=4` → JSON des 4 dernières années ; frise affichée avec une ligne par légume, groupes d'années séparés, scroll horizontal multi-années
3. `GET /garden/:id/stats/monthly?year=2024` → bar chart avec 12 groupes de barres colorées par action
4. `GET /garden/:id/stats/durations` → liste triée par durée croissante
5. Onglet "Données" : icône download dans AppBar → bottom sheet avec sélecteur d'année → `GET /garden/:id/export/excel?year=2025` → fichier `.xlsx` téléchargé, ouvrable dans Excel avec filtre auto activé
6. Changement d'année dans la frise → N-1 se recalcule automatiquement (N-2 pour la nouvelle sélection)
