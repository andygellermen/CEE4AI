# CEE4AI Go-Live Checkliste

Diese Checkliste ist fuer den **allerersten produktiven Test-Run** auf:

- `https://cpe.geller.men`

Sie ist bewusst fuer das aktuelle Live-MVP geschrieben:

- Go-API
- eigener Postgres-Container im CPE-Stack
- Rollout per Ansible
- Routing ueber den bestehenden zentralen Traefik

## Ziel dieses ersten Go-Live

Am Ende des Runs wollen wir belastbar sagen koennen:

- `cpe.geller.men` antwortet stabil
- das CPE-Backend ist ueber Traefik erreichbar
- Migrationen und Seeds sind sauber verfuegbar
- ein einfacher Snapshot-Flow funktioniert
- ein guarded Meaning-/Governance-Flow funktioniert

Das Ziel ist **nicht** schon ein vollwertiger Produkt-Launch, sondern ein ruhiger, sauberer, produktionsnaher Erstlauf.

## Erfolgskriterium

Der erste Go-Live gilt als erfolgreich, wenn alle folgenden Punkte stimmen:

- `GET /healthz` liefert `200` und `ok`
- der automatische Smoke-Test laeuft ohne Fehler durch
- `payload.governance.delivery_mode == "guarded"` im Meaning-Test
- keine Container sind in Restart-Schleifen
- keine offensichtlichen Traefik-, TLS- oder DB-Fehler in den Logs

## Abbruchkriterien

Den Run bitte bewusst abbrechen, wenn einer dieser Punkte eintritt:

- `cpe.geller.men` liefert dauerhaft `404`, `502` oder TLS-Fehler
- der API-Container startet wiederholt neu
- Postgres bleibt unhealthy
- Session-Start oder Fragen liefern reproduzierbar `500`
- der Smoke-Test scheitert mehrfach an derselben Stelle

## T-30 Minuten: Pre-Flight

- [ ] Es ist ein ruhiger Slot ohne parallele Infra-Aenderungen reserviert.
- [ ] Domain `cpe.geller.men` zeigt auf den richtigen Server.
- [ ] Zentraler Traefik laeuft auf dem Server.
- [ ] Docker und `docker compose` sind auf dem Server verfuegbar.
- [ ] In [cee4ai_live_mvp.yml](/Users/andygellermann/Documents/Projects/CEE4AI/infra/ansible/group_vars/cee4ai_live_mvp.yml:1) ist `cee4ai_postgres_password` gesetzt.
- [ ] Optionale Werte wie `cee4ai_domain_aliases` und `cee4ai_traefik_middlewares` sind bewusst geprueft.
- [ ] Das Repo ist auf dem gewuenschten Stand.

## T-20 Minuten: Infra-Check

Auf dem Server oder in eurer Infra-Konsole pruefen:

```bash
docker ps
docker network inspect traefik >/dev/null
```

Sollte passen:

- der zentrale Traefik-Container laeuft
- das Docker-Netzwerk `traefik` existiert

## T-15 Minuten: CPE-Deployment ausrollen

Im CEE4AI-Repo:

```bash
ansible-playbook -i ../infra/infra/ansible/inventory/hosts.ini infra/ansible/deploy-live-mvp.yml
```

Erwartung:

- das Playbook laeuft ohne Fehler durch
- am Ende erscheint die Meldung, dass der lokale Health-Endpunkt erreichbar ist

## T-10 Minuten: Container- und Log-Check

Auf dem Server:

```bash
docker ps --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}'
docker logs --tail 100 cpe-api
docker logs --tail 100 cpe-postgres
```

Sollte passen:

- `cpe-api` laeuft stabil
- `cpe-postgres` laeuft stabil
- keine offensichtlichen Fehler wie:
  - DB connection refused
  - migration failed
  - seed import failed
  - bind / routing errors

## T-8 Minuten: Seed-Status pruefen

Da Seeds standardmaessig nicht automatisch eingespielt werden, beim ersten echten Lauf bewusst pruefen:

```bash
docker exec cpe-api /app/cee4ai-seed
```

Erwartung:

- kein Fehler
- der Seed-Import beendet sich sauber

Hinweis:

- Wenn die Seeds bereits in derselben produktiven DB eingespielt wurden, ist der erneute Lauf durch die Upsert-Logik unkritisch.

## T-5 Minuten: Externe Health-Pruefung

Von lokal:

```bash
curl -fsS https://cpe.geller.men/healthz
```

Erwartung:

```text
ok
```

Wenn das nicht klappt:

- Traefik-Routing pruefen
- TLS pruefen
- Container-Logs pruefen
- nicht in den Smoke-Test weitergehen

## T-3 Minuten: Automatischen Smoke-Test fahren

Von lokal:

```bash
make smoke-test-live
```

Oder direkt:

```bash
BASE_URL=https://cpe.geller.men ./infra/ansible/smoke-test-live-mvp.sh
```

Wenn du die Antworten fuer spaeteres Debugging behalten willst:

```bash
BASE_URL=https://cpe.geller.men KEEP_ARTIFACTS=true ./infra/ansible/smoke-test-live-mvp.sh
```

## T-0: Ergebnis bewerten

Der Run ist gut, wenn:

- Healthcheck gruen ist
- Standard-Snapshot-Flow gruen ist
- Meaning-/Governance-Flow gruen ist
- keine Restart-Schleifen laufen
- Logs ruhig sind

Dann gilt:

- der erste produktive Test-Run auf `cpe.geller.men` ist bestanden

## Sofortmassnahmen bei Fehlern

### Fall A: TLS oder Routing kaputt

Pruefen:

```bash
docker logs --tail 100 infra-traefik
docker inspect cpe-api
```

Fragen:

- haengt `cpe-api` im Netzwerk `traefik`?
- stimmen Router-Regel und Domain?
- ist `letsEncrypt` verfuegbar?

### Fall B: API startet, aber `500` bei Sessions/Fragen

Pruefen:

```bash
docker logs --tail 200 cpe-api
docker logs --tail 200 cpe-postgres
```

Wahrscheinliche Ursachen:

- Migrationen fehlen
- Seed-Daten fehlen
- DB noch nicht sauber bereit

### Fall C: Meaning-Flow kommt nicht guarded zurueck

Pruefen:

- wurden die Governance-Seeds eingespielt?
- liegt wirklich Domain `5` vor?
- sind die Meaning-Fragen und Uebersetzungen aktiv und `approved`?

## Rollback fuer diesen MVP

Wenn der erste Test-Run sauber abgebrochen werden soll:

```bash
cd /srv/cpe/compose
docker compose --env-file cee4ai.env -f docker-compose.live-mvp.yml down
```

Wenn nur die API neu gebaut und neu gestartet werden soll:

```bash
ansible-playbook -i ../infra/infra/ansible/inventory/hosts.ini infra/ansible/deploy-live-mvp.yml
```

## Nach dem erfolgreichen Run

- [ ] Kurz notieren, wann der Run stattgefunden hat.
- [ ] Ergebnis in einem kurzen Ops-/Projekt-Log festhalten.
- [ ] Auffaellige Beobachtungen dokumentieren.
- [ ] Entscheiden, ob als naechstes ein zweiter stabiler Test-Run oder schon kleine Produktverfeinerungen folgen.

## Minimales Go-Live-Protokoll

```text
Datum:
Startzeit:
Endzeit:
Domain: cpe.geller.men
Deploy erfolgreich: ja/nein
Healthcheck erfolgreich: ja/nein
Smoke-Test erfolgreich: ja/nein
Meaning-Guardrail sichtbar: ja/nein
Auffaelligkeiten:
Naechster Schritt:
```
