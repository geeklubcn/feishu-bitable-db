package db

type FieldType int

const (
	String FieldType = iota + 1
	Int
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
