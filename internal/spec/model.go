package spec

// RelationKind enumerates the supported relationship kinds between models.
type RelationKind string

const (
	BelongsTo  RelationKind = "belongs_to"
	HasMany    RelationKind = "has_many"
	HasOne     RelationKind = "has_one"
	ManyToMany RelationKind = "many_to_many"
)

// IsRelation reports whether a raw type token denotes a relationship rather than
// a scalar field type.
func IsRelation(token string) bool {
	switch RelationKind(token) {
	case BelongsTo, HasMany, HasOne, ManyToMany:
		return true
	}
	return false
}

// Field is a single scalar attribute of a model.
type Field struct {
	Name     string    // logical name as written, e.g. "title"
	Type     FieldType // resolved type metadata
	Required bool      // NOT NULL + required in validation/zod
	Unique   bool      // UNIQUE constraint
}

// Relationship connects a model to another model.
type Relationship struct {
	Name   string       // field/accessor name, e.g. "author"
	Kind   RelationKind // belongs_to, has_many, ...
	Target string       // target model name, e.g. "User"
}

// Model is a single entity belonging to a module.
type Model struct {
	Module        string         // owning module, e.g. "blog"
	Name          string         // entity name in singular PascalCase, e.g. "Post"
	Fields        []Field        // scalar fields
	Relationships []Relationship // associations
}
