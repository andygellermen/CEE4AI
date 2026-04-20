# data-model-v3-final.md

## Zweck
Dieses Dokument ist die synchronisierte V3-Endfassung des Datenmodells für die CPE-Produktfamilie.

Es verbindet:
- Product Family
- Personality-Tags
- Meaning-/Spiritual-Tags
- Mehrsprachigkeit
- Paketlogik
- Eventmodell
- Credential-Layer
- Governance

---

## Grundsatz
Das Datenmodell V3 ist kein klassisches Quizschema mehr.  
Es ist ein **relationaler Kern mit mehrdimensionalen Zuordnungen, Übersetzungen, Ereignissen, Meaning-Achsen und Zertifikatslogik**.

---

## Hauptblöcke des Modells
### 1. Master Content
- Fragen-Master
- Optionen-Master
- Kategorien
- Subkategorien
- Domains

### 2. Translation & Localization
- Frageübersetzungen
- Antwortübersetzungen
- Ergebnistextübersetzungen
- Sprach- und Regionsmodelle
- Lokalisierungsstatus

### 3. Profiling Layer
- Denktypen
- Skills
- Personality Traits
- Interest Tags
- Trigger Groups
- Body References
- Development Paths
- Meaning Tags
- Worldview Frames

### 4. Runtime Layer
- Sessions
- Session Packages
- Answers
- Result Snapshots
- Profile Vectors

### 5. Review Layer
- Review Flags
- Question Reviews
- Localization Reviews
- Review Decisions

### 6. Credential Layer
- Assessment Programs
- Tracks
- Certificate Rules
- Templates
- Assertions
- Evidence
- Verifications
- Revocations

### 7. Governance Layer
- Policies
- Sensitivity Rules
- Ruleset Versions
- Audit Log

### 8. Event Layer
- Outbox
- Inbox
- Event Store

---

## Neue Kernfelder auf Frageebene
```text
question_family
confidence_tier
intended_use
estimated_time_seconds
cognitive_load_level
cultural_scope
region_scope
localization_status
meaning_depth
worldview_sensitivity
symbolic_interpretation_relevance
existential_load_level
```

### Bedeutung
- `question_family`: knowledge, skill, trait, interest, trigger, pathway, reflective, contemplative, comparative_worldview, symbolic_interpretation, existential
- `confidence_tier`: Vorhersage- und Aussagekraft einer Frage
- `intended_use`: Profil, Persönlichkeit, Coaching, Karriere, Zertifikat, Test, Meaning Journey
- `estimated_time_seconds`: wichtig für Paketsteuerung
- `cognitive_load_level`: mentale Belastung
- `cultural_scope`: universell, regional, lokal
- `region_scope`: optionale Beschränkung auf Regionen
- `localization_status`: Übersetzungs- und Freigabereife
- `meaning_depth`: spirituelle / existenzielle Tiefe
- `worldview_sensitivity`: Sensibilität für weltanschauliche Deutungsräume
- `symbolic_interpretation_relevance`: Relevanz für symbolische oder mythische Deutungen
- `existential_load_level`: Belastungspotenzial existenzieller Fragen

---

## Neue many-to-many Tag-Layer
### Zuordnungsebenen
- Frage ↔ Denktyp
- Frage ↔ Skill
- Frage ↔ Trait
- Frage ↔ Interest
- Frage ↔ Trigger
- Frage ↔ Body Reference
- Frage ↔ Pathway
- Frage ↔ Meaning Tag
- Frage ↔ Worldview Frame

### Jede Zuordnung kann zusätzlich tragen
- `weight`
- `rationale`
- `is_active`

---

## Mehrsprachigkeit im Datenmodell
Die Frage ist nicht der Text.  
Die Frage ist der fachliche Master.

### Deshalb getrennt modellieren
- `question_master`
- `question_translation`
- `question_option_master`
- `question_option_translation`
- `result_text_master`
- `result_text_translation`

### Zusätzliche sensible Marker
Für Meaning-/Spiritual-Inhalte sinnvoll:
- `requires_human_review`
- `worldview_sensitive`

---

## Session- und Paketlogik
V3 unterstützt zeitadaptive Verlaufsmuster.

### Neue Session-Felder
```text
session_goal
result_confidence
progress_state
```

### Neue Package-Felder
```text
package_index
package_size
estimated_time_seconds
actual_time_seconds
completion_quality
continuation_window_until
recommended_next_mode
```

### Sinn dieser Felder
Sie ermöglichen:
- Snapshot-Profile
- Guided Progression
- Appetizer-Profile
- Fortsetzungsfenster
- Confidence- und Abbruchmodelle
- Meaning Journeys

---

## Answer-Modell V3
### Antwortarten
- Single Choice
- Multiple Select
- Scale
- Reflection

### Antwortspeicherung enthält mindestens
```text
selected_option_ids
scale_value
free_text_answer
raw_score
evaluated_score
certainty_level
answered_at
```

---

## Result-Modell V3
### Ergebnisfelder
```text
result_type
profile_depth
certainty_level
snapshot_payload
ruleset_version
```

### Typische Ergebnisarten
- snapshot_profile
- deep_profile
- personality_extension
- coaching_reflection
- pathway_readiness
- certificate_basis_result
- meaning_journey_snapshot
- worldview_reflection_result

---

## Meaning-/Spiritual-Erweiterungen
### Neue Referenztabelle
`worldview_frames`

#### Beispielhafte Slugs
- `theistic_creation`
- `symbolic_creation`
- `evolutionary_naturalism`
- `integral_spiritual`
- `secular_existential`
- `comparative_open`

### Neue Meaning Tags
Beispiele:
- `meaning`
- `transcendence`
- `connectedness`
- `inner_truth`
- `creation_orientation`
- `sacredness`
- `service`
- `conscience`
- `mortality_awareness`
- `wonder`
- `worldview_reflection`
- `symbolic_depth`

### Neue spirituell relevante Profilachsen
- `meaning_orientation`
- `reflection_depth`
- `transcendence_sensitivity`
- `value_alignment`
- `connectedness_orientation`
- `existential_stability`
- `conscience_depth`
- `symbolic_openness`

---

## Credential-Layer im Datenmodell
### Wichtigste Objekte
- `assessment_programs`
- `assessment_tracks`
- `certificate_rules`
- `credential_assertions`
- `credential_evidence`
- `credential_verifications`
- `credential_revocations`

### Wichtige neue Ergebnis-/Credential-Felder
```text
certificate_level
issuer_id
evidence_hash
issued_under_ruleset_version
assessment_integrity_score
credential_standard
standards_alignment
verification_token
```

### Für Meaning-/Journey-Pfade
Sinnvoll eher:
- Journey Completion
- Reflective Development Track Completion
- Guided Pathway Completion

Nicht sinnvoll:
- objektivierende Stufenbehauptungen spiritueller Höherentwicklung

---

## Event-Layer im Datenmodell
### Kernobjekte
- `event_outbox`
- `event_inbox`
- `event_store`

### Event-Hauptfelder
```text
event_id
event_type
event_version
subject_type
subject_id
session_id
user_id
occurred_at
payload
correlation_id
causation_id
source
```

### Zusätzliche Meaning-/Spiritual-Ereignisse
- `meaning.path_started`
- `meaning.reflection_captured`
- `meaning.worldview_prompt_presented`
- `meaning.path_completed`
- `meaning.journey_certificate_issued`

---

## Review- und Governance-Daten
### Review
- `review_flags`
- `question_reviews`
- `review_decisions`
- `localization_reviews`

### Governance
- `content_policies`
- `sensitivity_rules`
- `ruleset_versions`
- `audit_log`

### Für Meaning-/Spiritual-Inhalte besonders wichtig
- dokumentierte Human Reviews
- weltanschauliche Sensitivität
- vorsichtige Ergebnistexte
- keine dogmatische Klassifikationslogik

---

## Sensitivity und Schutz
### Mindestfelder
```text
is_sensitive
age_gate
review_status
```

### Erweiterte sensible Marker
- `worldview_sensitivity`
- `meaning_depth`
- `existential_load_level`

### Sinn
- kontrollierte Ausspielung
- geschützte Ergebnisse
- zusätzliche Freigaben
- spätere Policylogik

---

## Abschluss
Das Datenmodell V3 Final schafft die Grundlage dafür, dass CPE zugleich präzise, flexibel, mehrsprachig, sinnorientiert, zertifikatsfähig und analyseoffen bleiben kann.

