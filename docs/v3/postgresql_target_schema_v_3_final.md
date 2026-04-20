# postgresql-target-schema-v3-final.md

## Zweck
Dieses Dokument beschreibt das strategische PostgreSQL-Zielschema in der V3-Endfassung für die CPE-/CEE4AI-Produktfamilie.

Es ist das fachlich-technische Zielmodell für:
- Inhalte und Übersetzungen
- Profiling und Meaning-Achsen
- Runtime und Events
- Zertifikate und Trust
- Governance und Audit

---

## Leitprinzip
PostgreSQL ist das **System of Record** für:
- kanonische Inhalte
- Profile und Resultate
- Sessionläufe und Pakete
- Zertifikatsregeln und ausgestellte Credentials
- Review-, Governance- und Audit-Daten

Analytische Massenauswertungen werden **ergänzend** gedacht, nicht als Ersatz des Wahrheitssystems.

---

## Logische Schemas
- `core`
- `content`
- `profiling`
- `runtime`
- `review`
- `credentials`
- `governance`
- `events`

---

## Schema `core`
### `languages`
- id
- code
- name
- script_code
- is_active

### `regions`
- id
- code
- name
- is_active

### `domains`
- id
- slug
- name
- description
- is_active

### `users`
- id
- external_ref
- locale_language_id
- locale_region_id
- created_at
- updated_at

### `roles`
- id
- slug
- name

### `user_roles`
- user_id
- role_id
- assigned_at

---

## Schema `content`
### `question_master`
- id
- external_id
- domain_id
- question_family
- category_id
- subcategory_id
- parent_question_id
- question_type
- scoring_mode
- intended_use
- confidence_tier
- estimated_time_seconds
- cognitive_load_level
- cultural_scope
- region_scope
- meaning_depth
- worldview_sensitivity
- symbolic_interpretation_relevance
- existential_load_level
- is_sensitive
- age_gate
- review_status
- is_active
- version
- created_at
- updated_at

### `question_translation`
- id
- question_id
- language_id
- region_id nullable
- localization_status
- requires_human_review
- worldview_sensitive
- title nullable
- question_text
- explanation_text
- reviewer_notes nullable
- is_active
- version
- created_at
- updated_at

### `question_option_master`
- id
- question_id
- option_key
- score_weight
- is_correct
- display_order
- is_active
- version

### `question_option_translation`
- id
- option_id
- language_id
- region_id nullable
- option_text
- localization_status
- is_active

### `categories`
- id
- domain_id
- slug
- name
- description
- is_sensitive
- is_active

### `subcategories`
- id
- category_id
- slug
- name
- description
- is_sensitive
- is_active

### `result_text_master`
- id
- domain_id
- result_type
- certainty_level
- profile_depth
- slug
- is_sensitive
- is_active

### `result_text_translation`
- id
- result_text_id
- language_id
- region_id nullable
- title
- body
- localization_status
- is_active

---

## Schema `profiling`
### `denktypes`
- id
- slug
- name
- description
- development_hint
- is_active

### `skills`
- id
- slug
- name
- description
- is_active

### `personality_traits`
- id
- slug
- name
- description
- polarity_model nullable
- is_active

### `interest_tags`
- id
- slug
- name
- description
- is_active

### `meaning_tags`
- id
- slug
- name
- description
- sensitivity_level
- is_active

### `worldview_frames`
- id
- slug
- name
- description
- sensitivity_level
- is_active

### `trigger_groups`
- id
- slug
- name
- description
- is_sensitive
- is_active

### `body_references`
- id
- slug
- name
- description
- is_sensitive
- is_active

### `development_paths`
- id
- slug
- name
- description
- target_domain_id nullable
- is_active

### mapping tables
- `question_denktype_tags`
- `question_skill_tags`
- `question_trait_tags`
- `question_interest_tags`
- `question_meaning_tags`
- `question_worldview_tags`
- `question_trigger_tags`
- `question_body_reference_tags`
- `question_path_tags`

Jede Mapping-Tabelle enthält mindestens:
- question_id
- referenced_id
- weight
- rationale nullable
- is_active

---

## Schema `runtime`
### `sessions`
- id uuid
- user_id nullable
- domain_id
- mode
- session_goal
- locale_language_id
- locale_region_id
- result_confidence
- progress_state
- started_at
- finished_at nullable
- created_at
- updated_at

### `session_packages`
- id
- session_id
- package_index
- package_size
- estimated_time_seconds
- actual_time_seconds nullable
- completion_quality nullable
- aborted_at nullable
- continuation_window_until nullable
- recommended_next_mode nullable
- created_at

### `answers`
- id
- session_id
- package_id nullable
- question_id
- answer_kind
- selected_option_ids jsonb nullable
- scale_value nullable
- free_text_answer nullable
- raw_score nullable
- evaluated_score nullable
- certainty_level nullable
- answered_at
- created_at

### `result_snapshots`
- id
- session_id
- result_type
- profile_depth
- certainty_level
- snapshot_payload jsonb
- ruleset_version
- created_at

### `profile_vectors`
- id
- session_id
- vector_type
- vector_payload jsonb
- created_at

---

## Schema `review`
### `review_flags`
- id
- slug
- name
- description
- default_severity
- is_active

### `question_reviews`
- id
- question_id
- session_id nullable
- reviewer_user_id nullable
- reviewer_role
- flag_id
- comment nullable
- severity
- created_at

### `review_decisions`
- id
- question_id
- decided_by_user_id
- old_status
- new_status
- reason
- created_at

### `localization_reviews`
- id
- translation_id
- reviewer_user_id
- status
- comment
- created_at

---

## Schema `credentials`
### `competency_frameworks`
- id
- slug
- name
- description
- version
- is_active

### `assessment_programs`
- id
- slug
- name
- description
- domain_id
- framework_id nullable
- issuer_id
- is_active

### `assessment_tracks`
- id
- program_id
- slug
- name
- description
- required_profile_depth
- required_certainty_level
- is_active

### `certificate_templates`
- id
- slug
- name
- template_payload jsonb
- output_format
- is_active

### `certificate_rules`
- id
- track_id
- ruleset_version
- minimum_completion
- minimum_confidence
- minimum_integrity_score nullable
- requires_identity_verification boolean
- requires_proctoring boolean
- validity_days nullable
- is_active

### `credential_issuers`
- id
- slug
- legal_name
- public_name
- verification_base_url
- trust_level
- is_active

### `credential_assertions`
- id uuid
- user_id
- session_id nullable
- track_id
- rule_id
- issuer_id
- certificate_level
- issued_under_ruleset_version
- result_confidence
- assessment_integrity_score nullable
- evidence_hash
- credential_standard
- standards_alignment jsonb
- verification_token
- issued_at
- expires_at nullable
- revoked_at nullable
- revocation_reason nullable

### `credential_evidence`
- id
- assertion_id
- evidence_type
- evidence_ref
- evidence_hash
- created_at

### `credential_verifications`
- id
- assertion_id
- verifier_ref nullable
- verification_result
- verified_at

### `credential_revocations`
- id
- assertion_id
- reason
- revoked_by_user_id nullable
- revoked_at

---

## Schema `governance`
### `content_policies`
- id
- slug
- name
- description
- applies_to_domain_id nullable
- applies_to_region_id nullable
- is_active

### `sensitivity_rules`
- id
- slug
- name
- description
- policy_payload jsonb
- is_active

### `ruleset_versions`
- id
- domain_id nullable
- slug
- version
- description
- effective_from
- effective_to nullable
- is_active

### `audit_log`
- id
- actor_user_id nullable
- entity_type
- entity_id
- action
- payload jsonb
- created_at

---

## Schema `events`
### `event_outbox`
- id uuid
- aggregate_type
- aggregate_id
- event_type
- event_version
- payload jsonb
- occurred_at
- published_at nullable

### `event_inbox`
- id uuid
- source
- event_type
- event_version
- payload jsonb
- received_at
- processed_at nullable

### `event_store`
- id bigserial
- event_id uuid
- event_type
- subject_type
- subject_id
- session_id nullable
- user_id nullable
- occurred_at
- payload jsonb
- partition_key

---

## Index- und Performance-Grundsätze
- FK-Indizes auf alle Verweise
- Indizes auf `slug`, `external_id`, `verification_token`
- GIN auf `jsonb`-Payloads, wo flexible Filter nötig sind
- Volltextindex für redaktionelle Suche auf Übersetzungen
- partielle Indizes auf aktive Datensätze

### Partitionierung zuerst für
- `event_store`
- `answers`
- `result_snapshots`
- `audit_log`
- später `credential_verifications`

---

## Besonderheiten für Meaning-/Spiritual-Inhalte
- Human Review bevorzugt verpflichtend
- `localization_status` nur nach Review auf `approved`
- weltanschauliche Sensitivität im Schema sichtbar
- Journey-/Completion-Zertifikate statt metaphysischer Stufenbehauptungen

---

## Abschluss
Diese V3-Endfassung des PostgreSQL-Zielschemas hält Profiling, Persönlichkeit, Meaning, Coaching, Mehrsprachigkeit, Zertifikate und Governance unter einem gemeinsamen, belastbaren relationalen Dach.

