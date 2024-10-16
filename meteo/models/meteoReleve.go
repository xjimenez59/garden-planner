package models

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type InfoClimaStation struct {
	Id        string            `json:"id"`
	Name      string            `json:"name"`
	Latitude  float64           `json:"latitude"`
	Longitude float64           `json:"longitude"`
	Elevation float64           `json:"elevation"`
	Type      string            `json:"type"`
	License   map[string]string `json:"license"`
}

type InfoClimatResponse struct {
	Status   string                 `json:"status"`
	Errors   []string               `json:"errors"`
	Data     []string               `json:"data"`
	Stations []InfoClimaStation     `json:"stations"`
	Metadata map[string]string      `json:"metadata"`
	Hourly   map[string]interface{} `json:"hourly"`
}

type MeteoReleve struct {
	Station        string `json:"id_station"`
	Date           string `json:"dh_utc"`
	Temperature    string `json:"temperature"`
	Pression       string `json:"pression"`
	Humidite       string `json:"humidite"`
	Point_de_rosee string `json:"point_de_rosee"`
	Vent_moyen     string `json:"vent_moyen"`
	Vent_rafales   string `json:"vent_rafales,omitempty"`
	Vent_direction string `json:"vent_direction"`
	Pluie_3h       string `json:"pluie_3h,omitempty"`
	Pluie_1h       string `json:"pluie_1h,omitempty"`
}

// Renvoie les infos météo pour la station et la journée passés
func GetMeteo(ctx context.Context, site string, date string) (result []MeteoReleve, err error) {

	responseData, err := callInfoClimatApi(site, date)
	if err == nil {
		data := InfoClimatResponse{}
		if err = json.Unmarshal(responseData, &data); err != nil {
			fmt.Printf("json.Unmarshal failed. Error:  %v\n", err)
		}
		hourly := data.Hourly
		delete(hourly, "_params") //--- on ne laisse que les relevés météo dans
		bytes, _ := json.Marshal(hourly["STATIC0095"])
		result = []MeteoReleve{}
		if err = json.Unmarshal(bytes, &result); err != nil {
			fmt.Printf("json.Unmarshal failed. Error:  %v\n", err)
		}
	}

	return result, err
}

func callInfoClimatApi(site string, date string) (responseData []byte, err error) {
	doActualCall := false //-- false en mode "developpement" ; true pour faire les vrais appels à l'api depuis la bonne IP
	if doActualCall {
		infoClimatUrl := "https://www.infoclimat.fr/opendata/?method=get&format=json&stations[]=STATIC0095&start=2024-10-13&end=2024-10-13&token=aG03ob8C35ysRMzetxIQAay57KojqXvEZrDicGfATATJVz1CSmpw"
		response, err := http.Get(infoClimatUrl)
		if err == nil {
			responseData, err = io.ReadAll(response.Body)
		}
	} else {
		responseData = getMockMeteo()
	}
	return responseData, err
}

func getMockMeteo() []byte {
	stringdata := `	{
        "status": "OK",
        "errors": [],
        "data": [],
        "stations": [
            {
                "id": "STATIC0095",
                "name": "Ile-d'Arz",
                "latitude": 47.589,
                "longitude": -2.804,
                "elevation": 15,
                "type": "static",
                "license": {
                    "license": "NON-COMMERCIAL ONLY: CC BY NC",
                    "url": "https:\/\/creativecommons.org\/licenses\/by-nc\/2.0\/fr\/",
                    "source": "infoclimat.fr",
                    "metadonnees": "https:\/\/www.infoclimat.fr\/stations\/metadonnees.php?id=STATIC0095"
                }
            }
        ],
        "metadata": {
            "temperature": "temperature,degC",
            "pression": "mean sea level pressure,hPa",
            "humidite": "relative humidity,%",
            "point_de_rosee": "dewpoint,degC",
            "visibilite": "horizontal visibility,m",
            "vent_moyen": "mean wind speed,km\/h",
            "vent_rafales": "wind gust,km\/h",
            "vent_direction": "wind direction,deg",
            "pluie_3h": "precipitation over 3h,mm",
            "pluie_1h": "precipitation over 1h,mm",
            "neige_au_sol": "snow depth,cm",
            "nebulosite": "Ncloud cover,octats",
            "temps_omm": "present weather,http:\/\/www.infoclimat.fr\/stations-meteo\/ww.php"
        },
        "hourly": {
            "STATIC0095": [
                {
                    "id_station": "STATIC0095",
                    "dh_utc": "2024-10-13 00:00:00",
                    "temperature": "14.3",
                    "pression": "1013.8",
                    "humidite": "89",
                    "point_de_rosee": "12.8",
                    "vent_moyen": "3.2",
                    "vent_rafales": null,
                    "vent_direction": "34",
                    "pluie_3h": null,
                    "pluie_1h": null
                },
                {
                    "id_station": "STATIC0095",
                    "dh_utc": "2024-10-13 00:15:00",
                    "temperature": "14.3",
                    "pression": "1013.8",
                    "humidite": "89",
                    "point_de_rosee": "12.8",
                    "vent_moyen": "3.2",
                    "vent_rafales": null,
                    "vent_direction": "34",
                    "pluie_3h": null,
                    "pluie_1h": null
                },
                {
                    "id_station": "STATIC0095",
                    "dh_utc": "2024-10-13 00:30:00",
                    "temperature": "14.3",
                    "pression": "1013.8",
                    "humidite": "89",
                    "point_de_rosee": "12.8",
                    "vent_moyen": "3.2",
                    "vent_rafales": null,
                    "vent_direction": "34",
                    "pluie_3h": null,
                    "pluie_1h": null
                }                
            ],
            "_params": [
                "temperature",
                "pression",
                "humidite",
                "point_de_rosee",
                "vent_moyen",
                "vent_rafales",
                "vent_direction",
                "pluie_3h",
                "pluie_1h"
            ]
        }
    }
    `
	result := []byte("{}")
	helper := make(map[string]interface{})
	err := json.Unmarshal([]byte(stringdata), &helper)
	if err != nil {
		fmt.Printf("json.Unmarshal([]byte(s), &helper) failed. Error:  %v\n", err)
		return result
	}
	bytes, err := json.Marshal(helper)
	if err != nil {
		fmt.Printf("json.Marshal(helper) failed. Error:  %v\n", err)
		return result
	}
	result = []byte(string(bytes))
	return result
}
