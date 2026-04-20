# repo-implementation-priorities-v3.md

## Zweck
Dieses Dokument leitet aus den kanonischen V3-Unterlagen eine konkrete Umsetzungs-Prioritaetenliste fuer das **aktuelle Repo** ab.

Es beantwortet:
- Was jetzt im bestehenden Projekt als Erstes gebaut werden sollte
- Welche V3-Bausteine fuer das Repo Pflicht sind
- Was bewusst noch nicht priorisiert werden sollte
- Woran wir erkennen, dass ein Ausbauschritt wirklich abgeschlossen ist

---

## Repo-Ausgangslage
Der aktuelle Stand des Repos ist bewusst frueh:

- `docs/v3/` enthaelt den neuen kanonischen Zielzustand
- `docs/v2/` enthaelt alte, fachlich nicht mehr leitende Planungsunterlagen
- `cmd/`, `internal/`, `web/` und `data/` existieren bereits als Struktur
- es gibt derzeit aber noch **keine tragende Go-Implementierung**
- `README.md` ist leer
- [cee4ai-Go-Projektstart.md](/Users/andygellermann/Documents/Projects/CEE4AI/docs/cee4ai-Go-Projektstart.md:1) spiegelt noch eine fruehere V2-/SQLite-Denke und sollte nicht mehr die weitere Umsetzung fuehren

Die wichtigste Konsequenz daraus:

> Dieses Repo braucht zuerst einen sauberen **V3-faehigen technischen Kern**.  
> Nicht die naechste Feature-Idee ist die erste Aufgabe, sondern die richtige technische und fachliche Basis.

---

## Kanonische V3-Basis
Diese Dokumente sollten die Umsetzung im Repo fuehren:

- [cee_4_ai_master_architecture_v_3_final.md](/Users/andygellermann/Documents/Projects/CEE4AI/docs/v3/cee_4_ai_master_architecture_v_3_final.md:21)
- [roadmap_v_3_final.md](/Users/andygellermann/Documents/Projects/CEE4AI/docs/v3/roadmap_v_3_final.md:39)
- [data_model_v_3_final.md](/Users/andygellermann/Documents/Projects/CEE4AI/docs/v3/data_model_v_3_final.md:18)
- [postgresql_target_schema_v_3_final.md](/Users/andygellermann/Documents/Projects/CEE4AI/docs/v3/postgresql_target_schema_v_3_final.md:15)
- [postgresql_ddl_and_go_repo_v_3_final.md](/Users/andygellermann/Documents/Projects/CEE4AI/docs/v3/postgresql_ddl_and_go_repo_v_3_final.md:804)
- [governance_model_v_3_final.md](/Users/andygellermann/Documents/Projects/CEE4AI/docs/v3/governance_model_v_3_final.md:17)
- [event_taxonomy_v_3_final.md](/Users/andygellermann/Documents/Projects/CEE4AI/docs/v3/event_taxonomy_v_3_final.md:26)
- [credential_trust_ladder_v_3_final.md](/Users/andygellermann/Documents/Projects/CEE4AI/docs/v3/credential_trust_ladder_v_3_final.md:27)
- [meaning_spirituality_integration_v_3.md](/Users/andygellermann/Documents/Projects/CEE4AI/docs/v3/meaning_spirituality_integration_v_3.md:1)
- [meaning_taxonomy_and_texts_v_3_final.md](/Users/andygellermann/Documents/Projects/CEE4AI/docs/v3/meaning_taxonomy_and_texts_v_3_final.md:19)

---

## Strategischer Grundsatz fuer dieses Repo
Die Reihenfolge fuer das Repo sollte der V3-Leitidee folgen:

1. Wahrheitsfaehigen Kern bauen
2. Runtime und Snapshot-/Progression-Logik aufsetzen
3. Governance und Review absichern
4. Meaning-/Spiritual-Basis kontrolliert integrieren
5. Trust-, Credential- und Event-Schicht aufsetzen
6. Erst danach staerker in Expansion, Analytics-Produkte und groessere UI-Ausbaustufen gehen

---

## Prioritaet 1 - V3-Repo-Baseline herstellen
### Ziel
Das Repo muss von einem leeren Geruest zu einem **startfaehigen Postgres-first V3-Kern** werden.

### Konkrete Arbeitspakete
- `go.mod` anlegen und auf eine stabile Modulstruktur festziehen
- `.env.example` anlegen
- `Makefile` fuer `run`, `migrate`, spaeter `seed`
- `cmd/api/main.go` anlegen
- `cmd/migrate/main.go` anlegen
- `internal/config/config.go` anlegen
- `internal/db/postgres.go` anlegen
- `internal/db/migrations.go` anlegen
- `internal/app/app.go` anlegen
- `internal/http/router.go` mit `GET /healthz`
- `migrations/` mit der V3-Sequenz `0001` bis `0009` anlegen

### Repo-Entscheidung
- **PostgreSQL-first**
- **pgx/v5**
- keine neue SQLite-Hauptlinie mehr aufbauen

### Definition of Done
- `go run ./cmd/api` startet lokal
- `go run ./cmd/migrate` laeuft gegen eine lokale Postgres-Datenbank
- die acht logischen Schemas aus V3 sind technisch vorhanden
- das Repo hat einen ersten lauffaehigen technischen Backbone

### Warum das zuerst kommt
Ohne diesen Schritt bleibt jede weitere Fachlogik an ein nicht belastbares Grundgeruest gekoppelt.

---

## Prioritaet 2 - Truth Core, Content und Profiling-Grundlage bauen
### Ziel
Die fachliche Wahrheit des Systems muss im Repo zuerst sauber modelliert und importierbar werden.

### Konkrete Arbeitspakete
- Tabellen fuer `core`, `content` und `profiling` vollstaendig migrieren
- Seed-Basis fuer `languages`, `regions`, `domains`, `roles`
- `questions/`, `profiling/`, `localization/` und `importer/` als echte Packages anlegen
- CSV-/Seed-Import fuer:
  - `categories`
  - `subcategories`
  - `question_master`
  - `question_translation`
  - `question_option_master`
  - `question_option_translation`
  - erste Profiling-Referenzen wie `denktypes`, `skills`, `personality_traits`, `meaning_tags`, `worldview_frames`
- V3-Felder direkt mitfuehren:
  - `question_family`
  - `intended_use`
  - `confidence_tier`
  - `cognitive_load_level`
  - `meaning_depth`
  - `worldview_sensitivity`

### Definition of Done
- mindestens ein kleiner seedbarer Content-Bestand ist in Postgres ladbar
- Inhalte sind sprachunabhaengig als Master-Objekte modelliert
- Uebersetzungen sind getrennt gespeichert
- Frage-zu-Tag-Mappings koennen importiert werden

### Warum das jetzt priorisiert ist
Ohne Truth Core gibt es keine sinnvolle Sessionlogik, keine belastbaren Ergebnisse und keine spaetere Governance.

---

## Prioritaet 3 - Runtime-Kern und erste Produktlogik
### Ziel
Aus dem Truth Core muss eine erste lauffaehige V3-Produktstrecke werden.

### Konkrete Arbeitspakete
- `sessions/`, `packages/`, `answers/`, `results/`, `scoring/` als echte Runtime-Pakete anlegen
- Tabellen fuer:
  - `runtime.sessions`
  - `runtime.session_packages`
  - `runtime.answers`
  - `runtime.result_snapshots`
  - `runtime.profile_vectors`
- erste API-Strecke fuer:
  - Session starten
  - naechstes Paket / naechste Frage liefern
  - Antwort speichern
  - Snapshot-Ergebnis abrufen
- zuerst nur zwei Modi produktiv denken:
  - `snapshot`
  - `guided_progression`
- Paketlogik aus V3 direkt beachten:
  - kleine Pakete
  - fruehe Vorprofile
  - `result_confidence`
  - `continuation_window_until`

### Definition of Done
- eine Session kann gestartet werden
- Fragen koennen paketweise ausgespielt werden
- Antworten werden sauber gespeichert
- ein erstes Snapshot-Profil kann erzeugt werden
- der Grundpfad ist V3-kompatibel und nicht mehr V2-testartig

### Wichtige Leitplanke
Nicht zuerst komplexe Persona-Texte oder grosse UI bauen. Erst die saubere Laufstrecke.

---

## Prioritaet 4 - Review, Sensitivity und Governance technisch verankern
### Ziel
Das Repo muss frueh governance-faehig werden, weil V3 sonst spaeter teuer umbauen wuerde.

### Konkrete Arbeitspakete
- `reviews/` und `governance/` als Packages anlegen
- Review-Tabellen und Audit-Log migrieren
- `ruleset_versions` nutzen, nicht nur definieren
- Rollenmodell im Seed aufnehmen:
  - `editor`
  - `reviewer`
  - `coach_reviewer`
  - `localization_reviewer`
  - `certification_admin`
  - `governance_admin`
- Workflows fuer:
  - Review-Flags
  - Statuswechsel
  - Localization Review
  - dokumentierte Freigaben
- sensible V3-Marker in Logik und APIs respektieren:
  - `is_sensitive`
  - `age_gate`
  - `worldview_sensitivity`
  - `requires_human_review`

### Definition of Done
- sensible Inhalte koennen technisch nicht mehr wie gewoehnliche Fragen behandelt werden
- Statuswechsel sind nachvollziehbar
- Audit-Eintraege entstehen bei relevanten Aenderungen
- die Basis fuer spaetere Meaning- und Zertifikatsgovernance ist gelegt

### Warum das nicht spaeter warten sollte
V3 macht Governance zur Voraussetzung von Glaubwuerdigkeit, nicht zum Admin-Nachtrag.

---

## Prioritaet 5 - Meaning / Spiritual Depth kontrolliert integrieren
### Ziel
Die neue V3-Kernverschiebung soll frueh vorbereitet, aber bewusst kontrolliert ausgerollt werden.

### Konkrete Arbeitspakete
- `internal/domains/meaning/` anlegen
- nur die **P1-Taxonomie** zuerst seedbar machen
- `meaning_tags` und `worldview_frames` aktiv nutzen
- neue Fragefamilien aufnehmbar machen:
  - `reflective`
  - `contemplative`
  - `comparative_worldview`
  - `symbolic_interpretation`
  - `existential`
- Ergebnistext-Basis fuer:
  - `meaning_journey_snapshot`
  - `worldview_reflection_result`
- Schutztexte und vorsichtige Ergebnislogik direkt mitdenken
- Human Review fuer worldview-/meaning-sensitive Inhalte technisch erzwingbar machen

### Definition of Done
- Meaning-/Spiritual-Inhalte sind als echte Domaene anschlussfaehig
- sie sind weder bloesser Freitext-Anhang noch dogmatische Sonderlogik
- Schutzsaetze und vorsichtige Ergebnislogik sind technisch vorbereitbar

### Strategische Empfehlung
Die kontrolliert gestufte Linie aus `meaning_taxonomy_and_texts_v_3_final.md` sollte gelten:

> erst sauber gefuehrte Kernintegration, nicht sofort maximale Durchdringung aller Produktteile

---

## Prioritaet 6 - Credential- und Trust-Layer als Repo-Baustein aufsetzen
### Ziel
Das Repo soll nicht erst spaeter kuenstlich zertifikatsfaehig gemacht werden.

### Konkrete Arbeitspakete
- `credentials/` als echtes Package anlegen
- Tabellen fuer:
  - `assessment_programs`
  - `assessment_tracks`
  - `certificate_rules`
  - `credential_issuers`
  - `credential_assertions`
  - `credential_evidence`
  - `credential_verifications`
  - `credential_revocations`
- Issuance- und Verification-Grundlogik vorbereiten
- `ruleset_version`, `evidence_hash`, `verification_token` ab erstem Credential mitfuehren
- anfangs nur niedrige bis mittlere Trust-Stufen wirklich aktiv ausgeben:
  - Participation
  - Snapshot
  - Assessed
  - spaeter Extended

### Definition of Done
- das Repo kann technisch ein regelbasiertes Credential erzeugen
- Regelversion und Evidenz sind rueckverfolgbar
- Widerruf und Verifikation sind modelliert

### Wichtige Leitplanke
Meaning-/Journey-Credentials nur als Completion-/Pathway-Nachweise, nicht als Ueberhoehungslogik.

---

## Prioritaet 7 - Event-Outbox und Analytics-MVP vorbereiten
### Ziel
V3 will Musterbildung frueh logisch vorbereiten, aber technologisch bewusst offen halten.

### Konkrete Arbeitspakete
- `events/` als Package anlegen
- `event_outbox`, `event_inbox`, `event_store` migrieren
- Prioritaet-A-Events zunaechst emittieren:
  - `session.started`
  - `package.started`
  - `question.presented`
  - `answer.submitted`
  - `answer.evaluated`
  - `profile.snapshot_created`
  - `review.flag_created`
  - `credential.issued`
- Event-Erzeugung in Runtime und Review nicht vergessen
- erste SQL-Sichten oder einfache Metriken fuer:
  - Session-Verlauf
  - Abbrueche
  - Snapshot-Haeufigkeit
  - Review-Signale

### Definition of Done
- zentrale Produktablaeufe erzeugen Events
- die Outbox ist lauffaehig
- erste operative KPIs koennen aus Postgres gelesen werden

### Wichtige Leitplanke
Noch keine schwere externe Analytics-Plattform bauen. Erst Postgres-first.

---

## Prioritaet 8 - Repo-Hygiene und Arbeitsfaehigkeit
### Ziel
Das Repo soll nicht nur inhaltlich, sondern auch als Arbeitsgrundlage sauber werden.

### Konkrete Arbeitspakete
- `README.md` mit V3-Realitaet fuellen
- lokale Startanleitung fuer Postgres und Migrationen dokumentieren
- `docs/cee4ai-Go-Projektstart.md` als historisch markieren oder spaeter durch V3-Startdokument ersetzen
- `.gitignore` pruefen
- Seed- und Migrationsablauf dokumentieren
- klare Benennung fuer API-, Domain- und Package-Begriffe festziehen

### Definition of Done
- ein neuer Entwickler kann das Repo lokal starten
- die V3-Zielrichtung ist im Repo selbst lesbar
- alte V2-Leitbilder fuehren die Implementierung nicht mehr versehentlich

---

## Nicht jetzt priorisieren
Diese Themen sind in V3 wichtig, aber fuer das aktuelle Repo **noch nicht zuerst** dran:

- grosse Frontend-Ausgestaltung
- White-Label
- Multi-Tenant
- Partner- und Lobby-Oekonomie
- externe OLAP-/Streaming-Infrastruktur
- breite internationale Rollouts
- komplexe Coach-Dashboards
- hochstufige Verified- oder Partner-Endorsed-Credentials
- spirituelle Feindifferenzierung jenseits der P1-/P2-Kernlogik

---

## Empfohlene Reihenfolge fuer die naechsten Umsetzungswellen
### Welle 1
Repo-Baseline, Postgres-Anbindung, Migrationslauf

### Welle 2
Truth Core, Content-Import, Profiling-Referenzen

### Welle 3
Session-, Package-, Answer- und Snapshot-Kern

### Welle 4
Review, Governance, Audit und Localization Review

### Welle 5
kontrollierte Meaning-Basis

### Welle 6
Credentials, Events und erste operative Analytics-Sichten

---

## Praktische naechste 10 Tasks fuer dieses Repo
1. `go.mod` und `.env.example` anlegen
2. `cmd/api` und `cmd/migrate` lauffaehig machen
3. `internal/config`, `internal/db`, `internal/app`, `internal/http` anlegen
4. `migrations/0001` bis `0009` aus V3 ins Repo uebernehmen
5. lokale Postgres-Datenbank anbinden und Migrationslauf pruefen
6. Seeds fuer `languages`, `regions`, `domains`, `roles` anlegen
7. `questions`-, `profiling`- und `importer`-Pakete anlegen
8. `runtime.sessions`, `runtime.session_packages`, `runtime.answers` produktiv anbinden
9. ersten Snapshot-Flow per API herstellen
10. Review-/Governance-Grundpfad technisch verankern

---

## Abschluss
Fuer dieses Repo ist die wichtigste Einsicht:

CEE4AI braucht jetzt **keine lose Sammlung naechster Features**, sondern eine entschlossene V3-Umbasis aus Postgres, Truth Core, Runtime, Governance, Meaning-Basis, Events und Trust.

Wenn diese Reihenfolge eingehalten wird, waechst das Projekt nicht chaotisch, sondern in genau der Ordnung, die die V3-Dokumente nun als finalen Kern vorgeben.
