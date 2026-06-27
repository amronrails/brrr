package spec

import (
	"fmt"
	"strings"
)

// ParseModel builds a Model from CLI-style field arguments. Each argument is of
// the form:
//
//	name:type[:modifier...]        scalar field, e.g. title:string:required
//	name:belongs_to:Target         relationship, e.g. author:belongs_to:User
//
// Supported scalar modifiers are "required" and "unique".
func ParseModel(module, name string, args []string) (*Model, error) {
	m := &Model{Module: module, Name: name}
	for _, arg := range args {
		parts := strings.Split(arg, ":")
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid field %q: expected name:type", arg)
		}
		fieldName, typeTok := parts[0], parts[1]
		if fieldName == "" {
			return nil, fmt.Errorf("invalid field %q: empty name", arg)
		}

		if IsRelation(typeTok) {
			if len(parts) < 3 || parts[2] == "" {
				return nil, fmt.Errorf("relationship %q requires a target: name:%s:Target", arg, typeTok)
			}
			m.Relationships = append(m.Relationships, Relationship{
				Name:   fieldName,
				Kind:   RelationKind(typeTok),
				Target: parts[2],
			})
			continue
		}

		ft, err := Lookup(typeTok)
		if err != nil {
			return nil, fmt.Errorf("field %q: %w", arg, err)
		}
		f := Field{Name: fieldName, Type: ft, Required: false}
		for _, mod := range parts[2:] {
			switch mod {
			case "required":
				f.Required = true
			case "unique":
				f.Unique = true
			case "":
				// ignore trailing colon
			default:
				return nil, fmt.Errorf("field %q: unknown modifier %q", arg, mod)
			}
		}
		m.Fields = append(m.Fields, f)
	}
	return m, nil
}
