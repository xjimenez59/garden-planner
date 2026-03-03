package config

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func ConnectDatabase() *sql.DB {
	dbPath := os.Getenv("SQLITE_PATH")
	if dbPath == "" {
		dbPath = "./garden-planner.db"
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	if _, err = db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		log.Fatal(err)
	}
	if _, err = db.Exec("PRAGMA journal_mode = WAL"); err != nil {
		log.Fatal(err)
	}

	if err = createTables(db); err != nil {
		log.Fatal(err)
	}

	DB = db
	return db
}

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
