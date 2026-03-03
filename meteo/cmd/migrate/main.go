// Script de migration des données météo de MongoDB vers SQLite.
//
// Usage :
//
//	go run . \
//	  -mongo "mongodb://root:password@localhost:27017/?authSource=admin" \
//	  -db garden-planner \
//	  -sqlite ./meteo.db
//
// Ou via variables d'environnement :
//
//	MONGO_HOST=localhost MONGO_PORT=27017 MONGO_USER=root MONGO_PWD=secret \
//	MONGO_DBNAME=garden-planner SQLITE_PATH=./meteo.db go run .
//
// Note : les IDs MongoDB (ObjectID hex) sont conservés comme clés primaires SQLite,
// ce qui garantit la cohérence des clés étrangères location_id dans weather_hourly.
package main

import (
	"context"
	"database/sql"
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
	mongoURI := flag.String("mongo", buildMongoURI(), "URI de connexion MongoDB")
	mongoDBName := flag.String("db", envOr("MONGO_DBNAME", "garden-planner"), "Nom de la base MongoDB")
	sqlitePath := flag.String("sqlite", envOr("SQLITE_PATH", "./meteo.db"), "Chemin du fichier SQLite")
	flag.Parse()

	log.Printf("Connexion MongoDB : %s / %s", *mongoURI, *mongoDBName)
	log.Printf("Destination SQLite : %s", *sqlitePath)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
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

	// --- migration des locations (doit être fait en premier pour les FK)
	idMap, locations, err := migrateLocations(ctx, mongoDB, sqliteDB)
	if err != nil {
		log.Fatalf("Erreur migration locations : %v", err)
	}
	log.Printf("Locations migrées : %d", locations)

	// --- migration des données horaires
	hourly, err := migrateHourly(ctx, mongoDB, sqliteDB, idMap)
	if err != nil {
		log.Fatalf("Erreur migration données horaires : %v", err)
	}
	log.Printf("Relevés horaires migrés : %d", hourly)

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

// ---- création des tables ----------------------------------------------------

func createTables(db *sql.DB) error {
	statements := []string{
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

// ---- types MongoDB (pour décodage BSON) -------------------------------------

type mongoWeatherValue struct {
	Value    float32 `bson:"value"`
	Unit     string  `bson:"unit"`
	UnitType int     `bson:"unit_type"`
	Phrase   string  `bson:"phrase"`
}

type mongoWeatherValueRange struct {
	Min mongoWeatherValue `bson:"minimum"`
	Max mongoWeatherValue `bson:"maximum"`
}

type mongoWindData struct {
	Degrees   int               `bson:"degrees"`
	Localized string            `bson:"localized"`
	English   string            `bson:"english"`
	Speed     mongoWeatherValue `bson:"speed"`
}

type mongoPressureTendency struct {
	LocalizedText string `bson:"localizedText"`
	Code          string `bson:"code"`
}

type mongoPrecipitationSummary struct {
	Precipitation mongoWeatherValue `bson:"precipitation"`
	PastHour      mongoWeatherValue `bson:"pastHour"`
	Past3Hours    mongoWeatherValue `bson:"past3Hours"`
	Past6Hours    mongoWeatherValue `bson:"past6Hours"`
	Past9Hours    mongoWeatherValue `bson:"past9Hours"`
	Past12Hours   mongoWeatherValue `bson:"past12Hours"`
	Past18Hours   mongoWeatherValue `bson:"past18Hours"`
	Past24Hours   mongoWeatherValue `bson:"past24Hours"`
}

type mongoTemperatureSummary struct {
	Past6HourRange  mongoWeatherValueRange `bson:"past6HourRange"`
	Past12HourRange mongoWeatherValueRange `bson:"past12HourRange"`
	Past24HourRange mongoWeatherValueRange `bson:"past24HourRange"`
}

type mongoWeatherHourly struct {
	ID                             primitive.ObjectID        `bson:"_id,omitempty"`
	LocationID                     primitive.ObjectID        `bson:"location_id"`
	LocalObservationDateTime       string                    `bson:"localObservationDateTime"`
	WeatherText                    string                    `bson:"weatherText"`
	WeatherIcon                    int                       `bson:"weatherIcon"`
	HasPrecipitation               bool                      `bson:"hasPrecipitation"`
	PrecipitationType              string                    `bson:"precipitationType"`
	IsDayTime                      bool                      `bson:"isDayTime"`
	Temperature                    mongoWeatherValue         `bson:"temperature"`
	RealFeelTemperature            mongoWeatherValue         `bson:"realFeelTemperature"`
	RealFeelTemperatureShade       mongoWeatherValue         `bson:"realFeelTemperatureShade"`
	RelativeHumidity               int                       `bson:"relativeHumidity"`
	IndoorRelativeHumidity         int                       `bson:"indoorRelativeHumidity"`
	DewPoint                       mongoWeatherValue         `bson:"dewPoint"`
	Wind                           mongoWindData             `bson:"wind"`
	WindGustSpeed                  mongoWeatherValue         `bson:"windGustSpeed"`
	UVIndex                        int                       `bson:"uvIndex"`
	UVIndexText                    string                    `bson:"uvIndexText"`
	Visibility                     mongoWeatherValue         `bson:"visibility"`
	ObstructionsToVisibility       string                    `bson:"obstructionsToVisibility"`
	CloudCover                     int                       `bson:"cloudCover"`
	Ceiling                        mongoWeatherValue         `bson:"ceiling"`
	Pressure                       mongoWeatherValue         `bson:"pressure"`
	PressureTendency               mongoPressureTendency     `bson:"pressureTendency"`
	Past24HourTemperatureDeparture mongoWeatherValue         `bson:"past24HourTemperatureDeparture"`
	ApparentTemperature            mongoWeatherValue         `bson:"apparentTemperature"`
	WindChillTemperature           mongoWeatherValue         `bson:"windChillTemperature"`
	WetBulbTemperature             mongoWeatherValue         `bson:"wetBulb"`
	WetBulbGlobeTemperature        mongoWeatherValue         `bson:"wetBulbGlobeTemperature"`
	Precip1hr                      mongoWeatherValue         `bson:"precip1hr"`
	PrecipitationSummary           mongoPrecipitationSummary `bson:"precipitationSummary"`
	TemperatureSummary             mongoTemperatureSummary   `bson:"temperatureSummary"`
	MobileLink                     string                    `bson:"mobileLink"`
	Link                           string                    `bson:"link"`
}

// ---- migration des locations ------------------------------------------------

// migrateLocations migre la collection weather_location.
// Retourne une map MongoDB ObjectID → SQLite id (hex string),
// le nombre de documents migrés et une éventuelle erreur.
func migrateLocations(ctx context.Context, mongoDB *mongo.Database, sqliteDB *sql.DB) (map[primitive.ObjectID]string, int, error) {
	type MongoWeatherValue struct {
		Value    float32 `bson:"value"`
		Unit     string  `bson:"unit"`
		UnitType int     `bson:"unit_type"`
	}
	type MongoWeatherLocation struct {
		ID            primitive.ObjectID `bson:"_id,omitempty"`
		Latitude      float32            `bson:"latitude"`
		Longitude     float32            `bson:"longitude"`
		Elevation     MongoWeatherValue  `bson:"elevation"`
		Key           string             `bson:"key"`
		KeyType       string             `bson:"key_type"`
		LocalizedName string             `bson:"localizedName"`
		PostalCode    string             `bson:"postalCode"`
	}

	cursor, err := mongoDB.Collection("weather_location").Find(ctx, bson.D{})
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	tx, err := sqliteDB.BeginTx(ctx, nil)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback()

	idMap := make(map[primitive.ObjectID]string)
	count := 0

	for cursor.Next(ctx) {
		var loc MongoWeatherLocation
		if err := cursor.Decode(&loc); err != nil {
			return nil, 0, fmt.Errorf("décodage location : %w", err)
		}

		// Conserver l'ObjectID hex comme ID SQLite pour cohérence des FK
		sqliteID := loc.ID.Hex()
		idMap[loc.ID] = sqliteID

		_, err = tx.ExecContext(ctx, `
			INSERT OR REPLACE INTO weather_location
				(id, key, key_type, latitude, longitude,
				 elevation_value, elevation_unit, elevation_unit_type,
				 localized_name, postal_code)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			sqliteID, loc.Key, loc.KeyType, loc.Latitude, loc.Longitude,
			loc.Elevation.Value, loc.Elevation.Unit, loc.Elevation.UnitType,
			loc.LocalizedName, loc.PostalCode)
		if err != nil {
			return nil, 0, fmt.Errorf("insert location %s : %w", sqliteID, err)
		}
		count++
	}
	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}

	return idMap, count, tx.Commit()
}

// ---- migration des relevés horaires -----------------------------------------

func migrateHourly(ctx context.Context, mongoDB *mongo.Database, sqliteDB *sql.DB, idMap map[primitive.ObjectID]string) (int, error) {
	cursor, err := mongoDB.Collection("weather_hourly").Find(ctx, bson.D{})
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
		var doc mongoWeatherHourly
		if err := cursor.Decode(&doc); err != nil {
			return 0, fmt.Errorf("décodage hourly : %w", err)
		}

		sqliteID := doc.ID.Hex()
		sqliteLocID, ok := idMap[doc.LocationID]
		if !ok {
			log.Printf("Attention : location_id %s introuvable pour le relevé %s, ignoré", doc.LocationID.Hex(), sqliteID)
			continue
		}

		_, err = tx.ExecContext(ctx, `
			INSERT OR REPLACE INTO weather_hourly (
				id, location_id, local_observation_datetime,
				weather_text, weather_icon, has_precipitation, precipitation_type, is_day_time,
				temperature_value, temperature_unit, temperature_unit_type, temperature_phrase,
				real_feel_temperature_value, real_feel_temperature_unit, real_feel_temperature_unit_type, real_feel_temperature_phrase,
				real_feel_temperature_shade_value, real_feel_temperature_shade_unit, real_feel_temperature_shade_unit_type, real_feel_temperature_shade_phrase,
				relative_humidity, indoor_relative_humidity,
				dew_point_value, dew_point_unit, dew_point_unit_type, dew_point_phrase,
				wind_degrees, wind_localized, wind_english,
				wind_speed_value, wind_speed_unit, wind_speed_unit_type, wind_speed_phrase,
				wind_gust_speed_value, wind_gust_speed_unit, wind_gust_speed_unit_type, wind_gust_speed_phrase,
				uv_index, uv_index_text,
				visibility_value, visibility_unit, visibility_unit_type, visibility_phrase,
				obstructions_to_visibility, cloud_cover,
				ceiling_value, ceiling_unit, ceiling_unit_type, ceiling_phrase,
				pressure_value, pressure_unit, pressure_unit_type, pressure_phrase,
				pressure_tendency_localized_text, pressure_tendency_code,
				past24h_temp_departure_value, past24h_temp_departure_unit, past24h_temp_departure_unit_type, past24h_temp_departure_phrase,
				apparent_temperature_value, apparent_temperature_unit, apparent_temperature_unit_type, apparent_temperature_phrase,
				wind_chill_temperature_value, wind_chill_temperature_unit, wind_chill_temperature_unit_type, wind_chill_temperature_phrase,
				wet_bulb_temperature_value, wet_bulb_temperature_unit, wet_bulb_temperature_unit_type, wet_bulb_temperature_phrase,
				wet_bulb_globe_temperature_value, wet_bulb_globe_temperature_unit, wet_bulb_globe_temperature_unit_type, wet_bulb_globe_temperature_phrase,
				precip1hr_value, precip1hr_unit, precip1hr_unit_type, precip1hr_phrase,
				precip_summary_precipitation_value, precip_summary_precipitation_unit, precip_summary_precipitation_unit_type, precip_summary_precipitation_phrase,
				precip_summary_past_hour_value, precip_summary_past_hour_unit, precip_summary_past_hour_unit_type, precip_summary_past_hour_phrase,
				precip_summary_past3h_value, precip_summary_past3h_unit, precip_summary_past3h_unit_type, precip_summary_past3h_phrase,
				precip_summary_past6h_value, precip_summary_past6h_unit, precip_summary_past6h_unit_type, precip_summary_past6h_phrase,
				precip_summary_past9h_value, precip_summary_past9h_unit, precip_summary_past9h_unit_type, precip_summary_past9h_phrase,
				precip_summary_past12h_value, precip_summary_past12h_unit, precip_summary_past12h_unit_type, precip_summary_past12h_phrase,
				precip_summary_past18h_value, precip_summary_past18h_unit, precip_summary_past18h_unit_type, precip_summary_past18h_phrase,
				precip_summary_past24h_value, precip_summary_past24h_unit, precip_summary_past24h_unit_type, precip_summary_past24h_phrase,
				temp_summary_past6h_min_value, temp_summary_past6h_min_unit, temp_summary_past6h_min_unit_type, temp_summary_past6h_min_phrase,
				temp_summary_past6h_max_value, temp_summary_past6h_max_unit, temp_summary_past6h_max_unit_type, temp_summary_past6h_max_phrase,
				temp_summary_past12h_min_value, temp_summary_past12h_min_unit, temp_summary_past12h_min_unit_type, temp_summary_past12h_min_phrase,
				temp_summary_past12h_max_value, temp_summary_past12h_max_unit, temp_summary_past12h_max_unit_type, temp_summary_past12h_max_phrase,
				temp_summary_past24h_min_value, temp_summary_past24h_min_unit, temp_summary_past24h_min_unit_type, temp_summary_past24h_min_phrase,
				temp_summary_past24h_max_value, temp_summary_past24h_max_unit, temp_summary_past24h_max_unit_type, temp_summary_past24h_max_phrase,
				mobile_link, link
			) VALUES (
				?, ?, ?,
				?, ?, ?, ?, ?,
				?, ?, ?, ?,
				?, ?, ?, ?,
				?, ?, ?, ?,
				?, ?,
				?, ?, ?, ?,
				?, ?, ?,
				?, ?, ?, ?,
				?, ?, ?, ?,
				?, ?,
				?, ?, ?, ?,
				?, ?,
				?, ?, ?, ?,
				?, ?, ?, ?,
				?, ?,
				?, ?, ?, ?,
				?, ?, ?, ?,
				?, ?, ?, ?,
				?, ?, ?, ?,
				?, ?, ?, ?,
				?, ?, ?, ?,
				?, ?, ?, ?,
				?, ?, ?, ?,
				?, ?, ?, ?,
				?, ?, ?, ?,
				?, ?, ?, ?,
				?, ?, ?, ?,
				?, ?, ?, ?,
				?, ?, ?, ?,
				?, ?, ?, ?,
				?, ?, ?, ?,
				?, ?, ?, ?,
				?, ?, ?, ?,
				?, ?, ?, ?,
				?, ?
			)`,
			sqliteID, sqliteLocID, doc.LocalObservationDateTime,
			doc.WeatherText, doc.WeatherIcon, doc.HasPrecipitation, doc.PrecipitationType, doc.IsDayTime,
			doc.Temperature.Value, doc.Temperature.Unit, doc.Temperature.UnitType, doc.Temperature.Phrase,
			doc.RealFeelTemperature.Value, doc.RealFeelTemperature.Unit, doc.RealFeelTemperature.UnitType, doc.RealFeelTemperature.Phrase,
			doc.RealFeelTemperatureShade.Value, doc.RealFeelTemperatureShade.Unit, doc.RealFeelTemperatureShade.UnitType, doc.RealFeelTemperatureShade.Phrase,
			doc.RelativeHumidity, doc.IndoorRelativeHumidity,
			doc.DewPoint.Value, doc.DewPoint.Unit, doc.DewPoint.UnitType, doc.DewPoint.Phrase,
			doc.Wind.Degrees, doc.Wind.Localized, doc.Wind.English,
			doc.Wind.Speed.Value, doc.Wind.Speed.Unit, doc.Wind.Speed.UnitType, doc.Wind.Speed.Phrase,
			doc.WindGustSpeed.Value, doc.WindGustSpeed.Unit, doc.WindGustSpeed.UnitType, doc.WindGustSpeed.Phrase,
			doc.UVIndex, doc.UVIndexText,
			doc.Visibility.Value, doc.Visibility.Unit, doc.Visibility.UnitType, doc.Visibility.Phrase,
			doc.ObstructionsToVisibility, doc.CloudCover,
			doc.Ceiling.Value, doc.Ceiling.Unit, doc.Ceiling.UnitType, doc.Ceiling.Phrase,
			doc.Pressure.Value, doc.Pressure.Unit, doc.Pressure.UnitType, doc.Pressure.Phrase,
			doc.PressureTendency.LocalizedText, doc.PressureTendency.Code,
			doc.Past24HourTemperatureDeparture.Value, doc.Past24HourTemperatureDeparture.Unit, doc.Past24HourTemperatureDeparture.UnitType, doc.Past24HourTemperatureDeparture.Phrase,
			doc.ApparentTemperature.Value, doc.ApparentTemperature.Unit, doc.ApparentTemperature.UnitType, doc.ApparentTemperature.Phrase,
			doc.WindChillTemperature.Value, doc.WindChillTemperature.Unit, doc.WindChillTemperature.UnitType, doc.WindChillTemperature.Phrase,
			doc.WetBulbTemperature.Value, doc.WetBulbTemperature.Unit, doc.WetBulbTemperature.UnitType, doc.WetBulbTemperature.Phrase,
			doc.WetBulbGlobeTemperature.Value, doc.WetBulbGlobeTemperature.Unit, doc.WetBulbGlobeTemperature.UnitType, doc.WetBulbGlobeTemperature.Phrase,
			doc.Precip1hr.Value, doc.Precip1hr.Unit, doc.Precip1hr.UnitType, doc.Precip1hr.Phrase,
			doc.PrecipitationSummary.Precipitation.Value, doc.PrecipitationSummary.Precipitation.Unit, doc.PrecipitationSummary.Precipitation.UnitType, doc.PrecipitationSummary.Precipitation.Phrase,
			doc.PrecipitationSummary.PastHour.Value, doc.PrecipitationSummary.PastHour.Unit, doc.PrecipitationSummary.PastHour.UnitType, doc.PrecipitationSummary.PastHour.Phrase,
			doc.PrecipitationSummary.Past3Hours.Value, doc.PrecipitationSummary.Past3Hours.Unit, doc.PrecipitationSummary.Past3Hours.UnitType, doc.PrecipitationSummary.Past3Hours.Phrase,
			doc.PrecipitationSummary.Past6Hours.Value, doc.PrecipitationSummary.Past6Hours.Unit, doc.PrecipitationSummary.Past6Hours.UnitType, doc.PrecipitationSummary.Past6Hours.Phrase,
			doc.PrecipitationSummary.Past9Hours.Value, doc.PrecipitationSummary.Past9Hours.Unit, doc.PrecipitationSummary.Past9Hours.UnitType, doc.PrecipitationSummary.Past9Hours.Phrase,
			doc.PrecipitationSummary.Past12Hours.Value, doc.PrecipitationSummary.Past12Hours.Unit, doc.PrecipitationSummary.Past12Hours.UnitType, doc.PrecipitationSummary.Past12Hours.Phrase,
			doc.PrecipitationSummary.Past18Hours.Value, doc.PrecipitationSummary.Past18Hours.Unit, doc.PrecipitationSummary.Past18Hours.UnitType, doc.PrecipitationSummary.Past18Hours.Phrase,
			doc.PrecipitationSummary.Past24Hours.Value, doc.PrecipitationSummary.Past24Hours.Unit, doc.PrecipitationSummary.Past24Hours.UnitType, doc.PrecipitationSummary.Past24Hours.Phrase,
			doc.TemperatureSummary.Past6HourRange.Min.Value, doc.TemperatureSummary.Past6HourRange.Min.Unit, doc.TemperatureSummary.Past6HourRange.Min.UnitType, doc.TemperatureSummary.Past6HourRange.Min.Phrase,
			doc.TemperatureSummary.Past6HourRange.Max.Value, doc.TemperatureSummary.Past6HourRange.Max.Unit, doc.TemperatureSummary.Past6HourRange.Max.UnitType, doc.TemperatureSummary.Past6HourRange.Max.Phrase,
			doc.TemperatureSummary.Past12HourRange.Min.Value, doc.TemperatureSummary.Past12HourRange.Min.Unit, doc.TemperatureSummary.Past12HourRange.Min.UnitType, doc.TemperatureSummary.Past12HourRange.Min.Phrase,
			doc.TemperatureSummary.Past12HourRange.Max.Value, doc.TemperatureSummary.Past12HourRange.Max.Unit, doc.TemperatureSummary.Past12HourRange.Max.UnitType, doc.TemperatureSummary.Past12HourRange.Max.Phrase,
			doc.TemperatureSummary.Past24HourRange.Min.Value, doc.TemperatureSummary.Past24HourRange.Min.Unit, doc.TemperatureSummary.Past24HourRange.Min.UnitType, doc.TemperatureSummary.Past24HourRange.Min.Phrase,
			doc.TemperatureSummary.Past24HourRange.Max.Value, doc.TemperatureSummary.Past24HourRange.Max.Unit, doc.TemperatureSummary.Past24HourRange.Max.UnitType, doc.TemperatureSummary.Past24HourRange.Max.Phrase,
			doc.MobileLink, doc.Link)
		if err != nil {
			return 0, fmt.Errorf("insert hourly %s : %w", sqliteID, err)
		}
		count++
	}
	if err := cursor.Err(); err != nil {
		return 0, err
	}

	return count, tx.Commit()
}
