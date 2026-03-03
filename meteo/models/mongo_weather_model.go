package models

import (
	"context"
	"garden-planner/meteo/config"
	"sort"
	"time"

	"github.com/google/uuid"
)

type WeatherValue struct {
	Value    float32 `json:"value"`
	Unit     string  `json:"unit"`
	UnitType int     `json:"unit_type"`
	Phrase   string  `json:"phrase,omitempty"`
}

type WeatherValueRange struct {
	Min WeatherValue `json:"min"`
	Max WeatherValue `json:"max"`
}

type WeatherLocation struct {
	ID            string       `json:"id"`
	Latitude      float32      `json:"latitude"`
	Longitude     float32      `json:"longitude"`
	Elevation     WeatherValue `json:"elevation"`
	Key           string       `json:"key"`
	KeyType       string       `json:"key_type"`
	LocalizedName string       `json:"localizedName"`
	PostalCode    string       `json:"postalCode"`
}

type WindData struct {
	Degrees   int          `json:"degrees"`
	Localized string       `json:"localized"`
	English   string       `json:"english"`
	Speed     WeatherValue `json:"speed"`
}

type PressureTendencyData struct {
	LocalizedText string `json:"localizedText"`
	Code          string `json:"code"`
}

type PrecipitationSummaryData struct {
	Precipitation WeatherValue `json:"precipitation"`
	PastHour      WeatherValue `json:"pastHour"`
	Past3Hours    WeatherValue `json:"past3Hours"`
	Past6Hours    WeatherValue `json:"past6Hours"`
	Past9Hours    WeatherValue `json:"past9Hours"`
	Past12Hours   WeatherValue `json:"past12Hours"`
	Past18Hours   WeatherValue `json:"past18Hours"`
	Past24Hours   WeatherValue `json:"past24Hours"`
}

type TemperatureSummaryData struct {
	Past6HourRange  WeatherValueRange `json:"past6HourRange"`
	Past12HourRange WeatherValueRange `json:"past12HourRange"`
	Past24HourRange WeatherValueRange `json:"past24HourRange"`
}

type WeatherHourlyData struct {
	ID                             string                   `json:"id"`
	LocationID                     string                   `json:"location_id"`
	LocalObservationDateTime       string                   `json:"localObservationDateTime"`
	WeatherText                    string                   `json:"weatherText"`
	WeatherIcon                    int                      `json:"weatherIcon"`
	HasPrecipitation               bool                     `json:"hasPrecipitation"`
	PrecipitationType              string                   `json:"precipitationType"`
	IsDayTime                      bool                     `json:"isDayTime"`
	Temperature                    WeatherValue             `json:"temperature"`
	RealFeelTemperature            WeatherValue             `json:"realFeelTemperature"`
	RealFeelTemperatureShade       WeatherValue             `json:"realFeelTemperatureShade"`
	RelativeHumidity               int                      `json:"relativeHumidity"`
	IndoorRelativeHumidity         int                      `json:"indoorRelativeHumidity"`
	DewPoint                       WeatherValue             `json:"dewPoint"`
	Wind                           WindData                 `json:"wind"`
	WindGustSpeed                  WeatherValue             `json:"windGustSpeed"`
	UVIndex                        int                      `json:"uvIndex"`
	UVIndexText                    string                   `json:"uvIndexText"`
	Visibility                     WeatherValue             `json:"visibility"`
	ObstructionsToVisibility       string                   `json:"obstructionsToVisibility"`
	CloudCover                     int                      `json:"cloudCover"`
	Ceiling                        WeatherValue             `json:"ceiling"`
	Pressure                       WeatherValue             `json:"pressure"`
	PressureTendency               PressureTendencyData     `json:"pressureTendency"`
	Past24HourTemperatureDeparture WeatherValue             `json:"past24HourTemperatureDeparture"`
	ApparentTemperature            WeatherValue             `json:"apparentTemperature"`
	WindChillTemperature           WeatherValue             `json:"windChillTemperature"`
	WetBulbTemperature             WeatherValue             `json:"wetBulb"`
	WetBulbGlobeTemperature        WeatherValue             `json:"wetBulbGlobeTemperature"`
	Precip1hr                      WeatherValue             `json:"precip1hr"`
	PrecipitationSummary           PrecipitationSummaryData `json:"precipitationSummary"`
	TemperatureSummary             TemperatureSummaryData   `json:"temperatureSummary"`
	MobileLink                     string                   `json:"mobileLink"`
	Link                           string                   `json:"link"`
}

type WeatherDaySummary struct {
	ID                     string            `json:"id"`
	LocationID             string            `json:"location_id"`
	Date                   string            `json:"date"`
	WeatherText            string            `json:"weatherText"`
	WeatherIcon            int               `json:"weatherIcon"`
	Precipitation          WeatherValue      `json:"precipitationType"`
	NightTemperature       WeatherValueRange `json:"nightTemperature"`
	DayTemperature         WeatherValueRange `json:"dayTemperature"`
	DayRealFeelTemperature WeatherValueRange `json:"dayRealFeelTemperature"`
	WindDirection          string            `json:"windDirection"`
	WindSpeed              WeatherValue      `json:"windSpeed"`
	WindGustSpeed          WeatherValue      `json:"windGustSpeed"`
	CloudCover             int               `json:"cloudCover"`
}

// ---- persistence ------------------------------------------------------------

func (wh *WeatherHourlyData) Save(ctx context.Context) error {
	if wh.ID == "" {
		wh.ID = uuid.New().String()
	}
	_, err := config.DB.ExecContext(ctx, `
		INSERT INTO weather_hourly (
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
		)
		ON CONFLICT(location_id, local_observation_datetime) DO UPDATE SET
			weather_text = excluded.weather_text,
			weather_icon = excluded.weather_icon,
			has_precipitation = excluded.has_precipitation,
			precipitation_type = excluded.precipitation_type,
			is_day_time = excluded.is_day_time,
			temperature_value = excluded.temperature_value,
			temperature_unit = excluded.temperature_unit,
			temperature_unit_type = excluded.temperature_unit_type,
			temperature_phrase = excluded.temperature_phrase,
			real_feel_temperature_value = excluded.real_feel_temperature_value,
			real_feel_temperature_unit = excluded.real_feel_temperature_unit,
			real_feel_temperature_unit_type = excluded.real_feel_temperature_unit_type,
			real_feel_temperature_phrase = excluded.real_feel_temperature_phrase,
			real_feel_temperature_shade_value = excluded.real_feel_temperature_shade_value,
			real_feel_temperature_shade_unit = excluded.real_feel_temperature_shade_unit,
			real_feel_temperature_shade_unit_type = excluded.real_feel_temperature_shade_unit_type,
			real_feel_temperature_shade_phrase = excluded.real_feel_temperature_shade_phrase,
			relative_humidity = excluded.relative_humidity,
			indoor_relative_humidity = excluded.indoor_relative_humidity,
			dew_point_value = excluded.dew_point_value,
			dew_point_unit = excluded.dew_point_unit,
			dew_point_unit_type = excluded.dew_point_unit_type,
			dew_point_phrase = excluded.dew_point_phrase,
			wind_degrees = excluded.wind_degrees,
			wind_localized = excluded.wind_localized,
			wind_english = excluded.wind_english,
			wind_speed_value = excluded.wind_speed_value,
			wind_speed_unit = excluded.wind_speed_unit,
			wind_speed_unit_type = excluded.wind_speed_unit_type,
			wind_speed_phrase = excluded.wind_speed_phrase,
			wind_gust_speed_value = excluded.wind_gust_speed_value,
			wind_gust_speed_unit = excluded.wind_gust_speed_unit,
			wind_gust_speed_unit_type = excluded.wind_gust_speed_unit_type,
			wind_gust_speed_phrase = excluded.wind_gust_speed_phrase,
			uv_index = excluded.uv_index,
			uv_index_text = excluded.uv_index_text,
			visibility_value = excluded.visibility_value,
			visibility_unit = excluded.visibility_unit,
			visibility_unit_type = excluded.visibility_unit_type,
			visibility_phrase = excluded.visibility_phrase,
			obstructions_to_visibility = excluded.obstructions_to_visibility,
			cloud_cover = excluded.cloud_cover,
			ceiling_value = excluded.ceiling_value,
			ceiling_unit = excluded.ceiling_unit,
			ceiling_unit_type = excluded.ceiling_unit_type,
			ceiling_phrase = excluded.ceiling_phrase,
			pressure_value = excluded.pressure_value,
			pressure_unit = excluded.pressure_unit,
			pressure_unit_type = excluded.pressure_unit_type,
			pressure_phrase = excluded.pressure_phrase,
			pressure_tendency_localized_text = excluded.pressure_tendency_localized_text,
			pressure_tendency_code = excluded.pressure_tendency_code,
			past24h_temp_departure_value = excluded.past24h_temp_departure_value,
			past24h_temp_departure_unit = excluded.past24h_temp_departure_unit,
			past24h_temp_departure_unit_type = excluded.past24h_temp_departure_unit_type,
			past24h_temp_departure_phrase = excluded.past24h_temp_departure_phrase,
			apparent_temperature_value = excluded.apparent_temperature_value,
			apparent_temperature_unit = excluded.apparent_temperature_unit,
			apparent_temperature_unit_type = excluded.apparent_temperature_unit_type,
			apparent_temperature_phrase = excluded.apparent_temperature_phrase,
			wind_chill_temperature_value = excluded.wind_chill_temperature_value,
			wind_chill_temperature_unit = excluded.wind_chill_temperature_unit,
			wind_chill_temperature_unit_type = excluded.wind_chill_temperature_unit_type,
			wind_chill_temperature_phrase = excluded.wind_chill_temperature_phrase,
			wet_bulb_temperature_value = excluded.wet_bulb_temperature_value,
			wet_bulb_temperature_unit = excluded.wet_bulb_temperature_unit,
			wet_bulb_temperature_unit_type = excluded.wet_bulb_temperature_unit_type,
			wet_bulb_temperature_phrase = excluded.wet_bulb_temperature_phrase,
			wet_bulb_globe_temperature_value = excluded.wet_bulb_globe_temperature_value,
			wet_bulb_globe_temperature_unit = excluded.wet_bulb_globe_temperature_unit,
			wet_bulb_globe_temperature_unit_type = excluded.wet_bulb_globe_temperature_unit_type,
			wet_bulb_globe_temperature_phrase = excluded.wet_bulb_globe_temperature_phrase,
			precip1hr_value = excluded.precip1hr_value,
			precip1hr_unit = excluded.precip1hr_unit,
			precip1hr_unit_type = excluded.precip1hr_unit_type,
			precip1hr_phrase = excluded.precip1hr_phrase,
			precip_summary_precipitation_value = excluded.precip_summary_precipitation_value,
			precip_summary_precipitation_unit = excluded.precip_summary_precipitation_unit,
			precip_summary_precipitation_unit_type = excluded.precip_summary_precipitation_unit_type,
			precip_summary_precipitation_phrase = excluded.precip_summary_precipitation_phrase,
			precip_summary_past_hour_value = excluded.precip_summary_past_hour_value,
			precip_summary_past_hour_unit = excluded.precip_summary_past_hour_unit,
			precip_summary_past_hour_unit_type = excluded.precip_summary_past_hour_unit_type,
			precip_summary_past_hour_phrase = excluded.precip_summary_past_hour_phrase,
			precip_summary_past3h_value = excluded.precip_summary_past3h_value,
			precip_summary_past3h_unit = excluded.precip_summary_past3h_unit,
			precip_summary_past3h_unit_type = excluded.precip_summary_past3h_unit_type,
			precip_summary_past3h_phrase = excluded.precip_summary_past3h_phrase,
			precip_summary_past6h_value = excluded.precip_summary_past6h_value,
			precip_summary_past6h_unit = excluded.precip_summary_past6h_unit,
			precip_summary_past6h_unit_type = excluded.precip_summary_past6h_unit_type,
			precip_summary_past6h_phrase = excluded.precip_summary_past6h_phrase,
			precip_summary_past9h_value = excluded.precip_summary_past9h_value,
			precip_summary_past9h_unit = excluded.precip_summary_past9h_unit,
			precip_summary_past9h_unit_type = excluded.precip_summary_past9h_unit_type,
			precip_summary_past9h_phrase = excluded.precip_summary_past9h_phrase,
			precip_summary_past12h_value = excluded.precip_summary_past12h_value,
			precip_summary_past12h_unit = excluded.precip_summary_past12h_unit,
			precip_summary_past12h_unit_type = excluded.precip_summary_past12h_unit_type,
			precip_summary_past12h_phrase = excluded.precip_summary_past12h_phrase,
			precip_summary_past18h_value = excluded.precip_summary_past18h_value,
			precip_summary_past18h_unit = excluded.precip_summary_past18h_unit,
			precip_summary_past18h_unit_type = excluded.precip_summary_past18h_unit_type,
			precip_summary_past18h_phrase = excluded.precip_summary_past18h_phrase,
			precip_summary_past24h_value = excluded.precip_summary_past24h_value,
			precip_summary_past24h_unit = excluded.precip_summary_past24h_unit,
			precip_summary_past24h_unit_type = excluded.precip_summary_past24h_unit_type,
			precip_summary_past24h_phrase = excluded.precip_summary_past24h_phrase,
			temp_summary_past6h_min_value = excluded.temp_summary_past6h_min_value,
			temp_summary_past6h_min_unit = excluded.temp_summary_past6h_min_unit,
			temp_summary_past6h_min_unit_type = excluded.temp_summary_past6h_min_unit_type,
			temp_summary_past6h_min_phrase = excluded.temp_summary_past6h_min_phrase,
			temp_summary_past6h_max_value = excluded.temp_summary_past6h_max_value,
			temp_summary_past6h_max_unit = excluded.temp_summary_past6h_max_unit,
			temp_summary_past6h_max_unit_type = excluded.temp_summary_past6h_max_unit_type,
			temp_summary_past6h_max_phrase = excluded.temp_summary_past6h_max_phrase,
			temp_summary_past12h_min_value = excluded.temp_summary_past12h_min_value,
			temp_summary_past12h_min_unit = excluded.temp_summary_past12h_min_unit,
			temp_summary_past12h_min_unit_type = excluded.temp_summary_past12h_min_unit_type,
			temp_summary_past12h_min_phrase = excluded.temp_summary_past12h_min_phrase,
			temp_summary_past12h_max_value = excluded.temp_summary_past12h_max_value,
			temp_summary_past12h_max_unit = excluded.temp_summary_past12h_max_unit,
			temp_summary_past12h_max_unit_type = excluded.temp_summary_past12h_max_unit_type,
			temp_summary_past12h_max_phrase = excluded.temp_summary_past12h_max_phrase,
			temp_summary_past24h_min_value = excluded.temp_summary_past24h_min_value,
			temp_summary_past24h_min_unit = excluded.temp_summary_past24h_min_unit,
			temp_summary_past24h_min_unit_type = excluded.temp_summary_past24h_min_unit_type,
			temp_summary_past24h_min_phrase = excluded.temp_summary_past24h_min_phrase,
			temp_summary_past24h_max_value = excluded.temp_summary_past24h_max_value,
			temp_summary_past24h_max_unit = excluded.temp_summary_past24h_max_unit,
			temp_summary_past24h_max_unit_type = excluded.temp_summary_past24h_max_unit_type,
			temp_summary_past24h_max_phrase = excluded.temp_summary_past24h_max_phrase,
			mobile_link = excluded.mobile_link,
			link = excluded.link`,
		wh.ID, wh.LocationID, wh.LocalObservationDateTime,
		wh.WeatherText, wh.WeatherIcon, wh.HasPrecipitation, wh.PrecipitationType, wh.IsDayTime,
		wh.Temperature.Value, wh.Temperature.Unit, wh.Temperature.UnitType, wh.Temperature.Phrase,
		wh.RealFeelTemperature.Value, wh.RealFeelTemperature.Unit, wh.RealFeelTemperature.UnitType, wh.RealFeelTemperature.Phrase,
		wh.RealFeelTemperatureShade.Value, wh.RealFeelTemperatureShade.Unit, wh.RealFeelTemperatureShade.UnitType, wh.RealFeelTemperatureShade.Phrase,
		wh.RelativeHumidity, wh.IndoorRelativeHumidity,
		wh.DewPoint.Value, wh.DewPoint.Unit, wh.DewPoint.UnitType, wh.DewPoint.Phrase,
		wh.Wind.Degrees, wh.Wind.Localized, wh.Wind.English,
		wh.Wind.Speed.Value, wh.Wind.Speed.Unit, wh.Wind.Speed.UnitType, wh.Wind.Speed.Phrase,
		wh.WindGustSpeed.Value, wh.WindGustSpeed.Unit, wh.WindGustSpeed.UnitType, wh.WindGustSpeed.Phrase,
		wh.UVIndex, wh.UVIndexText,
		wh.Visibility.Value, wh.Visibility.Unit, wh.Visibility.UnitType, wh.Visibility.Phrase,
		wh.ObstructionsToVisibility, wh.CloudCover,
		wh.Ceiling.Value, wh.Ceiling.Unit, wh.Ceiling.UnitType, wh.Ceiling.Phrase,
		wh.Pressure.Value, wh.Pressure.Unit, wh.Pressure.UnitType, wh.Pressure.Phrase,
		wh.PressureTendency.LocalizedText, wh.PressureTendency.Code,
		wh.Past24HourTemperatureDeparture.Value, wh.Past24HourTemperatureDeparture.Unit, wh.Past24HourTemperatureDeparture.UnitType, wh.Past24HourTemperatureDeparture.Phrase,
		wh.ApparentTemperature.Value, wh.ApparentTemperature.Unit, wh.ApparentTemperature.UnitType, wh.ApparentTemperature.Phrase,
		wh.WindChillTemperature.Value, wh.WindChillTemperature.Unit, wh.WindChillTemperature.UnitType, wh.WindChillTemperature.Phrase,
		wh.WetBulbTemperature.Value, wh.WetBulbTemperature.Unit, wh.WetBulbTemperature.UnitType, wh.WetBulbTemperature.Phrase,
		wh.WetBulbGlobeTemperature.Value, wh.WetBulbGlobeTemperature.Unit, wh.WetBulbGlobeTemperature.UnitType, wh.WetBulbGlobeTemperature.Phrase,
		wh.Precip1hr.Value, wh.Precip1hr.Unit, wh.Precip1hr.UnitType, wh.Precip1hr.Phrase,
		wh.PrecipitationSummary.Precipitation.Value, wh.PrecipitationSummary.Precipitation.Unit, wh.PrecipitationSummary.Precipitation.UnitType, wh.PrecipitationSummary.Precipitation.Phrase,
		wh.PrecipitationSummary.PastHour.Value, wh.PrecipitationSummary.PastHour.Unit, wh.PrecipitationSummary.PastHour.UnitType, wh.PrecipitationSummary.PastHour.Phrase,
		wh.PrecipitationSummary.Past3Hours.Value, wh.PrecipitationSummary.Past3Hours.Unit, wh.PrecipitationSummary.Past3Hours.UnitType, wh.PrecipitationSummary.Past3Hours.Phrase,
		wh.PrecipitationSummary.Past6Hours.Value, wh.PrecipitationSummary.Past6Hours.Unit, wh.PrecipitationSummary.Past6Hours.UnitType, wh.PrecipitationSummary.Past6Hours.Phrase,
		wh.PrecipitationSummary.Past9Hours.Value, wh.PrecipitationSummary.Past9Hours.Unit, wh.PrecipitationSummary.Past9Hours.UnitType, wh.PrecipitationSummary.Past9Hours.Phrase,
		wh.PrecipitationSummary.Past12Hours.Value, wh.PrecipitationSummary.Past12Hours.Unit, wh.PrecipitationSummary.Past12Hours.UnitType, wh.PrecipitationSummary.Past12Hours.Phrase,
		wh.PrecipitationSummary.Past18Hours.Value, wh.PrecipitationSummary.Past18Hours.Unit, wh.PrecipitationSummary.Past18Hours.UnitType, wh.PrecipitationSummary.Past18Hours.Phrase,
		wh.PrecipitationSummary.Past24Hours.Value, wh.PrecipitationSummary.Past24Hours.Unit, wh.PrecipitationSummary.Past24Hours.UnitType, wh.PrecipitationSummary.Past24Hours.Phrase,
		wh.TemperatureSummary.Past6HourRange.Min.Value, wh.TemperatureSummary.Past6HourRange.Min.Unit, wh.TemperatureSummary.Past6HourRange.Min.UnitType, wh.TemperatureSummary.Past6HourRange.Min.Phrase,
		wh.TemperatureSummary.Past6HourRange.Max.Value, wh.TemperatureSummary.Past6HourRange.Max.Unit, wh.TemperatureSummary.Past6HourRange.Max.UnitType, wh.TemperatureSummary.Past6HourRange.Max.Phrase,
		wh.TemperatureSummary.Past12HourRange.Min.Value, wh.TemperatureSummary.Past12HourRange.Min.Unit, wh.TemperatureSummary.Past12HourRange.Min.UnitType, wh.TemperatureSummary.Past12HourRange.Min.Phrase,
		wh.TemperatureSummary.Past12HourRange.Max.Value, wh.TemperatureSummary.Past12HourRange.Max.Unit, wh.TemperatureSummary.Past12HourRange.Max.UnitType, wh.TemperatureSummary.Past12HourRange.Max.Phrase,
		wh.TemperatureSummary.Past24HourRange.Min.Value, wh.TemperatureSummary.Past24HourRange.Min.Unit, wh.TemperatureSummary.Past24HourRange.Min.UnitType, wh.TemperatureSummary.Past24HourRange.Min.Phrase,
		wh.TemperatureSummary.Past24HourRange.Max.Value, wh.TemperatureSummary.Past24HourRange.Max.Unit, wh.TemperatureSummary.Past24HourRange.Max.UnitType, wh.TemperatureSummary.Past24HourRange.Max.Phrase,
		wh.MobileLink, wh.Link)
	return err
}

func Get_Location_byKey(ctx context.Context, accuweatherKey string) (WeatherLocation, error) {
	var loc WeatherLocation
	err := config.DB.QueryRowContext(ctx, `
		SELECT id, key, key_type, latitude, longitude,
		       elevation_value, elevation_unit, elevation_unit_type,
		       localized_name, postal_code
		FROM weather_location WHERE key = ?`, accuweatherKey).
		Scan(&loc.ID, &loc.Key, &loc.KeyType, &loc.Latitude, &loc.Longitude,
			&loc.Elevation.Value, &loc.Elevation.Unit, &loc.Elevation.UnitType,
			&loc.LocalizedName, &loc.PostalCode)
	return loc, err
}

func (loc *WeatherLocation) Save(ctx context.Context) error {
	if loc.ID == "" {
		loc.ID = uuid.New().String()
	}
	_, err := config.DB.ExecContext(ctx, `
		INSERT INTO weather_location (id, key, key_type, latitude, longitude,
		                              elevation_value, elevation_unit, elevation_unit_type,
		                              localized_name, postal_code)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(key) DO UPDATE SET
			key_type          = excluded.key_type,
			latitude          = excluded.latitude,
			longitude         = excluded.longitude,
			elevation_value   = excluded.elevation_value,
			elevation_unit    = excluded.elevation_unit,
			elevation_unit_type = excluded.elevation_unit_type,
			localized_name    = excluded.localized_name,
			postal_code       = excluded.postal_code`,
		loc.ID, loc.Key, loc.KeyType, loc.Latitude, loc.Longitude,
		loc.Elevation.Value, loc.Elevation.Unit, loc.Elevation.UnitType,
		loc.LocalizedName, loc.PostalCode)
	return err
}

// ---- data transformation ----------------------------------------------------

func FromAccuweatherMeasure(awv AccuweatherMeasure) (wv WeatherValue) {
	wv = WeatherValue{
		Value:    awv.Metric.Value,
		Unit:     awv.Metric.Unit,
		UnitType: awv.Metric.UnitType,
		Phrase:   awv.Metric.Phrase,
	}
	return wv
}

func FromAccuweatherRange(awr struct {
	Minimum AccuweatherMeasure
	Maximum AccuweatherMeasure
}) (wr WeatherValueRange) {
	wr.Min = FromAccuweatherMeasure(awr.Minimum)
	wr.Max = FromAccuweatherMeasure(awr.Maximum)
	return wr
}

func (wh *WeatherHourlyData) FromAccuweatherCurrent(awc AccuweatherCurrent) error {
	wh.ID = ""
	wh.LocationID = ""
	wh.LocalObservationDateTime = awc.LocalObservationDateTime
	wh.WeatherText = awc.WeatherText
	wh.WeatherIcon = awc.WeatherIcon
	wh.HasPrecipitation = awc.HasPrecipitation
	wh.PrecipitationType = awc.PrecipitationType
	wh.IsDayTime = awc.IsDayTime
	wh.Temperature = FromAccuweatherMeasure(awc.Temperature)
	wh.RealFeelTemperature = FromAccuweatherMeasure(awc.RealFeelTemperature)
	wh.RealFeelTemperatureShade = FromAccuweatherMeasure(awc.RealFeelTemperatureShade)
	wh.RelativeHumidity = awc.RelativeHumidity
	wh.IndoorRelativeHumidity = awc.IndoorRelativeHumidity
	wh.DewPoint = FromAccuweatherMeasure(awc.DewPoint)
	wh.Wind = WindData{
		Degrees:   awc.Wind.Direction.Degrees,
		Localized: awc.Wind.Direction.Localized,
		English:   awc.Wind.Direction.English,
		Speed:     FromAccuweatherMeasure(awc.Wind.Speed),
	}
	wh.WindGustSpeed = FromAccuweatherMeasure(awc.WindGust.Speed)
	wh.UVIndex = awc.UVIndex
	wh.UVIndexText = awc.UVIndexText
	wh.Visibility = FromAccuweatherMeasure(awc.Visibility)
	wh.ObstructionsToVisibility = awc.ObstructionsToVisibility
	wh.CloudCover = awc.CloudCover
	wh.Ceiling = FromAccuweatherMeasure(awc.Ceiling)
	wh.Pressure = FromAccuweatherMeasure(awc.Pressure)
	wh.PressureTendency = PressureTendencyData{
		LocalizedText: awc.PressureTendency.LocalizedText,
		Code:          awc.PressureTendency.Code,
	}
	wh.Past24HourTemperatureDeparture = FromAccuweatherMeasure(awc.Past24HourTemperatureDeparture)
	wh.ApparentTemperature = FromAccuweatherMeasure(awc.ApparentTemperature)
	wh.WindChillTemperature = FromAccuweatherMeasure(awc.WindChillTemperature)
	wh.WetBulbTemperature = FromAccuweatherMeasure(awc.WetBulbTemperature)
	wh.WetBulbGlobeTemperature = FromAccuweatherMeasure(awc.WetBulbGlobeTemperature)
	wh.Precip1hr = FromAccuweatherMeasure(awc.Precip1hr)
	wh.PrecipitationSummary = PrecipitationSummaryData{
		Precipitation: FromAccuweatherMeasure(awc.PrecipitationSummary.Precipitation),
		PastHour:      FromAccuweatherMeasure(awc.PrecipitationSummary.PastHour),
		Past3Hours:    FromAccuweatherMeasure(awc.PrecipitationSummary.Past3Hours),
		Past6Hours:    FromAccuweatherMeasure(awc.PrecipitationSummary.Past6Hours),
		Past9Hours:    FromAccuweatherMeasure(awc.PrecipitationSummary.Past9Hours),
		Past12Hours:   FromAccuweatherMeasure(awc.PrecipitationSummary.Past12Hours),
		Past18Hours:   FromAccuweatherMeasure(awc.PrecipitationSummary.Past18Hours),
		Past24Hours:   FromAccuweatherMeasure(awc.PrecipitationSummary.Past24Hours),
	}
	wh.TemperatureSummary = TemperatureSummaryData{
		Past6HourRange:  FromAccuweatherRange(awc.TemperatureSummary.Past6HourRange),
		Past12HourRange: FromAccuweatherRange(awc.TemperatureSummary.Past12HourRange),
		Past24HourRange: FromAccuweatherRange(awc.TemperatureSummary.Past24HourRange),
	}
	wh.MobileLink = awc.MobileLink
	wh.Link = awc.Link
	return nil
}

func FromAccuweatherLocation(awl AccuweatherLocation) (WeatherLocation, error) {
	return WeatherLocation{
		ID:            "",
		Latitude:      awl.GeoPosition.Latitude,
		Longitude:     awl.GeoPosition.Longitude,
		Elevation:     WeatherValue(awl.GeoPosition.Elevation.Metric),
		Key:           awl.Key,
		KeyType:       awl.Type,
		LocalizedName: awl.LocalizedName,
		PostalCode:    awl.PrimaryPostalCode,
	}, nil
}

// ---- aggregation ------------------------------------------------------------

/*
Retrieves weather data for the given location and date, and the day before.
Then builds a WeatherDaySummary from the data :
- night = data from the day before after 15:00, till sunrise today
- day = data from today (daytime only)
*/
func (loc *WeatherLocation) GetDaySummary(ctx context.Context, dayStr string) (WeatherDaySummary, error) {
	var summ WeatherDaySummary
	data, err := loc.GetHourlyData(ctx, dayStr)
	if err != nil {
		return summ, nil
	}
	summ.LocationID = loc.ID
	summ.Date = dayStr
	summ.WeatherText = GetMostFrequentValue(data, func(item WeatherHourlyData) any { return item.WeatherText }).(string)
	summ.WeatherIcon = GetMostFrequentValue(data, func(item WeatherHourlyData) any { return item.WeatherIcon }).(int)
	summ.Precipitation = reduce(data, func(acc WeatherValue, item WeatherHourlyData) WeatherValue {
		if item.HasPrecipitation {
			acc.Value += item.Precip1hr.Value
			acc.Unit = item.Precip1hr.Unit
		}
		return acc
	}, WeatherValue{Value: 0})
	return summ, err
}

/*
Récupère les données météo pour la date donnée,
ainsi que celles de la veille après 15h.
*/
func (loc *WeatherLocation) GetHourlyData(ctx context.Context, dayStr string) ([]WeatherHourlyData, error) {
	day, _ := time.Parse("2006-01-02", dayStr)
	dayBefore := day.AddDate(0, 0, -1)
	dayBeforeStr := dayBefore.Format("2006-01-02")
	dayAfterStr := day.AddDate(0, 0, 1).Format("2006-01-02")

	rows, err := config.DB.QueryContext(ctx, `
		SELECT
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
		FROM weather_hourly
		WHERE location_id = ?
		  AND local_observation_datetime >= ?
		  AND local_observation_datetime < ?
		ORDER BY local_observation_datetime ASC`,
		loc.ID, dayBeforeStr, dayAfterStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allData []WeatherHourlyData
	for rows.Next() {
		var wh WeatherHourlyData
		if err := rows.Scan(
			&wh.ID, &wh.LocationID, &wh.LocalObservationDateTime,
			&wh.WeatherText, &wh.WeatherIcon, &wh.HasPrecipitation, &wh.PrecipitationType, &wh.IsDayTime,
			&wh.Temperature.Value, &wh.Temperature.Unit, &wh.Temperature.UnitType, &wh.Temperature.Phrase,
			&wh.RealFeelTemperature.Value, &wh.RealFeelTemperature.Unit, &wh.RealFeelTemperature.UnitType, &wh.RealFeelTemperature.Phrase,
			&wh.RealFeelTemperatureShade.Value, &wh.RealFeelTemperatureShade.Unit, &wh.RealFeelTemperatureShade.UnitType, &wh.RealFeelTemperatureShade.Phrase,
			&wh.RelativeHumidity, &wh.IndoorRelativeHumidity,
			&wh.DewPoint.Value, &wh.DewPoint.Unit, &wh.DewPoint.UnitType, &wh.DewPoint.Phrase,
			&wh.Wind.Degrees, &wh.Wind.Localized, &wh.Wind.English,
			&wh.Wind.Speed.Value, &wh.Wind.Speed.Unit, &wh.Wind.Speed.UnitType, &wh.Wind.Speed.Phrase,
			&wh.WindGustSpeed.Value, &wh.WindGustSpeed.Unit, &wh.WindGustSpeed.UnitType, &wh.WindGustSpeed.Phrase,
			&wh.UVIndex, &wh.UVIndexText,
			&wh.Visibility.Value, &wh.Visibility.Unit, &wh.Visibility.UnitType, &wh.Visibility.Phrase,
			&wh.ObstructionsToVisibility, &wh.CloudCover,
			&wh.Ceiling.Value, &wh.Ceiling.Unit, &wh.Ceiling.UnitType, &wh.Ceiling.Phrase,
			&wh.Pressure.Value, &wh.Pressure.Unit, &wh.Pressure.UnitType, &wh.Pressure.Phrase,
			&wh.PressureTendency.LocalizedText, &wh.PressureTendency.Code,
			&wh.Past24HourTemperatureDeparture.Value, &wh.Past24HourTemperatureDeparture.Unit, &wh.Past24HourTemperatureDeparture.UnitType, &wh.Past24HourTemperatureDeparture.Phrase,
			&wh.ApparentTemperature.Value, &wh.ApparentTemperature.Unit, &wh.ApparentTemperature.UnitType, &wh.ApparentTemperature.Phrase,
			&wh.WindChillTemperature.Value, &wh.WindChillTemperature.Unit, &wh.WindChillTemperature.UnitType, &wh.WindChillTemperature.Phrase,
			&wh.WetBulbTemperature.Value, &wh.WetBulbTemperature.Unit, &wh.WetBulbTemperature.UnitType, &wh.WetBulbTemperature.Phrase,
			&wh.WetBulbGlobeTemperature.Value, &wh.WetBulbGlobeTemperature.Unit, &wh.WetBulbGlobeTemperature.UnitType, &wh.WetBulbGlobeTemperature.Phrase,
			&wh.Precip1hr.Value, &wh.Precip1hr.Unit, &wh.Precip1hr.UnitType, &wh.Precip1hr.Phrase,
			&wh.PrecipitationSummary.Precipitation.Value, &wh.PrecipitationSummary.Precipitation.Unit, &wh.PrecipitationSummary.Precipitation.UnitType, &wh.PrecipitationSummary.Precipitation.Phrase,
			&wh.PrecipitationSummary.PastHour.Value, &wh.PrecipitationSummary.PastHour.Unit, &wh.PrecipitationSummary.PastHour.UnitType, &wh.PrecipitationSummary.PastHour.Phrase,
			&wh.PrecipitationSummary.Past3Hours.Value, &wh.PrecipitationSummary.Past3Hours.Unit, &wh.PrecipitationSummary.Past3Hours.UnitType, &wh.PrecipitationSummary.Past3Hours.Phrase,
			&wh.PrecipitationSummary.Past6Hours.Value, &wh.PrecipitationSummary.Past6Hours.Unit, &wh.PrecipitationSummary.Past6Hours.UnitType, &wh.PrecipitationSummary.Past6Hours.Phrase,
			&wh.PrecipitationSummary.Past9Hours.Value, &wh.PrecipitationSummary.Past9Hours.Unit, &wh.PrecipitationSummary.Past9Hours.UnitType, &wh.PrecipitationSummary.Past9Hours.Phrase,
			&wh.PrecipitationSummary.Past12Hours.Value, &wh.PrecipitationSummary.Past12Hours.Unit, &wh.PrecipitationSummary.Past12Hours.UnitType, &wh.PrecipitationSummary.Past12Hours.Phrase,
			&wh.PrecipitationSummary.Past18Hours.Value, &wh.PrecipitationSummary.Past18Hours.Unit, &wh.PrecipitationSummary.Past18Hours.UnitType, &wh.PrecipitationSummary.Past18Hours.Phrase,
			&wh.PrecipitationSummary.Past24Hours.Value, &wh.PrecipitationSummary.Past24Hours.Unit, &wh.PrecipitationSummary.Past24Hours.UnitType, &wh.PrecipitationSummary.Past24Hours.Phrase,
			&wh.TemperatureSummary.Past6HourRange.Min.Value, &wh.TemperatureSummary.Past6HourRange.Min.Unit, &wh.TemperatureSummary.Past6HourRange.Min.UnitType, &wh.TemperatureSummary.Past6HourRange.Min.Phrase,
			&wh.TemperatureSummary.Past6HourRange.Max.Value, &wh.TemperatureSummary.Past6HourRange.Max.Unit, &wh.TemperatureSummary.Past6HourRange.Max.UnitType, &wh.TemperatureSummary.Past6HourRange.Max.Phrase,
			&wh.TemperatureSummary.Past12HourRange.Min.Value, &wh.TemperatureSummary.Past12HourRange.Min.Unit, &wh.TemperatureSummary.Past12HourRange.Min.UnitType, &wh.TemperatureSummary.Past12HourRange.Min.Phrase,
			&wh.TemperatureSummary.Past12HourRange.Max.Value, &wh.TemperatureSummary.Past12HourRange.Max.Unit, &wh.TemperatureSummary.Past12HourRange.Max.UnitType, &wh.TemperatureSummary.Past12HourRange.Max.Phrase,
			&wh.TemperatureSummary.Past24HourRange.Min.Value, &wh.TemperatureSummary.Past24HourRange.Min.Unit, &wh.TemperatureSummary.Past24HourRange.Min.UnitType, &wh.TemperatureSummary.Past24HourRange.Min.Phrase,
			&wh.TemperatureSummary.Past24HourRange.Max.Value, &wh.TemperatureSummary.Past24HourRange.Max.Unit, &wh.TemperatureSummary.Past24HourRange.Max.UnitType, &wh.TemperatureSummary.Past24HourRange.Max.Phrase,
			&wh.MobileLink, &wh.Link,
		); err != nil {
			return nil, err
		}
		allData = append(allData, wh)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	nightBefore := time.Date(dayBefore.Year(), dayBefore.Month(), dayBefore.Day(), 15, 0, 0, 0, time.Local)
	var data []WeatherHourlyData
	for _, result := range allData {
		resDate, _ := time.Parse("2006-01-02T15:04:05-0700", result.LocalObservationDateTime)
		if resDate.After(day) || (resDate.After(nightBefore) && !result.IsDayTime) {
			data = append(data, result)
		}
	}
	return data, nil
}

// ---- generic helpers --------------------------------------------------------

func reduce[T, M any](s []T, f func(M, T) M, initValue M) M {
	acc := initValue
	for _, v := range s {
		acc = f(acc, v)
	}
	return acc
}

func slice_find[T any](sl []T, found func(item T) bool) int {
	for i, v := range sl {
		if found(v) {
			return i
		}
	}
	return -1
}

func GetMostFrequentValue[W any, V comparable](data []W, getValue func(item W) V) V {
	type keyCount struct {
		key   V
		count int
	}
	var values []keyCount
	for _, item := range data {
		value := getValue(item)
		index := slice_find(values, func(v keyCount) bool { return v.key == value })
		if index == -1 {
			values = append(values, keyCount{key: value, count: 0})
		} else {
			values[index].count += 1
		}
	}
	sort.SliceStable(values, func(i, j int) bool {
		return values[i].count < values[j].count
	})
	return values[len(values)-1].key
}

func ValueCountIncrement[K comparable](count map[K]int, value K) {
	count[value] += 1
}

func (item WeatherValue) increment(newitem WeatherValue) {
	item.Unit = newitem.Unit
	item.Value += newitem.Value
}

func (item WeatherValueRange) updateWith(newitem WeatherValueRange) {
	if newitem.Min.Value < item.Min.Value {
		item.Min = newitem.Min
	}
	if newitem.Max.Value > item.Max.Value {
		item.Max = newitem.Max
	}
}

func (summ WeatherDaySummary) FromWeatherHourly(data []WeatherHourlyData) {
	var weatherTextValues map[string]int
	var weatherIconValues map[int]int
	for _, hourly := range data {
		ValueCountIncrement(weatherTextValues, hourly.WeatherText)
		ValueCountIncrement(weatherIconValues, hourly.WeatherIcon)
		summ.Precipitation.increment(hourly.Precip1hr)
		if hourly.IsDayTime {
			summ.DayTemperature.updateWith(hourly.TemperatureSummary.Past12HourRange)
		}
	}
}
