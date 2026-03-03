// Script de migration des données de MongoDB vers SQLite.
//
// Usage :
//
//	go run . \
//	  -mongo "mongodb://root:password@localhost:27017/?authSource=admin" \
//	  -db garden-planner \
//	  -sqlite ./garden-planner.db
//
// Ou via variables d'environnement :
//
//	MONGO_HOST=localhost MONGO_PORT=27017 MONGO_USER=root MONGO_PWD=secret \
//	MONGO_DBNAME=garden-planner SQLITE_PATH=./garden-planner.db go run .
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	_ "modernc.org/sqlite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// --- paramètres
	mongoURI := flag.String("mongo", buildMongoURI(), "URI de connexion MongoDB")
	mongoDBName := flag.String("db", envOr("MONGO_DBNAME", "garden-planner"), "Nom de la base MongoDB")
	sqlitePath := flag.String("sqlite", envOr("SQLITE_PATH", "./garden-planner.db"), "Chemin du fichier SQLite")
	flag.Parse()

	log.Printf("Connexion MongoDB : %s / %s", *mongoURI, *mongoDBName)
	log.Printf("Destination SQLite : %s", *sqlitePath)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// --- connexion MongoDB
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(*mongoURI))
	if err != nil {
		log.Fatalf("Erreur connexion MongoDB : %v", err)
	}
	defer mongoClient.Disconnect(ctx)
	if err = mongoClient.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB ne répond pas : %v", err)
	}
	mongoDB := mongoClient.Database(*mongoDBName)
	log.Println("Connexion MongoDB OK")

	// --- connexion SQLite
	sqliteDB, err := sql.Open("sqlite", *sqlitePath)
	if err != nil {
		log.Fatalf("Erreur ouverture SQLite : %v", err)
	}
	defer sqliteDB.Close()
	sqliteDB.Exec("PRAGMA foreign_keys = ON")
	sqliteDB.Exec("PRAGMA journal_mode = WAL")

	if err = createTables(sqliteDB); err != nil {
		log.Fatalf("Erreur création tables : %v", err)
	}
	log.Println("Tables SQLite prêtes")

	// --- migration
	gardens, err := migrateGardens(ctx, mongoDB, sqliteDB)
	if err != nil {
		log.Fatalf("Erreur migration jardins : %v", err)
	}
	log.Printf("Jardins migrés : %d", gardens)

	logs, err := migrateActionLogs(ctx, mongoDB, sqliteDB)
	if err != nil {
		log.Fatalf("Erreur migration logs : %v", err)
	}
	log.Printf("Logs d'action migrés : %d", logs)

	log.Println("Migration terminée avec succès.")
}

// ---- helpers ----------------------------------------------------------------

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func buildMongoURI() string {
	host := envOr("MONGO_HOST", "localhost")
	port := envOr("MONGO_PORT", "27017")
	user := envOr("MONGO_USER", "root")
	pwd := envOr("MONGO_PWD", "")
	return fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=admin&directConnection=true", user, pwd, host, port)
}

func objectIDToStr(id primitive.ObjectID) string {
	return id.Hex()
}

func dateTimeToStr(dt primitive.DateTime) string {
	return dt.Time().UTC().Format("2006-01-02")
}

// ---- création des tables ----------------------------------------------------

func createTables(db *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS garden (
			id TEXT PRIMARY KEY,
			nom TEXT NOT NULL DEFAULT '',
			notes TEXT NOT NULL DEFAULT '',
			mois_fin_recolte INTEGER NOT NULL DEFAULT 0,
			mois_fin_semis INTEGER NOT NULL DEFAULT 0,
			localisation TEXT NOT NULL DEFAULT '',
			surface INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS garden_jardinier (
			garden_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			role TEXT NOT NULL DEFAULT '',
			PRIMARY KEY (garden_id, user_id),
			FOREIGN KEY (garden_id) REFERENCES garden(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS action_log (
			id TEXT PRIMARY KEY,
			parent_id TEXT NOT NULL DEFAULT '',
			jardin_id TEXT NOT NULL DEFAULT '',
			date_action TEXT NOT NULL DEFAULT '',
			action TEXT NOT NULL DEFAULT '',
			statut TEXT NOT NULL DEFAULT '',
			lieu TEXT NOT NULL DEFAULT '',
			legume TEXT NOT NULL DEFAULT '',
			variete TEXT NOT NULL DEFAULT '',
			qte INTEGER NOT NULL DEFAULT 0,
			poids INTEGER NOT NULL DEFAULT 0,
			notes TEXT NOT NULL DEFAULT '',
			photos TEXT NOT NULL DEFAULT '[]',
			tags TEXT NOT NULL DEFAULT '[]'
		)`,
	}
	for _, stmt := range statements {
		if _, err := db.Exec(stmt); err != nil {
			return err
		}
	}
	return nil
}

// ---- migration des jardins --------------------------------------------------

func migrateGardens(ctx context.Context, mongoDB *mongo.Database, sqliteDB *sql.DB) (int, error) {
	type MongoGardenRole struct {
		UserID string `bson:"userId"`
		Role   string `bson:"role"`
	}
	type MongoGarden struct {
		ID             primitive.ObjectID `bson:"_id"`
		Nom            string             `bson:"nom"`
		Notes          string             `bson:"notes"`
		MoisFinRecolte int                `bson:"moisFinRecolte"`
		MoisFinSemis   int                `bson:"moisFinSemis"`
		Localisation   string             `bson:"localisation"`
		Surface        int                `bson:"surface"`
		Jardiniers     []MongoGardenRole  `bson:"jardiniers"`
	}

	cursor, err := mongoDB.Collection("garden").Find(ctx, bson.D{})
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	tx, err := sqliteDB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	count := 0
	for cursor.Next(ctx) {
		var g MongoGarden
		if err := cursor.Decode(&g); err != nil {
			return 0, fmt.Errorf("décodage jardin : %w", err)
		}

		id := objectIDToStr(g.ID)
		_, err = tx.ExecContext(ctx, `
			INSERT OR REPLACE INTO garden (id, nom, notes, mois_fin_recolte, mois_fin_semis, localisation, surface)
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
			id, g.Nom, g.Notes, g.MoisFinRecolte, g.MoisFinSemis, g.Localisation, g.Surface)
		if err != nil {
			return 0, fmt.Errorf("insert jardin %s : %w", id, err)
		}

		// supprimer les anciens jardiniers au cas où on rejoue la migration
		tx.ExecContext(ctx, `DELETE FROM garden_jardinier WHERE garden_id = ?`, id)

		for _, j := range g.Jardiniers {
			_, err = tx.ExecContext(ctx, `INSERT INTO garden_jardinier (garden_id, user_id, role) VALUES (?, ?, ?)`,
				id, j.UserID, j.Role)
			if err != nil {
				return 0, fmt.Errorf("insert jardinier %s/%s : %w", id, j.UserID, err)
			}
		}
		count++
	}
	if err := cursor.Err(); err != nil {
		return 0, err
	}

	return count, tx.Commit()
}

// ---- migration des logs d'action --------------------------------------------

func migrateActionLogs(ctx context.Context, mongoDB *mongo.Database, sqliteDB *sql.DB) (int, error) {
	type MongoActionLog struct {
		ID         primitive.ObjectID `bson:"_id,omitempty"`
		ParentId   primitive.ObjectID `bson:"_parentId,omitempty"`
		JardinId   primitive.ObjectID `bson:"jardinId"`
		DateAction primitive.DateTime `bson:"dateAction"`
		Action     string             `bson:"action"`
		Statut     string             `bson:"statut"`
		Lieu       string             `bson:"lieu"`
		Legume     string             `bson:"legume"`
		Variete    string             `bson:"variete"`
		Qte        int                `bson:"qte"`
		Poids      int                `bson:"poids"`
		Notes      string             `bson:"notes"`
		Photos     []string           `bson:"photos"`
		Tags       []string           `bson:"tags"`
	}

	cursor, err := mongoDB.Collection("actionLog").Find(ctx, bson.D{})
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	tx, err := sqliteDB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	count := 0
	for cursor.Next(ctx) {
		var a MongoActionLog
		if err := cursor.Decode(&a); err != nil {
			return 0, fmt.Errorf("décodage log : %w", err)
		}

		id := objectIDToStr(a.ID)
		parentId := ""
		if !a.ParentId.IsZero() {
			parentId = objectIDToStr(a.ParentId)
		}
		jardinId := ""
		if !a.JardinId.IsZero() {
			jardinId = objectIDToStr(a.JardinId)
		}
		dateAction := dateTimeToStr(a.DateAction)

		if a.Photos == nil {
			a.Photos = []string{}
		}
		if a.Tags == nil {
			a.Tags = []string{}
		}
		photosJSON, _ := json.Marshal(a.Photos)
		tagsJSON, _ := json.Marshal(a.Tags)

		_, err = tx.ExecContext(ctx, `
			INSERT OR REPLACE INTO action_log
				(id, parent_id, jardin_id, date_action, action, statut, lieu, legume, variete, qte, poids, notes, photos, tags)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			id, parentId, jardinId, dateAction,
			a.Action, a.Statut, a.Lieu, a.Legume, a.Variete,
			a.Qte, a.Poids, a.Notes,
			string(photosJSON), string(tagsJSON))
		if err != nil {
			return 0, fmt.Errorf("insert log %s : %w", id, err)
		}
		count++
	}
	if err := cursor.Err(); err != nil {
		return 0, err
	}

	return count, tx.Commit()
}
