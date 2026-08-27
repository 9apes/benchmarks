#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"

echo "==> pnpm-bench"
cd "$ROOT/pnpm-bench"
if command -v pnpm >/dev/null; then
  pnpm install
else
  echo "pnpm not found; install pnpm 9+ or run: corepack enable && pnpm install"
  exit 1
fi

echo "==> go-bench"
cd "$ROOT/go-bench"
go mod tidy
go build -o /dev/null ./...

echo "==> maven-bench"
cd "$ROOT/maven-bench"
mvn -B dependency:go-offline package -DskipTests

echo "Done. Commit pnpm-lock.yaml and go.sum (maven resolves from pom.xml on CI)."
