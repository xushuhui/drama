package main

var template = `package schema

import (
    "time"

    "entgo.io/ent"
    "entgo.io/ent/schema/field"
)


type %s struct {
    ent.Schema
}

// Fields of the user.
func (%s) Fields() []ent.Field {
    return []ent.Field
}`

var fields = `
	{
        field.Int("age"),
        field.String("name"),
        field.String("username").Unique(),
        field.Time("created_at").Default(time.Now),
    }
`
