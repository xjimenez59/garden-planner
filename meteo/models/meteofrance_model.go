package models

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"garden-planner/meteo/config"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	mfTokenURL         = "https://portail-api.meteofrance.fr/token"
	mfDefaultAPIBase   = "https://public-api.meteofrance.fr/public/DPClim/v1"
	mfDefaultBasicAuth = "SWI4OTRRVFRRbzlLZnBMU3RoUjNhRzlvWXhJYTpDWjhFMVZzUnRZTUtFQzJpUnQ2V29Ec0FjUEVh"
)

// ---- Token manager ----------------------------------------------------------

var mfTokenStore struct {
	sync.Mutex
	value     string
	expiresAt time.Time
}

type mfTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func getMFToken() (string, error) {
	mfTokenStore.Lock()
	defer mfTokenStore.Unlock()

	if mfTokenStore.value != "" && time.Now().Before(mfTokenStore.expiresAt.Add(-2*time.Minute)) {
		return mfTokenStore.value, nil
	}

	basicAuth := os.Getenv("METEOFRANCE_BASIC_AUTH")
	if basicAuth == "" {
		basicAuth = mfDefaultBasicAuth
	}

	req, err := http.NewRequest("POST", mfTokenURL, strings.NewReader("grant_type=client_credentials"))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+basicAuth)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("token HTTP %d: %s", resp.StatusCode, string(body))
	}

	var tr mfTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return "", fmt.Errorf("token decode: %w", err)
	}

	mfTokenStore.value = tr.AccessToken
	mfTokenStore.expiresAt = time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second)
	return mfTokenStore.value, nil
}

func mfAPIBase() string {
	if v := os.Getenv("METEOFRANCE_API_BASE"); v != "" {
		return v
	}
	return mfDefaultAPIBase
}

// ---- MFQuotidien struct (one row per station per day) -----------------------

// Column names in the exact order used for INSERT and CSV parsing.
// SQL column names are lowercase; CSV column names are uppercase.
var mfCols = []string{
	"poste", "date",
	"rr", "qrr", "drr", "qdrr",
	"tn", "qtn", "htn", "qhtn",
	"tx", "qtx", "htx", "qhtx",
	"tm", "qtm",
	"tmnx", "qtmnx",
	"tnsol", "qtnsol",
	"tn50", "qtn50",
	"dg", "qdg",
	"tampli", "qtampli",
	"tntxm", "qtntxm",
	"pmerm", "qpmerm",
	"pmermin", "qpmermin",
	"ffm", "qffm",
	"fxi", "qfxi", "dxi", "qdxi", "hxi", "qhxi",
	"fxy", "qfxy", "dxy", "qdxy", "hxy", "qhxy",
	"ff2m", "qff2m",
	"fxi2", "qfxi2", "dxi2", "qdxi2", "hxi2", "qhxi2",
	"fxi3s", "qfxi3s", "dxi3s", "qdxi3s", "hxi3s", "qhxi3s",
	"un", "qun", "hun", "qhun",
	"ux", "qux", "hux", "qhux",
	"dhumi40", "qdhumi40",
	"dhumi80", "qdhumi80",
	"tsvm", "qtsvm",
	"dhumec", "qdhumec",
	"um", "qum",
	"inst", "qinst",
	"glot", "qglot",
	"dift", "qdift",
	"dirt", "qdirt",
	"sigma", "qsigma",
	"infrart", "qinfrart",
	"uv_indicex", "quv_indicex",
	"nb300", "qnb300",
	"ba300", "qba300",
	"neig", "qneig",
	"brou", "qbrou",
	"orag", "qorag",
	"gresil", "qgresil",
	"grele", "qgrele",
	"rosee", "qrosee",
	"verglas", "qverglas",
	"solneige", "qsolneige",
	"gelee", "qgelee",
	"fumee", "qfumee",
	"brume", "qbrume",
	"eclair", "qeclair",
	"etpmon", "qetpmon",
	"etpgrille", "qetpgrille",
	"uv", "quv",
	"tmermax", "qtmermax",
	"tmermin", "qtmermin",
	"hneigef", "qhneigef",
	"neigetotx", "qneigetotx",
	"neigetot06", "qneigetot06",
}

type MFQuotidien struct {
	POSTE       string `json:"POSTE"`
	DATE        string `json:"DATE"`
	RR          string `json:"RR,omitempty"`
	QRR         string `json:"QRR,omitempty"`
	DRR         string `json:"DRR,omitempty"`
	QDRR        string `json:"QDRR,omitempty"`
	TN          string `json:"TN,omitempty"`
	QTN         string `json:"QTN,omitempty"`
	HTN         string `json:"HTN,omitempty"`
	QHTN        string `json:"QHTN,omitempty"`
	TX          string `json:"TX,omitempty"`
	QTX         string `json:"QTX,omitempty"`
	HTX         string `json:"HTX,omitempty"`
	QHTX        string `json:"QHTX,omitempty"`
	TM          string `json:"TM,omitempty"`
	QTM         string `json:"QTM,omitempty"`
	TMNX        string `json:"TMNX,omitempty"`
	QTMNX       string `json:"QTMNX,omitempty"`
	TNSOL       string `json:"TNSOL,omitempty"`
	QTNSOL      string `json:"QTNSOL,omitempty"`
	TN50        string `json:"TN50,omitempty"`
	QTN50       string `json:"QTN50,omitempty"`
	DG          string `json:"DG,omitempty"`
	QDG         string `json:"QDG,omitempty"`
	TAMPLI      string `json:"TAMPLI,omitempty"`
	QTAMPLI     string `json:"QTAMPLI,omitempty"`
	TNTXM       string `json:"TNTXM,omitempty"`
	QTNTXM      string `json:"QTNTXM,omitempty"`
	PMERM       string `json:"PMERM,omitempty"`
	QPMERM      string `json:"QPMERM,omitempty"`
	PMERMIN     string `json:"PMERMIN,omitempty"`
	QPMERMIN    string `json:"QPMERMIN,omitempty"`
	FFM         string `json:"FFM,omitempty"`
	QFFM        string `json:"QFFM,omitempty"`
	FXI         string `json:"FXI,omitempty"`
	QFXI        string `json:"QFXI,omitempty"`
	DXI         string `json:"DXI,omitempty"`
	QDXI        string `json:"QDXI,omitempty"`
	HXI         string `json:"HXI,omitempty"`
	QHXI        string `json:"QHXI,omitempty"`
	FXY         string `json:"FXY,omitempty"`
	QFXY        string `json:"QFXY,omitempty"`
	DXY         string `json:"DXY,omitempty"`
	QDXY        string `json:"QDXY,omitempty"`
	HXY         string `json:"HXY,omitempty"`
	QHXY        string `json:"QHXY,omitempty"`
	FF2M        string `json:"FF2M,omitempty"`
	QFF2M       string `json:"QFF2M,omitempty"`
	FXI2        string `json:"FXI2,omitempty"`
	QFXI2       string `json:"QFXI2,omitempty"`
	DXI2        string `json:"DXI2,omitempty"`
	QDXI2       string `json:"QDXI2,omitempty"`
	HXI2        string `json:"HXI2,omitempty"`
	QHXI2       string `json:"QHXI2,omitempty"`
	FXI3S       string `json:"FXI3S,omitempty"`
	QFXI3S      string `json:"QFXI3S,omitempty"`
	DXI3S       string `json:"DXI3S,omitempty"`
	QDXI3S      string `json:"QDXI3S,omitempty"`
	HXI3S       string `json:"HXI3S,omitempty"`
	QHXI3S      string `json:"QHXI3S,omitempty"`
	UN          string `json:"UN,omitempty"`
	QUN         string `json:"QUN,omitempty"`
	HUN         string `json:"HUN,omitempty"`
	QHUN        string `json:"QHUN,omitempty"`
	UX          string `json:"UX,omitempty"`
	QUX         string `json:"QUX,omitempty"`
	HUX         string `json:"HUX,omitempty"`
	QHUX        string `json:"QHUX,omitempty"`
	DHUMI40     string `json:"DHUMI40,omitempty"`
	QDHUMI40    string `json:"QDHUMI40,omitempty"`
	DHUMI80     string `json:"DHUMI80,omitempty"`
	QDHUMI80    string `json:"QDHUMI80,omitempty"`
	TSVM        string `json:"TSVM,omitempty"`
	QTSVM       string `json:"QTSVM,omitempty"`
	DHUMEC      string `json:"DHUMEC,omitempty"`
	QDHUMEC     string `json:"QDHUMEC,omitempty"`
	UM          string `json:"UM,omitempty"`
	QUM         string `json:"QUM,omitempty"`
	INST        string `json:"INST,omitempty"`
	QINST       string `json:"QINST,omitempty"`
	GLOT        string `json:"GLOT,omitempty"`
	QGLOT       string `json:"QGLOT,omitempty"`
	DIFT        string `json:"DIFT,omitempty"`
	QDIFT       string `json:"QDIFT,omitempty"`
	DIRT        string `json:"DIRT,omitempty"`
	QDIRT       string `json:"QDIRT,omitempty"`
	SIGMA       string `json:"SIGMA,omitempty"`
	QSIGMA      string `json:"QSIGMA,omitempty"`
	INFRART     string `json:"INFRART,omitempty"`
	QINFRART    string `json:"QINFRART,omitempty"`
	UV_INDICEX  string `json:"UV_INDICEX,omitempty"`
	QUV_INDICEX string `json:"QUV_INDICEX,omitempty"`
	NB300       string `json:"NB300,omitempty"`
	QNB300      string `json:"QNB300,omitempty"`
	BA300       string `json:"BA300,omitempty"`
	QBA300      string `json:"QBA300,omitempty"`
	NEIG        string `json:"NEIG,omitempty"`
	QNEIG       string `json:"QNEIG,omitempty"`
	BROU        string `json:"BROU,omitempty"`
	QBROU       string `json:"QBROU,omitempty"`
	ORAG        string `json:"ORAG,omitempty"`
	QORAG       string `json:"QORAG,omitempty"`
	GRESIL      string `json:"GRESIL,omitempty"`
	QGRESIL     string `json:"QGRESIL,omitempty"`
	GRELE       string `json:"GRELE,omitempty"`
	QGRELE      string `json:"QGRELE,omitempty"`
	ROSEE       string `json:"ROSEE,omitempty"`
	QROSEE      string `json:"QROSEE,omitempty"`
	VERGLAS     string `json:"VERGLAS,omitempty"`
	QVERGLAS    string `json:"QVERGLAS,omitempty"`
	SOLNEIGE    string `json:"SOLNEIGE,omitempty"`
	QSOLNEIGE   string `json:"QSOLNEIGE,omitempty"`
	GELEE       string `json:"GELEE,omitempty"`
	QGELEE      string `json:"QGELEE,omitempty"`
	FUMEE       string `json:"FUMEE,omitempty"`
	QFUMEE      string `json:"QFUMEE,omitempty"`
	BRUME       string `json:"BRUME,omitempty"`
	QBRUME      string `json:"QBRUME,omitempty"`
	ECLAIR      string `json:"ECLAIR,omitempty"`
	QECLAIR     string `json:"QECLAIR,omitempty"`
	ETPMON      string `json:"ETPMON,omitempty"`
	QETPMON     string `json:"QETPMON,omitempty"`
	ETPGRILLE   string `json:"ETPGRILLE,omitempty"`
	QETPGRILLE  string `json:"QETPGRILLE,omitempty"`
	UV          string `json:"UV,omitempty"`
	QUV         string `json:"QUV,omitempty"`
	TMERMAX     string `json:"TMERMAX,omitempty"`
	QTMERMAX    string `json:"QTMERMAX,omitempty"`
	TMERMIN     string `json:"TMERMIN,omitempty"`
	QTMERMIN    string `json:"QTMERMIN,omitempty"`
	HNEIGEF     string `json:"HNEIGEF,omitempty"`
	QHNEIGEF    string `json:"QHNEIGEF,omitempty"`
	NEIGETOTX   string `json:"NEIGETOTX,omitempty"`
	QNEIGETOTX  string `json:"QNEIGETOTX,omitempty"`
	NEIGETOT06  string `json:"NEIGETOT06,omitempty"`
	QNEIGETOT06 string `json:"QNEIGETOT06,omitempty"`
}

// fieldMap returns a map from CSV column name (uppercase) to the struct field pointer.
func (r *MFQuotidien) fieldMap() map[string]*string {
	return map[string]*string{
		"POSTE": &r.POSTE, "DATE": &r.DATE,
		"RR": &r.RR, "QRR": &r.QRR, "DRR": &r.DRR, "QDRR": &r.QDRR,
		"TN": &r.TN, "QTN": &r.QTN, "HTN": &r.HTN, "QHTN": &r.QHTN,
		"TX": &r.TX, "QTX": &r.QTX, "HTX": &r.HTX, "QHTX": &r.QHTX,
		"TM": &r.TM, "QTM": &r.QTM,
		"TMNX": &r.TMNX, "QTMNX": &r.QTMNX,
		"TNSOL": &r.TNSOL, "QTNSOL": &r.QTNSOL,
		"TN50": &r.TN50, "QTN50": &r.QTN50,
		"DG": &r.DG, "QDG": &r.QDG,
		"TAMPLI": &r.TAMPLI, "QTAMPLI": &r.QTAMPLI,
		"TNTXM": &r.TNTXM, "QTNTXM": &r.QTNTXM,
		"PMERM": &r.PMERM, "QPMERM": &r.QPMERM,
		"PMERMIN": &r.PMERMIN, "QPMERMIN": &r.QPMERMIN,
		"FFM": &r.FFM, "QFFM": &r.QFFM,
		"FXI": &r.FXI, "QFXI": &r.QFXI, "DXI": &r.DXI, "QDXI": &r.QDXI, "HXI": &r.HXI, "QHXI": &r.QHXI,
		"FXY": &r.FXY, "QFXY": &r.QFXY, "DXY": &r.DXY, "QDXY": &r.QDXY, "HXY": &r.HXY, "QHXY": &r.QHXY,
		"FF2M": &r.FF2M, "QFF2M": &r.QFF2M,
		"FXI2": &r.FXI2, "QFXI2": &r.QFXI2, "DXI2": &r.DXI2, "QDXI2": &r.QDXI2, "HXI2": &r.HXI2, "QHXI2": &r.QHXI2,
		"FXI3S": &r.FXI3S, "QFXI3S": &r.QFXI3S, "DXI3S": &r.DXI3S, "QDXI3S": &r.QDXI3S, "HXI3S": &r.HXI3S, "QHXI3S": &r.QHXI3S,
		"UN": &r.UN, "QUN": &r.QUN, "HUN": &r.HUN, "QHUN": &r.QHUN,
		"UX": &r.UX, "QUX": &r.QUX, "HUX": &r.HUX, "QHUX": &r.QHUX,
		"DHUMI40": &r.DHUMI40, "QDHUMI40": &r.QDHUMI40,
		"DHUMI80": &r.DHUMI80, "QDHUMI80": &r.QDHUMI80,
		"TSVM": &r.TSVM, "QTSVM": &r.QTSVM,
		"DHUMEC": &r.DHUMEC, "QDHUMEC": &r.QDHUMEC,
		"UM": &r.UM, "QUM": &r.QUM,
		"INST": &r.INST, "QINST": &r.QINST,
		"GLOT": &r.GLOT, "QGLOT": &r.QGLOT,
		"DIFT": &r.DIFT, "QDIFT": &r.QDIFT,
		"DIRT": &r.DIRT, "QDIRT": &r.QDIRT,
		"SIGMA": &r.SIGMA, "QSIGMA": &r.QSIGMA,
		"INFRART": &r.INFRART, "QINFRART": &r.QINFRART,
		"UV_INDICEX": &r.UV_INDICEX, "QUV_INDICEX": &r.QUV_INDICEX,
		"NB300": &r.NB300, "QNB300": &r.QNB300,
		"BA300": &r.BA300, "QBA300": &r.QBA300,
		"NEIG": &r.NEIG, "QNEIG": &r.QNEIG,
		"BROU": &r.BROU, "QBROU": &r.QBROU,
		"ORAG": &r.ORAG, "QORAG": &r.QORAG,
		"GRESIL": &r.GRESIL, "QGRESIL": &r.QGRESIL,
		"GRELE": &r.GRELE, "QGRELE": &r.QGRELE,
		"ROSEE": &r.ROSEE, "QROSEE": &r.QROSEE,
		"VERGLAS": &r.VERGLAS, "QVERGLAS": &r.QVERGLAS,
		"SOLNEIGE": &r.SOLNEIGE, "QSOLNEIGE": &r.QSOLNEIGE,
		"GELEE": &r.GELEE, "QGELEE": &r.QGELEE,
		"FUMEE": &r.FUMEE, "QFUMEE": &r.QFUMEE,
		"BRUME": &r.BRUME, "QBRUME": &r.QBRUME,
		"ECLAIR": &r.ECLAIR, "QECLAIR": &r.QECLAIR,
		"ETPMON": &r.ETPMON, "QETPMON": &r.QETPMON,
		"ETPGRILLE": &r.ETPGRILLE, "QETPGRILLE": &r.QETPGRILLE,
		"UV": &r.UV, "QUV": &r.QUV,
		"TMERMAX": &r.TMERMAX, "QTMERMAX": &r.QTMERMAX,
		"TMERMIN": &r.TMERMIN, "QTMERMIN": &r.QTMERMIN,
		"HNEIGEF": &r.HNEIGEF, "QHNEIGEF": &r.QHNEIGEF,
		"NEIGETOTX": &r.NEIGETOTX, "QNEIGETOTX": &r.QNEIGETOTX,
		"NEIGETOT06": &r.NEIGETOT06, "QNEIGETOT06": &r.QNEIGETOT06,
	}
}

// args returns INSERT arguments in the same order as mfCols.
// POSTE and DATE are required; all other fields store nil for empty strings.
func (r MFQuotidien) args() []any {
	ns := func(s string) any {
		if s == "" {
			return nil
		}
		return s
	}
	return []any{
		r.POSTE, r.DATE,
		ns(r.RR), ns(r.QRR), ns(r.DRR), ns(r.QDRR),
		ns(r.TN), ns(r.QTN), ns(r.HTN), ns(r.QHTN),
		ns(r.TX), ns(r.QTX), ns(r.HTX), ns(r.QHTX),
		ns(r.TM), ns(r.QTM),
		ns(r.TMNX), ns(r.QTMNX),
		ns(r.TNSOL), ns(r.QTNSOL),
		ns(r.TN50), ns(r.QTN50),
		ns(r.DG), ns(r.QDG),
		ns(r.TAMPLI), ns(r.QTAMPLI),
		ns(r.TNTXM), ns(r.QTNTXM),
		ns(r.PMERM), ns(r.QPMERM),
		ns(r.PMERMIN), ns(r.QPMERMIN),
		ns(r.FFM), ns(r.QFFM),
		ns(r.FXI), ns(r.QFXI), ns(r.DXI), ns(r.QDXI), ns(r.HXI), ns(r.QHXI),
		ns(r.FXY), ns(r.QFXY), ns(r.DXY), ns(r.QDXY), ns(r.HXY), ns(r.QHXY),
		ns(r.FF2M), ns(r.QFF2M),
		ns(r.FXI2), ns(r.QFXI2), ns(r.DXI2), ns(r.QDXI2), ns(r.HXI2), ns(r.QHXI2),
		ns(r.FXI3S), ns(r.QFXI3S), ns(r.DXI3S), ns(r.QDXI3S), ns(r.HXI3S), ns(r.QHXI3S),
		ns(r.UN), ns(r.QUN), ns(r.HUN), ns(r.QHUN),
		ns(r.UX), ns(r.QUX), ns(r.HUX), ns(r.QHUX),
		ns(r.DHUMI40), ns(r.QDHUMI40),
		ns(r.DHUMI80), ns(r.QDHUMI80),
		ns(r.TSVM), ns(r.QTSVM),
		ns(r.DHUMEC), ns(r.QDHUMEC),
		ns(r.UM), ns(r.QUM),
		ns(r.INST), ns(r.QINST),
		ns(r.GLOT), ns(r.QGLOT),
		ns(r.DIFT), ns(r.QDIFT),
		ns(r.DIRT), ns(r.QDIRT),
		ns(r.SIGMA), ns(r.QSIGMA),
		ns(r.INFRART), ns(r.QINFRART),
		ns(r.UV_INDICEX), ns(r.QUV_INDICEX),
		ns(r.NB300), ns(r.QNB300),
		ns(r.BA300), ns(r.QBA300),
		ns(r.NEIG), ns(r.QNEIG),
		ns(r.BROU), ns(r.QBROU),
		ns(r.ORAG), ns(r.QORAG),
		ns(r.GRESIL), ns(r.QGRESIL),
		ns(r.GRELE), ns(r.QGRELE),
		ns(r.ROSEE), ns(r.QROSEE),
		ns(r.VERGLAS), ns(r.QVERGLAS),
		ns(r.SOLNEIGE), ns(r.QSOLNEIGE),
		ns(r.GELEE), ns(r.QGELEE),
		ns(r.FUMEE), ns(r.QFUMEE),
		ns(r.BRUME), ns(r.QBRUME),
		ns(r.ECLAIR), ns(r.QECLAIR),
		ns(r.ETPMON), ns(r.QETPMON),
		ns(r.ETPGRILLE), ns(r.QETPGRILLE),
		ns(r.UV), ns(r.QUV),
		ns(r.TMERMAX), ns(r.QTMERMAX),
		ns(r.TMERMIN), ns(r.QTMERMIN),
		ns(r.HNEIGEF), ns(r.QHNEIGEF),
		ns(r.NEIGETOTX), ns(r.QNEIGETOTX),
		ns(r.NEIGETOT06), ns(r.QNEIGETOT06),
	}
}

// Save persists the row to SQLite (INSERT OR REPLACE).
func (r MFQuotidien) Save(ctx context.Context) error {
	placeholders := strings.Repeat("?,", len(mfCols))
	placeholders = placeholders[:len(placeholders)-1]
	sql := "INSERT OR REPLACE INTO meteofrance_quotidien (" +
		strings.Join(mfCols, ",") +
		") VALUES (" + placeholders + ")"
	_, err := config.DB.ExecContext(ctx, sql, r.args()...)
	return err
}

// ---- CSV parsing ------------------------------------------------------------

// ParseMFCSV parses the semicolon-delimited CSV returned by the MF fichier API.
// The first line must be the header row.
func ParseMFCSV(body io.Reader) ([]MFQuotidien, error) {
	scanner := bufio.NewScanner(body)

	if !scanner.Scan() {
		return nil, fmt.Errorf("CSV vide ou inaccessible")
	}
	headers := strings.Split(scanner.Text(), ";")

	var rows []MFQuotidien
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		values := strings.Split(line, ";")
		var row MFQuotidien
		fm := row.fieldMap()
		for i, h := range headers {
			if ptr, ok := fm[h]; ok && i < len(values) {
				*ptr = strings.TrimSpace(values[i])
			}
		}
		if row.POSTE == "" || row.DATE == "" {
			continue
		}
		rows = append(rows, row)
	}
	return rows, scanner.Err()
}

// ---- Lecture DB -------------------------------------------------------------

// MeteoQuotidienSummary contient les champs météo utiles pour l'application jardin.
type MeteoQuotidienSummary struct {
	POSTE  string `json:"POSTE"`
	DATE   string `json:"DATE"`
	RR     string `json:"RR,omitempty"`
	DRR    string `json:"DRR,omitempty"`
	DXY    string `json:"DXY,omitempty"`
	FFM    string `json:"FFM,omitempty"`
	FXI    string `json:"FXI,omitempty"`
	TM     string `json:"TM,omitempty"`
	TN     string `json:"TN,omitempty"`
	TX     string `json:"TX,omitempty"`
	INST   string `json:"INST,omitempty"`
	QINST  string `json:"QINST,omitempty"`
	SIGMA  string `json:"SIGMA,omitempty"`
	QSIGMA string `json:"QSIGMA,omitempty"`
	NB300  string `json:"NB300,omitempty"`
	QNB300 string `json:"QNB300,omitempty"`
	DG     string `json:"DG,omitempty"`
}

// GetMeteoQuotidien retourne les données journalières pour une station et une plage de dates.
// dateDeb et dateFin sont au format YYYYMMDD (correspond à la colonne 'date' de la table).
func GetMeteoQuotidien(ctx context.Context, station, dateDeb, dateFin string) ([]MeteoQuotidienSummary, error) {
	const query = `
		SELECT poste, date, rr, drr, dxy, ffm, fxi, tm, tn, tx, inst, qinst, sigma, qsigma, nb300, qnb300, dg
		FROM meteofrance_quotidien
		WHERE poste = ? AND date >= ? AND date <= ?
		ORDER BY date ASC`

	rows, err := config.DB.QueryContext(ctx, query, station, dateDeb, dateFin)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ns := func(n sql.NullString) string {
		if n.Valid {
			return n.String
		}
		return ""
	}

	var result []MeteoQuotidienSummary
	for rows.Next() {
		var r MeteoQuotidienSummary
		var rr, drr, dxy, ffm, fxi, tm, tn, tx, inst, qinst, sigma, qsigma, nb300, qnb300, dg sql.NullString
		if err := rows.Scan(
			&r.POSTE, &r.DATE,
			&rr, &drr, &dxy, &ffm, &fxi,
			&tm, &tn, &tx,
			&inst, &qinst, &sigma, &qsigma,
			&nb300, &qnb300, &dg,
		); err != nil {
			return nil, err
		}
		r.RR, r.DRR = ns(rr), ns(drr)
		r.DXY, r.FFM, r.FXI = ns(dxy), ns(ffm), ns(fxi)
		r.TM, r.TN, r.TX = ns(tm), ns(tn), ns(tx)
		r.INST, r.QINST = ns(inst), ns(qinst)
		r.SIGMA, r.QSIGMA = ns(sigma), ns(qsigma)
		r.NB300, r.QNB300 = ns(nb300), ns(qnb300)
		r.DG = ns(dg)
		result = append(result, r)
	}
	return result, rows.Err()
}

// ---- MétéoFrance API calls --------------------------------------------------

// MFCommandeQuotidienne orders daily climatological data for a station and date range.
// It returns the raw JSON body and HTTP status code from the MF API.
func MFCommandeQuotidienne(station, dateDeb, dateFin string) ([]byte, int, error) {
	token, err := getMFToken()
	if err != nil {
		return nil, 0, fmt.Errorf("authentification MF: %w", err)
	}

	url := fmt.Sprintf("%s/commande-station/quotidienne?id-station=%s&date-deb-periode=%s&date-fin-periode=%s",
		mfAPIBase(), station, dateDeb, dateFin)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("commande MF: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	return body, resp.StatusCode, err
}

// MFGetFichier retrieves the result file for the given command ID.
// Returns the parsed rows and HTTP status code.
// Status 202 means the file is not yet ready.
func MFGetFichier(idCmde string) ([]MFQuotidien, int, error) {
	token, err := getMFToken()
	if err != nil {
		return nil, 0, fmt.Errorf("authentification MF: %w", err)
	}

	url := fmt.Sprintf("%s/commande/fichier?id-cmde=%s", mfAPIBase(), idCmde)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "text/plain")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("fichier MF: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated { //-- 201 si les données sont dispo ; 202 s'il faut recommencer plus tard
		body, _ := io.ReadAll(resp.Body)
		return nil, resp.StatusCode, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	rows, err := ParseMFCSV(resp.Body)
	return rows, resp.StatusCode, err
}
