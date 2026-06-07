package models

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// ---- Station coords cache --------------------------------------------------

type StationInfo struct {
	ID  string
	Lat float64
	Lon float64
}

var (
	stationCacheMu sync.RWMutex
	stationCache   = map[string]StationInfo{}
)

// ---- Previsions cache (TTL 1h) ---------------------------------------------

type previsionsCacheEntry struct {
	data      []HourlyForecast
	fetchedAt time.Time
}

const previsionsCacheTTL = 1 * time.Hour

var (
	previsionsCacheMu sync.RWMutex
	previsionsCache   = map[string]previsionsCacheEntry{}
)

type mfStation struct {
	ID  string  `json:"id"`
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// GetStationCoords returns lat/lon for a MF station code.
// The MF API requires filtering by department (first 2 chars of station code).
// Results are cached in memory for the lifetime of the process.
func GetStationCoords(stationID string) (StationInfo, error) {
	stationCacheMu.RLock()
	if s, ok := stationCache[stationID]; ok {
		stationCacheMu.RUnlock()
		return s, nil
	}
	stationCacheMu.RUnlock()

	if len(stationID) < 2 {
		return StationInfo{}, fmt.Errorf("code station invalide: %s", stationID)
	}
	dept := stationID[:2]

	token, err := getMFToken()
	if err != nil {
		return StationInfo{}, fmt.Errorf("token MF: %w", err)
	}

	url := fmt.Sprintf("%s/liste-stations/quotidienne?id-departement=%s", mfAPIBase(), dept)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return StationInfo{}, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return StationInfo{}, fmt.Errorf("liste-stations: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return StationInfo{}, fmt.Errorf("liste-stations HTTP %d: %s", resp.StatusCode, string(body))
	}

	var stations []mfStation
	if err := json.Unmarshal(body, &stations); err != nil {
		return StationInfo{}, fmt.Errorf("parse liste-stations: %w", err)
	}

	for _, s := range stations {
		if s.ID == stationID {
			info := StationInfo{ID: stationID, Lat: s.Lat, Lon: s.Lon}
			stationCacheMu.Lock()
			stationCache[stationID] = info
			stationCacheMu.Unlock()
			return info, nil
		}
	}
	return StationInfo{}, fmt.Errorf("station %s introuvable dans le département %s", stationID, dept)
}

// ---- Open-Meteo forecast ---------------------------------------------------

type HourlyForecast struct {
	Time         string  `json:"time"`          // "2026-06-05T06:00"
	Temperature  float64 `json:"temperature"`   // °C
	Precipitation float64 `json:"precipitation"` // mm
	WindSpeed    float64 `json:"wind_speed"`    // km/h
	WindDir      int     `json:"wind_dir"`      // degrés 0-360
	WeatherCode  int     `json:"weather_code"`  // WMO code
}

type openMeteoResponse struct {
	Hourly struct {
		Time        []string  `json:"time"`
		Temperature []float64 `json:"temperature_2m"`
		Precip      []float64 `json:"precipitation"`
		WindSpeed   []float64 `json:"windspeed_10m"`
		WindDir     []int     `json:"winddirection_10m"`
		WeatherCode []int     `json:"weathercode"`
	} `json:"hourly"`
}

// GetPrevisions fetches today's hourly forecast from Open-Meteo (MétéoFrance model)
// for the given coordinates. Returns one entry per 3-hour slot.
// Results are cached for previsionsCacheTTL (1h) to avoid hammering the API.
func GetPrevisions(lat, lon float64) ([]HourlyForecast, error) {
	cacheKey := fmt.Sprintf("%.4f:%.4f", lat, lon)

	previsionsCacheMu.RLock()
	if e, ok := previsionsCache[cacheKey]; ok && time.Since(e.fetchedAt) < previsionsCacheTTL {
		previsionsCacheMu.RUnlock()
		return e.data, nil
	}
	previsionsCacheMu.RUnlock()

	url := fmt.Sprintf(
		"https://api.open-meteo.com/v1/meteofrance?latitude=%.4f&longitude=%.4f"+
			"&hourly=temperature_2m,precipitation,windspeed_10m,winddirection_10m,weathercode"+
			"&forecast_days=2&timezone=Europe%%2FParis",
		lat, lon,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("open-meteo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("open-meteo HTTP %d: %s", resp.StatusCode, string(body))
	}

	var raw openMeteoResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("open-meteo parse: %w", err)
	}

	paris, _ := time.LoadLocation("Europe/Paris")
	today := time.Now().In(paris).Format("2006-01-02")
	result := make([]HourlyForecast, 0)

	for i, t := range raw.Hourly.Time {
		if len(t) < 10 || t[:10] != today {
			continue
		}
		// Garder uniquement les créneaux de 3h (0, 3, 6, 9, 12, 15, 18, 21)
		hour := 0
		if len(t) >= 16 {
			fmt.Sscanf(t[11:13], "%d", &hour)
		}
		if hour%3 != 0 {
			continue
		}

		f := HourlyForecast{Time: t}
		if i < len(raw.Hourly.Temperature) {
			f.Temperature = raw.Hourly.Temperature[i]
		}
		if i < len(raw.Hourly.Precip) {
			f.Precipitation = raw.Hourly.Precip[i]
		}
		if i < len(raw.Hourly.WindSpeed) {
			f.WindSpeed = raw.Hourly.WindSpeed[i]
		}
		if i < len(raw.Hourly.WindDir) {
			f.WindDir = raw.Hourly.WindDir[i]
		}
		if i < len(raw.Hourly.WeatherCode) {
			f.WeatherCode = raw.Hourly.WeatherCode[i]
		}
		result = append(result, f)
	}

	previsionsCacheMu.Lock()
	previsionsCache[cacheKey] = previsionsCacheEntry{data: result, fetchedAt: time.Now()}
	previsionsCacheMu.Unlock()

	return result, nil
}
