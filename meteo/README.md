

===================================

exemple appel api infoclimat (données brutes de station): 
		https://www.infoclimat.fr/opendata/
			?method=get
			&format=json
			&stations[]=STATIC0095     (c'est le site sur l'ile d'ars)
			&start=2024-10-13
			&end=2024-10-13
			&token=aG03ob8C35ysRMzetxIQAay57KojqXvEZrDicGfATATJVz1CSmpw

	renvoie :
	{
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
			...
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

====================================================

exemple appel api "previsions à 7 jours" pour Le Tour du Parc :

    https://www.infoclimat.fr/public-api/gfs/json
        ?_ll=47.52521,-2.64577
        &_auth=ABoEEwZ4BCZfclptBnAHLgBoBzIOeAcgA39WNV81USwFY144VjBUNABuVSgALwMqUGEHeVpjBjpRN1Y2WzFeIgB8BGIGbARvXzlaOQY1BzQALAd4DjkHNgNiVjNfNlEtBXleMlY9VCgAblUwADUDKVBhB2VaewY%2FUTdWOVspXiIAYgRnBmcEbl84Wj4GNAcwADoHZg4uByADZlY2XzJRYQVjXjJWZ1QxAG9VYgAwA2RQZwdhWnsGOlE1VjJbMF4%2BAGsEaQZlBHlfL1pBBkUHLgBzByUOZAd5A31WZF9sUWY%3D
        &_c=731d30dcebf19e2fe0e4bd85298e3f9a

renvoie :

