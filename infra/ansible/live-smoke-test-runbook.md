# CEE4AI Live Smoke Test Runbook

Dieses Runbook beschreibt den ersten echten Smoke-Test fuer das Live-MVP auf:

- `https://cpe.geller.men`

Es deckt zwei bewusst unterschiedliche Faelle ab:

1. einen einfachen Snapshot-Flow
2. einen guarded Meaning-/Governance-Flow

Damit pruefen wir nicht nur, ob die API antwortet, sondern ob Runtime, Governance und Audit im Live-Betrieb zusammenhaengen.

## Zielbild

Nach dem Smoke-Test sollten wir belastbar sagen koennen:

- die Domain antwortet sauber
- die API ist ueber Traefik erreichbar
- Sessions koennen erzeugt werden
- Fragen werden ausgeliefert
- Antworten werden gespeichert
- Snapshots werden erzeugt
- sensitive Meaning-/Worldview-Inhalte laufen im guarded Modus

## Voraussetzungen

- das Live-MVP ist per Ansible ausgerollt
- `https://cpe.geller.men/healthz` antwortet
- Seed-Daten sind eingespielt
- lokal stehen `curl` und `jq` zur Verfuegung

## Seed-Voraussetzung pruefen

Beim ersten Server-Start sind Seeds standardmaessig **nicht** automatisch aktiv. Wenn noch keine Seed-Daten in der Live-DB liegen, einmalig auf dem Server ausfuehren:

```bash
docker exec cpe-api /app/cee4ai-seed
```

Anschliessend den Healthcheck pruefen:

```bash
curl -fsS https://cpe.geller.men/healthz
```

Erwartung:

```text
ok
```

## Automatischer Smoke-Test

Fuer den schnellen Gesamtlauf liegt ein kleiner Helper unter [smoke-test-live-mvp.sh](/Users/andygellermann/Documents/Projects/CEE4AI/infra/ansible/smoke-test-live-mvp.sh:1).

Beispiel:

```bash
BASE_URL=https://cpe.geller.men ./infra/ansible/smoke-test-live-mvp.sh
```

Wenn die JSON-Antworten fuer Debugging erhalten bleiben sollen:

```bash
BASE_URL=https://cpe.geller.men KEEP_ARTIFACTS=true ./infra/ansible/smoke-test-live-mvp.sh
```

Der Script prueft:

- `GET /healthz`
- Standard-Snapshot fuer Domain `1`
- guarded Meaning-Snapshot fuer Domain `5`
- zentrale Ergebnis- und Governance-Signale

## Manueller Test A: Standard Snapshot

### 1. Session starten

```bash
curl -fsS \
  -H 'Content-Type: application/json' \
  -d '{"domain_id":1,"mode":"snapshot","session_goal":"live smoke test","locale_language_id":1}' \
  https://cpe.geller.men/api/v1/sessions | jq
```

Erwartung:

- `session.id` ist gesetzt
- `session.mode == "snapshot"`
- `first_package` ist vorhanden

### 2. Naechste Frage holen

`SESSION_ID` aus dem ersten Schritt uebernehmen:

```bash
curl -fsS "https://cpe.geller.men/api/v1/questions/next?session_id=${SESSION_ID}" | jq
```

Erwartung:

- `question.id` ist gesetzt
- `question.question_type == "single_choice"`
- `governance.delivery_mode == "standard"`

### 3. Antwort senden

Bei den aktuellen Seed-Daten ist fuer die erste Domain die Antwortoption mit `option_key = "B"` die saubere positive Testantwort.

```bash
curl -fsS \
  -H 'Content-Type: application/json' \
  -d "{\"session_id\":\"${SESSION_ID}\",\"question_id\":1,\"selected_option_ids\":[2],\"certainty_level\":\"high\"}" \
  https://cpe.geller.men/api/v1/answers | jq
```

Erwartung:

- `answer.id` ist gesetzt
- `progress_state == "completed"`
- `governance.delivery_mode == "standard"`

### 4. Snapshot holen

```bash
curl -fsS "https://cpe.geller.men/api/v1/results?session_id=${SESSION_ID}" | jq
```

Erwartung:

- `result_type == "snapshot_profile"`
- `payload.top_signals.denktype == "analytical_logical"`
- `payload.top_signals.skill == "pattern_recognition"`
- `payload.governance.delivery_mode == "standard"`

## Manueller Test B: Guarded Meaning Snapshot

Dieser Test prueft die sensitive Linie mit Meaning-/Worldview-Governance.

### 1. Session starten

```bash
curl -fsS \
  -H 'Content-Type: application/json' \
  -d '{"domain_id":5,"mode":"snapshot","session_goal":"guarded meaning smoke test","locale_language_id":1}' \
  https://cpe.geller.men/api/v1/sessions | jq
```

Erwartung:

- `session.id` ist gesetzt
- `session.mode == "snapshot"`

### 2. Meaning-Frage holen

```bash
curl -fsS "https://cpe.geller.men/api/v1/questions/next?session_id=${SESSION_ID}" | jq
```

Erwartung:

- `question.question_type == "reflection"`
- `question.is_sensitive == true`
- `question.age_gate == 16`
- `governance.delivery_mode == "guarded"`
- `governance.review_required == true`
- `governance.sensitivity_flags` enthaelt mindestens:
  - `sensitive_content`
  - `age_gated`
  - `worldview_sensitive`
  - `human_review_required`

### 3. Freitext-Antwort senden

```bash
curl -fsS \
  -H 'Content-Type: application/json' \
  -d "{\"session_id\":\"${SESSION_ID}\",\"question_id\":3,\"free_text_answer\":\"Ich erlebe Sinn besonders in stillen Momenten, in Verbundenheit mit anderen Menschen und wenn mein Handeln mit meinen Werten zusammenpasst.\",\"certainty_level\":\"medium\"}" \
  https://cpe.geller.men/api/v1/answers | jq
```

Erwartung:

- `answer.id` ist gesetzt
- `governance.delivery_mode == "guarded"`

### 4. Guarded Snapshot holen

```bash
curl -fsS "https://cpe.geller.men/api/v1/results?session_id=${SESSION_ID}" | jq
```

Erwartung:

- `payload.governance.delivery_mode == "guarded"`
- `payload.governance.review_flag_count >= 1`
- `payload.governance.answered_sensitive_questions >= 1`
- `payload.governance.answered_worldview_sensitive_questions >= 1`
- `payload.governance.guardrails` ist nicht leer

## Optionaler Spot Check fuer Review-Endpunkte

Die Review-Endpunkte sollten fuer Admin-/Ops-Zwecke bevorzugt lokal gegen den gebundenen Port getestet werden, also auf dem Server oder per SSH-Tunnel, nicht als oeffentlicher Browser-Klicktest.

Basis:

- `http://127.0.0.1:18080`

### Manuelles Flag setzen

```bash
curl -fsS \
  -H 'Content-Type: application/json' \
  -d '{"question_id":3,"flag_slug":"worldview_sensitive","reviewer_role":"reviewer","comment":"manual smoke spot check","severity":"high"}' \
  http://127.0.0.1:18080/api/v1/reviews/questions/flag | jq
```

### Review-Status aendern

```bash
curl -fsS \
  -H 'Content-Type: application/json' \
  -d '{"question_id":3,"new_status":"review","reason":"manual smoke spot check"}' \
  http://127.0.0.1:18080/api/v1/reviews/questions/decisions | jq
```

Hinweis:

- Dieser Schritt veraendert echten Review-Zustand in der Live-DB.
- Deshalb nur bewusst und vorzugsweise auf einer frischen MVP-Instanz ausfuehren.

## Typische Fehlerbilder

### `404` oder `502` ueber `cpe.geller.men`

Pruefen:

- Traefik-Container laeuft
- CPE-Container haengt im Docker-Netzwerk `traefik`
- Domain zeigt auf den richtigen Server

### `500` bei Session oder Fragen

Pruefen:

- Migrationen sind gelaufen
- Seed-Daten sind vorhanden
- Postgres-Container ist gesund

### `question not found` oder `has_more: false` direkt nach Session-Start

Pruefen:

- Seeds wirklich eingespielt
- `question_translation.localization_status = approved` liegt fuer die verwendete Sprache vor

## Erfolgsbild

Ein guter erster Live-Smoke-Test ist erreicht, wenn:

- beide Flows ohne manuelle Reparatur durchlaufen
- der Standardfall `standard` bleibt
- der Meaning-Fall sichtbar `guarded` wird
- der Snapshot Governance-Informationen zurueckliefert

Dann darf das Kind wirklich anfangen, die ersten stabilen Schritte zu gehen.
