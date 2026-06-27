# brrr 🖨️

A CRUD code generator for **Go (modular monolith) + React**. Think *JHipster, for a
Go/chi/sqlc backend and a Vite/React/shadcn frontend*.

## Features

1. **`brrr init <app>`** — scaffold a new project: a Go modular monolith with a
   secure email/password + JWT **auth module** and a React admin frontend under
   `web/`. ✅ _built_
2. **`brrr generate <module> <Model> field:type ...`** — Rails-style CRUD
   generation (backend + frontend), wired in automatically. ✅ _built_
3. **`brrr init <app> --spec <file.yaml>`** — scaffold a project and generate a
   whole model graph (modules, models, relationships) from a YAML spec, in
   dependency order. ✅ _built_

## Generated stack

| Layer    | Choice                                                              |
| -------- | ------------------------------------------------------------------ |
| Backend  | Go · chi · sqlc + pgx · Postgres · goose · golang-jwt v5 · bcrypt   |
| Frontend | Vite · TypeScript · TanStack Query · React Router · shadcn-style UI + Tailwind |

The Go backend is a **modular monolith**: each feature is a vertical slice under
`internal/modules/<name>` with a layered interior — `domain → repository →
service → transport/http`. The built-in `user` module ships secure auth (the
first account to register becomes the admin).

## Usage

```sh
go build -o bin/brrr ./cmd/brrr

bin/brrr init myapp --module github.com/me/myapp
cd myapp
make db-up        # start Postgres (Docker)
make setup        # go mod tidy + sqlc generate + npm install
make migrate      # apply migrations
make dev          # backend (:8080) + frontend (:5173)
```

Open http://localhost:5173, register the first user (→ admin), and you have a
working dashboard with an admin-gated Users page.

## How it works

- Templates live under `internal/templates/init/**` and are embedded with
  `//go:embed`. Every template file has a `.tmpl` suffix so template Go/TS code
  is never compiled into the `brrr` binary.
- `internal/engine` renders both file paths and bodies with `text/template`,
  using helpers in `funcs.go` (pascal/camel/snake/kebab/plural/singular).
- `internal/spec` defines the field-type registry (Go/SQL/TS/zod mappings) and
  relationship model that features 2 and 3 will build on.
- `internal/project` holds the project `Context` and the `brrr.yaml` manifest
  that tracks generated modules and models.

## Growing an app with `generate`

```sh
brrr g blog Post title:string:required body:text published:bool author:belongs_to:User
brrr g blog Comment body:text:required post:belongs_to:Post
make sqlc && make migrate   # regenerate db code + apply the new migration
```

Each `generate` writes a layered backend slice (`domain → repository → service →
transport/http`) plus sqlc queries and a goose migration, writes the frontend
feature (types, API client, TanStack Query hooks, list + form pages), and wires
everything in: the module registry, the regenerated `module.go`, the React
router, and the sidebar nav. Supported field types include
`string, text, int, int64, float, decimal, bool, uuid, date, datetime, json`
plus `belongs_to` relationships.

## Whole-project spec

Describe modules, models, fields and relationships in YAML and generate them all
at creation time (see `spec.example.yaml` in any generated project):

```yaml
modules:
  blog:
    models:
      Post:
        fields:
          title: string required
          body: text
          published: bool
        relationships:
          author: belongs_to User
      Comment:
        fields:
          body: text required
        relationships:
          post: belongs_to Post
          author: belongs_to User
```

```sh
brrr init myapp --module github.com/me/myapp --spec spec.yaml
```

Models are generated in dependency order (a model is created after the models it
`belongs_to`), so cross-references resolve and migrations are ordered correctly.

## Development

```sh
go test ./...                   # engine unit tests
bash scripts/verify-init.sh     # e2e: init -> build backend + frontend
bash scripts/verify-generate.sh # e2e: init + generate 2 models -> build both
bash scripts/verify-spec.sh     # e2e: init --spec (whole graph) -> build both
```
