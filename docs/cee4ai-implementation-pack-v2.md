cee4ai-implementation-pack-v2.md
Zweck
Dieses Dokument schließt die aktuell offenen Punkte der V2-Checkliste für CEE4AI und bildet den Übergang von Vision und Datenmodell zu einer belastbaren Umsetzungsstruktur.
Es bündelt:
	•	die finale CSV-/Google-Sheets-Struktur V2
	•	die echte Go-Repo-Struktur V2
	•	den technischen Review-Workflow
	•	die Bewertungslogik für multiple_select, scale und reflection
	•	Schutzbausteine für sensible Domänen

1. Finale CSV-/Google-Sheets-Struktur V2
Leitprinzip
Google Sheets dient als redaktionelle Master-Datenquelle.Die Anwendung selbst arbeitet nicht direkt auf frei geschriebenen Labels, sondern auf stabilen IDs, Slugs und Referenzen.
Das schützt uns vor:
	•	Schreibfehlern
	•	inkonsistenten Kategorienamen
	•	unstimmigen Imports
	•	schwer nachvollziehbaren Datenfehlern

Empfohlene Tabellenblätter
Blatt 1 – domains
id
slug
name
description
is_active
Blatt 2 – categories
id
domain_id
slug
name
description
is_sensitive
is_active
Blatt 3 – subcategories
id
category_id
slug
name
description
is_sensitive
is_active
Blatt 4 – denktypes
id
slug
name
description
development_hint
is_active
Blatt 5 – trigger_groups
id
domain_id
slug
name
description
is_sensitive
is_active
Blatt 6 – body_references
id
slug
name
description
is_sensitive
is_active
Blatt 7 – questions
id
external_id
domain_id
category_id
subcategory_id
parent_id
question_type
scoring_mode
denktyp_id
level
difficulty
trigger_group_id
body_reference_id
question
explanation
weight
review_status
is_sensitive
age_gate
is_active
created_by
notes
Blatt 8 – question_options
id
question_id
option_key
option_label
is_correct
score_weight
display_order
is_active
Blatt 9 – review_flags
id
slug
name
description
default_severity
is_active
Blatt 10 – result_texts
id
domain_id
result_type
slug
title
text
hint_level
is_sensitive
is_active

Warum question_options als separates Blatt?
Weil wir uns damit bewusst von der starren Vier-Spalten-Struktur lösen.
Vorteile
	•	mehrere richtige Antworten sind sauber abbildbar
	•	variable Anzahl an Antwortoptionen möglich
	•	Teilgewichtungen werden möglich
	•	spätere UI-Varianten werden leichter

Empfohlene ENUM-Werte für V2
question_type
single_choice
multiple_select
scale
reflection
scoring_mode
exact
partial
weighted
non_scored
path_only
review_status
draft
active
review
flagged
disabled
archived

Minimale Referenzwerte für domains
1 | cognitive_profile | Cognitive Profile | Denkprofil, Interessen und Entwicklung | true
2 | coaching_analyzer | Coaching Analyzer | Belastungs-, Druck- und Trigger-Navigation | true

Minimaler Importprozess
Schritt 1
Redaktion pflegt Google Sheets.
Schritt 2
Google Sheets wird als CSV exportiert.
Schritt 3
Ein Importer liest die CSV-Dateien und validiert:
	•	Pflichtfelder
	•	Referenzen
	•	ENUM-Werte
	•	Aktivitätsstatus
Schritt 4
Valide Datensätze werden in SQLite oder später PostgreSQL übernommen.
Vermeidungsstrategie
Bitte zunächst nicht direkt live aus Google Sheets lesen.Für den Start ist ein kontrollierter CSV-Import robuster und besser nachvollziehbar.

2. Echte Go-Repo-Struktur V2
Zielbild
Die Repo-Struktur soll den Engine-Gedanken tragen und spätere Domänen nicht blockieren.
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
│   │   ├── migrations.go
│   │   └── seed.go
│   ├── domains/
│   │   ├── cognitive/
│   │   │   ├── service.go
│   │   │   └── result_texts.go
│   │   └── coaching/
│   │       ├── service.go
│   │       └── safeguards.go
│   ├── questions/
│   │   ├── model.go
│   │   ├── repository.go
│   │   ├── service.go
│   │   ├── options.go
│   │   └── validation.go
│   ├── sessions/
│   │   ├── model.go
│   │   ├── repository.go
│   │   └── service.go
│   ├── answers/
│   │   ├── model.go
│   │   ├── repository.go
│   │   └── service.go
│   ├── scoring/
│   │   ├── model.go
│   │   ├── service.go
│   │   ├── exact.go
│   │   ├── partial.go
│   │   ├── weighted.go
│   │   ├── scale.go
│   │   └── reflection.go
│   ├── reviews/
│   │   ├── model.go
│   │   ├── repository.go
│   │   └── service.go
│   ├── importer/
│   │   ├── csv.go
│   │   ├── mapper.go
│   │   └── validator.go
│   ├── results/
│   │   ├── model.go
│   │   └── service.go
│   └── http/
│       ├── dto.go
│       ├── router.go
│       ├── sessions_handler.go
│       ├── questions_handler.go
│       ├── answers_handler.go
│       ├── reviews_handler.go
│       └── results_handler.go
├── data/
│   ├── imports/
│   └── cee4ai_v2.db
├── web/
│   └── ui/
├── .env.example
├── Makefile
├── README.md
└── go.mod

Verantwortlichkeiten je Paket
domains/
Domänenspezifische Regeln und Ergebnislogik.
questions/
Frageobjekte, Optionen, Validierung, Auswahlregeln.
answers/
Speichern und Lesen von Antworten unabhängig von der Auswertung.
scoring/
Bewertungslogik je scoring_mode.
reviews/
Test-Flags, Kommentare, Review-Status.
results/
Zusammenführung von Scores, Hinweisen und Ergebnistexten.
importer/
Import aus CSV/Google-Sheets-Exporten.
domains/coaching/
Safeguards und spezielle Pfadlogik für das Angst-/Druck-Modul.

Warum diese Struktur sinnvoll ist
Sie trennt früh:
	•	Datenhaltung
	•	Domänenlogik
	•	Bewertungslogik
	•	Test-/Review-Logik
	•	API-Ausgabe
Damit verhindern wir, dass alles im Session- oder Questions-Service zusammenklebt.

3. Review-Workflow technisch sauber trennen
Ziel
Review und Testbetrieb dürfen nicht als Nebenbemerkung im Fragefluss untergehen. Sie brauchen einen eigenständigen Ablauf.

Technischer Review-Ablauf
Zustand 1 – draft
Neue oder noch nicht geprüfte Frage.
Zustand 2 – active
Frage ist freigegeben und regulär nutzbar.
Zustand 3 – review
Frage bleibt sichtbar, soll aber geprüft werden.
Zustand 4 – flagged
Frage zeigt deutliche Auffälligkeiten und sollte priorisiert geprüft werden.
Zustand 5 – disabled
Frage wird vorerst nicht mehr ausgeliefert.
Zustand 6 – archived
Frage bleibt historisch erhalten, wird aber nicht mehr verwendet.

Empfohlene Review-Flags
unclear_question
multiple_answers_seem_correct
no_answer_fits
too_easy
too_hard
linguistically_unclear
factually_problematic
emotionally_triggering

Entscheidungslogik für automatische Eskalation
Beispielregel 1
Wenn eine Frage in Testmodus innerhalb kurzer Zeit 3-mal mit no_answer_fits markiert wird:
	•	setze review_status auf flagged
Beispielregel 2
Wenn emotionally_triggering bei sensiblen Domänen mehrfach vorkommt:
	•	leite Frage an Coach/Admin-Review
Beispielregel 3
Wenn multiple_answers_seem_correct gehäuft auftaucht:
	•	prüfe, ob question_type auf multiple_select umgestellt werden sollte

Rollen im Review-Workflow
user
tester
coach
admin
expert_reviewer
Empfehlung
Für den MVP reichen zunächst:
	•	tester
	•	coach
	•	admin

Vermeidungsstrategie
Reviewdaten nie nur als Freitext speichern.Immer zuerst strukturierte Flags, dann optional Kommentar.

4. Bewertungslogik für multiple_select, scale und reflection
4.1 single_choice mit exact
Logik
	•	genau eine Antwort ist richtig
	•	nur exakte Übereinstimmung zählt als korrekt
Scoreformel
score = difficulty * weight
bei korrekter Antwort, sonst 0.

4.2 multiple_select mit partial
Logik
	•	mehrere Antworten können richtig sein
	•	Teilpunkte sind möglich
	•	falsche Mit-Auswahl sollte Punkte mindern oder neutralisieren
Vorschlag für eine faire MVP-Formel
score = (richtig_getroffen - falsch_mitgewählt * 0.5) / anzahl_richtiger_optionen
Dann clampen auf Minimum 0.
Beispiel
Richtige Antworten: B, DNutzer wählt: B, C, D
richtig_getroffen = 2
falsch_mitgewählt = 1
score = (2 - 0.5) / 2 = 0.75
Dann:
raw_score = 0.75 * difficulty * weight
Vermeidungsstrategie
Nicht jede falsche Auswahl maximal bestrafen. Sonst wird multiple_select unnötig frustrierend.

4.3 multiple_select mit weighted
Logik
Jede Option trägt ein eigenes Gewicht. Das ist hilfreich, wenn Antworten unterschiedlich stark zur Deutung beitragen.
Beispiel
	•	B = 1.0
	•	D = 0.5
	•	C = -0.25
Dann wird die Summe der gewählten Optionsgewichte gebildet und auf einen Mindestwert von 0 begrenzt.
Einsatzbereich
	•	feinere Denkprofilfragen
	•	Coaching-Pfade mit unterschiedlich starker Relevanz
Empfehlung
Für den frühen MVP nur sparsam einsetzen.

4.4 scale
Logik
Antwort erfolgt auf einer Skala, z. B. 1 bis 5.
Nutzungsarten
	•	Selbsteinschätzung
	•	Intensität von Erleben
	•	Häufigkeit von Zuständen
	•	Zustimmung zu Aussagen
Bewertungsansatz A – Profilwert
Skalenwert fließt direkt als Stärkeindikator in einen Denktyp oder Triggerpfad ein.
Bewertungsansatz B – nicht-diagnostische Verdichtung
Mehrere Skalenwerte werden zu Clustern aggregiert, z. B.:
	•	innere Unruhe
	•	Kontrollbedürfnis
	•	soziale Unsicherheit
MVP-Empfehlung
scale zunächst nicht als richtig/falsch, sondern als profilsignalgebend auswerten.

4.5 reflection
Logik
Keine objektiv richtige Antwort, sondern Selbstbeobachtung oder offene Reflexion.
Geeignete Speicherung
	•	free_text_answer
	•	optional Tags oder Muster später ergänzen
MVP-Empfehlung
reflection zunächst mit scoring_mode = non_scored
Nutzen
	•	wertvoll für Coaching und persönliche Erkenntnis
	•	nicht als klassische Punktequelle behandeln
Vermeidungsstrategie
Reflexionsfragen nicht künstlich numerisieren, wenn dafür noch keine belastbare Methodik existiert.

5. Schutzbausteine für sensible Domänen
Warum das wichtig ist
Nicht jede Frage darf gleich behandelt oder gleich ausgewertet werden.Sensible Domänen brauchen sprachliche, methodische und technische Schutzschichten.

5.1 Technische Marker
Jede sensible Frage kann mindestens tragen:
is_sensitive = true
age_gate = 16 oder 18
review_status = active nur nach zusätzlicher Prüfung
Zusätzlich kann domänenseitig definiert werden:
	•	nur im coaching-Mode ausspielbar
	•	nicht im offenen öffentlichen Testmodus
	•	Ergebnis nur mit Schutztexten

5.2 Sprachliche Schutzbausteine
Diese Textlogik sollte für sensible Resultate vorbereitet werden.
Schutztext A – Orientierung statt Diagnose
Diese Auswertung dient der Selbstklärung und Orientierung. Sie ersetzt keine medizinische, psychotherapeutische oder sonstige fachliche Diagnose.
Schutztext B – Beobachtung statt Festlegung
Die erkannten Muster sind Hinweise auf mögliche Belastungs- oder Themenräume, keine endgültige Aussage über Ursachen oder Störungen.
Schutztext C – freiwillige Vertiefung
Bitte gehe nur so weit in die Reflexion, wie es sich für dich aktuell sicher und stimmig anfühlt.
Schutztext D – professionelle Hilfe
Wenn dich Belastung, Angst oder innere Anspannung stark beeinträchtigen, ist es sinnvoll, fachliche Unterstützung durch qualifizierte Ansprechpartner einzubeziehen.

5.3 Sensible Domänen mit erhöhter Vorsicht
	•	Sexualität
	•	Coaching Analyzer / Angst- und Druckthemen
	•	Medizinische Unsicherheit
	•	starke Verlust- oder Krisenthemen
	•	politisch oder weltanschaulich aufgeladene Konfliktfelder

5.4 Freigaberegel für sensible Fragen
Vorschlag
Eine sensible Frage darf nur auf active gehen, wenn:
	•	Sprache geprüft wurde
	•	Optionen keine Suggestion enthalten
	•	Ergebnistext Schutzlogik enthält
	•	mindestens ein Coach/Admin sie gesehen hat

5.5 UI-Hinweise für sensible Pfade
Im UI sollten sensible Pfade ruhiger und sicherer wirken:
	•	weniger spielerische Animationen
	•	zurückhaltendere Farben
	•	keine scherzhafte Tonalität
	•	klare Rückzugsmöglichkeiten
	•	„Überspringen“-Funktion sichtbar

6. Priorisierte Umsetzungsreihenfolge ab jetzt
- Schritt 1: Google-Sheets-Struktur V2 exakt anlegen.
- Schritt 2: CSV-Importer und Validator bauen.
- Schritt 3: Go-Repo-Struktur V2 erzeugen.
- Schritt 4: Fragen-, Optionen- und Review-Modelle umsetzen.
- Schritt 5: Scoring-Pakete für exact, partial, scale, non_scored einführen.
- Schritt 6: Erste sensible Textbausteine in result_texts pflegen.
- Schritt 7: Danach erst UI und Pfadlogik enger mit dem Backend verheiraten.

7. Nächste konkrete TODOs
Erledigt mit diesem Dokument
	•	CSV-/Google-Sheets-Struktur V2 finalisiert
	•	Go-Repo-Struktur V2 konkretisiert
	•	Review-Workflow getrennt beschrieben
	•	Bewertungslogik für multiple_select, scale, reflection formuliert
	•	Schutzbausteine für sensible Domänen festgehalten
Nächste echte Umsetzungsaufgaben
	•	Tabellenblätter im Google Sheet anlegen
	•	Referenz-IDs festziehen
	•	CSV-Seed-Dateien erzeugen
	•	Go-Pakete anlegen
	•	Importer schreiben
	•	Review-Tabellen implementieren
	•	Scoring-Services codieren

8. Abschluss
Mit diesem Paket ist der nächste Evolutionsschritt für CEE4AI sauber vorbereitet.
Wir haben jetzt nicht nur Vision und grobe Architektur, sondern eine belastbare Brücke zwischen:
	•	redaktioneller Datenpflege
	•	technischer Engine-Struktur
	•	Bewertungslogik
	•	Coaching-Ausbau
	•	Schutzmechanismen für sensible Domänen
Genau dadurch bekommt unser Projekt die Substanz, die du dir wünschst – und die es verdient.
