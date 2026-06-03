# Service API — Garden Planner

Service back-end en Go + Gin. Gère les jardins, les journaux d'actions, les récoltes et le stockage de photos.

- **Port** : `8081`
- **Base de données** : SQLite (chemin configuré via `SQLITE_PATH`)
- **Bucket photos** : Google Cloud Storage (`jactez01`)

---

## Variables d'environnement

| Variable | Défaut | Description |
|----------|--------|-------------|
| `SQLITE_PATH` | `./garden-planner.db` | Chemin vers le fichier SQLite |

---

## Routes

### Jardins

| Méthode | Route | Description |
|---------|-------|-------------|
| `GET` | `/gardens` | Liste les jardins de l'utilisateur |
| `POST` | `/garden` | Crée un jardin |
| `DELETE` | `/garden/:id` | Supprime un jardin |

**Authentification** : l'en-tête `Authorization` doit contenir le `userId` brut (pas de JWT).

**Corps `POST /garden`** :
```json
{
  "nom": "Mon potager",
  "notes": "...",
  "localisation": "Bretagne",
  "surface": 50,
  "moisFinRecolte": 10,
  "moisFinSemis": 4,
  "meteofSite": "56243001",
  "jardiniers": [
    { "userId": "user-uuid", "role": "admin" }
  ]
}
```

---

### Journaux d'actions

| Méthode | Route | Description |
|---------|-------|-------------|
| `GET` | `/logs` | Tous les logs (15 derniers mois) |
| `GET` | `/garden/:gardenId/logs` | Logs d'un jardin spécifique |
| `POST` | `/log` | Crée un log |
| `POST` | `/logs` | Import en masse (tableau JSON) |
| `DELETE` | `/log/:id` | Supprime un log (et ses photos GCS) |
| `PUT` | `/logs/garden?value=<id>` | Rattache tous les logs à un jardin |

**Filtre temporel** : `GET /logs` ne retourne que les logs des 15 derniers mois.

**Corps `POST /log`** :
```json
{
  "jardinId": "uuid-jardin",
  "dateAction": "2024-06-01",
  "action": "Semis",
  "statut": "Fait",
  "lieu": "Serre",
  "legume": "Tomate",
  "variete": "Cœur de bœuf",
  "qte": 6,
  "poids": 0,
  "notes": "...",
  "photos": [],
  "tags": ["printemps"]
}
```

**Valeurs `action` courantes** : `Semis`, `Plantation`, `Récolte`, `Arrosage`, `Traitement`, `Taille`

---

### Référentiels

| Méthode | Route | Description |
|---------|-------|-------------|
| `GET` | `/legumes` | Liste des légumes (extraits des logs) |
| `GET` | `/tags` | Tags distincts présents dans les logs |
| `GET` | `/lieux` | Lieux distincts présents dans les logs |
| `GET` | `/recoltes` | Cumuls de récolte par légume et année |
| `GET` | `/recoltes/lieux` | Cumuls de récolte par lieu, légume et année |

> Les données de récolte rattachent les mois janvier–mars à l'année N-1 (cycle maraîcher).

---

### Photos (Google Cloud Storage)

| Méthode | Route | Description |
|---------|-------|-------------|
| `POST` | `/photo` | Upload un fichier vers le bucket `jactez01` |
| `DELETE` | `/photo/:id` | Supprime un objet du bucket |

**Upload** : multipart/form-data, champ `file`.

**Réponse** :
```json
{ "message": "file uploaded successfully", "pathname": "/jactez01/nom-fichier.jpg" }
```

> Les credentials GCS sont actuellement codés en dur dans `controllers/cloudBucket.go`.

---

## Structure des tables SQLite

### `garden`

| Colonne | Type | Description |
|---------|------|-------------|
| `id` | TEXT (UUID) | Clé primaire |
| `nom` | TEXT | Nom du jardin |
| `localisation` | TEXT | Localisation géographique |
| `surface` | INTEGER | Surface en m² |
| `mois_fin_recolte` | INTEGER | Mois de fin de récolte |
| `mois_fin_semis` | INTEGER | Mois de fin de semis |
| `meteofrance_site` | TEXT | Code station MétéoFrance |
| `notes` | TEXT | |

### `garden_jardinier`

| Colonne | Type | Description |
|---------|------|-------------|
| `garden_id` | TEXT | FK → garden.id |
| `user_id` | TEXT | Identifiant utilisateur |
| `role` | TEXT | `admin` ou `lecteur` |

### `action_log`

| Colonne | Type | Description |
|---------|------|-------------|
| `id` | TEXT (UUID) | Clé primaire |
| `parent_id` | TEXT | ID parent (pour sous-actions) |
| `jardin_id` | TEXT | FK → garden.id |
| `date_action` | TEXT | Format `YYYY-MM-DD` |
| `action` | TEXT | Type d'action |
| `statut` | TEXT | Statut de l'action |
| `lieu` | TEXT | Emplacement dans le jardin |
| `legume` | TEXT | Nom du légume |
| `variete` | TEXT | Variété |
| `qte` | INTEGER | Quantité |
| `poids` | INTEGER | Poids en grammes |
| `notes` | TEXT | |
| `photos` | TEXT | Tableau JSON d'URLs |
| `tags` | TEXT | Tableau JSON de tags |

---

## Lancer en local

```bash
cd api
go run .
```

L'API écoute sur `http://0.0.0.0:8081`.
