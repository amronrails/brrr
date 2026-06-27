#!/usr/bin/env bash
#
# verify-spec.sh — end-to-end check that `brrr init --spec` scaffolds a project
# and generates a whole model graph (with cross-model relationships and
# dependency ordering) that compiles (backend) and builds (frontend).
#
# Usage:
#   scripts/verify-spec.sh [workdir]
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
WORKDIR="${1:-$(mktemp -d)}"
APP="specapp"
APP_DIR="${WORKDIR}/${APP}"
SPEC="${WORKDIR}/spec.yaml"

# See verify-init.sh for why DWARF is disabled.
export GOFLAGS="${GOFLAGS:--gcflags=all=-dwarf=false}"

step() { printf '\n\033[1;36m==> %s\033[0m\n' "$1"; }

# Comment is intentionally declared before Post to exercise topological ordering.
cat > "${SPEC}" <<'YAML'
modules:
  blog:
    models:
      Comment:
        fields:
          body: text required
        relationships:
          post: belongs_to Post
          author: belongs_to User
      Post:
        fields:
          title: string required
          body: text
          published: bool
        relationships:
          author: belongs_to User
  shop:
    models:
      Product:
        fields:
          name: string required
          price: decimal
          sku: string unique
          in_stock: bool
YAML

step "Building brrr"
go build -o "${WORKDIR}/brrr" ./cmd/brrr

step "init --spec ${SPEC}"
rm -rf "${APP_DIR}"
"${WORKDIR}/brrr" init "${APP}" --module "github.com/example/${APP}" --dir "${APP_DIR}" --spec "${SPEC}"
cd "${APP_DIR}"

step "Backend: sqlc generate / tidy / build / vet"
go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate
go mod tidy
go build ./...
go vet ./...

step "Frontend: npm install + build"
cd web
npm install --no-audit --no-fund
npm run build

printf '\n\033[1;32m✓ verify-spec passed\033[0m  (%s)\n' "${APP_DIR}"
