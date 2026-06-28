// Package gen turns a parsed model specification into the rich, fully-derived
// "views" consumed by the generate templates, and orchestrates writing and
// wiring the resulting files into an existing project.
package gen

import (
	"fmt"
	"strings"

	"github.com/amronrails/brrr/internal/engine"
	"github.com/amronrails/brrr/internal/spec"
)

// Column is a single persisted, user-supplied attribute of a model (a scalar
// field or a belongs-to foreign key). It carries everything the backend and
// frontend templates need, so the templates stay logic-light.
type Column struct {
	Name       string // snake_case column / json key, e.g. "title", "author_id"
	GoName     string // PascalCase Go field, e.g. "Title", "AuthorID"
	GoType     string // Go type, e.g. "string", "bool", "uuid.UUID"
	SQLType    string // Postgres column type
	Default    string // SQL default expression for optional columns, else ""
	Required   bool
	Unique     bool   // emit a UNIQUE constraint
	IsRelation bool   // true for belongs-to foreign keys
	RefTable   string // referenced table for relation columns, e.g. "users"

	TSType      string // TypeScript type: "string" | "number" | "boolean"
	TSOptional  bool   // render the TS property as optional (e.g. json -> z.unknown())
	TSDefault   string // TypeScript default form value: "\"\"" | "0" | "false"
	Zod         string // zod schema fragment
	InputType   string // HTML input type: "text" | "number" | "checkbox"
	ValidateTag string // go-playground validate tag, e.g. "required" or ""
	Label       string // UI label
}

// ModelView is the template data for generating one model.
type ModelView struct {
	ModulePath string // Go module path of the project
	Module     string // raw module name, e.g. "blog"
	ModulePkg  string // Go package name for the module, e.g. "blog"
	ModuleKebab string
	HTTPPkg    string // transport package name, e.g. "bloghttp"

	Model        string // raw model name
	Pascal       string // "Post"
	Camel        string // "post"
	Snake        string // "post"
	Plural       string // "posts" (snake)
	PluralPascal string // "Posts"
	KebabPlural  string // "posts"
	Table        string // "posts"
	RoutePrefix  string // "blog/posts"

	Columns []Column // user columns (fields + relation FKs), in declaration order
	HasJSON bool     // whether any column needs encoding/json

	MigrationName string // e.g. "0003_create_posts"
}

// ModelRef is a lightweight reference used when regenerating a module's wiring.
type ModelRef struct {
	Pascal string
	Camel  string
}

// ModuleView is the template data for (re)generating a module's module.go.
type ModuleView struct {
	ModulePath string
	Module     string
	ModulePkg  string
	HTTPPkg    string
	Models     []ModelRef
}

// BuildModelView derives a ModelView from a parsed spec.Model.
func BuildModelView(modulePath string, m *spec.Model, migrationSeq int) (*ModelView, error) {
	pkg := pkgName(m.Module)
	v := &ModelView{
		ModulePath:  modulePath,
		Module:      m.Module,
		ModulePkg:   pkg,
		ModuleKebab: engine.Kebab(m.Module),
		HTTPPkg:     pkg + "http",
		Model:       m.Name,
		Pascal:      engine.Pascal(m.Name),
		Camel:       engine.Camel(m.Name),
		Snake:       engine.Snake(m.Name),
		Plural:      engine.Plural(engine.Snake(m.Name)),
		PluralPascal: engine.Plural(engine.Pascal(m.Name)),
		KebabPlural: engine.Plural(engine.Kebab(m.Name)),
		Table:       engine.Plural(engine.Snake(m.Name)),
	}
	v.RoutePrefix = v.ModuleKebab + "/" + v.KebabPlural
	v.MigrationName = fmt.Sprintf("%04d_create_%s", migrationSeq, v.Table)

	for _, f := range m.Fields {
		col := Column{
			Name:      engine.Snake(f.Name),
			GoName:    engine.Pascal(f.Name),
			GoType:    f.Type.GoType,
			SQLType:   f.Type.SQLType,
			Required:  f.Required,
			Unique:    f.Unique,
			TSType:    f.Type.TSType,
			TSOptional: f.Type.TSType == "unknown", // z.unknown()/z.any() infer optional
			TSDefault: tsDefault(f.Type.TSType),
			Zod:       zodFor(f.Type),
			InputType: inputTypeFor(f.Type),
			Label:     engine.Title(f.Name),
		}
		col.Default = sqlDefault(col)
		col.ValidateTag = validateTag(f.Required, f.Type.TSType)
		if f.Type.GoImport == "encoding/json" {
			v.HasJSON = true
		}
		v.Columns = append(v.Columns, col)
	}

	for _, r := range m.Relationships {
		if r.Kind != spec.BelongsTo {
			return nil, fmt.Errorf("relationship %q (%s): only belongs_to is supported by generate today", r.Name, r.Kind)
		}
		col := Column{
			Name:       engine.Snake(r.Name) + "_id",
			GoName:     engine.Pascal(r.Name) + "ID",
			GoType:     "uuid.UUID",
			SQLType:    "uuid",
			Required:   true,
			IsRelation: true,
			RefTable:   engine.Plural(engine.Snake(r.Target)),
			TSType:      "string",
			TSDefault:   `""`,
			Zod:         "z.string().uuid()",
			InputType:   "text",
			ValidateTag: "required",
			Label:       engine.Title(r.Name) + " ID",
		}
		v.Columns = append(v.Columns, col)
	}

	return v, nil
}

// pkgName produces a valid lowercase Go package identifier from a module name.
func pkgName(module string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(module) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		}
	}
	out := b.String()
	if out == "" {
		out = "module"
	}
	return out
}

func zodFor(ft spec.FieldType) string {
	switch ft.TSType {
	case "number":
		if strings.Contains(ft.Zod, "int") {
			return "z.coerce.number().int()"
		}
		return "z.coerce.number()"
	case "boolean":
		return "z.boolean()"
	default:
		return ft.Zod
	}
}

func tsDefault(tsType string) string {
	switch tsType {
	case "number":
		return "0"
	case "boolean":
		return "false"
	default:
		return `""`
	}
}

func validateTag(required bool, tsType string) string {
	if required && tsType != "boolean" {
		return "required"
	}
	return ""
}

func inputTypeFor(ft spec.FieldType) string {
	switch ft.TSType {
	case "number":
		return "number"
	case "boolean":
		return "checkbox"
	default:
		return "text"
	}
}

func sqlDefault(c Column) string {
	if c.Required || c.IsRelation {
		return ""
	}
	switch c.SQLType {
	case "boolean":
		return "false"
	case "integer", "bigint", "double precision", "numeric":
		return "0"
	case "text":
		return "''"
	case "timestamptz":
		return "now()"
	case "date":
		return "CURRENT_DATE"
	case "time":
		return "CURRENT_TIME"
	case "jsonb":
		return "'{}'::jsonb"
	default:
		return ""
	}
}
