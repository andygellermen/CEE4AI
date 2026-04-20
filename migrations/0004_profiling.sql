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
    id              BIGSERIAL PRIMARY KEY,
    slug            CITEXT NOT NULL UNIQUE,
    name            TEXT NOT NULL,
    description     TEXT,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE profiling.personality_traits (
    id              BIGSERIAL PRIMARY KEY,
    slug            CITEXT NOT NULL UNIQUE,
    name            TEXT NOT NULL,
    description     TEXT,
    polarity_model  TEXT,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE profiling.interest_tags (
    id              BIGSERIAL PRIMARY KEY,
    slug            CITEXT NOT NULL UNIQUE,
    name            TEXT NOT NULL,
    description     TEXT,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE
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
    id              BIGSERIAL PRIMARY KEY,
    slug            CITEXT NOT NULL UNIQUE,
    name            TEXT NOT NULL,
    description     TEXT,
    is_sensitive    BOOLEAN NOT NULL DEFAULT FALSE,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE profiling.body_references (
    id              BIGSERIAL PRIMARY KEY,
    slug            CITEXT NOT NULL UNIQUE,
    name            TEXT NOT NULL,
    description     TEXT,
    is_sensitive    BOOLEAN NOT NULL DEFAULT FALSE,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE
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
    question_id BIGINT NOT NULL REFERENCES content.question_master(id) ON DELETE CASCADE,
    denktype_id BIGINT NOT NULL REFERENCES profiling.denktypes(id) ON DELETE CASCADE,
    weight      NUMERIC(8,4) NOT NULL DEFAULT 1.0,
    rationale   TEXT,
    is_active   BOOLEAN NOT NULL DEFAULT TRUE,
    PRIMARY KEY (question_id, denktype_id)
);

CREATE TABLE profiling.question_skill_tags (
    question_id BIGINT NOT NULL REFERENCES content.question_master(id) ON DELETE CASCADE,
    skill_id    BIGINT NOT NULL REFERENCES profiling.skills(id) ON DELETE CASCADE,
    weight      NUMERIC(8,4) NOT NULL DEFAULT 1.0,
    rationale   TEXT,
    is_active   BOOLEAN NOT NULL DEFAULT TRUE,
    PRIMARY KEY (question_id, skill_id)
);

CREATE TABLE profiling.question_trait_tags (
    question_id BIGINT NOT NULL REFERENCES content.question_master(id) ON DELETE CASCADE,
    trait_id    BIGINT NOT NULL REFERENCES profiling.personality_traits(id) ON DELETE CASCADE,
    weight      NUMERIC(8,4) NOT NULL DEFAULT 1.0,
    rationale   TEXT,
    is_active   BOOLEAN NOT NULL DEFAULT TRUE,
    PRIMARY KEY (question_id, trait_id)
);

CREATE TABLE profiling.question_interest_tags (
    question_id     BIGINT NOT NULL REFERENCES content.question_master(id) ON DELETE CASCADE,
    interest_tag_id BIGINT NOT NULL REFERENCES profiling.interest_tags(id) ON DELETE CASCADE,
    weight          NUMERIC(8,4) NOT NULL DEFAULT 1.0,
    rationale       TEXT,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    PRIMARY KEY (question_id, interest_tag_id)
);

CREATE TABLE profiling.question_meaning_tags (
    question_id     BIGINT NOT NULL REFERENCES content.question_master(id) ON DELETE CASCADE,
    meaning_tag_id  BIGINT NOT NULL REFERENCES profiling.meaning_tags(id) ON DELETE CASCADE,
    weight          NUMERIC(8,4) NOT NULL DEFAULT 1.0,
    rationale       TEXT,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
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
    question_id      BIGINT NOT NULL REFERENCES content.question_master(id) ON DELETE CASCADE,
    trigger_group_id BIGINT NOT NULL REFERENCES profiling.trigger_groups(id) ON DELETE CASCADE,
    weight           NUMERIC(8,4) NOT NULL DEFAULT 1.0,
    rationale        TEXT,
    is_active        BOOLEAN NOT NULL DEFAULT TRUE,
    PRIMARY KEY (question_id, trigger_group_id)
);

CREATE TABLE profiling.question_body_reference_tags (
    question_id        BIGINT NOT NULL REFERENCES content.question_master(id) ON DELETE CASCADE,
    body_reference_id  BIGINT NOT NULL REFERENCES profiling.body_references(id) ON DELETE CASCADE,
    weight             NUMERIC(8,4) NOT NULL DEFAULT 1.0,
    rationale          TEXT,
    is_active          BOOLEAN NOT NULL DEFAULT TRUE,
    PRIMARY KEY (question_id, body_reference_id)
);

CREATE TABLE profiling.question_path_tags (
    question_id          BIGINT NOT NULL REFERENCES content.question_master(id) ON DELETE CASCADE,
    development_path_id  BIGINT NOT NULL REFERENCES profiling.development_paths(id) ON DELETE CASCADE,
    weight               NUMERIC(8,4) NOT NULL DEFAULT 1.0,
    rationale            TEXT,
    is_active            BOOLEAN NOT NULL DEFAULT TRUE,
    PRIMARY KEY (question_id, development_path_id)
);

COMMIT;
