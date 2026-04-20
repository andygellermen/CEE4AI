# CEE4AI Infra Live MVP

Dieses Verzeichnis enthaelt die erste Infra-Basis fuer das Live-MVP auf einem Server mit:

- Docker + Docker Compose Plugin
- Ansible
- Traefik als bestehendem Reverse Proxy

## Annahmen fuer diesen MVP

- Die API wird als Container auf dem Zielserver gebaut.
- Der Stack bringt fuer den ersten Live-MVP eine eigene Postgres-Instanz mit.
- Traefik laeuft bereits auf dem Server und nutzt in eurer Infra das Docker-Netzwerk `traefik`.
- Die API bleibt zusaetzlich lokal auf `127.0.0.1:<bind-port>` erreichbar, damit Ansible den Healthcheck pruefen kann.
- Der aktuelle Default ist auf `cpe.geller.men` sowie den Resolver `letsEncrypt` ausgerichtet.

## Dateien

- `deploy-live-mvp.yml`: Playbook fuer Build und Rollout
- `group_vars/cee4ai_live_mvp.yml`: zentrale Variablen fuer Domain, Traefik, DB und Bootverhalten
- `inventory.example.ini`: einfaches Inventory-Beispiel
- `templates/docker-compose.live-mvp.yml.j2`: Compose-Template mit Traefik-Labels
- `templates/cee4ai.env.j2`: Laufzeit-Env fuer App und Postgres
- `live-smoke-test-runbook.md`: erster End-to-End-Test fuer `cpe.geller.men`
- `smoke-test-live-mvp.sh`: kleiner Curl-/JQ-Helper fuer den schnellen Lauf
- `go-live-checklist-cpe.md`: ruhige Checkliste fuer den ersten produktiven Test-Run

## Erstes Setup

1. Wenn ihr wie in eurer zentralen Infra lokal auf dem Server deployt, kann das Beispiel-Inventory direkt so bleiben.
2. `group_vars/cee4ai_live_mvp.yml` anpassen:
   - `cee4ai_postgres_password`
   - optional `cee4ai_domain_aliases`
   - optional `cee4ai_traefik_middlewares`
   - optional `cee4ai_run_seed_on_boot`
3. Sicherstellen, dass Docker und `docker compose` auf dem Zielserver verfuegbar sind.
4. Sicherstellen, dass Traefik auf demselben Docker-Netzwerk erreichbar ist.

## Deployment

```bash
ansible-playbook -i infra/ansible/inventory.example.ini infra/ansible/deploy-live-mvp.yml
```

Wenn ihr den Rollout innerhalb eurer bestehenden `infra`-Umgebung fahrt, ist die naheliegende Variante:

```bash
ansible-playbook -i ../infra/infra/ansible/inventory/hosts.ini infra/ansible/deploy-live-mvp.yml
```

## Reale Infra-Ausrichtung

Die CPE-MVP-Basis ist jetzt bewusst an eure bestehende Infra angenaehert:

- Docker-Netzwerk: `traefik`
- TLS-Resolver: `letsEncrypt`
- Standard-Domain: `cpe.geller.men`
- HTTP -> HTTPS Redirect ueber eigene Router/Middlewares
- optionale Alias-Domains und optionale Traefik-Middlewares

## Verhalten beim Start

- Die API fuehrt standardmaessig Migrationen beim Containerstart aus.
- Seeds werden nur eingespielt, wenn `cee4ai_run_seed_on_boot: true` gesetzt ist.
- Danach startet die API auf Port `8080` im Container.

## Smoke Test

Nach dem ersten Rollout ist der naechste Schritt das Runbook unter [live-smoke-test-runbook.md](/Users/andygellermann/Documents/Projects/CEE4AI/infra/ansible/live-smoke-test-runbook.md:1).

Direkter Schnelllauf:

```bash
BASE_URL=https://cpe.geller.men ./infra/ansible/smoke-test-live-mvp.sh
```

Mit behaltenen Artefakten:

```bash
BASE_URL=https://cpe.geller.men KEEP_ARTIFACTS=true ./infra/ansible/smoke-test-live-mvp.sh
```

## Erster Go-Live

Fuer den allerersten produktiven Test-Run liegt jetzt zusaetzlich die Checkliste unter [go-live-checklist-cpe.md](/Users/andygellermann/Documents/Projects/CEE4AI/infra/ansible/go-live-checklist-cpe.md:1).
