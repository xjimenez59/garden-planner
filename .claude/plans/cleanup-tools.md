# Plan — Outils de nettoyage des données (légumes, variétés, lieux, tags)

> **Statut : IMPLÉMENTÉ** (branche `feature/cleanup-tools`, commit `c7ceaae`)
> Ce plan est conservé pour référence et pour les futures évolutions du module nettoyage.

## Contexte

Les légumes, variétés, lieux et tags sont saisis librement dans les `action_log`, sans table de référence. Des doublons apparaissent (casse, fautes d'orthographe). Ce module permet de visualiser ces valeurs avec leur fréquence d'utilisation et d'effectuer des opérations de nettoyage : renommer, fusionner, supprimer.

**Architecture** : `legume`, `variete`, `lieu` sont des champs `TEXT` libres dans `action_log`. `tags` est un tableau JSON sérialisé. Renommer = `UPDATE action_log SET legume = 'new' WHERE legume = 'old'`.

---

## Décisions d'architecture

- **Scoping per-jardin** pour toutes les opérations (légumes, variétés, lieux, tags) — pas d'effet de bord inter-jardins
- **Autocomplétion cross-jardins préservée** — `/legumes`, `/tags` lisent sur tous les jardins de l'utilisateur
- **Catalogue de référence statique** (`api/data/legumes_reference.json`, 59 légumes, embarqué via `//go:embed`) — endpoint `GET /legumes/reference`
- **Suggestions par trigrammes** (seuil 0.30) — signale les valeurs proches d'un nom de référence
- **Barre alphabétique verticale** — affichée dès ≥ 20 items, scroll vers la première lettre
- **Tri Flutter-side** — SQLite sort les accents après Z ; le tri utilise `_norm()` (suppression diacritiques)

---

## Routes implémentées

| Méthode | Route | Contrôleur |
|---------|-------|------------|
| `GET` | `/legumes/reference` | `GetLegumesReference` |
| `GET` | `/garden/:gardenId/cleanup/:field` | `GetCleanupList` |
| `POST` | `/garden/:gardenId/cleanup/rename` | `RenameCleanupValue` |
| `DELETE` | `/garden/:gardenId/cleanup/value` | `DeleteCleanupValue` |

---

## Fichiers clés

| Fichier | Rôle |
|---------|------|
| `api/data/legumes_reference.json` | Catalogue 59 légumes avec variétés |
| `api/data/embed.go` | `//go:embed legumes_reference.json` |
| `api/models/cleanup.go` | SQL : comptages, rename, delete (y compris tags JSON) |
| `api/controllers/cleanup.go` | Handlers Gin |
| `app/lib/cleanup_view.dart` | UI nettoyage : 4 tabs, alpha-bar, suggestions trigrammes, actions bottom sheet |
| `app/lib/cleanup_model.dart` | `CleanupItem`, `LegumeReference` |
| `app/lib/list_selector.dart` | Support `getSecondaryOptions` (items référence en gris) |
| `app/lib/action_detail.dart` | Autocomplétion fusionnée : items utilisateur (noir) + référence (gris) |

---

## Points d'attention pour évolutions futures

- **Groupement par famille** : ajouter un champ `famille` dans `legumes_reference.json` (ex. "Solanacées", "Cucurbitacées") pour regrouper dans la `CleanupView` et la frise temporelle des stats
- **Seuil trigrammes** : 0.30 calibré sur "arico-vers" ↔ "Haricot vert" (~33%). Ajuster si trop de faux positifs.
- **Route DELETE /garden/:gardenId** : a nécessité de renommer le param de `:id` → `:gardenId` dans `gardens.go` pour éviter le conflit httprouter avec `DELETE /garden/:gardenId/cleanup/value`
- **Tags** : reconstruction JSON via `json_group_array` / `json_each` — tester sur des logs avec beaucoup de tags
