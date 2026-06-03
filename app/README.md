# Application Web — Garden Planner

Front-end Flutter compilé en web, servi par Nginx. Communique avec le service `api` (port 8081) et le service `meteo` (port 8082).

---

## Configuration

L'URL des services est configurable via des variables `--dart-define` au moment du build Flutter.

| Variable | Défaut | Description |
|----------|--------|-------------|
| `METEO_BASE_URL` | `http://localhost:8082` | URL du service météo |

L'URL de l'API est définie dans [lib/constants.dart](lib/constants.dart) (`ApiConstants.baseUrl`). En production, le trafic `/api/` est proxié par Nginx vers `http://api:8081`.

**Authentification** : l'en-tête `Authorization` est renseigné avec l'ID utilisateur Firebase (`ApiConstants.defaultUser`). L'intégration Firebase Auth est commentée dans le code — l'ID est actuellement en dur.

---

## Build de production

```bash
# Build standard (URL meteo par défaut)
flutter build web

# Build avec URL meteo custom
flutter build web --dart-define=METEO_BASE_URL=https://gardenplanner.app.jactez.com/meteo
```

Les fichiers compilés sont générés dans `build/web/` et copiés dans l'image Docker vers `/usr/share/nginx/html`.

---

## Configuration Nginx

Le fichier [nginx.conf](nginx.conf) est embarqué dans l'image Docker (cibles `nas` et `raspberry`).

- `GET /` → fichiers statiques Flutter
- `GET /api/*` → proxy vers `http://api:8081/`

En développement local, la configuration étendue se trouve dans [../nginx/conf/garden-planner.conf](../nginx/conf/garden-planner.conf) (montée en volume, supporte HTTP + HTTPS).

---

## Services consommés

### API (`api_service.dart`)

| Méthode | Appel Flutter | Route API |
|---------|---------------|-----------|
| Lister les jardins | `getGardens()` | `GET /gardens` |
| Créer un jardin | `postGarden(garden)` | `POST /garden` |
| Lister les logs | `getLogs(jardin)` | `GET /garden/:id/logs` |
| Créer un log | `postLog(log)` | `POST /log` |
| Import en masse | `postLogs(logs)` | `POST /logs` |
| Supprimer un log | `deleteLog(id)` | `DELETE /log/:id` |
| Lister les légumes | `getLegumes()` | `GET /legumes` |
| Lister les tags | `getTags()` | `GET /tags` |
| Lister les lieux | `getLieux()` | `GET /lieux` |
| Récoltes par légume | `getRecoltes()` | `GET /recoltes` |
| Récoltes par lieu | `getRecolteAnnuelle()` | `GET /recoltes/lieux` |
| Uploader une photo | `postPicture(bytes)` | `POST /photo` |
| Supprimer une photo | `deletePicture(url)` | `DELETE /photo/:id` |

### Météo (`meteo_service.dart`)

| Méthode | Description | Route Meteo |
|---------|-------------|-------------|
| `getMeteo(site, dateDeb, dateFin)` | Relevés journaliers MétéoFrance | `GET /meteo?station=&date_deb=&date_fin=` |

Le code station (`site`) est lu depuis le champ `meteofSite` du jardin sélectionné.

**Champs météo utilisés par l'app** :

| Champ | Description |
|-------|-------------|
| `RR` | Précipitations cumulées (mm) |
| `TM` / `TN` / `TX` | Températures moyenne / min / max (°C) |
| `FFM` / `FXI` | Vent moyen / rafale max (m/s) |
| `DXY` | Direction vent dominant (°) |
| `INST` / `SIGMA` | Durée et fraction d'insolation |
| `DG` | Durée de gel (min) |
| `NB300` | Nb d'heures sous 3°C |

---

## Développement local

```bash
cd app

# Web (navigateur)
flutter run -d chrome

# Avec URL meteo custom
flutter run -d chrome --dart-define=METEO_BASE_URL=http://localhost:8082
```

L'app attend l'API sur `http://localhost:8081` par défaut.
