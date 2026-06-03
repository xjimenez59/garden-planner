# Service Meteo — Garden Planner

Service back-end en Go + Gin. Agrège des données météo depuis trois fournisseurs (MétéoFrance, AccuWeather, InfoClimat) et calcule les cycles lunaires/biodynamiques.

- **Port** : `8082`
- **Base de données** : SQLite (chemin configuré via `SQLITE_PATH`)

---

## Variables d'environnement

| Variable | Défaut | Description |
|----------|--------|-------------|
| `SQLITE_PATH` | `./meteo.db` | Chemin vers le fichier SQLite |
| `METEOFRANCE_BASIC_AUTH` | *(valeur encodée)* | Credential Basic Auth pour l'API MétéoFrance |
| `METEOFRANCE_API_BASE` | `https://public-api.meteofrance.fr/public/DPClim/v1` | URL de base de l'API MétéoFrance |

---

## Routes

### MétéoFrance (source principale recommandée)

| Méthode | Route | Description |
|---------|-------|-------------|
| `GET` | `/meteofrance/quotidien` | Commande des données quotidiennes à MF |
| `GET` | `/meteofrance/resultats` | Récupère le fichier CSV commandé et le sauvegarde |
| `GET` | `/meteo` | Lit les données météo sauvegardées en base |

#### Import des données (2 étapes asynchrones)

**Étape 1 — Commander les données** :
```
GET /meteofrance/quotidien?station=56243001&date_deb=20240601&date_fin=20240630
```

Réponse (200) : JSON contenant le numéro de commande
```json
{ "elaboreProduitAvecDemandeResponse": { "return": "2026017935560" } }
```

**Étape 2 — Récupérer le fichier** :
```
GET /meteofrance/resultats?id_cmde=abc123xyz
```

- Réponse `202` : fichier en cours de préparation côté MF, relancer dans quelques secondes
- Réponse `200` : données sauvegardées en base
```json
{ "saved": 30 }
```

**Lecture des données sauvegardées** :
```
GET /meteo?station=56243001&date_deb=20240601&date_fin=20240630
```

**Format des dates** : `YYYYMMDD` (ex : `20240601`)

**Code station** : identifiant MétéoFrance à 8 chiffres (ex : `56243001` pour Le Tour-du-Parc).  
Le code est renseigné dans la configuration du jardin (`meteofSite`).

---

### AccuWeather

| Méthode | Route | Description |
|---------|-------|-------------|
| `GET` | `/meteo/accuweather/location/search?cp=<code_postal>` | Recherche de locations |
| `GET` | `/meteo/accuweather/location/import?cp=<code_postal>` | Import des locations en base |
| `GET` | `/meteo/accuweather/:location/past24h` | Données brutes des 24h (sans sauvegarde) |
| `GET` | `/meteo/accuweather/:location/past24h/import` | Import des données 24h en base |

**Flux d'initialisation (à faire une seule fois par zone)** :
```
# 1. Chercher la location par code postal
GET /meteo/accuweather/location/search?cp=56370

# 2. Importer la location en base (même paramètre)
GET /meteo/accuweather/location/import?cp=56370
```

**Flux d'import quotidien** :
```
# Remplacer <location> par la clé AccuWeather (ex: 166808_PC)
GET /meteo/accuweather/166808_PC/past24h/import
```

> La clé AccuWeather pour Tropark est `166808_PC`.

---

### InfoClimat

| Méthode | Route | Description |
|---------|-------|-------------|
| `GET` | `/meteo/infoclimat/:site/:date` | Relevés horaires d'une station publique |

**Exemple** :
```
GET /meteo/infoclimat/STATIC0095/2024-10-13
```

> Ce fournisseur est actuellement en mode mock (pas d'appel réel à l'API InfoClimat).

---

### Lune & calendrier biodynamique

| Méthode | Route | Description |
|---------|-------|-------------|
| `GET` | `/lune` | Infos lunaires du jour |
| `GET` | `/lune?date=YYYY-MM-DD` | Infos lunaires pour une date donnée |

**Réponse** :
```json
{
  "date": "2024-06-15",
  "revolution_periodique": "lune_montante",
  "revolution_cyclique": "lune_croissante",
  "jour_biodynamique": "Fruit",
  "signe_zodiaque": "Lion",
  "noeud_ascendant_longitude": 12.34,
  "prochain_noeud_ascendant": "2024-07-01",
  "prochain_noeud_descendant": "2024-07-15",
  "prochain_perigee": "2024-06-22",
  "prochain_apogee": "2024-07-06",
  "etat_orbite": "apogee"
}
```

**Champs** :

| Champ | Valeurs possibles | Description |
|-------|-------------------|-------------|
| `revolution_periodique` | `lune_montante` / `lune_descendante` | Cycle tropical (~27,32 j) — utilisé en biodynamie |
| `revolution_cyclique` | `lune_croissante` / `lune_decroissante` | Cycle synodique (~29,53 j) — phases lunaires |
| `jour_biodynamique` | `Fruit` / `Racine` / `Fleur` / `Feuille` | Calendrier Maria Thun |
| `signe_zodiaque` | Bélier … Poissons | Signe sidéral (ayanamsha de Lahiri) |
| `etat_orbite` | `perigee` / `apogee` / `entre_perigee_apogee` | Position sur l'orbite elliptique |

---

## Automatisation depuis crontab-ui (ou tout scheduler)

Le service ne possède pas de tâche planifiée interne. Les mises à jour doivent être déclenchées par un scheduler externe (ex : conteneur `crontab-ui`, `ofelia`, ou `cron` système).

### Mise à jour quotidienne MétéoFrance

Les deux étapes doivent être chaînées car MF prépare les fichiers de manière asynchrone.

**Script shell recommandé** (`/scripts/update-meteo-mf.sh`) :

> Important : créer ce fichier depuis un environnement Linux ou le convertir avec `dos2unix` avant utilisation.
> Les fins de ligne Windows (CRLF) cassent l'exécution sous BusyBox/Alpine.

```sh
#!/bin/sh
METEO_HOST="http://meteo:8082"
STATION="56243001"

# Hier au format YYYYMMDD - compatible BusyBox/Alpine (date -d "yesterday" non supporté)
DATE_DEB=$(date -d @$(($(date +%s) - 86400)) +%Y%m%d)
DATE_FIN=$DATE_DEB

# Étape 1 : Commander les données
RESPONSE=$(curl -sf "$METEO_HOST/meteofrance/quotidien?station=$STATION&date_deb=$DATE_DEB&date_fin=$DATE_FIN")
echo "Réponse commande : $RESPONSE"

ID_CMDE=$(echo "$RESPONSE" | grep -o '"return": *"[^"]*"' | cut -d'"' -f4)

if [ -z "$ID_CMDE" ]; then
  echo "Erreur : id de commande introuvable dans : $RESPONSE"
  exit 1
fi

echo "Commande MétéoFrance : $ID_CMDE"

# Étape 2 : Attendre et récupérer le fichier (10 tentatives)
i=1
while [ $i -le 10 ]; do
  STATUS=$(curl -s -o /dev/null -w "%{http_code}" \
    "$METEO_HOST/meteofrance/resultats?id_cmde=$ID_CMDE")
  if [ "$STATUS" = "200" ]; then
    echo "Données importées avec succès"
    exit 0
  fi
  echo "Tentative $i : en attente (HTTP $STATUS), retry dans 10s..."
  sleep 10
  i=$((i + 1))
done

echo "Échec : fichier non disponible après 10 tentatives"
exit 1
```

**Entrée crontab** (chaque jour à 6h) :
```
0 6 * * * sh /scripts/update-meteo-mf.sh
```

### Mise à jour quotidienne AccuWeather

```bash
#!/bin/sh
METEO_HOST="http://meteo:8082"
LOCATION="166808_PC"

curl -sf "$METEO_HOST/meteo/accuweather/$LOCATION/past24h/import"
echo "Import AccuWeather terminé"
```

**Entrée crontab** (chaque jour à 7h) :
```
0 7 * * * curl -sf http://meteo:8082/meteo/accuweather/166808_PC/past24h/import
```

> Dans un réseau Docker Compose, `meteo` est le nom du service et `8082` le port interne.  
> Depuis l'hôte ou un autre réseau, remplacer par `http://localhost:8082` ou l'IP du serveur.

---

## Structure des tables SQLite

### `weather_location`

| Colonne | Description |
|---------|-------------|
| `id` | TEXT (UUID ou ObjectID Mongo migré) |
| `key` | Clé AccuWeather (ex : `166808_PC`) |
| `key_type` | Type de clé |
| `localized_name` | Nom de la localité |
| `postal_code` | Code postal |
| `latitude` / `longitude` | Coordonnées |
| `elevation` | Altitude (m) |

### `weather_hourly`

Table large (~100 colonnes). Contient les relevés horaires AccuWeather :
température, humidité, pression, précipitations, vent, UV, visibilité, etc.

Clé : `(id)` — INSERT OR REPLACE.

### `meteofrance_quotidien`

Données quotidiennes MétéoFrance (CSV parsé) :

| Colonne clé | Description |
|-------------|-------------|
| `POSTE` | Code station (ex : `56243001`) |
| `DATE` | Format `YYYYMMDD` |
| `RR` | Précipitations (mm) |
| `TN` / `TX` / `TM` | Températures min / max / moyenne (°C) |
| `FFM` | Vitesse vent moyen (m/s) |
| `FXI` | Vitesse vent max (m/s) |
| `DXY` | Direction vent dominant (°) |

---

## Lancer en local

```bash
cd meteo
go run .
```

Le service écoute sur `http://0.0.0.0:8082`.
