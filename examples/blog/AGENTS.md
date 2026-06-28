# AGENTS.md

Guidance for AI coding agents (and humans) working in **Blog**.
Read this before making changes — it explains how the project is generated, how
it is organized, and the conventions you must preserve.

## What this project is

`blog` (`github.com/example/blog`) is a Go **modular monolith** with a React
admin frontend, scaffolded by [brrr](https://github.com/amronrails/brrr) — a
CRUD code generator. brrr produced the initial structure; you keep building on
it either by hand or by running brrr again:

```sh
brrr generate <module> <Model> field:type ...   # add a CRUD slice (backend + frontend)
brrr init <app> --spec spec.yaml                 # scaffold a whole app from a YAML spec
```

`brrr.yaml` at the repo root is the manifest brrr maintains: it lists every
module and model. Treat it as generator-owned state.

## Stack

- **Backend:** Go · chi · sqlc + pgx · Postgres · goose · JWT (golang-jwt v5) · bcrypt · go-playground/validator
- **Frontend (`web/`):** Vite · TypeScript · TanStack Query · React Router · shadcn-style UI + Tailwind

## Layout

```
cmd/api/                     # entrypoint: config → db → modules → HTTP server
internal/
  db/                        # sqlc-GENERATED query layer — DO NOT EDIT
  modules/<name>/            # one vertical slice per feature
    module.go                # PUBLIC facade: wires the slice, mounts routes
    internal/                # private to the module (Go enforces this)
      domain/                # entities, value objects, errors — pure, no I/O
      ports/                 # interfaces the services depend on (repos, gateways)
      services/              # use cases / application logic
      adapters/
        http/                # chi handlers, DTOs, route registration
        postgres/            # repository: maps internal/db rows ↔ domain entities
  modules/registry.go        # lists every module for the server to mount
  platform/                  # cross-cutting: auth, config, database, httpx, server
db/
  migrations/                # goose SQL migrations (also the sqlc schema)
  queries/                   # hand-written sqlc queries (*.sql)
sqlc.yaml                    # one centralized config → internal/db
web/                         # React admin (see web/ section below)
```

### The dependency rule

Dependencies point **inward**: `adapters → services → ports → domain`. The
domain knows nothing about HTTP or SQL. Services depend on the `ports`
interfaces, never on a concrete adapter. The `postgres` adapter implements a
`ports` interface and is the only place that imports `internal/db`.

Each module's `internal/` directory is enforced by the Go compiler: nothing
outside `internal/modules/<name>/` can import its inner packages. The only public
surface is the module's root package (`module.go` + `api.go`), consumed by
`internal/modules/registry.go`.

## Inter-module communication

Modules talk to each other **only through public interfaces**, never by reaching
into another module's `internal/` (the compiler forbids it). There is no event
bus — calls are direct, synchronous, and compile-time checked.

- **Every module publishes an `API`** in its root package (`api.go`): an
  interface plus public DTO types (e.g. `user.API` with `user.User`). `*Module`
  implements it. DTOs are deliberately separate from domain entities and HTTP
  responses — they are the stable cross-module surface.
- **A consumer depends on the interface**, not the concrete module. Add the
  dependency as a field on the consumer's `Deps` (optionally narrow it behind an
  interface in the consumer's `ports` package).
- **Wiring happens in the composition root**, `internal/modules/registry.go`,
  which constructs modules in two phases: providers first (named variables),
  then consumers receive them.

Example — make the `blog` module look up authors via the user module:

```go
// internal/modules/blog/module.go   (consumer accepts the interface)
type Deps struct {
    Queries   *db.Queries
    Tokens    *auth.TokenService
    Validator *validator.Validate
    Users     user.API            // <- dependency on user's public interface
}

// internal/modules/registry.go      (composition root wires it)
func New(d Deps) []Module {
    userMod := user.New(user.Deps{...})           // provider
    blogMod := blog.New(blog.Deps{..., Users: userMod}) // userMod satisfies user.API
    return []Module{userMod, blogMod}
}
```

`brrr generate` already builds each module into a named `<name>Mod` variable and
publishes its `API`, so wiring a new edge is just adding the field and passing
the variable. Keep providers constructed before their consumers.

## Conventions to preserve

- **sqlc is centralized.** All generated DB code lives in one package,
  `internal/db`, generated from `db/migrations` (schema) + `db/queries/*.sql`.
  Never edit `internal/db` by hand. Add a migration + a `db/queries/<model>.sql`,
  then run `make sqlc`. A `postgres` adapter wraps `*db.Queries` and maps rows to
  domain types.
- **Marker comments are injection points.** brrr appends generated wiring at
  these markers; keep them intact and in place:
  - `// brrr:modules` and `// brrr:module-imports` in `internal/modules/registry.go`
  - `// brrr:imports-fe` and `// brrr:routes-fe` in `web/src/router.tsx`
  - `{/* brrr:nav */}` in `web/src/components/layout/DashboardLayout.tsx`
- **`module.go` is regenerated** by `brrr generate`. Put per-model code in the
  layer files under `internal/`, not in `module.go`.
- **Build order matters:** run `sqlc generate` (`make sqlc`) before `go build` /
  `go mod tidy`, because the adapters import the not-yet-generated `internal/db`.

## web/ (frontend)

Feature-sliced: each feature owns its slice under `web/src/features/<feature>/`
and exposes a public API through an `index.ts` barrel. Import a feature from its
barrel (`@/features/users`), not via deep paths; use relative imports *inside* a
feature.

```
web/src/
  features/<feature>/        # api.ts, hooks.ts (TanStack Query), types.ts, *Page.tsx, index.ts
  components/ui/             # shadcn-style primitives
  components/layout/         # app shell (DashboardLayout)
  lib/                       # api client, query client, utils
  router.tsx                 # routes (brrr injects generated routes here)
```

## Common commands

```sh
make db-up        # start Postgres (docker compose)
make setup        # go mod tidy + sqlc generate + npm install
make migrate      # apply migrations (goose)
make dev          # run backend (:8080) and frontend (:5173) together
make sqlc         # regenerate internal/db after changing migrations/queries
make test         # go test ./...
```

The first account to register becomes the **admin**.

## Checklist before finishing a change

- [ ] `make sqlc` run if you touched `db/migrations` or `db/queries`
- [ ] `go build ./...` and `go vet ./...` pass
- [ ] `cd web && npm run build` type-checks and builds
- [ ] Marker comments and `brrr.yaml` left intact
- [ ] New code respects the inward dependency rule (no `internal/db` outside the postgres adapter)
