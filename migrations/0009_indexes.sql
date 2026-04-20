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
