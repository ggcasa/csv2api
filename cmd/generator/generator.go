package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"csv2api" // Importă pachetul din rădăcină

	"gopkg.in/yaml.v3"
)

func random(min, max int) int {
	if min >= max {
		return min
	}
	return rand.Intn(max-min) + min
}

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
	rand.Seed(time.Now().UnixNano())

	fmt.Printf("Generare CSV: %d linii -> %s\n", cfg.Generator.NumarLinii, cfg.Generator.CaleIesireCsv)

	file, err := os.Create(cfg.Generator.CaleIesireCsv)
	if err != nil {
		fmt.Printf("Eroare la creare fisier: %v\n", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	_ = writer.Write([]string{"ID_Bula", "Marca", "Culoare", "Marime_Volum", "Greutate_Kg", "Location_Block"})

	startChar := byte('A')
	var linie int64
	for linie = 0; linie < cfg.Generator.NumarLinii; linie++ {
		idAleatoriu := ""
		for i := int64(0); i < cfg.Generator.LungimeId; i++ {
			idAleatoriu += string(startChar + byte(random(0, 26)))
		}

		marca := cfg.Atribute.Marci[random(0, len(cfg.Atribute.Marci))]
		culoare := cfg.Atribute.Culori[random(0, len(cfg.Atribute.Culori))]
		bloc := cfg.Atribute.LocatieBlocuri[random(0, len(cfg.Atribute.LocatieBlocuri))]
		greutate := random(cfg.Atribute.GreutateInterval.Min, cfg.Atribute.GreutateInterval.Max)
		marime := random(cfg.Atribute.MarimeInterval.Min, cfg.Atribute.MarimeInterval.Max)

		_ = writer.Write([]string{
			idAleatoriu, marca, culoare,
			strconv.Itoa(marime), strconv.Itoa(greutate), bloc,
		})
	}
	fmt.Println("CSV generat cu succes!")
}
