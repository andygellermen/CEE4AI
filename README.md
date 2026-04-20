# CEE4AI

CEE4AI ist derzeit ein Go-basiertes V3-Backend mit:

- Postgres-first Datenmodell
- Migrations- und Seed-Pfad
- Runtime fuer Sessions, Fragen, Antworten und Snapshots
- Review-, Governance- und Audit-Flow

## Lokaler Kern

### Voraussetzungen

- Go `1.23.4`
- Postgres

### Wichtige Commands

```bash
make migrate
make seed
make run
make test
```

### Standard-Endpunkte

- `GET /healthz`
- `POST /api/v1/sessions`
- `GET /api/v1/questions/next?session_id=...`
- `POST /api/v1/answers`
- `GET /api/v1/results?session_id=...`
- `POST /api/v1/reviews/questions/flag`
- `POST /api/v1/reviews/questions/decisions`
- `POST /api/v1/reviews/localizations`

## Container-Basis

Das Repo enthaelt jetzt eine erste Container-Basis fuer ein Live-MVP:

- [Dockerfile](/Users/andygellermann/Documents/Projects/CEE4AI/Dockerfile:1)
- [docker/entrypoint.sh](/Users/andygellermann/Documents/Projects/CEE4AI/docker/entrypoint.sh:1)

Der Container:

- baut `api`, `migrate` und `seed`
- fuehrt optional Migrationen und Seeds beim Start aus
- startet anschliessend die API

### Docker lokal

```bash
make docker-build
```

Der Image-Build ist bewusst der erste lokale Container-Schritt. Fuer einen echten Lauf braucht der Container eine erreichbare Postgres-Instanz oder den infra-aehnlichen Compose-Stack aus `infra/ansible`.

## Infra Live MVP

Fuer die Einspielung auf einen Infra-Server in eurer Docker-/Ansible-/Traefik-Umgebung gibt es jetzt ein erstes Geruest unter [infra/ansible/README.md](/Users/andygellermann/Documents/Projects/CEE4AI/infra/ansible/README.md:1).

Wesentliche Dateien:

- [infra/ansible/deploy-live-mvp.yml](/Users/andygellermann/Documents/Projects/CEE4AI/infra/ansible/deploy-live-mvp.yml:1)
- [infra/ansible/group_vars/cee4ai_live_mvp.yml](/Users/andygellermann/Documents/Projects/CEE4AI/infra/ansible/group_vars/cee4ai_live_mvp.yml:1)
- [infra/ansible/templates/docker-compose.live-mvp.yml.j2](/Users/andygellermann/Documents/Projects/CEE4AI/infra/ansible/templates/docker-compose.live-mvp.yml.j2:1)
- [infra/ansible/templates/cee4ai.env.j2](/Users/andygellermann/Documents/Projects/CEE4AI/infra/ansible/templates/cee4ai.env.j2:1)
- [infra/ansible/live-smoke-test-runbook.md](/Users/andygellermann/Documents/Projects/CEE4AI/infra/ansible/live-smoke-test-runbook.md:1)
- [infra/ansible/smoke-test-live-mvp.sh](/Users/andygellermann/Documents/Projects/CEE4AI/infra/ansible/smoke-test-live-mvp.sh:1)
- [infra/ansible/go-live-checklist-cpe.md](/Users/andygellermann/Documents/Projects/CEE4AI/infra/ansible/go-live-checklist-cpe.md:1)

### Wichtige MVP-Annahme

Der erste Infra-Stack bringt absichtlich eine eigene Postgres-Instanz mit, damit wir den ersten Server-Lauf ohne externen Registry- oder DB-Zwang stabil hinbekommen. Das ist bewusst ein Live-MVP-Schritt und keine finale Produktionsarchitektur.

### Aktuelle Infra-Ausrichtung

Die Defaults sind jetzt auf eure bestehende Linie gehaertet:

- Domain: `cpe.geller.men`
- Traefik-Docker-Netzwerk: `traefik`
- TLS-Resolver: `letsEncrypt`
- lokaler Health-Port fuer den Server-Check: `127.0.0.1:18080`

### Smoke Test Schnelllauf

```bash
BASE_URL=https://cpe.geller.men ./infra/ansible/smoke-test-live-mvp.sh
```
