BEGIN;

CREATE TABLE content.categories (
    id                          BIGSERIAL PRIMARY KEY,
    domain_id                   BIGINT NOT NULL REFERENCES core.domains(id),
    slug                        CITEXT NOT NULL UNIQUE,
    name                        TEXT NOT NULL,
    description                 TEXT,
    default_sensitivity         TEXT,
    cultural_scope              TEXT,
    spiritual_relevance         TEXT,
    worldview_sensitivity       TEXT,
    meaning_pathway_relevance   TEXT,
    is_sensitive                BOOLEAN NOT NULL DEFAULT FALSE,
    is_active                   BOOLEAN NOT NULL DEFAULT TRUE,
    created_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE content.subcategories (
    id              BIGSERIAL PRIMARY KEY,
    category_id     BIGINT NOT NULL REFERENCES content.categories(id) ON DELETE CASCADE,
    slug            CITEXT NOT NULL UNIQUE,
    name            TEXT NOT NULL,
    description     TEXT,
    is_sensitive    BOOLEAN NOT NULL DEFAULT FALSE,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
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
    CONSTRAINT question_family_chk CHECK (
        question_family IN (
            'knowledge', 'skill', 'trait', 'interest', 'trigger', 'pathway',
            'reflective', 'contemplative', 'comparative_worldview',
            'symbolic_interpretation', 'existential'
        )
    ),
    CONSTRAINT question_type_chk CHECK (
        question_type IN ('single_choice', 'multiple_select', 'scale', 'reflection')
    ),
    CONSTRAINT scoring_mode_chk CHECK (
        scoring_mode IN ('exact', 'partial', 'weighted', 'non_scored', 'path_only')
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
    id              BIGSERIAL PRIMARY KEY,
    question_id     BIGINT NOT NULL REFERENCES content.question_master(id) ON DELETE CASCADE,
    option_key      TEXT NOT NULL,
    score_weight    NUMERIC(8,4) NOT NULL DEFAULT 1.0,
    is_correct      BOOLEAN NOT NULL DEFAULT FALSE,
    display_order   INTEGER NOT NULL,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    version         INTEGER NOT NULL DEFAULT 1,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
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
    id              BIGSERIAL PRIMARY KEY,
    domain_id       BIGINT NOT NULL REFERENCES core.domains(id),
    result_type     TEXT NOT NULL,
    certainty_level TEXT,
    profile_depth   TEXT,
    slug            CITEXT NOT NULL UNIQUE,
    is_sensitive    BOOLEAN NOT NULL DEFAULT FALSE,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
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
