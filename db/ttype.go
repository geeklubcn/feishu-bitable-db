package db

type FieldType int

const (
	String      FieldType = 1
	Int         FieldType = 2
	Radio       FieldType = 3
	MultiSelect FieldType = 4
	Date        FieldType = 5
	People      FieldType = 11
)

const (
	ID   = "id"
	Name = "Databases"
)

type Database struct {
	Name   string
	Tables []Table
}

type Table struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name string
	Type FieldType
}

type SearchCmd struct {
	Key, Operator string
	Val           interface{}
}
