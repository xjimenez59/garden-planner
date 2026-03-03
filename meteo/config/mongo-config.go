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
		dbPath = "./meteo.db"
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

func CloseDatabase() {
	if DB != nil {
		DB.Close()
	}
}

func createTables(db *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS meteofrance_quotidien (
			poste TEXT NOT NULL,
			date TEXT NOT NULL,
			rr TEXT, qrr TEXT, drr TEXT, qdrr TEXT,
			tn TEXT, qtn TEXT, htn TEXT, qhtn TEXT,
			tx TEXT, qtx TEXT, htx TEXT, qhtx TEXT,
			tm TEXT, qtm TEXT,
			tmnx TEXT, qtmnx TEXT,
			tnsol TEXT, qtnsol TEXT,
			tn50 TEXT, qtn50 TEXT,
			dg TEXT, qdg TEXT,
			tampli TEXT, qtampli TEXT,
			tntxm TEXT, qtntxm TEXT,
			pmerm TEXT, qpmerm TEXT,
			pmermin TEXT, qpmermin TEXT,
			ffm TEXT, qffm TEXT,
			fxi TEXT, qfxi TEXT, dxi TEXT, qdxi TEXT, hxi TEXT, qhxi TEXT,
			fxy TEXT, qfxy TEXT, dxy TEXT, qdxy TEXT, hxy TEXT, qhxy TEXT,
			ff2m TEXT, qff2m TEXT,
			fxi2 TEXT, qfxi2 TEXT, dxi2 TEXT, qdxi2 TEXT, hxi2 TEXT, qhxi2 TEXT,
			fxi3s TEXT, qfxi3s TEXT, dxi3s TEXT, qdxi3s TEXT, hxi3s TEXT, qhxi3s TEXT,
			un TEXT, qun TEXT, hun TEXT, qhun TEXT,
			ux TEXT, qux TEXT, hux TEXT, qhux TEXT,
			dhumi40 TEXT, qdhumi40 TEXT,
			dhumi80 TEXT, qdhumi80 TEXT,
			tsvm TEXT, qtsvm TEXT,
			dhumec TEXT, qdhumec TEXT,
			um TEXT, qum TEXT,
			inst TEXT, qinst TEXT,
			glot TEXT, qglot TEXT,
			dift TEXT, qdift TEXT,
			dirt TEXT, qdirt TEXT,
			sigma TEXT, qsigma TEXT,
			infrart TEXT, qinfrart TEXT,
			uv_indicex TEXT, quv_indicex TEXT,
			nb300 TEXT, qnb300 TEXT,
			ba300 TEXT, qba300 TEXT,
			neig TEXT, qneig TEXT,
			brou TEXT, qbrou TEXT,
			orag TEXT, qorag TEXT,
			gresil TEXT, qgresil TEXT,
			grele TEXT, qgrele TEXT,
			rosee TEXT, qrosee TEXT,
			verglas TEXT, qverglas TEXT,
			solneige TEXT, qsolneige TEXT,
			gelee TEXT, qgelee TEXT,
			fumee TEXT, qfumee TEXT,
			brume TEXT, qbrume TEXT,
			eclair TEXT, qeclair TEXT,
			etpmon TEXT, qetpmon TEXT,
			etpgrille TEXT, qetpgrille TEXT,
			uv TEXT, quv TEXT,
			tmermax TEXT, qtmermax TEXT,
			tmermin TEXT, qtmermin TEXT,
			hneigef TEXT, qhneigef TEXT,
			neigetotx TEXT, qneigetotx TEXT,
			neigetot06 TEXT, qneigetot06 TEXT,
			PRIMARY KEY (poste, date)
		)`,
		`CREATE TABLE IF NOT EXISTS weather_location (
			id TEXT PRIMARY KEY,
			key TEXT NOT NULL UNIQUE,
			key_type TEXT NOT NULL DEFAULT '',
			latitude REAL NOT NULL DEFAULT 0,
			longitude REAL NOT NULL DEFAULT 0,
			elevation_value REAL NOT NULL DEFAULT 0,
			elevation_unit TEXT NOT NULL DEFAULT '',
			elevation_unit_type INTEGER NOT NULL DEFAULT 0,
			localized_name TEXT NOT NULL DEFAULT '',
			postal_code TEXT NOT NULL DEFAULT ''
		)`,
		`CREATE TABLE IF NOT EXISTS weather_hourly (
			id TEXT PRIMARY KEY,
			location_id TEXT NOT NULL,
			local_observation_datetime TEXT NOT NULL,
			weather_text TEXT NOT NULL DEFAULT '',
			weather_icon INTEGER NOT NULL DEFAULT 0,
			has_precipitation INTEGER NOT NULL DEFAULT 0,
			precipitation_type TEXT NOT NULL DEFAULT '',
			is_day_time INTEGER NOT NULL DEFAULT 0,
			temperature_value REAL NOT NULL DEFAULT 0,
			temperature_unit TEXT NOT NULL DEFAULT '',
			temperature_unit_type INTEGER NOT NULL DEFAULT 0,
			temperature_phrase TEXT NOT NULL DEFAULT '',
			real_feel_temperature_value REAL NOT NULL DEFAULT 0,
			real_feel_temperature_unit TEXT NOT NULL DEFAULT '',
			real_feel_temperature_unit_type INTEGER NOT NULL DEFAULT 0,
			real_feel_temperature_phrase TEXT NOT NULL DEFAULT '',
			real_feel_temperature_shade_value REAL NOT NULL DEFAULT 0,
			real_feel_temperature_shade_unit TEXT NOT NULL DEFAULT '',
			real_feel_temperature_shade_unit_type INTEGER NOT NULL DEFAULT 0,
			real_feel_temperature_shade_phrase TEXT NOT NULL DEFAULT '',
			relative_humidity INTEGER NOT NULL DEFAULT 0,
			indoor_relative_humidity INTEGER NOT NULL DEFAULT 0,
			dew_point_value REAL NOT NULL DEFAULT 0,
			dew_point_unit TEXT NOT NULL DEFAULT '',
			dew_point_unit_type INTEGER NOT NULL DEFAULT 0,
			dew_point_phrase TEXT NOT NULL DEFAULT '',
			wind_degrees INTEGER NOT NULL DEFAULT 0,
			wind_localized TEXT NOT NULL DEFAULT '',
			wind_english TEXT NOT NULL DEFAULT '',
			wind_speed_value REAL NOT NULL DEFAULT 0,
			wind_speed_unit TEXT NOT NULL DEFAULT '',
			wind_speed_unit_type INTEGER NOT NULL DEFAULT 0,
			wind_speed_phrase TEXT NOT NULL DEFAULT '',
			wind_gust_speed_value REAL NOT NULL DEFAULT 0,
			wind_gust_speed_unit TEXT NOT NULL DEFAULT '',
			wind_gust_speed_unit_type INTEGER NOT NULL DEFAULT 0,
			wind_gust_speed_phrase TEXT NOT NULL DEFAULT '',
			uv_index INTEGER NOT NULL DEFAULT 0,
			uv_index_text TEXT NOT NULL DEFAULT '',
			visibility_value REAL NOT NULL DEFAULT 0,
			visibility_unit TEXT NOT NULL DEFAULT '',
			visibility_unit_type INTEGER NOT NULL DEFAULT 0,
			visibility_phrase TEXT NOT NULL DEFAULT '',
			obstructions_to_visibility TEXT NOT NULL DEFAULT '',
			cloud_cover INTEGER NOT NULL DEFAULT 0,
			ceiling_value REAL NOT NULL DEFAULT 0,
			ceiling_unit TEXT NOT NULL DEFAULT '',
			ceiling_unit_type INTEGER NOT NULL DEFAULT 0,
			ceiling_phrase TEXT NOT NULL DEFAULT '',
			pressure_value REAL NOT NULL DEFAULT 0,
			pressure_unit TEXT NOT NULL DEFAULT '',
			pressure_unit_type INTEGER NOT NULL DEFAULT 0,
			pressure_phrase TEXT NOT NULL DEFAULT '',
			pressure_tendency_localized_text TEXT NOT NULL DEFAULT '',
			pressure_tendency_code TEXT NOT NULL DEFAULT '',
			past24h_temp_departure_value REAL NOT NULL DEFAULT 0,
			past24h_temp_departure_unit TEXT NOT NULL DEFAULT '',
			past24h_temp_departure_unit_type INTEGER NOT NULL DEFAULT 0,
			past24h_temp_departure_phrase TEXT NOT NULL DEFAULT '',
			apparent_temperature_value REAL NOT NULL DEFAULT 0,
			apparent_temperature_unit TEXT NOT NULL DEFAULT '',
			apparent_temperature_unit_type INTEGER NOT NULL DEFAULT 0,
			apparent_temperature_phrase TEXT NOT NULL DEFAULT '',
			wind_chill_temperature_value REAL NOT NULL DEFAULT 0,
			wind_chill_temperature_unit TEXT NOT NULL DEFAULT '',
			wind_chill_temperature_unit_type INTEGER NOT NULL DEFAULT 0,
			wind_chill_temperature_phrase TEXT NOT NULL DEFAULT '',
			wet_bulb_temperature_value REAL NOT NULL DEFAULT 0,
			wet_bulb_temperature_unit TEXT NOT NULL DEFAULT '',
			wet_bulb_temperature_unit_type INTEGER NOT NULL DEFAULT 0,
			wet_bulb_temperature_phrase TEXT NOT NULL DEFAULT '',
			wet_bulb_globe_temperature_value REAL NOT NULL DEFAULT 0,
			wet_bulb_globe_temperature_unit TEXT NOT NULL DEFAULT '',
			wet_bulb_globe_temperature_unit_type INTEGER NOT NULL DEFAULT 0,
			wet_bulb_globe_temperature_phrase TEXT NOT NULL DEFAULT '',
			precip1hr_value REAL NOT NULL DEFAULT 0,
			precip1hr_unit TEXT NOT NULL DEFAULT '',
			precip1hr_unit_type INTEGER NOT NULL DEFAULT 0,
			precip1hr_phrase TEXT NOT NULL DEFAULT '',
			precip_summary_precipitation_value REAL NOT NULL DEFAULT 0,
			precip_summary_precipitation_unit TEXT NOT NULL DEFAULT '',
			precip_summary_precipitation_unit_type INTEGER NOT NULL DEFAULT 0,
			precip_summary_precipitation_phrase TEXT NOT NULL DEFAULT '',
			precip_summary_past_hour_value REAL NOT NULL DEFAULT 0,
			precip_summary_past_hour_unit TEXT NOT NULL DEFAULT '',
			precip_summary_past_hour_unit_type INTEGER NOT NULL DEFAULT 0,
			precip_summary_past_hour_phrase TEXT NOT NULL DEFAULT '',
			precip_summary_past3h_value REAL NOT NULL DEFAULT 0,
			precip_summary_past3h_unit TEXT NOT NULL DEFAULT '',
			precip_summary_past3h_unit_type INTEGER NOT NULL DEFAULT 0,
			precip_summary_past3h_phrase TEXT NOT NULL DEFAULT '',
			precip_summary_past6h_value REAL NOT NULL DEFAULT 0,
			precip_summary_past6h_unit TEXT NOT NULL DEFAULT '',
			precip_summary_past6h_unit_type INTEGER NOT NULL DEFAULT 0,
			precip_summary_past6h_phrase TEXT NOT NULL DEFAULT '',
			precip_summary_past9h_value REAL NOT NULL DEFAULT 0,
			precip_summary_past9h_unit TEXT NOT NULL DEFAULT '',
			precip_summary_past9h_unit_type INTEGER NOT NULL DEFAULT 0,
			precip_summary_past9h_phrase TEXT NOT NULL DEFAULT '',
			precip_summary_past12h_value REAL NOT NULL DEFAULT 0,
			precip_summary_past12h_unit TEXT NOT NULL DEFAULT '',
			precip_summary_past12h_unit_type INTEGER NOT NULL DEFAULT 0,
			precip_summary_past12h_phrase TEXT NOT NULL DEFAULT '',
			precip_summary_past18h_value REAL NOT NULL DEFAULT 0,
			precip_summary_past18h_unit TEXT NOT NULL DEFAULT '',
			precip_summary_past18h_unit_type INTEGER NOT NULL DEFAULT 0,
			precip_summary_past18h_phrase TEXT NOT NULL DEFAULT '',
			precip_summary_past24h_value REAL NOT NULL DEFAULT 0,
			precip_summary_past24h_unit TEXT NOT NULL DEFAULT '',
			precip_summary_past24h_unit_type INTEGER NOT NULL DEFAULT 0,
			precip_summary_past24h_phrase TEXT NOT NULL DEFAULT '',
			temp_summary_past6h_min_value REAL NOT NULL DEFAULT 0,
			temp_summary_past6h_min_unit TEXT NOT NULL DEFAULT '',
			temp_summary_past6h_min_unit_type INTEGER NOT NULL DEFAULT 0,
			temp_summary_past6h_min_phrase TEXT NOT NULL DEFAULT '',
			temp_summary_past6h_max_value REAL NOT NULL DEFAULT 0,
			temp_summary_past6h_max_unit TEXT NOT NULL DEFAULT '',
			temp_summary_past6h_max_unit_type INTEGER NOT NULL DEFAULT 0,
			temp_summary_past6h_max_phrase TEXT NOT NULL DEFAULT '',
			temp_summary_past12h_min_value REAL NOT NULL DEFAULT 0,
			temp_summary_past12h_min_unit TEXT NOT NULL DEFAULT '',
			temp_summary_past12h_min_unit_type INTEGER NOT NULL DEFAULT 0,
			temp_summary_past12h_min_phrase TEXT NOT NULL DEFAULT '',
			temp_summary_past12h_max_value REAL NOT NULL DEFAULT 0,
			temp_summary_past12h_max_unit TEXT NOT NULL DEFAULT '',
			temp_summary_past12h_max_unit_type INTEGER NOT NULL DEFAULT 0,
			temp_summary_past12h_max_phrase TEXT NOT NULL DEFAULT '',
			temp_summary_past24h_min_value REAL NOT NULL DEFAULT 0,
			temp_summary_past24h_min_unit TEXT NOT NULL DEFAULT '',
			temp_summary_past24h_min_unit_type INTEGER NOT NULL DEFAULT 0,
			temp_summary_past24h_min_phrase TEXT NOT NULL DEFAULT '',
			temp_summary_past24h_max_value REAL NOT NULL DEFAULT 0,
			temp_summary_past24h_max_unit TEXT NOT NULL DEFAULT '',
			temp_summary_past24h_max_unit_type INTEGER NOT NULL DEFAULT 0,
			temp_summary_past24h_max_phrase TEXT NOT NULL DEFAULT '',
			mobile_link TEXT NOT NULL DEFAULT '',
			link TEXT NOT NULL DEFAULT '',
			UNIQUE (location_id, local_observation_datetime),
			FOREIGN KEY (location_id) REFERENCES weather_location(id)
		)`,
	}
	for _, stmt := range statements {
		if _, err := db.Exec(stmt); err != nil {
			return err
		}
	}
	return nil
}
