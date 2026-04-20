# analytics-plane-v3-final.md

## Zweck
Dieses Dokument modelliert die Analytics Plane der CPE-Produktfamilie in der V3-Endfassung.

Sie ist die Ebene, auf der aus großen Mengen von Sessions, Antworten, Zeitmustern, Abbrüchen, Reviews, Lokalisierungen, Meaning-Journeys und Zertifikaten **echte Musterintelligenz** entsteht.

---

## Leitprinzip
Die Analytics Plane ist **ergänzend** zur operativen Produktdatenbank zu verstehen.

### Nicht ihre Aufgabe
- Primärwahrheit verwalten
- operative Geschäftsregeln erzwingen
- Zertifikate ausstellen

### Ihre Aufgabe
- Muster erkennen
- Verdichtungen berechnen
- Confidence-Modelle verbessern
- Paketlogik optimieren
- Benchmarking erzeugen
- Qualitäts- und Lokalisierungsentscheidungen unterstützen
- Meaning-/Journey-Verläufe analysieren

---

## Architekturhaltung: bewusst offen
Die Analytics Plane soll logisch früh angelegt, technologisch aber bewusst offen gehalten werden.

### Warum?
Weil der spätere Bedarf offen ist:
- einfache SQL-Aggregationen
- Eventauswertungen
- Segmentierungen
- Recommender-/Heuristik-Modelle
- ML-/AI-Modelle
- regionale Datenprodukte
- Journey- und Meaning-Musterbildung

---

## Empfohlene Schichten der Analytics Plane
### 1. Raw Event Layer
Unveränderte oder nahezu unveränderte Ereignisse aus der operativen Ebene.

### 2. Curated / Modeled Layer
Bereinigte, normalisierte, joinbare Analyseansichten.

### 3. Feature Layer
Abgeleitete Kennzahlen und Merkmale für:
- Confidence
- Profilstabilität
- Abbruchwahrscheinlichkeit
- Paketoptimierung
- Zertifikatswürdigkeit
- Meaning-/Journey-Tiefe

### 4. Insight / Product Layer
Dashboards, Reports, Benchmarks, Empfehlungen und Modelle.

---

## Datenquellen
### Operative Hauptquellen
- Sessions
- Session Packages
- Answers
- Result Snapshots
- Reviews
- Localizations
- Credential Assertions
- Verifications

### Eventquellen
- Session-Events
- Package-Events
- Question-Events
- Answer-Events
- Review-Events
- Localization-Events
- Credential-Events
- Meaning-/Journey-Events

---

## Analytische Kernfragen
### A. Profil- und Produktfragen
- Welche frühen Pakete haben hohe Vorhersagekraft?
- Welche Fragefamilien erzeugen stabile Profile?
- Welche Traits oder Skills korrelieren mit welchen Denktypen?

### B. Experience-Fragen
- Wo brechen Nutzer ab?
- Welche Sprachen oder Regionen zeigen besondere Hürden?
- Welche Fragepakete sind zu lang, zu schwer oder unklar?

### C. Qualitätsfragen
- Welche Inhalte werden häufig geflaggt?
- Welche Übersetzungen erzeugen Verständnisprobleme?
- Welche sensiblen Fragen brauchen stärkere Prüfung?

### D. Trust- und Zertifikatsfragen
- Welche Pfade führen zu verifizierbaren Zertifikaten?
- Welche Integritätsmuster unterscheiden belastbare von schwächeren Durchläufen?
- Welche Zertifikatsstufen werden tatsächlich genutzt oder verifiziert?

### E. Markt- und Business-Fragen
- Welche Produkte, Tracks oder Sprachräume konvertieren?
- Welche Pathways werden häufig vertieft?
- Welche Module stiften den höchsten Folgewert?

### F. Meaning-/Journey-Fragen
- Welche Meaning-/Spiritual-Pfade werden freiwillig vertieft?
- Welche Fragefamilien erzeugen hohe Reflexionsbereitschaft?
- Welche Regionen oder Sprachen reagieren sensibler auf bestimmte Weltdeutungsfragen?
- Welche Reflexionspfade fördern stabile Rückkehr und vertiefte Sessionverläufe?

---

## Offene Technologieoptionen
### Option A – SQL-first / Postgres-first
Start mit analytischen Views, Materialized Views und Exporten direkt aus PostgreSQL.

### Option B – Separates OLAP-System
Spätere Ausleitung in ein analytisches System.

### Option C – Eigene Entwicklungsroute
Eigene Analytics-Schicht mit:
- eventbasierten Exporten
- eigenen Snapshots
- eigenen Modellierungsjobs
- domänenspezifischen Confidence- und Benchmarking-Komponenten

### Strategische Empfehlung
Kurzfristig Postgres-first, mittelfristig logisch getrennte Analytics Plane, langfristig offene Entscheidung zwischen spezialisierten Tools und eigener Weiterentwicklung.

---

## Feature-Ideen für die spätere Musterbildung
- Profilstabilität
- Confidence Engine
- Early Prediction
- Drift Detection
- Certificate Readiness
- Localization Quality Signals
- Meaning Journey Depth Signals
- Worldview Sensitivity Patterns

---

## Datenprodukte der Analytics Plane
### Interne Datenprodukte
- Qualitätsdashboard
- Confidence Dashboard
- Profil- und Pfadbenchmarking
- Zertifikats- und Verifikationsmonitor
- Localization Health Monitor
- Meaning / Journey Resonance Dashboard

### Externe / monetarisierbare Datenprodukte
- Benchmark Reports
- Markt- und Sprachraumanalysen
- B2B-Insights
- Coaching-Qualitätsreports
- anonymisierte Branchenvergleiche

---

## Datenschutz und Trennung
Die Analytics Plane darf nie einfach die gesamte operative Wahrheit 1:1 spiegeln.

### Grundsätze
- Pseudonymisierung, wo möglich
- sensible Freitexte vorsichtig behandeln
- Coaching- und Körperdaten besonders schützen
- Meaning-/Worldview-Daten nicht für degradierende oder vereinfachende Rankings missbrauchen
- Zertifikats- und Identitätsdaten von breiten Produktanalysen trennen

---

## MVP für die Analytics Plane
### Reicht für den Start
- Event-Schema definieren
- Event-Outbox in Postgres
- erste materialisierte Sicht auf Sessions, Packages, Answers, Reviews
- einfache KPI-Dashboards
- Confidence-Grundmetriken
- erste Journey-/Meaning-Aggregate

### Noch nicht nötig
- komplexe Echtzeitarchitektur
- globale Stream-Plattform
- massive verteilte Datenverarbeitung

---

## Abschluss
Die Analytics Plane V3 Final ist der Ort, an dem CPE von einem starken Profil-Tool zu einer lernenden Muster- und Entscheidungsplattform wird.  
Sie bleibt technisch offen, aber logisch klar genug, um Profiling, Zertifikate, Mehrsprachigkeit und Meaning-/Journey-Pfade gemeinsam auswertbar zu machen.

