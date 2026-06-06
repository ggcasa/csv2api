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
	Culoare       string `json:"carma"` // Păstrăm structura ta existentă
	MarimeVolum   int    `json:"marime_volum"`
	GreutateKg    int    `json:"greutate_kg"`
	LocationBlock string `json:"location_block"`
}

// Structura nouă pentru a decoda jurnalul minuscul trimis de browser
type PunctTraseu struct {
	IDBula         string `json:"id_bula"`
	LocationBlock  string `json:"location_block"`
	Latitudine     string `json:"latitudine"`
	Longitudine    string `json:"longitudine"`
	DataModificare string `json:"data_modificare"`
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

	// INTEGRARE STRUCTURĂ: Creăm tabelul nou de tracking în "cars.db" dacă nu există
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tracking_bule (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		id_bula TEXT,
		location_block TEXT,
		latitudine TEXT,
		longitudine TEXT,
		data_modificare TEXT
	);`)
	if err != nil {
		log.Fatalf("Eroare la crearea tabelului de tracking: %v", err)
	}

	// Endpoint-ul tău existent
	http.HandleFunc(cfg.Server.ApiEndpoint, masiniHandler)

	// ENDPOINT-UL NOU pentru salvarea istoricului trimis de pe laptop
	http.HandleFunc("/salveaza-tracking", salveazaTrackingHandler)

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

// HANDLERUL NOU: Primește JSON-ul compact din IndexedDB și îl pune în SQLite3
func salveazaTrackingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Metoda nu este permisa", http.StatusMethodNotAllowed)
		return
	}

	var puncte []PunctTraseu
	if err := json.NewDecoder(r.Body).Decode(&puncte); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Folosim o tranzacție asincronă rapidă, exact cum ai făcut tu în migrator.go
	tx, err := db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stmt, err := tx.Prepare("INSERT INTO tracking_bule (id_bula, location_block, latitudine, longitudine, data_modificare) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	for _, p := range puncte {
		_, err = stmt.Exec(p.IDBula, p.LocationBlock, p.Latitudine, p.Longitudine, p.DataModificare)
		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	_ = tx.Commit()
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "succes"})
}
