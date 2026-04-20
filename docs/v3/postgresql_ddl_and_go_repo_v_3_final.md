# postgresql-ddl-and-go-repo-v3-final.md

## Zweck
Dieses Dokument liefert die beiden angefragten technischen Endbausteine für die CPE-/CEE4AI-V3-Architektur:

1. **PostgreSQL-Zielschema als echte DDL-/Migrationsfassung**  
2. **Go-Repo-Struktur auf Final-V3 umgestellt**

Es ist bewusst so aufgebaut, dass es als direkte Grundlage für:
- ein Repo-Setup
- echte Migrationsdateien
- die erste Postgres-Integration mit Go
- die modulare Weiterentwicklung der Produktfamilie

dienen kann.

---

# Teil A – PostgreSQL-DLL / Migrationsfassung

## Migrationsstrategie
Empfohlen ist eine sequenzielle Struktur, z. B.:

```text
migrations/
  0001_extensions_and_schemas.sql
  0002_core.sql
  0003_content.sql
  0004_profiling.sql
  0005_runtime.sql
  0006_review.sql
  0007_credentials.sql
  0008_governance_and_events.sql
  0009_indexes.sql
```

---

## `0001_extensions_and_schemas.sql`
```sql
BEGIN;

CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE SCHEMA IF NOT EXISTS core;
CREATE SCHEMA IF NOT EXISTS content;
CREATE SCHEMA IF NOT EXISTS profiling;
CREATE SCHEMA IF NOT EXISTS runtime;
CREATE SCHEMA IF NOT EXISTS review;
CREATE SCHEMA IF NOT EXISTS credentials;
CREATE SCHEMA IF NOT EXISTS governance;
CREATE SCHEMA IF NOT EXISTS events;

COMMIT;
```

---

## `0002_core.sql`
```sql
BEGIN;

CREATE TABLE core.languages (
    id              BIGSERIAL PRIMARY KEY,
    code            TEXT NOT NULL UNIQUE,
    name            TEXT NOT NULL,
    script_code     TEXT,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE core.regions (
    id              BIGSERIAL PRIMARY KEY,
    code            TEXT NOT NULL UNIQUE,
    name            TEXT NOT NULL,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE core.domains (
    id              BIGSERIAL PRIMARY KEY,
    slug            CITEXT NOT NULL UNIQUE,
    name            TEXT NOT NULL,
    description     TEXT,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE core.users (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    external_ref            TEXT UNIQUE,
    locale_language_id      BIGINT REFERENCES core.languages(id),
    locale_region_id        BIGINT REFERENCES core.regions(id),
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE core.roles (
    id              BIGSERIAL PRIMARY KEY,
    slug            CITEXT NOT NULL UNIQUE,
    name            TEXT NOT NULL
);

CREATE TABLE core.user_roles (
    user_id         UUID NOT NULL REFERENCES core.users(id) ON DELETE CASCADE,
    role_id         BIGINT NOT NULL REFERENCES core.roles(id) ON DELETE CASCADE,
    assigned_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, role_id)
);

COMMIT;
```

---

## `0003_content.sql`
```sql
BEGIN;

CREATE TABLE content.categories (
    id                              BIGSERIAL PRIMARY KEY,
    domain_id                       BIGINT NOT NULL REFERENCES core.domains(id),
    slug                            CITEXT NOT NULL UNIQUE,
    name                            TEXT NOT NULL,
    description                     TEXT,
    default_sensitivity             TEXT,
    cultural_scope                  TEXT,
    spiritual_relevance             TEXT,
    worldview_sensitivity           TEXT,
    meaning_pathway_relevance       TEXT,
    is_sensitive                    BOOLEAN NOT NULL DEFAULT FALSE,
    is_active                       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at                      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE content.subcategories (
    id                  BIGSERIAL PRIMARY KEY,
    category_id         BIGINT NOT NULL REFERENCES content.categories(id) ON DELETE CASCADE,
    slug                CITEXT NOT NULL UNIQUE,
    name                TEXT NOT NULL,
    description         TEXT,
    is_sensitive        BOOLEAN NOT NULL DEFAULT FALSE,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE content.question_master (
    id                                  BIGSERIAL PRIMARY KEY,
    external_id                         TEXT NOT NULL UNIQUE,
    domain_id                           BIGINT NOT NULL REFERENCES core.domains(id),
    question_family                     TEXT NOT NULL,
    category_id                         BIGINT NOT NULL REFERENCES content.categories(id),
    subcategory_id                      BIGINT REFERENCES content.subcategories(id),
    parent_question_id                  BIGINT REFERENCES content.question_master(id),
    question_type                       TEXT NOT NULL,
    scoring_mode                        TEXT NOT NULL,
    intended_use                        TEXT,
    confidence_tier                     TEXT,
    estimated_time_seconds              INTEGER,
    cognitive_load_level                TEXT,
    cultural_scope                      TEXT,
    region_scope                        TEXT,
    meaning_depth                       TEXT,
    worldview_sensitivity               TEXT,
    symbolic_interpretation_relevance   TEXT,
    existential_load_level              TEXT,
    is_sensitive                        BOOLEAN NOT NULL DEFAULT FALSE,
    age_gate                            INTEGER NOT NULL DEFAULT 0,
    review_status                       TEXT NOT NULL DEFAULT 'draft',
    is_active                           BOOLEAN NOT NULL DEFAULT TRUE,
    version                             INTEGER NOT NULL DEFAULT 1,
    created_at                          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT question_type_chk CHECK (
        question_type IN (
            'single_choice','multiple_select','scale','reflection',
            'reflective','contemplative','comparative_worldview',
            'symbolic_interpretation','existential'
        )
    ),
    CONSTRAINT scoring_mode_chk CHECK (
        scoring_mode IN ('exact','partial','weighted','non_scored','path_only')
    )
);

CREATE TABLE content.question_translation (
    id                      BIGSERIAL PRIMARY KEY,
    question_id             BIGINT NOT NULL REFERENCES content.question_master(id) ON DELETE CASCADE,
    language_id             BIGINT NOT NULL REFERENCES core.languages(id),
    region_id               BIGINT REFERENCES core.regions(id),
    localization_status     TEXT NOT NULL DEFAULT 'draft',
    requires_human_review   BOOLEAN NOT NULL DEFAULT FALSE,
    worldview_sensitive     BOOLEAN NOT NULL DEFAULT FALSE,
    title                   TEXT,
    question_text           TEXT NOT NULL,
    explanation_text        TEXT,
    reviewer_notes          TEXT,
    is_active               BOOLEAN NOT NULL DEFAULT TRUE,
    version                 INTEGER NOT NULL DEFAULT 1,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (question_id, language_id, region_id, version)
);

CREATE TABLE content.question_option_master (
    id                  BIGSERIAL PRIMARY KEY,
    question_id         BIGINT NOT NULL REFERENCES content.question_master(id) ON DELETE CASCADE,
    option_key          TEXT NOT NULL,
    score_weight        NUMERIC(8,4) NOT NULL DEFAULT 1.0,
    is_correct          BOOLEAN NOT NULL DEFAULT FALSE,
    display_order       INTEGER NOT NULL,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,
    version             INTEGER NOT NULL DEFAULT 1,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (question_id, option_key, version)
);

CREATE TABLE content.question_option_translation (
    id                      BIGSERIAL PRIMARY KEY,
    option_id               BIGINT NOT NULL REFERENCES content.question_option_master(id) ON DELETE CASCADE,
    language_id             BIGINT NOT NULL REFERENCES core.languages(id),
    region_id               BIGINT REFERENCES core.regions(id),
    option_text             TEXT NOT NULL,
    localization_status     TEXT NOT NULL DEFAULT 'draft',
    is_active               BOOLEAN NOT NULL DEFAULT TRUE,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (option_id, language_id, region_id)
);

CREATE TABLE content.result_text_master (
    id                  BIGSERIAL PRIMARY KEY,
    domain_id           BIGINT NOT NULL REFERENCES core.domains(id),
    result_type         TEXT NOT NULL,
    certainty_level     TEXT,
    profile_depth       TEXT,
    slug                CITEXT NOT NULL UNIQUE,
    is_sensitive        BOOLEAN NOT NULL DEFAULT FALSE,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE content.result_text_translation (
    id                      BIGSERIAL PRIMARY KEY,
    result_text_id          BIGINT NOT NULL REFERENCES content.result_text_master(id) ON DELETE CASCADE,
    language_id             BIGINT NOT NULL REFERENCES core.languages(id),
    region_id               BIGINT REFERENCES core.regions(id),
    title                   TEXT NOT NULL,
    body                    TEXT NOT NULL,
    localization_status     TEXT NOT NULL DEFAULT 'draft',
    is_active               BOOLEAN NOT NULL DEFAULT TRUE,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (result_text_id, language_id, region_id)
);

COMMIT;
```

---

## `0004_profiling.sql`
```sql
BEGIN;

CREATE TABLE profiling.denktypes (
    id                  BIGSERIAL PRIMARY KEY,
    slug                CITEXT NOT NULL UNIQUE,
    name                TEXT NOT NULL,
    description         TEXT,
    development_hint    TEXT,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE profiling.skills (
    id                  BIGSERIAL PRIMARY KEY,
    slug                CITEXT NOT NULL UNIQUE,
    name                TEXT NOT NULL,
    description         TEXT,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE profiling.personality_traits (
    id                  BIGSERIAL PRIMARY KEY,
    slug                CITEXT NOT NULL UNIQUE,
    name                TEXT NOT NULL,
    description         TEXT,
    polarity_model      TEXT,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE profiling.interest_tags (
    id                  BIGSERIAL PRIMARY KEY,
    slug                CITEXT NOT NULL UNIQUE,
    name                TEXT NOT NULL,
    description         TEXT,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE profiling.meaning_tags (
    id                  BIGSERIAL PRIMARY KEY,
    slug                CITEXT NOT NULL UNIQUE,
    name                TEXT NOT NULL,
    description         TEXT,
    sensitivity_level   TEXT,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE profiling.worldview_frames (
    id                  BIGSERIAL PRIMARY KEY,
    slug                CITEXT NOT NULL UNIQUE,
    name                TEXT NOT NULL,
    description         TEXT,
    sensitivity_level   TEXT,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE profiling.trigger_groups (
    id                  BIGSERIAL PRIMARY KEY,
    slug                CITEXT NOT NULL UNIQUE,
    name                TEXT NOT NULL,
    description         TEXT,
    is_sensitive        BOOLEAN NOT NULL DEFAULT FALSE,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE profiling.body_references (
    id                  BIGSERIAL PRIMARY KEY,
    slug                CITEXT NOT NULL UNIQUE,
    name                TEXT NOT NULL,
    description         TEXT,
    is_sensitive        BOOLEAN NOT NULL DEFAULT FALSE,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE profiling.development_paths (
    id                  BIGSERIAL PRIMARY KEY,
    slug                CITEXT NOT NULL UNIQUE,
    name                TEXT NOT NULL,
    description         TEXT,
    target_domain_id    BIGINT REFERENCES core.domains(id),
    is_active           BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE profiling.question_denktype_tags (
    question_id         BIGINT NOT NULL REFERENCES content.question_master(id) ON DELETE CASCADE,
    denktype_id         BIGINT NOT NULL REFERENCES profiling.denktypes(id) ON DELETE CASCADE,
    weight              NUMERIC(8,4) NOT NULL DEFAULT 1.0,
    rationale           TEXT,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,
    PRIMARY KEY (question_id, denktype_id)
);

CREATE TABLE profiling.question_skill_tags (
    question_id         BIGINT NOT NULL REFERENCES content.question_master(id) ON DELETE CASCADE,
    skill_id            BIGINT NOT NULL REFERENCES profiling.skills(id) ON DELETE CASCADE,
    weight              NUMERIC(8,4) NOT NULL DEFAULT 1.0,
    rationale           TEXT,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,
    PRIMARY KEY (question_id, skill_id)
);

CREATE TABLE profiling.question_trait_tags (
    question_id         BIGINT NOT NULL REFERENCES content.question_master(id) ON DELETE CASCADE,
    trait_id            BIGINT NOT NULL REFERENCES profiling.personality_traits(id) ON DELETE CASCADE,
    weight              NUMERIC(8,4) NOT NULL DEFAULT 1.0,
    rationale           TEXT,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,
    PRIMARY KEY (question_id, trait_id)
);

CREATE TABLE profiling.question_interest_tags (
    question_id         BIGINT NOT NULL REFERENCES content.question_master(id) ON DELETE CASCADE,
    interest_tag_id     BIGINT NOT NULL REFERENCES profiling.interest_tags(id) ON DELETE CASCADE,
    weight              NUMERIC(8,4) NOT NULL DEFAULT 1.0,
    rationale           TEXT,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,
    PRIMARY KEY (question_id, interest_tag_id)
);

CREATE TABLE profiling.question_meaning_tags (
    question_id         BIGINT NOT NULL REFERENCES content.question_master(id) ON DELETE CASCADE,
    meaning_tag_id      BIGINT NOT NULL REFERENCES profiling.meaning_tags(id) ON DELETE CASCADE,
    weight              NUMERIC(8,4) NOT NULL DEFAULT 1.0,
    rationale           TEXT,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,
    PRIMARY KEY (question_id, meaning_tag_id)
);

CREATE TABLE profiling.question_worldview_tags (
    question_id         BIGINT NOT NULL REFERENCES content.question_master(id) ON DELETE CASCADE,
    worldview_frame_id  BIGINT NOT NULL REFERENCES profiling.worldview_frames(id) ON DELETE CASCADE,
    weight              NUMERIC(8,4) NOT NULL DEFAULT 1.0,
    rationale           TEXT,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,
    PRIMARY KEY (question_id, worldview_frame_id)
);

CREATE TABLE profiling.question_trigger_tags (
    question_id         BIGINT NOT NULL REFERENCES content.question_master(id) ON DELETE CASCADE,
    trigger_group_id    BIGINT NOT NULL REFERENCES profiling.trigger_groups(id) ON DELETE CASCADE,
    weight              NUMERIC(8,4) NOT NULL DEFAULT 1.0,
    rationale           TEXT,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,
    PRIMARY KEY (question_id, trigger_group_id)
);

CREATE TABLE profiling.question_body_reference_tags (
    question_id         BIGINT NOT NULL REFERENCES content.question_master(id) ON DELETE CASCADE,
    body_reference_id   BIGINT NOT NULL REFERENCES profiling.body_references(id) ON DELETE CASCADE,
    weight              NUMERIC(8,4) NOT NULL DEFAULT 1.0,
    rationale           TEXT,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,
    PRIMARY KEY (question_id, body_reference_id)
);

CREATE TABLE profiling.question_path_tags (
    question_id         BIGINT NOT NULL REFERENCES content.question_master(id) ON DELETE CASCADE,
    development_path_id BIGINT NOT NULL REFERENCES profiling.development_paths(id) ON DELETE CASCADE,
    weight              NUMERIC(8,4) NOT NULL DEFAULT 1.0,
    rationale           TEXT,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,
    PRIMARY KEY (question_id, development_path_id)
);

COMMIT;
```

---

## `0005_runtime.sql`
```sql
BEGIN;

CREATE TABLE runtime.sessions (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id                 UUID REFERENCES core.users(id) ON DELETE SET NULL,
    domain_id               BIGINT NOT NULL REFERENCES core.domains(id),
    mode                    TEXT NOT NULL,
    session_goal            TEXT,
    locale_language_id      BIGINT REFERENCES core.languages(id),
    locale_region_id        BIGINT REFERENCES core.regions(id),
    result_confidence       NUMERIC(8,4),
    progress_state          TEXT,
    started_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    finished_at             TIMESTAMPTZ,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE runtime.session_packages (
    id                          BIGSERIAL PRIMARY KEY,
    session_id                  UUID NOT NULL REFERENCES runtime.sessions(id) ON DELETE CASCADE,
    package_index               INTEGER NOT NULL,
    package_size                INTEGER NOT NULL,
    estimated_time_seconds      INTEGER,
    actual_time_seconds         INTEGER,
    completion_quality          NUMERIC(8,4),
    aborted_at                  TIMESTAMPTZ,
    continuation_window_until   TIMESTAMPTZ,
    recommended_next_mode       TEXT,
    created_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (session_id, package_index)
);

CREATE TABLE runtime.answers (
    id                      BIGSERIAL PRIMARY KEY,
    session_id              UUID NOT NULL REFERENCES runtime.sessions(id) ON DELETE CASCADE,
    package_id              BIGINT REFERENCES runtime.session_packages(id) ON DELETE SET NULL,
    question_id             BIGINT NOT NULL REFERENCES content.question_master(id),
    answer_kind             TEXT NOT NULL,
    selected_option_ids     JSONB,
    scale_value             INTEGER,
    free_text_answer        TEXT,
    raw_score               NUMERIC(10,4),
    evaluated_score         NUMERIC(10,4),
    certainty_level         TEXT,
    answered_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE runtime.result_snapshots (
    id                      BIGSERIAL PRIMARY KEY,
    session_id              UUID NOT NULL REFERENCES runtime.sessions(id) ON DELETE CASCADE,
    result_type             TEXT NOT NULL,
    profile_depth           TEXT,
    certainty_level         TEXT,
    snapshot_payload        JSONB NOT NULL,
    ruleset_version         TEXT,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE runtime.profile_vectors (
    id                      BIGSERIAL PRIMARY KEY,
    session_id              UUID NOT NULL REFERENCES runtime.sessions(id) ON DELETE CASCADE,
    vector_type             TEXT NOT NULL,
    vector_payload          JSONB NOT NULL,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
```

---

## `0006_review.sql`
```sql
BEGIN;

CREATE TABLE review.review_flags (
    id                  BIGSERIAL PRIMARY KEY,
    slug                CITEXT NOT NULL UNIQUE,
    name                TEXT NOT NULL,
    description         TEXT,
    default_severity    TEXT,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE review.question_reviews (
    id                  BIGSERIAL PRIMARY KEY,
    question_id         BIGINT NOT NULL REFERENCES content.question_master(id) ON DELETE CASCADE,
    session_id          UUID REFERENCES runtime.sessions(id) ON DELETE SET NULL,
    reviewer_user_id    UUID REFERENCES core.users(id) ON DELETE SET NULL,
    reviewer_role       TEXT,
    flag_id             BIGINT NOT NULL REFERENCES review.review_flags(id),
    comment             TEXT,
    severity            TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE review.review_decisions (
    id                  BIGSERIAL PRIMARY KEY,
    question_id         BIGINT NOT NULL REFERENCES content.question_master(id) ON DELETE CASCADE,
    decided_by_user_id  UUID REFERENCES core.users(id) ON DELETE SET NULL,
    old_status          TEXT,
    new_status          TEXT,
    reason              TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE review.localization_reviews (
    id                  BIGSERIAL PRIMARY KEY,
    translation_id      BIGINT NOT NULL REFERENCES content.question_translation(id) ON DELETE CASCADE,
    reviewer_user_id    UUID REFERENCES core.users(id) ON DELETE SET NULL,
    status              TEXT NOT NULL,
    comment             TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
```

---

## `0007_credentials.sql`
```sql
BEGIN;

CREATE TABLE credentials.competency_frameworks (
    id                  BIGSERIAL PRIMARY KEY,
    slug                CITEXT NOT NULL UNIQUE,
    name                TEXT NOT NULL,
    description         TEXT,
    version             TEXT NOT NULL,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE credentials.credential_issuers (
    id                      BIGSERIAL PRIMARY KEY,
    slug                    CITEXT NOT NULL UNIQUE,
    legal_name              TEXT NOT NULL,
    public_name             TEXT NOT NULL,
    verification_base_url   TEXT,
    trust_level             TEXT,
    is_active               BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE credentials.assessment_programs (
    id                  BIGSERIAL PRIMARY KEY,
    slug                CITEXT NOT NULL UNIQUE,
    name                TEXT NOT NULL,
    description         TEXT,
    domain_id           BIGINT NOT NULL REFERENCES core.domains(id),
    framework_id        BIGINT REFERENCES credentials.competency_frameworks(id),
    issuer_id           BIGINT NOT NULL REFERENCES credentials.credential_issuers(id),
    is_active           BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE credentials.assessment_tracks (
    id                          BIGSERIAL PRIMARY KEY,
    program_id                  BIGINT NOT NULL REFERENCES credentials.assessment_programs(id) ON DELETE CASCADE,
    slug                        CITEXT NOT NULL UNIQUE,
    name                        TEXT NOT NULL,
    description                 TEXT,
    required_profile_depth      TEXT,
    required_certainty_level    TEXT,
    is_active                   BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE credentials.certificate_templates (
    id                  BIGSERIAL PRIMARY KEY,
    slug                CITEXT NOT NULL UNIQUE,
    name                TEXT NOT NULL,
    template_payload    JSONB NOT NULL,
    output_format       TEXT NOT NULL,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE credentials.certificate_rules (
    id                              BIGSERIAL PRIMARY KEY,
    track_id                        BIGINT NOT NULL REFERENCES credentials.assessment_tracks(id) ON DELETE CASCADE,
    ruleset_version                 TEXT NOT NULL,
    minimum_completion              NUMERIC(8,4),
    minimum_confidence              NUMERIC(8,4),
    minimum_integrity_score         NUMERIC(8,4),
    requires_identity_verification  BOOLEAN NOT NULL DEFAULT FALSE,
    requires_proctoring             BOOLEAN NOT NULL DEFAULT FALSE,
    validity_days                   INTEGER,
    is_active                       BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE credentials.credential_assertions (
    id                              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id                         UUID REFERENCES core.users(id) ON DELETE SET NULL,
    session_id                      UUID REFERENCES runtime.sessions(id) ON DELETE SET NULL,
    track_id                        BIGINT NOT NULL REFERENCES credentials.assessment_tracks(id),
    rule_id                         BIGINT NOT NULL REFERENCES credentials.certificate_rules(id),
    issuer_id                       BIGINT NOT NULL REFERENCES credentials.credential_issuers(id),
    certificate_level               TEXT NOT NULL,
    issued_under_ruleset_version    TEXT NOT NULL,
    result_confidence               NUMERIC(8,4),
    assessment_integrity_score      NUMERIC(8,4),
    evidence_hash                   TEXT NOT NULL,
    credential_standard             TEXT,
    standards_alignment             JSONB,
    verification_token              TEXT NOT NULL UNIQUE,
    issued_at                       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at                      TIMESTAMPTZ,
    revoked_at                      TIMESTAMPTZ,
    revocation_reason               TEXT
);

CREATE TABLE credentials.credential_evidence (
    id                  BIGSERIAL PRIMARY KEY,
    assertion_id        UUID NOT NULL REFERENCES credentials.credential_assertions(id) ON DELETE CASCADE,
    evidence_type       TEXT NOT NULL,
    evidence_ref        TEXT,
    evidence_hash       TEXT NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE credentials.credential_verifications (
    id                  BIGSERIAL PRIMARY KEY,
    assertion_id        UUID NOT NULL REFERENCES credentials.credential_assertions(id) ON DELETE CASCADE,
    verifier_ref        TEXT,
    verification_result TEXT NOT NULL,
    verified_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE credentials.credential_revocations (
    id                  BIGSERIAL PRIMARY KEY,
    assertion_id        UUID NOT NULL REFERENCES credentials.credential_assertions(id) ON DELETE CASCADE,
    reason              TEXT NOT NULL,
    revoked_by_user_id  UUID REFERENCES core.users(id) ON DELETE SET NULL,
    revoked_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
```

---

## `0008_governance_and_events.sql`
```sql
BEGIN;

CREATE TABLE governance.content_policies (
    id                      BIGSERIAL PRIMARY KEY,
    slug                    CITEXT NOT NULL UNIQUE,
    name                    TEXT NOT NULL,
    description             TEXT,
    applies_to_domain_id    BIGINT REFERENCES core.domains(id),
    applies_to_region_id    BIGINT REFERENCES core.regions(id),
    is_active               BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE governance.sensitivity_rules (
    id                  BIGSERIAL PRIMARY KEY,
    slug                CITEXT NOT NULL UNIQUE,
    name                TEXT NOT NULL,
    description         TEXT,
    policy_payload      JSONB NOT NULL,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE governance.ruleset_versions (
    id                  BIGSERIAL PRIMARY KEY,
    domain_id           BIGINT REFERENCES core.domains(id),
    slug                CITEXT NOT NULL,
    version             TEXT NOT NULL,
    description         TEXT,
    effective_from      TIMESTAMPTZ NOT NULL,
    effective_to        TIMESTAMPTZ,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,
    UNIQUE (slug, version)
);

CREATE TABLE governance.audit_log (
    id                  BIGSERIAL PRIMARY KEY,
    actor_user_id       UUID REFERENCES core.users(id) ON DELETE SET NULL,
    entity_type         TEXT NOT NULL,
    entity_id           TEXT NOT NULL,
    action              TEXT NOT NULL,
    payload             JSONB,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE events.event_outbox (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    aggregate_type      TEXT NOT NULL,
    aggregate_id        TEXT NOT NULL,
    event_type          TEXT NOT NULL,
    event_version       INTEGER NOT NULL DEFAULT 1,
    payload             JSONB NOT NULL,
    occurred_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    published_at        TIMESTAMPTZ
);

CREATE TABLE events.event_inbox (
    id                  UUID PRIMARY KEY,
    source              TEXT NOT NULL,
    event_type          TEXT NOT NULL,
    event_version       INTEGER NOT NULL DEFAULT 1,
    payload             JSONB NOT NULL,
    received_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at        TIMESTAMPTZ
);

CREATE TABLE events.event_store (
    id                  BIGSERIAL PRIMARY KEY,
    event_id            UUID NOT NULL UNIQUE,
    event_type          TEXT NOT NULL,
    subject_type        TEXT NOT NULL,
    subject_id          TEXT NOT NULL,
    session_id          UUID,
    user_id             UUID,
    occurred_at         TIMESTAMPTZ NOT NULL,
    payload             JSONB NOT NULL,
    partition_key       TEXT NOT NULL
);

COMMIT;
```

---

## `0009_indexes.sql`
```sql
BEGIN;

CREATE INDEX idx_question_master_domain ON content.question_master(domain_id);
CREATE INDEX idx_question_master_category ON content.question_master(category_id);
CREATE INDEX idx_question_master_active_review ON content.question_master(is_active, review_status);
CREATE INDEX idx_question_translation_question_lang_region ON content.question_translation(question_id, language_id, region_id);
CREATE INDEX idx_result_text_translation_lang_region ON content.result_text_translation(result_text_id, language_id, region_id);

CREATE INDEX idx_sessions_user ON runtime.sessions(user_id);
CREATE INDEX idx_sessions_domain_mode ON runtime.sessions(domain_id, mode);
CREATE INDEX idx_session_packages_session_idx ON runtime.session_packages(session_id, package_index);
CREATE INDEX idx_answers_session ON runtime.answers(session_id);
CREATE INDEX idx_answers_question ON runtime.answers(question_id);
CREATE INDEX idx_answers_selected_options_gin ON runtime.answers USING GIN (selected_option_ids);
CREATE INDEX idx_result_snapshots_session ON runtime.result_snapshots(session_id);
CREATE INDEX idx_profile_vectors_session ON runtime.profile_vectors(session_id);

CREATE INDEX idx_question_reviews_question ON review.question_reviews(question_id);
CREATE INDEX idx_question_reviews_flag ON review.question_reviews(flag_id);
CREATE INDEX idx_localization_reviews_translation ON review.localization_reviews(translation_id);

CREATE INDEX idx_credential_assertions_user ON credentials.credential_assertions(user_id);
CREATE INDEX idx_credential_assertions_session ON credentials.credential_assertions(session_id);
CREATE INDEX idx_credential_assertions_track ON credentials.credential_assertions(track_id);
CREATE INDEX idx_credential_assertions_issuer ON credentials.credential_assertions(issuer_id);

CREATE INDEX idx_audit_log_entity ON governance.audit_log(entity_type, entity_id);

CREATE INDEX idx_event_outbox_unpublished ON events.event_outbox(published_at) WHERE published_at IS NULL;
CREATE INDEX idx_event_store_type_time ON events.event_store(event_type, occurred_at DESC);
CREATE INDEX idx_event_store_subject ON events.event_store(subject_type, subject_id);
CREATE INDEX idx_event_store_payload_gin ON events.event_store USING GIN (payload);

COMMIT;
```

---

# Teil B – Go-Repo-Struktur auf Final-V3 umstellen

## Zielbild
Die Repo-Struktur soll jetzt:
- PostgreSQL-first sein
- die Product Family tragen
- Event- und Credential-Layer nativ berücksichtigen
- Meaning-/Journey-Pfade nicht nachträglich hineinquetschen

```text
cee4ai/
├── cmd/
│   ├── api/
│   │   └── main.go
│   ├── migrate/
│   │   └── main.go
│   └── seed/
│       └── main.go
├── internal/
│   ├── app/
│   │   └── app.go
│   ├── config/
│   │   └── config.go
│   ├── db/
│   │   ├── postgres.go
│   │   ├── migrations.go
│   │   └── tx.go
│   ├── domains/
│   │   ├── cognitive/
│   │   │   ├── service.go
│   │   │   └── result_texts.go
│   │   ├── personality/
│   │   │   ├── service.go
│   │   │   └── traits.go
│   │   ├── coaching/
│   │   │   ├── service.go
│   │   │   └── safeguards.go
│   │   ├── pathways/
│   │   │   ├── service.go
│   │   │   └── recommendations.go
│   │   └── meaning/
│   │       ├── service.go
│   │       ├── worldview.go
│   │       ├── journey_texts.go
│   │       └── safeguards.go
│   ├── questions/
│   │   ├── model.go
│   │   ├── repository.go
│   │   ├── service.go
│   │   ├── translations.go
│   │   ├── options.go
│   │   └── validation.go
│   ├── profiling/
│   │   ├── model.go
│   │   ├── repository.go
│   │   ├── tags.go
│   │   ├── traits.go
│   │   ├── skills.go
│   │   └── meaning.go
│   ├── sessions/
│   │   ├── model.go
│   │   ├── repository.go
│   │   └── service.go
│   ├── packages/
│   │   ├── model.go
│   │   ├── repository.go
│   │   └── service.go
│   ├── answers/
│   │   ├── model.go
│   │   ├── repository.go
│   │   └── service.go
│   ├── scoring/
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
│   ├── credentials/
│   │   ├── model.go
│   │   ├── repository.go
│   │   ├── rules.go
│   │   ├── issuance.go
│   │   ├── verification.go
│   │   └── revocation.go
│   ├── governance/
│   │   ├── model.go
│   │   ├── repository.go
│   │   └── service.go
│   ├── events/
│   │   ├── model.go
│   │   ├── outbox.go
│   │   ├── publisher.go
│   │   └── taxonomy.go
│   ├── importer/
│   │   ├── csv.go
│   │   ├── mapper.go
│   │   ├── validator.go
│   │   └── upsert.go
│   ├── localization/
│   │   ├── model.go
│   │   ├── repository.go
│   │   └── service.go
│   ├── results/
│   │   ├── model.go
│   │   ├── service.go
│   │   └── snapshots.go
│   └── http/
│       ├── dto.go
│       ├── router.go
│       ├── sessions_handler.go
│       ├── questions_handler.go
│       ├── answers_handler.go
│       ├── reviews_handler.go
│       ├── credentials_handler.go
│       └── results_handler.go
├── migrations/
│   ├── 0001_extensions_and_schemas.sql
│   ├── 0002_core.sql
│   ├── 0003_content.sql
│   ├── 0004_profiling.sql
│   ├── 0005_runtime.sql
│   ├── 0006_review.sql
│   ├── 0007_credentials.sql
│   ├── 0008_governance_and_events.sql
│   └── 0009_indexes.sql
├── seeds/
│   ├── core/
│   ├── profiling/
│   └── content/
├── web/
│   └── ui/
├── .env.example
├── Makefile
├── go.mod
└── README.md
```

---

## Package-Prinzipien

### `db/`
Nur Datenbankzugang, Connection Pool, Transactions, Migrations-Handling.

### `domains/`
Domänenspezifische Logik. Keine generischen Repositories hineinmischen.

### `profiling/`
Tag-Mappings, Trait-/Skill-/Meaning-Vektoren, Zuordnungslogik.

### `credentials/`
Alles rund um Trust Ladder, Rules, Issuance, Verification, Revocation.

### `events/`
Outbox-Pattern, Event-Taxonomie, Publisher-Schnittstelle.

### `results/`
Zusammenführung von Scores, Vektoren, Deutungen und Snapshots.

### `importer/`
CSV-/Seed-Import getrennt von Runtime.

---

## Empfohlene technische Basis in Go
### Driver / DB-Zugriff
Empfohlen: **pgx/v5** mit `pgxpool` für PostgreSQL-first.

### Begründung
- PostgreSQL-spezifisch stark
- gute Performance
- sauber für JSONB, Arrays, Transactions
- offen für späteres `LISTEN/NOTIFY`, `COPY`, Bulk-Operationen

---

## Startdateien – minimale Final-V3-Skelette

## `cmd/api/main.go`
```go
package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/your-org/cee4ai/internal/app"
	"github.com/your-org/cee4ai/internal/config"
)

func main() {
	cfg := config.MustLoad()
	application, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: application.Router,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		_ = srv.Shutdown(context.Background())
	}()

	log.Printf("api listening on %s", cfg.HTTPAddr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
```

---

## `internal/config/config.go`
```go
package config

import (
	"log"
	"os"
)

type Config struct {
	HTTPAddr     string
	PostgresURL  string
	AppEnv       string
}

func MustLoad() Config {
	cfg := Config{
		HTTPAddr:    getenv("HTTP_ADDR", ":8080"),
		PostgresURL: getenv("POSTGRES_URL", "postgres://localhost:5432/cee4ai?sslmode=disable"),
		AppEnv:      getenv("APP_ENV", "development"),
	}
	if cfg.PostgresURL == "" {
		log.Fatal("POSTGRES_URL is required")
	}
	return cfg
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
```

---

## `internal/db/postgres.go`
```go
package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	cfg.MaxConns = 20
	cfg.MinConns = 2
	cfg.MaxConnLifetime = time.Hour
	cfg.MaxConnIdleTime = 15 * time.Minute
	return pgxpool.NewWithConfig(ctx, cfg)
}
```

---

## `internal/app/app.go`
```go
package app

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/cee4ai/internal/config"
	appdb "github.com/your-org/cee4ai/internal/db"
	apphttp "github.com/your-org/cee4ai/internal/http"
)

type App struct {
	Config config.Config
	DB     *pgxpool.Pool
	Router http.Handler
}

func New(cfg config.Config) (*App, error) {
	pool, err := appdb.NewPool(context.Background(), cfg.PostgresURL)
	if err != nil {
		return nil, err
	}

	router := apphttp.NewRouter(pool)

	return &App{
		Config: cfg,
		DB:     pool,
		Router: router,
	}, nil
}
```

---

## `internal/http/router.go`
```go
package http

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(pool *pgxpool.Pool) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// TODO: sessions, questions, answers, reviews, credentials, results
	return mux
}
```

---

## `internal/questions/model.go`
```go
package questions

type Question struct {
	ID                                int64
	ExternalID                        string
	DomainID                          int64
	QuestionFamily                    string
	CategoryID                        int64
	SubcategoryID                     *int64
	ParentQuestionID                  *int64
	QuestionType                      string
	ScoringMode                       string
	IntendedUse                       string
	ConfidenceTier                    string
	EstimatedTimeSeconds              *int
	CognitiveLoadLevel                string
	CulturalScope                     string
	RegionScope                       string
	MeaningDepth                      string
	WorldviewSensitivity              string
	SymbolicInterpretationRelevance   string
	ExistentialLoadLevel              string
	IsSensitive                       bool
	AgeGate                           int
	ReviewStatus                      string
	IsActive                          bool
	Version                           int
}
```

---

## `internal/credentials/model.go`
```go
package credentials

type Assertion struct {
	ID                          string
	UserID                      *string
	SessionID                   *string
	TrackID                     int64
	RuleID                      int64
	IssuerID                    int64
	CertificateLevel            string
	IssuedUnderRulesetVersion   string
	ResultConfidence            *float64
	AssessmentIntegrityScore    *float64
	EvidenceHash                string
	CredentialStandard          string
	VerificationToken           string
}
```

---

## `internal/events/model.go`
```go
package events

import "time"

type Event struct {
	ID           string
	EventType    string
	EventVersion int
	SubjectType  string
	SubjectID    string
	SessionID    *string
	UserID       *string
	OccurredAt   time.Time
	Payload      []byte
	PartitionKey string
}
```

---

## `internal/domains/meaning/service.go`
```go
package meaning

type JourneyHint struct {
	Title string
	Body  string
}

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) BuildJourneyHints() []JourneyHint {
	return []JourneyHint{
		{
			Title: "Erste Sinn- und Tiefenresonanz",
			Body:  "Deine bisherigen Antworten zeigen erste Hinweise auf Sinn- und Bedeutungsräume. Dieser Zwischenstand ist vorläufig und eher als Richtung als als fertiges Bild zu verstehen.",
		},
	}
}
```

---

## Migration Runner
### `cmd/migrate/main.go`
```go
package main

import (
	"context"
	"log"

	"github.com/your-org/cee4ai/internal/config"
	appdb "github.com/your-org/cee4ai/internal/db"
)

func main() {
	cfg := config.MustLoad()
	pool, err := appdb.NewPool(context.Background(), cfg.PostgresURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	if err := appdb.RunMigrations(context.Background(), pool, "./migrations"); err != nil {
		log.Fatal(err)
	}

	log.Println("migrations applied")
}
```

---

## `internal/db/migrations.go`
```go
package db

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunMigrations(ctx context.Context, pool *pgxpool.Pool, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	var files []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".sql" {
			files = append(files, filepath.Join(dir, e.Name()))
		}
	}
	sort.Strings(files)

	for _, file := range files {
		body, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		if _, err := pool.Exec(ctx, string(body)); err != nil {
			return fmt.Errorf("migration %s failed: %w", file, err)
		}
	}
	return nil
}
```

---

## `.env.example`
```dotenv
APP_ENV=development
HTTP_ADDR=:8080
POSTGRES_URL=postgres://cee4ai:cee4ai@localhost:5432/cee4ai?sslmode=disable
```

---

## `go.mod` – empfohlene Pakete
```text
github.com/jackc/pgx/v5
```

Optional später:
```text
github.com/pressly/goose/v3
```

---

## `Makefile`
```make
.PHONY: migrate run

migrate:
	go run ./cmd/migrate

run:
	go run ./cmd/api
```

---

# Teil C – Reihenfolge der tatsächlichen Umsetzung

## Schritt 1
PostgreSQL lokal starten und Datenbank anlegen.

## Schritt 2
Migrationsordner anlegen und die obigen SQL-Dateien übernehmen.

## Schritt 3
Go-Modul mit `pgx/v5` initialisieren.

## Schritt 4
`cmd/migrate` und `internal/db` zuerst lauffähig machen.

## Schritt 5
`questions`, `sessions`, `answers`, `results` als erste Runtime-Pakete umsetzen.

## Schritt 6
Danach `profiling`, `reviews`, `credentials`, `events`, `meaning` ergänzen.

---

# Teil D – Vermeidungsstrategien

## 1. Keine Monolith-Datei für alles
Domain-Logik, DB-Zugriff und HTTP nicht zusammenkleben.

## 2. Nicht zu früh alle Domänen komplett implementieren
Cognitive Runtime zuerst, aber Repo-Struktur schon für Personality, Meaning, Credentials und Events offen halten.

## 3. JSONB nur gezielt nutzen
Für Payloads, Snapshots und flexible Metadaten ja – nicht als Ersatz für saubere Kernrelationen.

## 4. Meaning-/Journey-Logik nicht als Sonderfall verstecken
Besser als eigenes Paket und saubere Tag-/Text-Logik führen.

---

# Abschluss
Mit dieser DDL-/Migrationsfassung und der Final-V3-Go-Repo-Struktur ist CPE jetzt technisch so geerdet, dass die zuvor geschärfte Vision sauber in echte Entwicklung übergehen kann.

