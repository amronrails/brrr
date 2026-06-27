#!/usr/bin/env bash
#
# verify-generate.sh — end-to-end check that `brrr generate` produces CRUD that
# compiles (backend) and type-checks/builds (frontend), for both a brand-new
# module and an additional model in an existing module.
#
# Usage:
#   scripts/verify-generate.sh [workdir]
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
WORKDIR="${1:-$(mktemp -d)}"
APP="genapp"
APP_DIR="${WORKDIR}/${APP}"

# See verify-init.sh for why DWARF is disabled.
export GOFLAGS="${GOFLAGS:--gcflags=all=-dwarf=false}"

step() { printf '\n\033[1;36m==> %s\033[0m\n' "$1"; }

step "Building brrr"
go build -o "${WORKDIR}/brrr" ./cmd/brrr
BRRR="${WORKDIR}/brrr"

step "Initialising ${APP_DIR}"
rm -rf "${APP_DIR}"
"${BRRR}" init "${APP}" --module "github.com/example/${APP}" --dir "${APP_DIR}"
cd "${APP_DIR}"

step "Generating Post (new module) and Comment (existing module)"
"${BRRR}" g blog Post title:string:required body:text published:bool author:belongs_to:User
"${BRRR}" g blog Comment body:text:required post:belongs_to:Post

step "Backend: sqlc generate"
go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate

step "Backend: go mod tidy / build / vet"
go mod tidy
go build ./...
go vet ./...

step "Frontend: npm install + build"
cd web
npm install --no-audit --no-fund
npm run build

printf '\n\033[1;32m✓ verify-generate passed\033[0m  (%s)\n' "${APP_DIR}"
