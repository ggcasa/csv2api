package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"csv2api" // Importă pachetul din rădăcină

	_ "github.com/glebarez/go-sqlite"
	"gopkg.in/yaml.v3"
)

type Masina struct {
	ID_Bula       string `json:"id_bula"`
	Marca         string `json:"marca"`
	Culoare       string `json:"culoare"`
	MarimeVolum   int    `json:"marime_volum"`
	GreutateKg    int    `json:"greutate_kg"`
	LocationBlock string `json:"location_block"`
}

var db *sql.DB
var tabel string

func incarcaConfig() csv2api.Config {
	file, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Eroare la citirea config.yaml: %v", err)
	}
	var cfg csv2api.Config
	_ = yaml.Unmarshal(file, &cfg)
	return cfg
}

func main() {
	cfg := incarcaConfig()
	tabel = cfg.Migrare.NumeTabel

	var err error
	db, err = sql.Open("sqlite", cfg.Migrare.CaleBazaDate)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, _ = db.Exec("PRAGMA journal_mode = WAL;")

	http.HandleFunc(cfg.Server.ApiEndpoint, masiniHandler)

	log.Printf("Serverul ruleaza pe http://localhost%s", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(cfg.Server.Port, nil))
}

func masiniHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	limita := r.URL.Query().Get("limit")
	if limita == "" {
		limita = "5000"
	}

	query := fmt.Sprintf("SELECT id_bula, marca, culoare, marime_volum, greutate_kg, location_block FROM %s LIMIT ?", tabel)
	rows, err := db.Query(query, limita)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var lista []Masina
	for rows.Next() {
		var m Masina
		_ = rows.Scan(&m.ID_Bula, &m.Marca, &m.Culoare, &m.MarimeVolum, &m.GreutateKg, &m.LocationBlock)
		lista = append(lista, m)
	}
	_ = json.NewEncoder(w).Encode(lista)
}
