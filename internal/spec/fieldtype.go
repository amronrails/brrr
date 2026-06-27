package spec

import (
	"fmt"
	"sort"
	"strings"
)

// FieldType describes how a logical field type maps onto each layer of the
// generated stack. A single registry entry drives the Go struct field, the
// Postgres column, the sqlc-facing type and the TypeScript/zod frontend code.
type FieldType struct {
	Name      string // logical name used in the CLI / spec, e.g. "string"
	GoType    string // Go type, e.g. "string", "int64", "time.Time"
	GoImport  string // import path required by GoType, if any
	SQLType   string // Postgres column type, e.g. "text", "timestamptz"
	TSType    string // TypeScript type, e.g. "string", "number", "boolean"
	Zod       string // zod schema fragment, e.g. "z.string()"
	Nullable  bool   // whether the type is rendered as a pointer / optional by default
}

// registry holds the built-in field types. Relationship "types"
// (belongs_to, has_many, ...) are handled separately in model.go.
var registry = map[string]FieldType{
	"string":   {Name: "string", GoType: "string", SQLType: "text", TSType: "string", Zod: "z.string()"},
	"text":     {Name: "text", GoType: "string", SQLType: "text", TSType: "string", Zod: "z.string()"},
	"int":      {Name: "int", GoType: "int32", SQLType: "integer", TSType: "number", Zod: "z.number().int()"},
	"int64":    {Name: "int64", GoType: "int64", SQLType: "bigint", TSType: "number", Zod: "z.number().int()"},
	"float":    {Name: "float", GoType: "float64", SQLType: "double precision", TSType: "number", Zod: "z.number()"},
	"decimal":  {Name: "decimal", GoType: "string", SQLType: "numeric", TSType: "string", Zod: "z.string()"},
	"bool":     {Name: "bool", GoType: "bool", SQLType: "boolean", TSType: "boolean", Zod: "z.boolean()"},
	"uuid":     {Name: "uuid", GoType: "uuid.UUID", GoImport: "github.com/google/uuid", SQLType: "uuid", TSType: "string", Zod: "z.string().uuid()"},
	"date":     {Name: "date", GoType: "time.Time", GoImport: "time", SQLType: "date", TSType: "string", Zod: "z.string()"},
	"datetime": {Name: "datetime", GoType: "time.Time", GoImport: "time", SQLType: "timestamptz", TSType: "string", Zod: "z.string().datetime()"},
	"time":     {Name: "time", GoType: "time.Time", GoImport: "time", SQLType: "time", TSType: "string", Zod: "z.string()"},
	"json":     {Name: "json", GoType: "json.RawMessage", GoImport: "encoding/json", SQLType: "jsonb", TSType: "unknown", Zod: "z.unknown()"},
}

// Lookup returns the FieldType registered under name. The "enum(...)" form is
// recognised and produces a text-backed string type.
func Lookup(name string) (FieldType, error) {
	name = strings.TrimSpace(name)
	if strings.HasPrefix(name, "enum(") && strings.HasSuffix(name, ")") {
		return FieldType{Name: name, GoType: "string", SQLType: "text", TSType: "string", Zod: "z.string()"}, nil
	}
	ft, ok := registry[name]
	if !ok {
		return FieldType{}, fmt.Errorf("unknown field type %q (known: %s)", name, strings.Join(KnownTypes(), ", "))
	}
	return ft, nil
}

// KnownTypes returns the sorted list of registered type names, for error
// messages and CLI help.
func KnownTypes() []string {
	names := make([]string, 0, len(registry))
	for n := range registry {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}
