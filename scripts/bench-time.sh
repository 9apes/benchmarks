#!/usr/bin/env bash
# Source from workflow steps: source scripts/bench-time.sh
now_ms() { date +%s%3N; }
dur_ms() { echo $(( $(now_ms) - "$1" )); }
log_step() { echo "STEP_MS_${1}=$(dur_ms "$2")" >> "${GITHUB_ENV:-/dev/null}"; }
