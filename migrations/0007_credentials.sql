BEGIN;

CREATE TABLE credentials.competency_frameworks (
    id              BIGSERIAL PRIMARY KEY,
    slug            CITEXT NOT NULL UNIQUE,
    name            TEXT NOT NULL,
    description     TEXT,
    version         TEXT NOT NULL,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE
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
    id              BIGSERIAL PRIMARY KEY,
    slug            CITEXT NOT NULL UNIQUE,
    name            TEXT NOT NULL,
    description     TEXT,
    domain_id       BIGINT NOT NULL REFERENCES core.domains(id),
    framework_id    BIGINT REFERENCES credentials.competency_frameworks(id),
    issuer_id       BIGINT NOT NULL REFERENCES credentials.credential_issuers(id),
    is_active       BOOLEAN NOT NULL DEFAULT TRUE
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
    id              BIGSERIAL PRIMARY KEY,
    assertion_id    UUID NOT NULL REFERENCES credentials.credential_assertions(id) ON DELETE CASCADE,
    evidence_type   TEXT NOT NULL,
    evidence_ref    TEXT,
    evidence_hash   TEXT NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
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
