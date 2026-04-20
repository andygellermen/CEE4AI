#!/bin/sh
set -eu

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "required command not found: $1" >&2
    exit 1
  fi
}

require_cmd curl
require_cmd jq
require_cmd mktemp

BASE_URL="${BASE_URL:-https://cpe.geller.men}"
LANGUAGE_ID="${LANGUAGE_ID:-1}"
STANDARD_DOMAIN_ID="${STANDARD_DOMAIN_ID:-1}"
GUARDED_DOMAIN_ID="${GUARDED_DOMAIN_ID:-5}"
KEEP_ARTIFACTS="${KEEP_ARTIFACTS:-false}"
WORKDIR="$(mktemp -d)"

cleanup() {
  if [ "$KEEP_ARTIFACTS" = "true" ]; then
    return 0
  fi
  rm -rf "$WORKDIR"
}

trap cleanup EXIT INT TERM

request() {
  method="$1"
  path="$2"
  body="${3:-}"

  if [ -n "$body" ]; then
    curl -fsS -X "$method" \
      -H 'Content-Type: application/json' \
      -d "$body" \
      "${BASE_URL}${path}"
  else
    curl -fsS -X "$method" "${BASE_URL}${path}"
  fi
}

assert_jq() {
  file="$1"
  expression="$2"
  if ! jq -e "$expression" "$file" >/dev/null; then
    echo "assertion failed for ${file}: ${expression}" >&2
    jq '.' "$file" >&2
    exit 1
  fi
}

print_step() {
  echo
  echo "== $1 =="
}

health_file="$WORKDIR/health.txt"
standard_session_file="$WORKDIR/standard-session.json"
standard_question_file="$WORKDIR/standard-question.json"
standard_answer_file="$WORKDIR/standard-answer.json"
standard_result_file="$WORKDIR/standard-result.json"
guarded_session_file="$WORKDIR/guarded-session.json"
guarded_question_file="$WORKDIR/guarded-question.json"
guarded_answer_file="$WORKDIR/guarded-answer.json"
guarded_result_file="$WORKDIR/guarded-result.json"

print_step "Health"
request GET /healthz >"$health_file"
if [ "$(cat "$health_file")" != "ok" ]; then
  echo "unexpected health response:" >&2
  cat "$health_file" >&2
  exit 1
fi
echo "health check ok at ${BASE_URL}/healthz"

print_step "Standard Snapshot Flow"
standard_session_payload="$(jq -nc \
  --argjson domain_id "$STANDARD_DOMAIN_ID" \
  --argjson locale_language_id "$LANGUAGE_ID" \
  '{domain_id:$domain_id, mode:"snapshot", session_goal:"live smoke test", locale_language_id:$locale_language_id}')"
request POST /api/v1/sessions "$standard_session_payload" >"$standard_session_file"
assert_jq "$standard_session_file" '.session.id != ""'
assert_jq "$standard_session_file" '.session.mode == "snapshot"'
standard_session_id="$(jq -r '.session.id' "$standard_session_file")"
echo "standard session: ${standard_session_id}"

request GET "/api/v1/questions/next?session_id=${standard_session_id}" >"$standard_question_file"
assert_jq "$standard_question_file" '.question.id > 0'
assert_jq "$standard_question_file" '.governance.delivery_mode == "standard"'
standard_question_id="$(jq -r '.question.id' "$standard_question_file")"
standard_option_id="$(jq -r '(.question.options[] | select(.option_key == "B") | .id) // .question.options[0].id' "$standard_question_file")"
echo "standard question: ${standard_question_id}, answer option: ${standard_option_id}"

standard_answer_payload="$(jq -nc \
  --arg session_id "$standard_session_id" \
  --argjson question_id "$standard_question_id" \
  --argjson option_id "$standard_option_id" \
  '{session_id:$session_id, question_id:$question_id, selected_option_ids:[$option_id], certainty_level:"high"}')"
request POST /api/v1/answers "$standard_answer_payload" >"$standard_answer_file"
assert_jq "$standard_answer_file" '.answer.id > 0'
assert_jq "$standard_answer_file" '.governance.delivery_mode == "standard"'

request GET "/api/v1/results?session_id=${standard_session_id}" >"$standard_result_file"
assert_jq "$standard_result_file" '.result_type == "snapshot_profile"'
assert_jq "$standard_result_file" '.payload.governance.delivery_mode == "standard"'
assert_jq "$standard_result_file" '.payload.top_signals.denktype != ""'
echo "standard snapshot top signals:"
jq '{denktype: .payload.top_signals.denktype, skill: .payload.top_signals.skill, trait: .payload.top_signals.trait}' "$standard_result_file"

print_step "Guarded Meaning Flow"
guarded_session_payload="$(jq -nc \
  --argjson domain_id "$GUARDED_DOMAIN_ID" \
  --argjson locale_language_id "$LANGUAGE_ID" \
  '{domain_id:$domain_id, mode:"snapshot", session_goal:"guarded meaning smoke test", locale_language_id:$locale_language_id}')"
request POST /api/v1/sessions "$guarded_session_payload" >"$guarded_session_file"
assert_jq "$guarded_session_file" '.session.id != ""'
guarded_session_id="$(jq -r '.session.id' "$guarded_session_file")"
echo "guarded session: ${guarded_session_id}"

request GET "/api/v1/questions/next?session_id=${guarded_session_id}" >"$guarded_question_file"
assert_jq "$guarded_question_file" '.question.id > 0'
assert_jq "$guarded_question_file" '.question.question_type == "reflection"'
assert_jq "$guarded_question_file" '.governance.delivery_mode == "guarded"'
assert_jq "$guarded_question_file" '.governance.review_required == true'
guarded_question_id="$(jq -r '.question.id' "$guarded_question_file")"
echo "guarded question: ${guarded_question_id}"

guarded_answer_payload="$(jq -nc \
  --arg session_id "$guarded_session_id" \
  --argjson question_id "$guarded_question_id" \
  '{session_id:$session_id, question_id:$question_id, free_text_answer:"Ich erlebe Sinn besonders in stillen Momenten, in Verbundenheit mit anderen Menschen und wenn mein Handeln mit meinen Werten zusammenpasst.", certainty_level:"medium"}')"
request POST /api/v1/answers "$guarded_answer_payload" >"$guarded_answer_file"
assert_jq "$guarded_answer_file" '.answer.id > 0'
assert_jq "$guarded_answer_file" '.governance.delivery_mode == "guarded"'

request GET "/api/v1/results?session_id=${guarded_session_id}" >"$guarded_result_file"
assert_jq "$guarded_result_file" '.payload.governance.delivery_mode == "guarded"'
assert_jq "$guarded_result_file" '.payload.governance.review_flag_count >= 1'
assert_jq "$guarded_result_file" '.payload.governance.answered_sensitive_questions >= 1'
assert_jq "$guarded_result_file" '.payload.governance.answered_worldview_sensitive_questions >= 1'
echo "guarded governance summary:"
jq '.payload.governance' "$guarded_result_file"

print_step "Smoke Test Complete"
echo "all live smoke checks passed for ${BASE_URL}"
if [ "$KEEP_ARTIFACTS" = "true" ]; then
  echo "artifacts kept at ${WORKDIR}"
else
  echo "artifacts were written to ${WORKDIR} during execution"
fi
