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
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    external_ref        TEXT UNIQUE,
    locale_language_id  BIGINT REFERENCES core.languages(id),
    locale_region_id    BIGINT REFERENCES core.regions(id),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE core.roles (
    id      BIGSERIAL PRIMARY KEY,
    slug    CITEXT NOT NULL UNIQUE,
    name    TEXT NOT NULL
);

CREATE TABLE core.user_roles (
    user_id     UUID NOT NULL REFERENCES core.users(id) ON DELETE CASCADE,
    role_id     BIGINT NOT NULL REFERENCES core.roles(id) ON DELETE CASCADE,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, role_id)
);

COMMIT;
