# brrr examples

Three projects scaffolded entirely with the `brrr` tool, each from a single YAML
spec (`<name>-spec.yaml`) via `brrr init <name> --spec <name>-spec.yaml`. They are
kept as pristine generator output — run `make setup` in any of them to install
deps and generate the sqlc layer, then `make migrate && make dev`.

| Project | Modules | Models | Exercises |
| ------- | ------- | ------ | --------- |
| **blog** | `blog` | Post, Comment | text/bool/int fields, intra-module FK (Comment→Post) + `belongs_to User` |
| **shop** | `catalog`, `sales` | Category, Product, Order, OrderItem | multi-module, `decimal`, `json`, unique `sku`, cross-module FK (OrderItem→Product) |
| **tasks** | `projects` | Project, Task, Label | `date` field, `int` priority, unique `key`/`name` |

Each project is a Go modular monolith (chi · sqlc + pgx · Postgres · goose ·
JWT) with a Vite/React admin in `web/`, following the ports & adapters layout
documented in the generated `AGENTS.md`.

## Regenerate from scratch

```sh
brrr init blog  --module github.com/example/blog  --spec blog-spec.yaml  --dir ./blog
brrr init shop  --module github.com/example/shop  --spec shop-spec.yaml  --dir ./shop
brrr init tasks --module github.com/example/tasks --spec tasks-spec.yaml --dir ./tasks
```

## Validity

All three were verified end-to-end: backend `sqlc generate` → `go build ./...` →
`go vet ./...`, frontend `tsc -b && vite build`, and (for shop and tasks) live
runtime against Postgres — migrations, auth, CRUD, decimal/json round-trips, the
`date` field, and unique-violation → HTTP 409.
