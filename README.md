# CSV2API - Data Pipeline & Generator Engine

Acest proiect este un motor autonom de simulare, migrare și distribuție a datelor prin API. Arhitectura urmează principiul Single Source of Truth (O singură sursă de adevăr), comportamentul tuturor serviciilor fiind dictat în întregime de un fișier central de configurare în format YAML (config.yaml).

Sistemul este optimizat pentru a gestiona seturi de date de până la câteva mii de linii (implicit 5.000), generând ID-uri unice, proprietăți fizice editabile (volum, greutate) și atribute de localizare, oferind suport ideal pentru aplicații cu arhitectură Offline-First.

---

## 🏗️ Structura Proiectului
```text
Proiectul este împărțit în module Go independente (cmd/) și un fișier comun de configurare:

csv2api/
├── config.yaml          # Configurarea centrală (Sursa unică de adevăr)
├── config.go            # Structura de date mapată peste YAML
├── go.mod               # Definirea modulului Go
├── go.sum               # Suma de control a dependențelor
└── cmd/
    ├── generator/       # Serviciul de generare fișier CSV aleatoriu
    ├── migrator/        # Serviciul de import ultra-rapid în SQLite3
    └── server/          # Serverul HTTP API (Livrează JSON din SQLite)

---
```
## ⚙️ Parametri Editabili (config.yaml)

Orice modificare asupra regulilor de generare, structurii bazei de date sau porturilor de rețea se face exclusiv din acest fișier, fără a fi necesară modificarea sau recompilarea codului sursă:

```text
generator:
  numar_linii: 5000         # Numărul de înregistrări generate în CSV
  lungime_id: 17            # Lungimea string-ului alfanumeric pentru ID
  cale_iesire_csv: "cars.csv"

migrare:
  cale_intrare_csv: "cars.csv"
  cale_baza_date: "cars.db"
  nume_tabel: "masini"      # Numele tabelului creat în SQLite

server:
  port: ":9977"             # Portul pe care va asculta serverul API
  api_endpoint: "/"         # Endpoint-ul principal al API-ului

atribute:                   # Nomenclatoare folosite pentru generarea aleatorie
  marci: [ "Dacia", "Trabant", "Aro", "DAC", "Roman" ]
  culori: [ "Alb", "Negru", "Rosu", "Albastru", "Verde", "Galben", "Gri" ]
  locatie_blocuri: [ "Zona A1", "Bloc Nord", "Terminal Sud", "Hangar 2" ]
  greutate_interval:
    min: 600
    max: 3500
  marime_interval:
    min: 3
    max: 15
```
---

## 🚀 Ghid de Compilare și Rulare Locală

Toate comenzile se execută din rădăcina proiectului (csv2api/).

### 1. Descărcarea dependențelor
```bash
go get gopkg.in/yaml.v3
go get github.com/glebarez/go-sqlite
go mod tidy
```
### 2. Executarea directă (Mod Dezvoltare)
#### generarea documentului cars.csv
```bash
go run cmd/generator/generator.go config.go
```
#### Migrarerea documentului csv la sqlite3
```bash
go run cmd/migrator/migrator.go config.go
```
#### Rulare server port:9977 
```bash
go run cmd/server/server.go config.go
```
### 3. Compilarea nativă în binare independente
```bash
go build -o bin/generator ./cmd/generator
go build -o bin/migrator ./cmd/migrator
go build -o bin/server ./cmd/server
```
Rularea binarelor se face în ordine:
```bash
./bin/generator && ./bin/migrator && ./bin/server
```
---

## 🐳 Rulare cu Docker, Podman sau CLI alternative

Deoarece proiectul utilizează driverul SQLite scris în pure Go, binarele nu depind de CGO. Acest lucru face imaginile de container extrem de mici și portabile.

### Fișierul Dockerfile de referință
Plasați acest text într-un fișier numit Dockerfile în rădăcină:
```text
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o server ./cmd/server

FROM scratch
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/config.yaml .
COPY --from=builder /app/cars.db . 
EXPOSE 9977
CMD ["./server"]
```

### Executare cu Docker

```bash
docker build -t csv2api-server .
docker run -d -p 9977:9977 --name api-engine csv2api-server
```

### Executare cu Podman (Rootless)

```bash
podman build -t csv2api-server .
podman run -d -p 9977:9977 --name api-engine csv2api-server
```
---

## 🗺️ Ghid de Integrare pe Platforme Target (Cross-Compilation)

Compilarea direct de pe mașina curentă de dezvoltare pentru alte medii:

### 1. Linux (x86_64 Server)
```bash
env GOOS=linux GOARCH=amd64 go build -o bin/server-linux ./cmd/server
```
### 2. Raspberry Pi (RPI 3, 4, 5, Zero 2 W - 64-bit)
```bash
env GOOS=linux GOARCH=arm64 go build -o bin/server-rpi ./cmd/server
```
### 3. FreeBSD (Sisteme Unix / FreeBSD 15.0+)
```bash
env GOOS=freebsd GOARCH=amd64 go build -o bin/server-freebsd ./cmd/server
```
Execuție persistentă în fundal pe FreeBSD:
```bash
chmod +x bin/server-freebsd
daemon -f -p /var/run/csv2api.pid ./bin/server-freebsd
```
---

## 🛠️ Administrare rapidă prin terminal (Bash)

Interogarea bazei de date din linia de comandă:
```bash
sqlite3 cars.db ".schema"
sqlite3 cars.db "SELECT count(*) FROM masini;"
sqlite3 -box cars.db "SELECT * FROM masini LIMIT 5;"
```
---

## 🐙 Ghid Git (Excludere fișiere mari)

Fișierul .gitignore pre-configurat din proiect:
```text
*.csv
*.db
bin/
```
Pentru a trimite modificările curat în repository-ul de pe GitHub:
```bash
git status
git add config.yaml config.go cmd/ go.mod go.sum README.md
git commit -m "Actualizare nomenclatoare în YAML și completare documentație tehnică"
git push origin master
```