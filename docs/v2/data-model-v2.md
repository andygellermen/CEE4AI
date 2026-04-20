# data-model-v2.md

## Zweck

Dieses Dokument erweitert das bisherige Datenmodell von CEE4AI auf den neuen Erkenntnisstand.

Es berücksichtigt jetzt ausdrücklich:

* mehrere Domänen
* verschiedene Fragetypen
* mehrere richtige Antworten
* hierarchische Fragepfade
* Coaching-orientierte Triggerlogik
* sensible Inhalte
* Test- und Review-Prozesse

---

## Grundprinzip

Das Datenmodell darf nicht länger nur auf einen simplen Multiple-Choice-Test zugeschnitten sein.

Es muss künftig mindestens zwei große Betriebsweisen tragen:

1. **Profil- und Entwicklungsmodus**
2. **Test-/Reviewmodus**
3. perspektivisch **Coaching- und Belastungs-Navigationsmodus**

---

## Kernfelder für Fragen

### Pflichtfelder

```text
id
external_id
domain
category_id
subcategory_id
parent_id
question_type
scoring_mode
denktyp_id
level
difficulty
question
explanation
weight
review_status
is_sensitive
age_gate
is_active
created_at
updated_at
```

---

## Bedeutung der neuen Felder

### `domain`

Ordnet eine Frage einem fachlichen Modul zu.

Beispielwerte:

* `cognitive_profile`
* `coaching_analyzer`
* später weitere Domänen

### `question_type`

Bestimmt, wie eine Frage dargestellt und beantwortet wird.

Empfohlene Startwerte:

* `single_choice`
* `multiple_select`
* `scale`
* `reflection`
* optional später `open_text`
* optional später `ranking`

### `parent_id`

Erlaubt hierarchische oder adaptive Folgefragen.

Beispiele:

* Subkategorie folgt auf Hauptkategorie
* Körperreaktionsfrage folgt auf Belastungsthema
* Vertiefungsfrage folgt auf auffällige Auswahl

### `trigger_group`

Dient besonders dem Coaching-Modul. Markiert Fragen, die zu einem gemeinsamen Belastungs- oder Themencluster gehören.

Beispielwerte:

* `krieg_weltlage`
* `arbeit_leistung`
* `gesundheit_unsicherheit`
* `verlust_beziehung`

### `body_reference`

Optionales Feld für körperbezogene Marker oder Reaktionsbezüge.

Beispielwerte:

* `brustenge`
* `bauchdruck`
* `herzklopfen`
* `innere_unruhe`
* `schlafprobleme`

### `scoring_mode`

Legt fest, wie Antworten gewertet werden.

Empfohlene Startwerte:

* `exact`
* `partial`
* `weighted`
* `non_scored`
* optional später `path_only`

### `review_status`

Status für Test-, Qualitäts- und Freigabeprozesse.

Empfohlene Werte:

* `draft`
* `active`
* `review`
* `flagged`
* `disabled`
* `archived`

### `is_sensitive`

Kennzeichnet sensible Inhalte.

Beispiele:

* Sexualität
* belastende Gesundheitsfragen
* Angst-/Druck-Themen
* politische oder weltanschauliche Reizthemen

### `age_gate`

Optionales Alters- oder Freigabefeld.

Beispiele:

* `0`
* `12`
* `16`
* `18`

---

## Antwortmodell

Da künftig nicht jede Frage genau eine richtige Antwort hat, muss das Antwortmodell flexibler werden.

### Tabelle `question_options`

```text
id
question_id
option_key
option_label
is_correct
score_weight
display_order
is_active
created_at
updated_at
```

### Bedeutung

* `option_key` z. B. `A`, `B`, `C`, `D`
* `is_correct` kann bei `multiple_select` mehrfach wahr sein
* `score_weight` erlaubt Teilgewichtungen

### Vorteil

Damit löst du dich von der starren Vier-Spalten-Logik `a`, `b`, `c`, `d`.

---

## Fragen-CSV V2

Für die Pflege per CSV oder Google Sheets ist eine flache Struktur weiterhin hilfreich.

### Vorschlag für das Blatt `questions`

```text
id
external_id
domain
category_id
subcategory_id
parent_id
question_type
scoring_mode
denktyp_id
level
difficulty
trigger_group
body_reference
question
explanation
weight
review_status
is_sensitive
age_gate
is_active
```

### Vorschlag für das Blatt `question_options`

```text
id
question_id
option_key
option_label
is_correct
score_weight
display_order
is_active
```

---

## Referenztabellen

### `categories`

```text
id
slug
name
description
is_sensitive
is_active
```

### `subcategories`

```text
id
category_id
slug
name
description
is_sensitive
is_active
```

### `denktypes`

```text
id
slug
name
description
development_hint
is_active
```

### `domains`

```text
id
slug
name
description
is_active
```

### `trigger_groups`

```text
id
slug
name
description
domain
is_sensitive
is_active
```

### `body_references`

```text
id
slug
name
description
is_sensitive
is_active
```

---

## Sessionmodell V2

### Tabelle `sessions`

```text
id
user_id
mode
domain
started_at
finished_at
status
created_at
updated_at
```

### Bedeutung neuer Felder

* `mode` z. B. `profile`, `test`, `coaching`
* `domain` trennt Domänenläufe
* `status` z. B. `active`, `completed`, `abandoned`

---

## Antwortspeicherung V2

### Tabelle `answers`

```text
id
session_id
question_id
selected_option_ids
free_text_answer
scale_value
is_correct
raw_score
evaluated_at
created_at
updated_at
```

### Hinweise

* `selected_option_ids` kann als JSON-Array gespeichert werden
* `free_text_answer` ist für Reflexionen oder offene Formate reserviert
* `scale_value` ist für Skalenfragen hilfreich
* `raw_score` erlaubt spätere flexible Auswertung

---

## Review- und Testmodell

### Tabelle `question_reviews`

```text
id
question_id
session_id
reviewer_role
flag_type
comment
severity
created_at
```

### Mögliche `flag_type`-Werte

* `unclear_question`
* `multiple_answers_seem_correct`
* `no_answer_fits`
* `too_easy`
* `too_hard`
* `linguistically_unclear`
* `factually_problematic`
* `emotionally_triggering`

### Mögliche `reviewer_role`-Werte

* `tester`
* `coach`
* `admin`
* später `expert_reviewer`

---

## Ergebnislogik V2

### Bereich A – kognitive Profile

Bewertung entlang von Denktypen, Interessen und Entwicklungsniveaus.

### Bereich B – Review

Aggregierte Fragequalität, Häufung von Flags, Prioritäten für Überarbeitung.

### Bereich C – Coaching-Modul

Pfad- und Themenverdichtungen, Triggercluster, Körperreaktionsbezüge, nicht-diagnostische Reflexionshinweise.

---

## Beispiel für Fragenobjekt in JSON

```json
{
  "id": 101,
  "external_id": "Q-COG-0001",
  "domain": "cognitive_profile",
  "category_id": 10,
  "subcategory_id": 101,
  "parent_id": null,
  "question_type": "single_choice",
  "scoring_mode": "exact",
  "denktyp_id": 1,
  "level": 1,
  "difficulty": 1,
  "trigger_group": null,
  "body_reference": null,
  "question": "Welche Zahl setzt die Reihe sinnvoll fort: 2, 4, 8, 16, ?",
  "explanation": "Die Zahlen verdoppeln sich jeweils.",
  "weight": 1.0,
  "review_status": "active",
  "is_sensitive": false,
  "age_gate": 0,
  "is_active": true
}
```

---

## Beispiel für Coaching-Frageobjekt

```json
{
  "id": 2001,
  "external_id": "Q-COA-0101",
  "domain": "coaching_analyzer",
  "category_id": 80,
  "subcategory_id": 801,
  "parent_id": null,
  "question_type": "single_choice",
  "scoring_mode": "path_only",
  "denktyp_id": null,
  "level": 1,
  "difficulty": 1,
  "trigger_group": "krieg_weltlage",
  "body_reference": "bauchdruck",
  "question": "Spürst du in Bezug auf Nachrichten zu Krieg und Weltlage häufiger Druck oder Enge im Körper?",
  "explanation": "Diese Frage dient der Selbstklärung, nicht der Diagnosestellung.",
  "weight": 1.0,
  "review_status": "draft",
  "is_sensitive": true,
  "age_gate": 16,
  "is_active": true
}
```

---

## Empfehlungen für Google Sheets V2

### Blattstruktur

* `questions`
* `question_options`
* `categories`
* `subcategories`
* `denktypes`
* `domains`
* `trigger_groups`
* `body_references`
* `review_flags`

### Vermeidungsstrategie

Freitext bei IDs und Slugs vermeiden.
Pflege lieber Referenzen sauber, sonst entstehen früh Inkonsistenzen.

---

## Minimalanforderung für den nächsten Go-Schritt

Der nächste Code-Schritt sollte mindestens bereits verstehen:

* `domain`
* `question_type`
* `scoring_mode`
* `review_status`
* `is_sensitive`
* `age_gate`
* `parent_id`

`trigger_group` und `body_reference` sollten ebenfalls vorbereitet, aber im ersten MVP noch nicht vollständig ausgewertet werden.

---

## Checkliste

### Neu aufgenommen

* neue Kernfelder integriert
* Antwortmodell flexibilisiert
* Review-Modell ergänzt
* Coaching-Felder vorbereitet

### Offen

* finale ENUM-Werte technisch festziehen
* SQLite-Schema ableiten
* CSV- und Sheets-Dateien auf V2 umstellen
* Go-Structs und Repository-Logik angleichen

---

## Abschluss

Das Datenmodell V2 macht CEE4AI deutlich zukunftsfähiger.
Es schützt uns davor, in der Anfangsphase an einem zu engen Quiz-Modell festzukleben, und öffnet zugleich den Weg zu Reviewfähigkeit, Coaching-Ausbau und echter adaptiver Domänenlogik.
