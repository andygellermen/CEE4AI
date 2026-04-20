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
    id              BIGSERIAL PRIMARY KEY,
    slug            CITEXT NOT NULL UNIQUE,
    name            TEXT NOT NULL,
    description     TEXT,
    policy_payload  JSONB NOT NULL,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE governance.ruleset_versions (
    id              BIGSERIAL PRIMARY KEY,
    domain_id       BIGINT REFERENCES core.domains(id),
    slug            CITEXT NOT NULL,
    version         TEXT NOT NULL,
    description     TEXT,
    effective_from  TIMESTAMPTZ NOT NULL,
    effective_to    TIMESTAMPTZ,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    UNIQUE (slug, version)
);

CREATE TABLE governance.audit_log (
    id              BIGSERIAL PRIMARY KEY,
    actor_user_id   UUID REFERENCES core.users(id) ON DELETE SET NULL,
    entity_type     TEXT NOT NULL,
    entity_id       TEXT NOT NULL,
    action          TEXT NOT NULL,
    payload         JSONB,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE events.event_outbox (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    aggregate_type  TEXT NOT NULL,
    aggregate_id    TEXT NOT NULL,
    event_type      TEXT NOT NULL,
    event_version   INTEGER NOT NULL DEFAULT 1,
    payload         JSONB NOT NULL,
    occurred_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    published_at    TIMESTAMPTZ
);

CREATE TABLE events.event_inbox (
    id              UUID PRIMARY KEY,
    source          TEXT NOT NULL,
    event_type      TEXT NOT NULL,
    event_version   INTEGER NOT NULL DEFAULT 1,
    payload         JSONB NOT NULL,
    received_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at    TIMESTAMPTZ
);

CREATE TABLE events.event_store (
    id              BIGSERIAL PRIMARY KEY,
    event_id        UUID NOT NULL UNIQUE,
    event_type      TEXT NOT NULL,
    subject_type    TEXT NOT NULL,
    subject_id      TEXT NOT NULL,
    session_id      UUID,
    user_id         UUID,
    occurred_at     TIMESTAMPTZ NOT NULL,
    payload         JSONB NOT NULL,
    partition_key   TEXT NOT NULL
);

COMMIT;
