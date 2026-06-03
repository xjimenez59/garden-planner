# Déploiement — Garden Planner

## Architecture

L'application est composée de trois services Docker :

| Service | Rôle | Port exposé |
|---------|------|-------------|
| `api` | Back-end Go + Gin, accès à la base SQLite | 8081 |
| `webapp` | Front-end Flutter Web servi par Nginx | 8080 (HTTP) / 8083 (HTTPS) |
| `meteo` | Service météo Go, base SQLite dédiée | 8082 |

Nginx (dans `webapp`) fait office de reverse proxy : les requêtes vers `/api/` sont proxiées vers le service `api:8081`.

---

## Cibles de build

Chaque `Dockerfile` contient plusieurs cibles (`target`) selon l'environnement :

| Cible | Plateforme | Usage |
|-------|-----------|-------|
| `debug` | linux/amd64 (image golang complète) | Développement local |
| `windows` | linux/amd64 | Dev sous Windows avec Docker Desktop |
| `nas` | linux/amd64 (alpine) | NAS Synology ou serveur x86 |
| `raspberry` | linux/arm64 (alpine) | Raspberry Pi 4 |

---

## Déploiement sur NAS (linux/amd64)

### 1. Prérequis sur l'hôte

- Docker + Docker Compose installés
- Répertoires de données créés :

```bash
mkdir -p /local/containers-data/garden-planner/api
mkdir -p /local/containers-data/garden-planner/meteo-sqlite
```

### 2. Construction des images

Depuis la racine du dépôt, sur une machine avec Docker disponible :

```bash
# API
docker build --target nas -t garden-planner-api:latest ./api

# Webapp
docker build --target nas -t garden-planner-webapp:latest ./app

# Meteo
docker build --target nas -t garden-planner-meteo:latest ./meteo
```

> Pour une cross-compilation depuis Windows ou Mac, ajouter `--platform linux/amd64`.

### 3. Transfert des images vers le NAS

```bash
docker save garden-planner-api:latest   | gzip > garden-planner-api.tar.gz
docker save garden-planner-webapp:latest | gzip > garden-planner-webapp.tar.gz
docker save garden-planner-meteo:latest  | gzip > garden-planner-meteo.tar.gz

scp *.tar.gz user@nas:/tmp/

# Sur le NAS :
docker load < /tmp/garden-planner-api.tar.gz
docker load < /tmp/garden-planner-webapp.tar.gz
docker load < /tmp/garden-planner-meteo.tar.gz
```

### 4. Démarrage

Copier `docker-compose.yaml` et le dossier `nginx/` sur le NAS, puis :

```bash
docker compose -f docker-compose.yaml up -d
```

### 5. Vérification

```bash
docker compose ps
curl http://localhost:8081/gardens    # API
curl http://localhost:8080            # Webapp
curl http://localhost:8082            # Meteo
```

---

## Déploiement sur Raspberry Pi (linux/arm64)

### 1. Prérequis

- Raspberry Pi 4 avec Raspberry Pi OS 64-bit
- Docker installé : `curl -fsSL https://get.docker.com | sh`

### 2. Construction des images

Depuis la racine du dépôt (sur le Pi ou en cross-compilation) :

```bash
docker build --target raspberry -t garden-planner-api:latest ./api
docker build --target raspberry -t garden-planner-webapp:latest ./app
docker build --target raspberry -t garden-planner-meteo:latest ./meteo
```

### 3. Démarrage

```bash
docker compose -f docker-compose-rpi.yaml up -d
```

> Le fichier `docker-compose-rpi.yaml` utilise encore MongoDB. Adapter les volumes et variables d'environnement si la migration SQLite a été effectuée.

---

## Configuration Nginx

La configuration Nginx est montée depuis `nginx/conf/garden-planner.conf` dans le conteneur `webapp`.

### HTTP (port 80)

- Sert les fichiers statiques Flutter depuis `/usr/share/nginx/html`
- Proxie `/api/` vers `http://api:8081/`

### HTTPS (port 443)

Pour activer HTTPS, décommenter le bloc `server { listen 443 ssl; }` dans `nginx/conf/garden-planner.conf` et monter les certificats :

```yaml
# dans docker-compose.yaml, service webapp :
volumes:
  - ./nginx/conf/:/etc/nginx/conf.d/:ro
  - /etc/letsencrypt/:/etc/nginx/ssl/:ro
```

Les certificats attendus (Let's Encrypt) :
- `/etc/nginx/ssl/live/gardenplanner.app.jactez.com/fullchain.pem`
- `/etc/nginx/ssl/live/gardenplanner.app.jactez.com/privkey.pem`

> Le domaine configuré est `gardenplanner.app.jactez.com`. Modifier `nginx/conf/garden-planner.conf` pour changer le `server_name`.

---

## Variables d'environnement

### Service `api`

| Variable | Défaut | Description |
|----------|--------|-------------|
| `SQLITE_PATH` | `./garden-planner.db` | Chemin du fichier SQLite |

### Service `meteo`

| Variable | Défaut | Description |
|----------|--------|-------------|
| `SQLITE_PATH` | `./meteo.db` | Chemin du fichier SQLite météo |

---

## Persistance des données

| Service | Chemin hôte (NAS) | Chemin conteneur | Contenu |
|---------|-------------------|-----------------|---------|
| `api` | `/local/containers-data/garden-planner/api` | `/opt/app/api` | `garden-planner.db` |
| `meteo` | `/local/containers-data/garden-planner/meteo-sqlite` | `/data` | `meteo.db` |

> Les bases SQLite sont des fichiers uniques. Sauvegarder ces répertoires suffit pour une sauvegarde complète des données.

---

## Migration depuis MongoDB

Si les données sont stockées dans MongoDB, des scripts de migration sont disponibles.

### Migration API (garden-planner.db)

```bash
cd api/cmd/migrate
go run . \
  -mongo "mongodb://user:password@host:27017/?authSource=admin" \
  -db garden-planner \
  -sqlite ./garden-planner.db
```

Ou via variables d'environnement :

```bash
export MONGO_HOST=host
export MONGO_PORT=27017
export MONGO_USER=user
export MONGO_PWD=password
export MONGO_DBNAME=garden-planner
export SQLITE_PATH=./garden-planner.db
go run .
```

### Migration Meteo (meteo.db)

```bash
cd meteo/cmd/migrate
go run . \
  -mongo "mongodb://user:password@host:27017/?authSource=admin" \
  -db garden-planner \
  -sqlite ./meteo.db
```

> Les scripts sont idempotents (`INSERT OR REPLACE`) : ils peuvent être relancés sans risque de doublon.

---

## Mise à jour de l'application

```bash
# 1. Reconstruire les images avec les nouvelles sources
docker build --target nas -t garden-planner-api:latest ./api
docker build --target nas -t garden-planner-webapp:latest ./app
docker build --target nas -t garden-planner-meteo:latest ./meteo

# 2. Redémarrer les conteneurs
docker compose up -d
```

Les données SQLite étant dans des volumes bind, elles sont conservées entre les mises à jour.

---

## Développement local

```bash
# Lancer tous les services en mode debug
docker compose -f .devcontainer/docker-compose.yml up
```

Ou directement sans Docker :

```bash
# API
cd api && go run .

# Meteo
cd meteo && go run .

# Webapp
cd app && flutter run -d chrome
```
