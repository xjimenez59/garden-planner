package models

import (
	"context"
	"garden-planner/meteo/config"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WeatherValue struct {
	Value    float32 `bson:"value" json:"value"`
	Unit     string  `bson:"unit" json:"unit"`
	UnitType int     `bson:"unit_type" json:"unit_type"`
	Phrase   string  `bson:"phrase,omitempty" json:"phrase,omitempty"`
}

type WeatherValueRange struct {
	Min WeatherValue `bson:"min" json:"min"`
	Max WeatherValue `bson:"max" json:"max"`
}

type WeatherLocation struct {
	ID            *primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Latitude      float32             `bson:"latitude" json:"latitude"`
	Longitude     float32             `bson:"longitude" json:"longitude"`
	Elevation     WeatherValue        `bson:"elevation" json:"elevation"`
	Key           string              `bson:"key" json:"key"`
	KeyType       string              `bson:"key_type" json:"key_type"`
	LocalizedName string              `bson:"localizedName" json:"localizedName"`
	PostalCode    string              `bson:"postalCode" json:"postalCode"`
}

type WindData struct {
	Degrees   int          `bson:"degrees" json:"degrees"`
	Localized string       `bson:"localized" json:"localized"`
	English   string       `bson:"english" json:"english"`
	Speed     WeatherValue `bson:"speed" json:"speed"`
}
type PressureTendencyData struct {
	LocalizedText string `bson:"localizedText" json:"localizedText"`
	Code          string `bson:"code" json:"code"`
}
type PrecipitationSummaryData struct {
	Precipitation WeatherValue `bson:"precipitation" json:"precipitation"`
	PastHour      WeatherValue `bson:"pastHour" json:"pastHour"`
	Past3Hours    WeatherValue `bson:"past3Hours" json:"past3Hours"`
	Past6Hours    WeatherValue `bson:"past6Hours" json:"past6Hours"`
	Past9Hours    WeatherValue `bson:"past9Hours" json:"past9Hours"`
	Past12Hours   WeatherValue `bson:"past12Hours" json:"past12Hours"`
	Past18Hours   WeatherValue `bson:"past18Hours" json:"past18Hours"`
	Past24Hours   WeatherValue `bson:"past24Hours" json:"past24Hours"`
}
type TemperatureSummaryData struct {
	Past6HourRange  WeatherValueRange `bson:"past6HourRange" json:"past6HourRange"`
	Past12HourRange WeatherValueRange `bson:"past12HourRange" json:"past12HourRange"`
	Past24HourRange WeatherValueRange `bson:"past24HourRange" json:"past24HourRange"`
}

type WeatherHourlyData struct {
	ID                             *primitive.ObjectID      `bson:"_id,omitempty" json:"id"`
	LocationID                     *primitive.ObjectID      `bson:"location_id" json:"location_id"`
	LocalObservationDateTime       string                   `bson:"localObservationDateTime" json:"localObservationDateTime"`
	WeatherText                    string                   `bson:"weatherText" json:"weatherText"`
	WeatherIcon                    int                      `bson:"weatherIcon" json:"weatherIcon"`
	HasPrecipitation               bool                     `bson:"hasPrecipitation" json:"hasPrecipitation"`
	PrecipitationType              string                   `bson:"precipitationType" json:"precipitationType"`
	IsDayTime                      bool                     `bson:"isDayTime" json:"isDayTime"`
	Temperature                    WeatherValue             `bson:"temperature" json:"temperature"`
	RealFeelTemperature            WeatherValue             `bson:"realFeelTemperature" json:"realFeelTemperature"`
	RealFeelTemperatureShade       WeatherValue             `bson:"realFeelTemperatureShade" json:"realFeelTemperatureShade"`
	RelativeHumidity               int                      `bson:"relativeHumidity" json:"relativeHumidity"`
	IndoorRelativeHumidity         int                      `bson:"indoorRelativeHumidity" json:"indoorRelativeHumidity"`
	DewPoint                       WeatherValue             `bson:"dewPoint" json:"dewPoint"`
	Wind                           WindData                 `bson:"wind" json:"wind"`
	WindGustSpeed                  WeatherValue             `bson:"windGustSpeed" json:"windGustSpeed"`
	UVIndex                        int                      `bson:"uvIndex" json:"uvIndex"`
	UVIndexText                    string                   `bson:"uvIndexText" json:"uvIndexText"`
	Visibility                     WeatherValue             `bson:"visibility" json:"visibility"`
	ObstructionsToVisibility       string                   `bson:"obstructionsToVisibility" json:"obstructionsToVisibility"`
	CloudCover                     int                      `bson:"cloudCover" json:"cloudCover"`
	Ceiling                        WeatherValue             `bson:"ceiling" json:"ceiling"`
	Pressure                       WeatherValue             `bson:"pressure" json:"pressure"`
	PressureTendency               PressureTendencyData     `bson:"pressureTendency" json:"pressureTendency"`
	Past24HourTemperatureDeparture WeatherValue             `bson:"past24HourTemperatureDeparture" json:"past24HourTemperatureDeparture"`
	ApparentTemperature            WeatherValue             `bson:"apparentTemperature" json:"apparentTemperature"`
	WindChillTemperature           WeatherValue             `bson:"windChillTemperature" json:"windChillTemperature"`
	WetBulbTemperature             WeatherValue             `bson:"wetBulbTemperature" json:"wetBulb"`
	WetBulbGlobeTemperature        WeatherValue             `bson:"wetBulbGlobeTemperature" json:"wetBulbGlobeTemperature"`
	Precip1hr                      WeatherValue             `bson:"precip1hr" json:"precip1hr"`
	PrecipitationSummary           PrecipitationSummaryData `bson:"precipitationSummary" json:"precipitationSummary"`
	TemperatureSummary             TemperatureSummaryData   `bson:"temperatureSummary" json:"temperatureSummary"`
	MobileLink                     string                   `bson:"mobileLink" json:"mobileLink"`
	Link                           string                   `bson:"link" json:"link"`
}

type WeatherDaySummary struct {
	ID                     *primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	LocationID             *primitive.ObjectID `bson:"location_id" json:"location_id"`
	Date                   string              `bson:"date"  json:"date"`
	WeatherText            string              `bson:"weatherText" json:"weatherText"`
	WeatherIcon            int                 `bson:"weatherIcon" json:"weatherIcon"`
	Precipitation          WeatherValue        `bson:"precipitationType" json:"precipitationType"`
	NightTemperature       WeatherValueRange   `bson:"nightTemperature" json:"nightTemperature"`
	DayTemperature         WeatherValueRange   `bson:"dayTemperature" json:"dayTemperature"`
	DayRealFeelTemperature WeatherValueRange   `bson:"dayRealFeelTemperature" json:"dayRealFeelTemperature"`
	WindDirection          string              `bson:"windDirection" json:"windDirection"`
	WindSpeed              WeatherValue        `bson:"windSpeed" json:"windSpeed"`
	WindGustSpeed          WeatherValue        `bson:"windGustSpeed" json:"windGustSpeed"`
	CloudCover             int                 `bson:"cloudCover" json:"cloudCover"`
}

func (wh *WeatherHourlyData) Save(ctx context.Context) (err error) {
	whCollection := config.DB.Collection("weather_hourly")
	filter := bson.D{{"location_id", wh.LocationID}, {"localObservationDateTime", wh.LocalObservationDateTime}}
	opts := options.Replace().SetUpsert(true)
	_, err = whCollection.ReplaceOne(ctx, filter, wh, opts)
	return err
}

func Get_Location_byKey(ctx context.Context, accuweatherKey string) (loc WeatherLocation, err error) {
	locationCollection := config.DB.Collection("weather_location")
	filter := bson.D{{"key", accuweatherKey}}
	err = locationCollection.FindOne(ctx, filter).Decode(&loc)
	return loc, err
}

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

func (wh *WeatherHourlyData) FromAccuweatherCurrent(awc AccuweatherCurrent) (err error) {

	wh.ID = nil
	wh.LocationID = nil
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

	return err
}

func FromAccuweatherLocation(awl AccuweatherLocation) (loc WeatherLocation, err error) {
	loc.ID = nil
	err = nil
	loc.Latitude = awl.GeoPosition.Latitude
	loc.Longitude = awl.GeoPosition.Longitude
	loc.Elevation = WeatherValue(awl.GeoPosition.Elevation.Metric)
	loc.Key = awl.Key
	loc.KeyType = awl.Type
	loc.LocalizedName = awl.LocalizedName
	loc.PostalCode = awl.PrimaryPostalCode

	return loc, err
}

func (loc *WeatherLocation) Save(ctx context.Context) (err error) {
	locationCollection := config.DB.Collection("weather_location")
	filter := bson.D{{"key", loc.Key}}
	opts := options.Replace().SetUpsert(true)
	_, err = locationCollection.ReplaceOne(ctx, filter, loc, opts)

	return err
}

/*
Retrieves weather data for the  given location and date, and the day before.
Then builds a WeatherDaySummary out of the retrieved data :
- nil  if no data is found
- night is data from the day before after sunset, till this day until sunrise.
- day is data since today at sunrise till today  at sunset.
- attributes are computed getting the mostly found values in the considered hourly data
*/
func (loc *WeatherLocation) GetDaySummary(ctx context.Context, dayStr string) (summ WeatherDaySummary, err error) {
	var data []WeatherHourlyData
	if data, err = loc.GetHourlyData(ctx, dayStr); err != nil {
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
On récupère les données météo pour la date donnée ,
ainsi que celles de la veille après le coucher de soleil
*/
func (loc *WeatherLocation) GetHourlyData(ctx context.Context, dayStr string) (data []WeatherHourlyData, err error) {
	day, _ := time.Parse("2006-01-02", dayStr)
	dayBefore := day.AddDate(0, 0, -1)
	dayBeforeStr := dayBefore.Format("2006-01-02")
	dayAfterStr := day.AddDate(0, 0, 1).Format("2006-01-02")

	whCollection := config.DB.Collection("weather_hourly")
	filter := bson.M{
		"location_id": loc.ID,
		"localObservationDateTime": bson.M{
			"$gte": dayBeforeStr,
			"lt":   dayAfterStr,
		},
	}
	opts := options.Find().SetSort(bson.D{{"localObservationDateTime", 1}})

	var results []WeatherHourlyData
	cursor, err := whCollection.Find(ctx, filter, opts)
	if err = cursor.All(ctx, cursor); err != nil {
		return nil, err
	}

	nightBefore := time.Date(dayBefore.Year(), dayBefore.Month(), dayBefore.Day(), 15, 0, 0, 0, time.Local)
	for _, result := range results {

		resDate, _ := time.Parse("2006-01-02T15:04:05-0700", result.LocalObservationDateTime)

		if resDate.After(day) || (resDate.After(nightBefore) && result.IsDayTime == false) {
			data = append(data, result)
		}
	}

	return data, err

}

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

	//-- count values
	for _, item := range data {
		value := getValue(item)
		index := slice_find(values, func(v keyCount) bool { return v.key == value })
		if index == -1 {
			values = append(values, keyCount{key: value, count: 0})
		} else {
			values[index].count += 1
		}
	}
	sort.SliceStable(values, func(i int, j int) bool {
		return values[i].count < values[j].count
	})

	return values[len(values)-1].key

}

func ValueCountIncrement[K comparable](count map[K]int, value K) {
	_, exists := count[value]
	if exists {
		count[value] += 1
	} else {
		count[value] = 1
	}
	return
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
	//var windDirectionValues map[string]int

	for _, hourly := range data {
		ValueCountIncrement(weatherTextValues, hourly.WeatherText)
		ValueCountIncrement(weatherIconValues, hourly.WeatherIcon)
		summ.Precipitation.increment(hourly.Precip1hr)
		if hourly.IsDayTime {
			summ.DayTemperature.updateWith(hourly.TemperatureSummary.Past12HourRange)
		}

	}
}
