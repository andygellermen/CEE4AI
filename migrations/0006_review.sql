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
