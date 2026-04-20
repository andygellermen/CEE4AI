# event-taxonomy-v3-final.md

## Zweck
Dieses Dokument definiert die Event-Taxonomie in der V3-Endfassung für die CPE-Produktfamilie.

Events sind der Rohstoff für:
- Musterbildung
- Nutzungsanalysen
- Qualitätsverbesserung
- Confidence-Modelle
- Zertifikats- und Integritätsnachweise
- Meaning-/Journey-Analysen

---

## Leitprinzip
Nicht jede Datenänderung ist ein Event – aber jedes relevante Nutzer-, Content-, Review-, Meaning- oder Zertifikatsereignis sollte als Event modellierbar sein.

Die Event-Taxonomie dient drei Zwecken:
1. operative Nachvollziehbarkeit
2. analytische Verdichtung
3. spätere Ableitung von Heuristiken und Modellen

---

## Eventfamilien
### A. Session Events
- `session.created`
- `session.started`
- `session.paused`
- `session.resumed`
- `session.completed`
- `session.abandoned`
- `session.expired`

### B. Package Events
- `package.created`
- `package.started`
- `package.completed`
- `package.aborted`
- `package.continuation_window_opened`
- `package.continuation_window_expired`

### C. Question Delivery Events
- `question.presented`
- `question.skipped`
- `question.timeout`
- `question.returned`
- `question.replaced`

### D. Answer Events
- `answer.submitted`
- `answer.updated`
- `answer.evaluated`
- `answer.scored`
- `answer.reflection_captured`
- `answer.scale_recorded`

### E. Profiling Events
- `profile.snapshot_created`
- `profile.updated`
- `profile.depth_increased`
- `profile.confidence_recalculated`
- `profile.extension_recommended`

### F. Review & Quality Events
- `review.flag_created`
- `review.flag_resolved`
- `review.status_changed`
- `review.translation_approved`
- `review.question_disabled`

### G. Localization Events
- `localization.created`
- `localization.machine_seeded`
- `localization.human_reviewed`
- `localization.approved`
- `localization.archived`

### H. Credential Events
- `credential.rule_matched`
- `credential.issued`
- `credential.viewed`
- `credential.downloaded`
- `credential.verified`
- `credential.revoked`
- `credential.reissued`

### I. Governance Events
- `ruleset.activated`
- `policy.updated`
- `sensitivity_rule.changed`
- `issuer.trust_level_changed`

### J. Meaning / Spiritual Events
- `meaning.path_started`
- `meaning.reflection_captured`
- `meaning.worldview_prompt_presented`
- `meaning.path_completed`
- `meaning.journey_certificate_issued`
- `meaning.review_required`
- `meaning.review_approved`
- `meaning.localization_human_reviewed`

---

## Event-Hauptfelder
```text
event_id
event_type
event_version
occurred_at
subject_type
subject_id
session_id nullable
user_id nullable
domain nullable
language nullable
region nullable
payload
correlation_id
causation_id nullable
source
```

### Bedeutungen
- `correlation_id`: verbindet zusammengehörige Abläufe, z. B. komplette Sessions
- `causation_id`: zeigt, welches Event ein Folgeevent ausgelöst hat
- `subject_type`: session, package, question, answer, profile, credential, translation, meaning_path etc.
- `payload`: versionierte Zusatzdaten

---

## Event-Taxonomie für Musterbildung
### Mustergruppe 1 – Verlauf
Wie bewegen sich Nutzer durch Sessions, Pakete und Modi?

### Mustergruppe 2 – Belastung
Wo entstehen Zeitdruck, Unsicherheit, Überforderung, Abbrüche?

### Mustergruppe 3 – Qualität
Wo zeigen Fragen, Übersetzungen oder Paketlogiken Schwächen?

### Mustergruppe 4 – Vertrauen
Welche Bedingungen führten zu belastbaren Zertifikaten?

### Mustergruppe 5 – Wachstum
Wie werden Profile tiefer, stabiler, differenzierter oder sicherer?

### Mustergruppe 6 – Meaning / Spiritual Depth
Welche Reflexionspfade werden vertieft? Welche Journey-Formen fördern Rückkehr, Stille, Resonanz oder Überforderung?

---

## Priorisierte Eventtypen für den ersten Ausbau
### Priorität A
- `session.started`
- `package.started`
- `question.presented`
- `answer.submitted`
- `answer.evaluated`
- `profile.snapshot_created`
- `review.flag_created`
- `credential.issued`

### Priorität B
- `session.abandoned`
- `package.aborted`
- `question.skipped`
- `localization.approved`
- `credential.verified`
- `meaning.path_started`
- `meaning.reflection_captured`

---

## Event-Retention und Speicherung
### Operativ in PostgreSQL
Kurz- bis mittelfristig als Event-Outbox und operatives Event-Archiv.

### Analytisch später separat
Zur Musterbildung und Verdichtung in eine Analytics-Ebene ausleiten.

### Retentionsidee
- operative Detail-Events: begrenzte Primärhaltung
- analytische Verdichtungen: langfristig
- Governance-/Credential-Events: langfristig und auditierbar

---

## Idempotenz und Deduplizierung
Besonders kritisch für:
- `credential.issued`
- `credential.revoked`
- `review.status_changed`
- `ruleset.activated`
- `meaning.journey_certificate_issued`

Darum:
- `event_id` eindeutig
- Konsumenten idempotent
- fachlich kritische Events zusätzlich absichern

---

## Sensible Events
### Besonders sensibel
- Coaching-Trigger
- Körperreaktionen
- freie Reflexionstexte
- Meaning-/Worldview-Inhalte
- Identitäts- und Zertifikatsdaten

### Vermeidungsstrategie
- sensible Payloads minimieren
- PII trennen
- pseudonymisieren, wo möglich
- Meaning-/Spiritual-Daten nicht in abwertende Rankingmodelle kippen lassen

---

## Strategische Bedeutung
Die Event-Taxonomie ist kein Logging-Zubehör, sondern das Fundament für:
- Confidence-Modelle
- Pfadoptimierung
- Zertifikatsintegrität
- Benchmarking
- Meaning-/Journey-Analysen
- spätere ML-/Heuristik-Layer

---

## Abschluss
Mit dieser V3-Endfassung kann CPE Ereignisse so sammeln, dass aus ihnen echte Musterintelligenz entsteht – ohne das operative Wahrheitssystem zu verwässern.

