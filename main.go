package main

import (
	"encoding/json"
	"errors"
	"log"
	"math"
	"net/http"
	"strings"
)

func main() {
	mux := http.NewServeMux()

	// Serve every file in the web directory (HTML, JS, CSS, assets, etc.).
	fileServer := http.FileServer(http.Dir("web"))
	mux.Handle("/", fileServer)
	mux.HandleFunc("GET /units", getUnits)
	mux.HandleFunc("POST /convert", convertUnit)

	log.Println("listening on http://localhost:9742")
	if err := http.ListenAndServe(":9742", mux); err != nil {
		log.Fatal(err)
	}
}

type convertUnitRequest struct {
	Value float64 `json:"value"`
	From  string  `json:"from"`
	To    string  `json:"to"`
}

func convertUnit(w http.ResponseWriter, r *http.Request) {
	t := r.URL.Query().Get("type")
	if len(t) == 0 {
		http.Error(w, "no type provided", http.StatusBadRequest)
		return
	}

	var req convertUnitRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req.From = strings.ToLower(strings.TrimSpace(req.From))
	req.To = strings.ToLower(strings.TrimSpace(req.To))

	if req.From == "" || req.To == "" {
		http.Error(w, "from/to unit is required", http.StatusBadRequest)
		return
	}
	if math.IsNaN(req.Value) || math.IsInf(req.Value, 0) {
		http.Error(w, "invalid value", http.StatusBadRequest)
		return
	}

	var (
		result float64
		err    error
	)

	switch t {
	case "length":
		// base: meter
		lengthToMeter := map[string]float64{
			"mm": 0.001,
			"cm": 0.01,
			"m":  1,
			"km": 1000,
			"in": 0.0254,
			"ft": 0.3048,
			"yd": 0.9144,
			"mi": 1609.344,
		}
		result, err = convertLinear(req.Value, req.From, req.To, lengthToMeter)

	case "weight", "mass":
		// base: kilogram
		weightToKg := map[string]float64{
			"mg": 0.000001,
			"g":  0.001,
			"kg": 1,
			"oz": 0.028349523125,
			"lb": 0.45359237,
		}
		result, err = convertLinear(req.Value, req.From, req.To, weightToKg)

	case "temperature":
		result, err = convertTemperature(req.Value, req.From, req.To)

	default:
		http.Error(w, "Wrong type provided", http.StatusForbidden)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return a single value (JSON number)
	json.NewEncoder(w).Encode(result)
}

func convertLinear(value float64, from, to string, toBase map[string]float64) (float64, error) {
	fromRatio, ok := toBase[from]
	if !ok {
		return 0, errors.New("invalid from unit: " + from)
	}
	toRatio, ok := toBase[to]
	if !ok {
		return 0, errors.New("invalid to unit: " + to)
	}

	// value_in_base = value * fromRatio
	// result = value_in_base / toRatio
	return (value * fromRatio) / toRatio, nil
}

func convertTemperature(value float64, from, to string) (float64, error) {
	if from == to {
		return value, nil
	}

	// Convert from -> Celsius
	var c float64
	switch from {
	case "c", "°c", "celsius":
		c = value
	case "f", "°f", "fahrenheit":
		c = (value - 32) * 5 / 9
	case "k", "kelvin":
		c = value - 273.15
	default:
		return 0, errors.New("invalid from unit: " + from)
	}

	// Convert Celsius -> to
	switch to {
	case "c", "°c", "celsius":
		return c, nil
	case "f", "°f", "fahrenheit":
		return c*9/5 + 32, nil
	case "k", "kelvin":
		return c + 273.15, nil
	default:
		return 0, errors.New("invalid to unit: " + to)
	}
}

func getUnits(w http.ResponseWriter, r *http.Request) {
	t := r.URL.Query().Get("type")
	if len(t) == 0 {
		http.Error(w, "no type provided", http.StatusBadRequest)
		return
	}

	units := []string{}

	switch t {
	case "length":
		units = []string{
			"mm",
			"cm",
			"m",
			"km",
			"in",
			"ft",
		}
	case "weight":
		units = []string{
			"g",
			"kg",
			"lb",
			"oz",
		}
	case "temperature":
		units = []string{
			"c",
			"f",
			"k",
		}
	default:
		http.Error(w, "Wrong type provided", http.StatusForbidden)
		return
	}

	b, err := json.Marshal(units)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
