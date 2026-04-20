BEGIN;

CREATE TABLE runtime.sessions (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID REFERENCES core.users(id) ON DELETE SET NULL,
    domain_id           BIGINT NOT NULL REFERENCES core.domains(id),
    mode                TEXT NOT NULL,
    session_goal        TEXT,
    locale_language_id  BIGINT REFERENCES core.languages(id),
    locale_region_id    BIGINT REFERENCES core.regions(id),
    result_confidence   NUMERIC(8,4),
    progress_state      TEXT,
    started_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    finished_at         TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT session_mode_chk CHECK (
        mode IN (
            'snapshot', 'guided_progression', 'deep_profile',
            'personality_extension', 'career_path', 'meaning_journey'
        )
    )
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
    id                  BIGSERIAL PRIMARY KEY,
    session_id          UUID NOT NULL REFERENCES runtime.sessions(id) ON DELETE CASCADE,
    package_id          BIGINT REFERENCES runtime.session_packages(id) ON DELETE SET NULL,
    question_id         BIGINT NOT NULL REFERENCES content.question_master(id),
    answer_kind         TEXT NOT NULL,
    selected_option_ids JSONB,
    scale_value         INTEGER,
    free_text_answer    TEXT,
    raw_score           NUMERIC(10,4),
    evaluated_score     NUMERIC(10,4),
    certainty_level     TEXT,
    answered_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT answer_kind_chk CHECK (
        answer_kind IN ('single_choice', 'multiple_select', 'scale', 'reflection')
    )
);

CREATE TABLE runtime.result_snapshots (
    id              BIGSERIAL PRIMARY KEY,
    session_id      UUID NOT NULL REFERENCES runtime.sessions(id) ON DELETE CASCADE,
    result_type     TEXT NOT NULL,
    profile_depth   TEXT,
    certainty_level TEXT,
    snapshot_payload JSONB NOT NULL,
    ruleset_version TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE runtime.profile_vectors (
    id              BIGSERIAL PRIMARY KEY,
    session_id      UUID NOT NULL REFERENCES runtime.sessions(id) ON DELETE CASCADE,
    vector_type     TEXT NOT NULL,
    vector_payload  JSONB NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
