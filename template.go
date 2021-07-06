package main

var TemplateMain = `package schema

import (
    "entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type %s struct {
    ent.Schema
}

func (%s) Fields() []ent.Field {
    return []ent.Field{
        %s
    }
}`

var DataTypeMap = map[string]string{
	"int": "Int", "tinyint": "Int8", "bigint": "Int", "varchar": "String", "char": "String", "datetime": "Time", "date": "Time",
	"text": "Text", "decimal": "Float", "longtext": "Text",
}
