package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// Serve every file in the web directory (HTML, JS, CSS, assets, etc.).
	fileServer := http.FileServer(http.Dir("web"))
	mux.Handle("/", fileServer)
	mux.HandleFunc("GET /units", getUnits)

	log.Println("listening on http://localhost:9742")
	if err := http.ListenAndServe(":9742", mux); err != nil {
		log.Fatal(err)
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
