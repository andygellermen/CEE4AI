# CEE4AI – Go-Projektstart

## Ziel des MVP

Der erste lauffähige Stand soll genau das leisten:

1. Session starten
2. Eine noch nicht beantwortete Frage liefern
3. Antwort speichern
4. Denktypen-Scores berechnen
5. Ein erstes Profil zurückgeben

---

## Empfohlene Projektstruktur

```text
cee4ai/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── app/
│   │   └── app.go
│   ├── config/
│   │   └── config.go
│   ├── db/
│   │   ├── sqlite.go
│   │   └── migrations.go
│   ├── http/
│   │   ├── handlers.go
│   │   ├── router.go
│   │   └── dto.go
│   ├── questions/
│   │   ├── model.go
│   │   ├── service.go
│   │   └── repository.go
│   ├── sessions/
│   │   ├── model.go
│   │   ├── service.go
│   │   └── repository.go
│   ├── scoring/
│   │   ├── model.go
│   │   └── service.go
│   └── sheets/
│       └── importer.go
├── web/
│   └── ui/
├── data/
│   └── cee4ai.db
├── .env.example
├── go.mod
├── Makefile
└── README.md
```

---

## Warum diese Struktur sinnvoll ist

### `cmd/api`

Startpunkt deiner API.

### `internal/config`

Saubere Konfiguration über Umgebungsvariablen.

### `internal/db`

SQLite-Verbindung und Migrationen an einer Stelle.

### `internal/questions`

Fragen laden, filtern, vorbereiten.

### `internal/sessions`

Session anlegen, Fortschritt speichern, doppelte Fragen vermeiden.

### `internal/scoring`

Denktypen-Scores und Ergebnisprofil berechnen.

### `internal/sheets`

Späterer Import aus Google Sheets oder CSV-Dateien.

### `web/ui`

Moderne Oberfläche, getrennt vom Backend.

---

## Minimaler Startplan

### Phase 1 – Backend-Grundlagen

* Go-Modul initialisieren
* HTTP-Server starten
* SQLite anbinden
* Migrationen ausführen

### Phase 2 – Fachlogik

* Session erstellen
* Frage laden
* Antwort speichern
* Scoring berechnen

### Phase 3 – Präsentation

* einfache UI
* Fortschrittsanzeige
* Ergebnisprofil

---

## Empfohlene erste Commands

```bash
git init
go mod init github.com/andygellermen/cee4ai
mkdir -p cmd/api internal/{app,config,db,http,questions,sessions,scoring,sheets} web/ui data
```

---

## Empfohlene Dependencies

### Pflicht

```bash
go get github.com/mattn/go-sqlite3
go get github.com/google/uuid
```

### Optional, aber praktisch

```bash
go get github.com/joho/godotenv
```

---

## API-Endpunkte für den MVP

### `POST /api/v1/sessions`

Legt eine neue Session an.

### `GET /api/v1/questions/next?session_id=...`

Liefert die nächste passende Frage, die in dieser Session noch nicht beantwortet wurde.

### `POST /api/v1/answers`

Speichert die Antwort des Probanden.

### `GET /api/v1/results?session_id=...`

Liefert Denktypen-Scores und ein erstes Profil.

---

## Fachliche Leitidee für die erste Version

Du baust noch **kein fertiges Intelligenz-Entwicklungsuniversum**.
Du baust zuerst den **stabilen Kern**:

* Fragenfluss
* Sessiontracking
* Bewertungslogik
* Profil-Rückgabe

Alles Weitere wächst darauf evolutionär weiter.

---

## Empfehlung für deine Google-Sheets-Struktur

Ja: **bitte trennen**.

### Blatt 1: `questions`

Enthält alle Fragen.

### Blatt 2: `categories`

Enthält Kategorien und optionale Beschreibungen.

### Blatt 3: `think_types`

Enthält Denktypen, Definitionen, Entwicklungsbeschreibung.

### Blatt 4: `levels`

Enthält Level 1 bis 5 und deren sprachliche Deutung.

### Blatt 5: `weights`

Optionale Feingewichte für Schwierigkeit, Kategorie, Denktyp.

---

## Warum die Trennung wichtig ist

### Vorteile

* konsistenter Datenbestand
* leichter pflegbar
* später besser validierbar
* ideal für Importlogik in Go

### Vermeidungsstrategie

Bitte die Labels nicht direkt an zig Stellen frei eintippen.
Sonst entstehen schnell Varianten wie:

* Spiritualität
* spirituell
* Spirituell
* Spiritualitaet

Sauberer ist eine Referenzstruktur.

---

## Nächste sinnvolle Arbeitsreihenfolge

1. Google Sheet aufsetzen
2. Repo-Struktur erzeugen
3. SQLite-Schema anlegen
4. API-Endpunkte implementieren
5. Fragenimport vorbereiten
6. UI anbinden

---

## Checkliste

### Erledigt

* Vision geschärft
* Denktypen definiert
* erste 30 Fragen vorbereitet
* Zielbild für den MVP gesetzt

### Offen

* Google-Sheets-Blätter anlegen
* Repo initialisieren
* SQLite-Schema schreiben
* erste Go-Dateien anlegen
* Endpunkte implementieren
* UI verbinden

---

## Zwei alternative Denkansätze

### Ansatz A – pragmatisch schnell

CSV zuerst lokal einlesen, Google Sheets erst im zweiten Schritt.

### Ansatz B – direkt sauber

Google Sheets sofort als pflegbare Master-Datenquelle nutzen.

Meine Empfehlung für dich:
**Starte mit Google Sheets als Master, importiere aber technisch zunächst über CSV-Export.**
Das ist robust, nachvollziehbar und für dein MVP sehr vernünftig.
