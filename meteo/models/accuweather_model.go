package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const BASE_URL string = "http://dataservice.accuweather.com"
const API_KEY string = "zOS65wGzjmkKuzunEFoJ8hiTjIAy5nBE"

type AccuweatherValue struct {
	Value    float32
	Unit     string
	UnitType int
	Phrase   string `json:"Phrase,omitempty"`
}

type AccuweatherMeasure struct {
	Metric   AccuweatherValue
	Imperial AccuweatherValue
}

type AccuweatherCurrent struct {
	LocalObservationDateTime string
	EpochTime                int64
	WeatherText              string
	WeatherIcon              int
	HasPrecipitation         bool
	PrecipitationType        string
	IsDayTime                bool
	Temperature              AccuweatherMeasure
	RealFeelTemperature      AccuweatherMeasure
	RealFeelTemperatureShade AccuweatherMeasure
	RelativeHumidity         int
	IndoorRelativeHumidity   int
	DewPoint                 AccuweatherMeasure
	Wind                     struct {
		Direction struct {
			Degrees   int
			Localized string
			English   string
		}
		Speed AccuweatherMeasure
	}
	WindGust struct {
		Speed AccuweatherMeasure
	}
	UVIndex                  int
	UVIndexText              string
	Visibility               AccuweatherMeasure
	ObstructionsToVisibility string
	CloudCover               int
	Ceiling                  AccuweatherMeasure
	Pressure                 AccuweatherMeasure
	PressureTendency         struct {
		LocalizedText string
		Code          string
	}
	Past24HourTemperatureDeparture AccuweatherMeasure
	ApparentTemperature            AccuweatherMeasure
	WindChillTemperature           AccuweatherMeasure
	WetBulbTemperature             AccuweatherMeasure
	WetBulbGlobeTemperature        AccuweatherMeasure
	Precip1hr                      AccuweatherMeasure
	PrecipitationSummary           struct {
		Precipitation AccuweatherMeasure
		PastHour      AccuweatherMeasure
		Past3Hours    AccuweatherMeasure
		Past6Hours    AccuweatherMeasure
		Past9Hours    AccuweatherMeasure
		Past12Hours   AccuweatherMeasure
		Past18Hours   AccuweatherMeasure
		Past24Hours   AccuweatherMeasure
	}
	TemperatureSummary struct {
		Past6HourRange struct {
			Minimum AccuweatherMeasure
			Maximum AccuweatherMeasure
		}
		Past12HourRange struct {
			Minimum AccuweatherMeasure
			Maximum AccuweatherMeasure
		}
		Past24HourRange struct {
			Minimum AccuweatherMeasure
			Maximum AccuweatherMeasure
		}
	}
	MobileLink string
	Link       string
}

type AccuweatherLocation struct {
	Version           int
	Key               string
	Type              string
	Rank              int
	LocalizedName     string
	EnglishName       string
	PrimaryPostalCode string
	Region            struct {
		ID            string
		LocalizedName string
		EnglishName   string
	}
	Country struct {
		ID            string
		LocalizedName string
		EnglishName   string
	}
	AdministrativeArea struct {
		ID            string
		LocalizedName string
		EnglishName   string
		Level         int
		LocalizedType string
		EnglishType   string
		CountryID     string
	}
	TimeZone struct {
		Code             string
		Name             string
		GmtOffset        float32
		IsDaylightSaving bool
		NextOffsetChange string
	}
	GeoPosition struct {
		Latitude  float32
		Longitude float32
		Elevation AccuweatherMeasure
	}
	IsAlias                bool
	SupplementalAdminAreas []string
	DataSets               []string
}

// Renvoie les infos météo pour la station et la journée passés
func Accuweather_GetCurrent24H(location string) (data []AccuweatherCurrent, err error) {
	doActualCall := true //-- false en mode "developpement" ; true pour faire les vrais appels à l'api depuis la bonne IP

	url := fmt.Sprintf("%s/currentconditions/v1/%s/historical/24?language=fr&details=true&apikey=%s", BASE_URL, location, API_KEY)

	var responseData []byte
	if doActualCall {
		responseData, err = callAccuweatherApi(url)
	} else {
		responseData = getAccuweatherMockData()
		err = nil
	}
	if err == nil {
		data = make([]AccuweatherCurrent, 0)
		if err = json.Unmarshal(responseData, &data); err != nil {
			fmt.Printf("json.Unmarshal failed. Error:  %v\n", err)
		}
	}

	return data, err
}

func Accuweather_GetLocations_from_CP(cp string) (data []AccuweatherLocation, err error) {

	url := fmt.Sprintf("%s/locations/v1/search?q=%s,FR&apikey=%s", BASE_URL, cp, API_KEY)
	responseData, err := callAccuweatherApi(url)
	if err == nil {
		data = make([]AccuweatherLocation, 0)
		if err = json.Unmarshal(responseData, &data); err != nil {
			fmt.Printf("json.Unmarshal failed. Error:  %v\n", err)
		}
	}
	return data, err
}

func callAccuweatherApi(url string) (responseData []byte, err error) {
	response, err := http.Get(url)
	if err == nil {
		responseData, err = io.ReadAll(response.Body)
	}
	return responseData, err
}

func getAccuweatherMockData() []byte {
	stringdata := `
	[
	{
		"LocalObservationDateTime": "2024-10-16T16:01:00+02:00",
		"EpochTime": 1729087260,
		"WeatherText": "Faible pluie",
		"WeatherIcon": 12,
		"HasPrecipitation": true,
		"PrecipitationType": "Rain",
		"IsDayTime": true,
		"Temperature": {
			"Metric": {
				"Value": 16.5,
				"Unit": "C",
				"UnitType": 17
			},
			"Imperial": {
				"Value": 62.0,
				"Unit": "F",
				"UnitType": 18
			}
		},
		"RealFeelTemperature": {
			"Metric": {
				"Value": 13.8,
				"Unit": "C",
				"UnitType": 17,
				"Phrase": "Frais"
			},
			"Imperial": {
				"Value": 57.0,
				"Unit": "F",
				"UnitType": 18,
				"Phrase": "Frais"
			}
		},
		"RealFeelTemperatureShade": {
			"Metric": {
				"Value": 13.8,
				"Unit": "C",
				"UnitType": 17,
				"Phrase": "Frais"
			},
			"Imperial": {
				"Value": 57.0,
				"Unit": "F",
				"UnitType": 18,
				"Phrase": "Frais"
			}
		},
		"RelativeHumidity": 99,
		"IndoorRelativeHumidity": 80,
		"DewPoint": {
			"Metric": {
				"Value": 16.4,
				"Unit": "C",
				"UnitType": 17
			},
			"Imperial": {
				"Value": 62.0,
				"Unit": "F",
				"UnitType": 18
			}
		},
		"Wind": {
			"Direction": {
				"Degrees": 203,
				"Localized": "SSO",
				"English": "SSW"
			},
			"Speed": {
				"Metric": {
					"Value": 10.5,
					"Unit": "km/h",
					"UnitType": 7
				},
				"Imperial": {
					"Value": 6.5,
					"Unit": "mi/h",
					"UnitType": 9
				}
			}
		},
		"WindGust": {
			"Speed": {
				"Metric": {
					"Value": 20.5,
					"Unit": "km/h",
					"UnitType": 7
				},
				"Imperial": {
					"Value": 12.7,
					"Unit": "mi/h",
					"UnitType": 9
				}
			}
		},
		"UVIndex": 0,
		"UVIndexText": "Minimum",
		"Visibility": {
			"Metric": {
				"Value": 3.2,
				"Unit": "km",
				"UnitType": 6
			},
			"Imperial": {
				"Value": 2.0,
				"Unit": "mi",
				"UnitType": 2
			}
		},
		"ObstructionsToVisibility": "R-",
		"CloudCover": 100,
		"Ceiling": {
			"Metric": {
				"Value": 457.0,
				"Unit": "m",
				"UnitType": 5
			},
			"Imperial": {
				"Value": 1500.0,
				"Unit": "ft",
				"UnitType": 0
			}
		},
		"Pressure": {
			"Metric": {
				"Value": 1000.3,
				"Unit": "mb",
				"UnitType": 14
			},
			"Imperial": {
				"Value": 29.54,
				"Unit": "inHg",
				"UnitType": 12
			}
		},
		"PressureTendency": {
			"LocalizedText": "Stationnaire",
			"Code": "S"
		},
		"Past24HourTemperatureDeparture": {
			"Metric": {
				"Value": -3.2,
				"Unit": "C",
				"UnitType": 17
			},
			"Imperial": {
				"Value": -6.0,
				"Unit": "F",
				"UnitType": 18
			}
		},
		"ApparentTemperature": {
			"Metric": {
				"Value": 19.4,
				"Unit": "C",
				"UnitType": 17
			},
			"Imperial": {
				"Value": 67.0,
				"Unit": "F",
				"UnitType": 18
			}
		},
		"WindChillTemperature": {
			"Metric": {
				"Value": 16.7,
				"Unit": "C",
				"UnitType": 17
			},
			"Imperial": {
				"Value": 62.0,
				"Unit": "F",
				"UnitType": 18
			}
		},
		"WetBulbTemperature": {
			"Metric": {
				"Value": 16.5,
				"Unit": "C",
				"UnitType": 17
			},
			"Imperial": {
				"Value": 62.0,
				"Unit": "F",
				"UnitType": 18
			}
		},
		"WetBulbGlobeTemperature": {
			"Metric": {
				"Value": 16.6,
				"Unit": "C",
				"UnitType": 17
			},
			"Imperial": {
				"Value": 62.0,
				"Unit": "F",
				"UnitType": 18
			}
		},
		"Precip1hr": {
			"Metric": {
				"Value": 0.5,
				"Unit": "mm",
				"UnitType": 3
			},
			"Imperial": {
				"Value": 0.02,
				"Unit": "in",
				"UnitType": 1
			}
		},
		"PrecipitationSummary": {
			"Precipitation": {
				"Metric": {
					"Value": 3.3,
					"Unit": "mm",
					"UnitType": 3
				},
				"Imperial": {
					"Value": 0.13,
					"Unit": "in",
					"UnitType": 1
				}
			},
			"PastHour": {
				"Metric": {
					"Value": 0.5,
					"Unit": "mm",
					"UnitType": 3
				},
				"Imperial": {
					"Value": 0.02,
					"Unit": "in",
					"UnitType": 1
				}
			},
			"Past3Hours": {
				"Metric": {
					"Value": 8.9,
					"Unit": "mm",
					"UnitType": 3
				},
				"Imperial": {
					"Value": 0.35,
					"Unit": "in",
					"UnitType": 1
				}
			},
			"Past6Hours": {
				"Metric": {
					"Value": 11.8,
					"Unit": "mm",
					"UnitType": 3
				},
				"Imperial": {
					"Value": 0.46,
					"Unit": "in",
					"UnitType": 1
				}
			},
			"Past9Hours": {
				"Metric": {
					"Value": 18.7,
					"Unit": "mm",
					"UnitType": 3
				},
				"Imperial": {
					"Value": 0.74,
					"Unit": "in",
					"UnitType": 1
				}
			},
			"Past12Hours": {
				"Metric": {
					"Value": 32.7,
					"Unit": "mm",
					"UnitType": 3
				},
				"Imperial": {
					"Value": 1.29,
					"Unit": "in",
					"UnitType": 1
				}
			},
			"Past18Hours": {
				"Metric": {
					"Value": 33.2,
					"Unit": "mm",
					"UnitType": 3
				},
				"Imperial": {
					"Value": 1.31,
					"Unit": "in",
					"UnitType": 1
				}
			},
			"Past24Hours": {
				"Metric": {
					"Value": 33.2,
					"Unit": "mm",
					"UnitType": 3
				},
				"Imperial": {
					"Value": 1.31,
					"Unit": "in",
					"UnitType": 1
				}
			}
		},
		"TemperatureSummary": {
			"Past6HourRange": {
				"Minimum": {
					"Metric": {
						"Value": 16.5,
						"Unit": "C",
						"UnitType": 17
					},
					"Imperial": {
						"Value": 62.0,
						"Unit": "F",
						"UnitType": 18
					}
				},
				"Maximum": {
					"Metric": {
						"Value": 18.8,
						"Unit": "C",
						"UnitType": 17
					},
					"Imperial": {
						"Value": 66.0,
						"Unit": "F",
						"UnitType": 18
					}
				}
			},
			"Past12HourRange": {
				"Minimum": {
					"Metric": {
						"Value": 16.5,
						"Unit": "C",
						"UnitType": 17
					},
					"Imperial": {
						"Value": 62.0,
						"Unit": "F",
						"UnitType": 18
					}
				},
				"Maximum": {
					"Metric": {
						"Value": 18.8,
						"Unit": "C",
						"UnitType": 17
					},
					"Imperial": {
						"Value": 66.0,
						"Unit": "F",
						"UnitType": 18
					}
				}
			},
			"Past24HourRange": {
				"Minimum": {
					"Metric": {
						"Value": 16.5,
						"Unit": "C",
						"UnitType": 17
					},
					"Imperial": {
						"Value": 62.0,
						"Unit": "F",
						"UnitType": 18
					}
				},
				"Maximum": {
					"Metric": {
						"Value": 20.9,
						"Unit": "C",
						"UnitType": 17
					},
					"Imperial": {
						"Value": 70.0,
						"Unit": "F",
						"UnitType": 18
					}
				}
			}
		},
		"MobileLink": "http://www.accuweather.com/fr/fr/le-tour-du-parc/56370/current-weather/166808_pc",
		"Link": "http://www.accuweather.com/fr/fr/le-tour-du-parc/56370/current-weather/166808_pc"
	}
	]
    `
	result := &bytes.Buffer{}
	json.Compact(result, []byte(stringdata))
	return result.Bytes()
}
