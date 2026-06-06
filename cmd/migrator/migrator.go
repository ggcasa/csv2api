package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"csv2api" // Importă pachetul din rădăcină

	_ "github.com/glebarez/go-sqlite"
	"gopkg.in/yaml.v3"
)

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

	csvFile, err := os.Open(cfg.Migrare.CaleIntrareCsv)
	if err != nil {
		log.Fatalf("Eroare la deschidere CSV: %v", err)
	}
	defer csvFile.Close()

	db, err := sql.Open("sqlite", cfg.Migrare.CaleBazaDate)
	if err != nil {
		log.Fatalf("Eroare la conectare DB: %v", err)
	}
	defer db.Close()

	_, _ = db.Exec("PRAGMA journal_mode = WAL;")
	_, _ = db.Exec("PRAGMA synchronous = OFF;")

	queryTabel := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id_bula TEXT PRIMARY KEY,
		marca TEXT,
		culoare TEXT,
		marime_volum INTEGER,
		greutate_kg INTEGER,
		location_block TEXT
	);`, cfg.Migrare.NumeTabel)
	_, _ = db.Exec(queryTabel)

	reader := csv.NewReader(csvFile)
	_, _ = reader.Read() // Sari peste antet

	tx, _ := db.Begin()
	queryInsert := fmt.Sprintf("INSERT OR REPLACE INTO %s VALUES (?, ?, ?, ?, ?, ?)", cfg.Migrare.NumeTabel)
	stmt, _ := tx.Prepare(queryInsert)
	defer stmt.Close()

	fmt.Println("Se importa datele in SQLite...")
	for {
		rec, err := reader.Read()
		if err == io.EOF {
			break
		}
		_, _ = stmt.Exec(rec[0], rec[1], rec[2], rec[3], rec[4], rec[5])
	}
	_ = tx.Commit()
	fmt.Println("Import finalizat cu succes!")
}
