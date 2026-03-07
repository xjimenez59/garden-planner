package models

import (
	"math"
	"time"
)

// luneEcliptiqueLongitude calcule la longitude écliptique approchée de la lune
// (précision ~1°) d'après Meeus, "Astronomical Algorithms", ch. 47.
// n = jours depuis J2000.0.
func luneEcliptiqueLongitude(n float64) float64 {
	toRad := math.Pi / 180

	L := normDeg(218.3164477 + 13.17639648*n)  // longitude moyenne
	M := normDeg(134.9633964 + 13.06499295*n)  // anomalie moyenne
	D := normDeg(297.8501921 + 12.19074912*n)  // élongation moyenne
	Ms := normDeg(357.5291092 + 0.98560028*n)  // anomalie moyenne du soleil
	F := normDeg(93.2720950 + 13.22935024*n)   // argument de latitude

	lambda := L +
		6.289*math.Sin(M*toRad) +
		1.274*math.Sin((2*D-M)*toRad) +
		0.658*math.Sin(2*D*toRad) +
		0.214*math.Sin(2*M*toRad) -
		0.186*math.Sin(Ms*toRad) -
		0.114*math.Sin(2*F*toRad)

	return normDeg(lambda)
}

// soleilEcliptiqueLongitude calcule la longitude écliptique approchée du soleil
// (précision ~1°). n = jours depuis J2000.0.
func soleilEcliptiqueLongitude(n float64) float64 {
	toRad := math.Pi / 180

	L := normDeg(280.4665 + 0.9856474*n)  // longitude moyenne du soleil
	Ms := normDeg(357.5291 + 0.9856003*n) // anomalie moyenne du soleil

	lambda := L +
		1.9148*math.Sin(Ms*toRad) +
		0.0200*math.Sin(2*Ms*toRad)

	return normDeg(lambda)
}

// RevolutionPeriodique détermine si la lune est montante ou descendante
// (cycle tropical ~27,32 jours, oscillation de la déclinaison).
//   - Lune montante  : déclinaison croissante → longitude écliptique ∈ [270°,360°[ ∪ [0°,90°[
//   - Lune descendante : déclinaison décroissante → longitude écliptique ∈ [90°,270°[
func RevolutionPeriodique(date time.Time) string {
	n := julianDay(date) - 2451545.0
	lambda := luneEcliptiqueLongitude(n)

	if lambda >= 90 && lambda < 270 {
		return "lune_descendante"
	}
	return "lune_montante"
}

// RevolutionCyclique détermine si la lune est croissante ou décroissante
// (cycle synodique ~29,53 jours, phases lunaires).
//   - Lune croissante  : nouvelle lune → pleine lune (élongation ∈ [0°, 180°[)
//   - Lune décroissante : pleine lune → nouvelle lune (élongation ∈ [180°, 360°[)
func RevolutionCyclique(date time.Time) string {
	n := julianDay(date) - 2451545.0
	lambdaLune := luneEcliptiqueLongitude(n)
	lambdaSoleil := soleilEcliptiqueLongitude(n)

	elongation := normDeg(lambdaLune - lambdaSoleil)

	if elongation < 180 {
		return "lune_croissante"
	}
	return "lune_decroissante"
}

// julianDay calcule le Jour Julien à partir d'une date UTC.
func julianDay(t time.Time) float64 {
	t = t.UTC()
	year := t.Year()
	month := int(t.Month())
	day := t.Day()
	hour := t.Hour()
	minute := t.Minute()
	second := t.Second()

	if month <= 2 {
		year--
		month += 12
	}
	A := year / 100
	B := 2 - A + A/4

	dayFrac := float64(day) +
		float64(hour)/24.0 +
		float64(minute)/1440.0 +
		float64(second)/86400.0

	return math.Floor(365.25*float64(year+4716)) +
		math.Floor(30.6001*float64(month+1)) +
		dayFrac + float64(B) - 1524.5
}

// normDeg normalise un angle en degrés dans [0, 360[.
func normDeg(d float64) float64 {
	d = math.Mod(d, 360)
	if d < 0 {
		d += 360
	}
	return d
}
