#!/usr/bin/env bash
#
# verify-init.sh — end-to-end check that `brrr init` produces a project whose
# backend compiles/vets and whose frontend type-checks and builds.
#
# Usage:
#   scripts/verify-init.sh [workdir]
#
# It generates a throwaway project, then runs:
#   backend  : sqlc generate -> go mod tidy -> go build -> go vet
#   frontend : npm install -> npm run build (tsc -b && vite build)
#
# The runtime auth smoke test (Postgres + register/login) is documented in the
# generated README and is not run here, so this script needs no Docker.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
WORKDIR="${1:-$(mktemp -d)}"
APP="verifyapp"
APP_DIR="${WORKDIR}/${APP}"

# Some Go 1.26 toolchains crash in DWARF generation when building large
# `go run` tools (sqlc/goose). Disabling DWARF sidesteps that without affecting
# the generated project; debug symbols are irrelevant to a verification build.
export GOFLAGS="${GOFLAGS:--gcflags=all=-dwarf=false}"

step() { printf '\n\033[1;36m==> %s\033[0m\n' "$1"; }

step "Building brrr"
go build -o "${WORKDIR}/brrr" ./cmd/brrr

step "Generating sample project at ${APP_DIR}"
rm -rf "${APP_DIR}"
"${WORKDIR}/brrr" init "${APP}" --module "github.com/example/${APP}" --dir "${APP_DIR}"

cd "${APP_DIR}"

step "Backend: sqlc generate"
go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate

step "Backend: go mod tidy"
go mod tidy

step "Backend: go build ./..."
go build ./...

step "Backend: go vet ./..."
go vet ./...

step "Frontend: npm install"
cd web
npm install --no-audit --no-fund

step "Frontend: npm run build"
npm run build

printf '\n\033[1;32m✓ verify-init passed\033[0m  (%s)\n' "${APP_DIR}"
