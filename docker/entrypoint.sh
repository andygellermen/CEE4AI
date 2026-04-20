#!/bin/sh
set -eu

run_with_retry() {
  command_name="$1"
  retries="${DB_STARTUP_RETRIES:-20}"
  sleep_seconds="${DB_STARTUP_SLEEP_SECONDS:-3}"
  attempt=1

  while [ "$attempt" -le "$retries" ]; do
    if "$command_name"; then
      return 0
    fi

    if [ "$attempt" -eq "$retries" ]; then
      echo "command failed after ${retries} attempts: ${command_name}" >&2
      return 1
    fi

    echo "command failed, retrying in ${sleep_seconds}s (${attempt}/${retries}): ${command_name}" >&2
    attempt=$((attempt + 1))
    sleep "$sleep_seconds"
  done

  return 1
}

if [ "${RUN_MIGRATIONS_ON_BOOT:-true}" = "true" ]; then
  run_with_retry /app/cee4ai-migrate
fi

if [ "${RUN_SEED_ON_BOOT:-false}" = "true" ]; then
  run_with_retry /app/cee4ai-seed
fi

exec /app/cee4ai-api
