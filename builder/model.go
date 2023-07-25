package builder

type Struct struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name string
	Kind Kind
	Tags string
}
