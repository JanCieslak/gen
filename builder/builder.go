package builder

import (
	"fmt"
	"strings"
)

//// Interfaces

//// Builder

func QueryStruct(name string) *StructBuilder {
	return &StructBuilder{name: name}
}

type StructBuilder struct {
	name   string
	fields []FieldBuilder
}

func (sb *StructBuilder) AsFieldBuilder(fieldName string, options ...option) *FieldBuilder {
	f := &FieldBuilder{
		Name: fieldName,
		Kind: sb.name,
		Tags: make(map[string][]string),
	}
	for _, opt := range options {
		opt(sb, f)
	}
	return f
}

//// Fields

type Kind = string

const (
	KindBoolPtr Kind = "*bool"
	KindBool    Kind = "bool"
)

type FieldBuilder struct {
	Name string
	Kind Kind
	Tags map[string][]string
}

func (sb *StructBuilder) addField(f FieldBuilder) {
	sb.fields = append(sb.fields, f)
}

//// Groups

type Group struct {
	*StructBuilder
}

func (sb *StructBuilder) OneOfGroupped(builders ...*StructBuilder) *StructBuilder {
	for _, b := range builders {
		for _, f := range b.fields {
			sb.addField(f)
		}
	}

	return sb
}

func FieldGroup(fb ...FieldBuilder) *Group {
	return &Group{
		StructBuilder: &StructBuilder{
			fields: fb,
		},
	}
}

//// Options

type option func(sb *StructBuilder, f *FieldBuilder)

func WithTag(tagName string, tagValue string) func(sb *StructBuilder, f *FieldBuilder) {
	return func(sb *StructBuilder, f *FieldBuilder) {
		if _, ok := f.Tags[tagName]; !ok {
			f.Tags[tagName] = make([]string, 0)
		}
		f.Tags[tagName] = append(f.Tags[tagName], tagValue)
	}
}

func WithSQL(sql string) option {
	return WithTag("sql", sql)
}

func WithType(typeKind Kind) func(sb *StructBuilder, f *FieldBuilder) {
	return func(sb *StructBuilder, f *FieldBuilder) {
		f.Kind = typeKind
	}
}

func WithTypeBoolPtr() option {
	return WithType(KindBoolPtr)
}

func SingleQuoted() option {
	return WithTag("dll", "single_quotes")
}

//// Building functions

func (sb *StructBuilder) Field(dllValue string, fieldName string, options ...option) *StructBuilder {
	f := FieldBuilder{Name: fieldName}
	f.Tags = make(map[string][]string)
	f.Tags["dll"] = make([]string, 0)
	f.Tags["dll"] = append(f.Tags["dll"], dllValue)
	for _, opt := range options {
		opt(sb, &f)
	}
	sb.addField(f)
	return sb
}

// Static
func (sb *StructBuilder) Static(fieldName string, options ...option) *StructBuilder {
	options = append(options, WithType(KindBool))
	return sb.Field("static", fieldName, options...)
}

func (sb *StructBuilder) Create() *StructBuilder {
	return sb.Static("create")
}

func (sb *StructBuilder) Alter() *StructBuilder {
	return sb.Static("alter")
}

// Keyword
func (sb *StructBuilder) Keyword(fieldName string, options ...option) *StructBuilder {
	return sb.Field("keyword", fieldName, options...)
}

func (sb *StructBuilder) OrReplace(options ...option) *StructBuilder {
	options = append(options, WithTypeBoolPtr())
	options = append(options, WithSQL("OR REPLACE"))
	return sb.Keyword("OrReplace", options...)
}

func (sb *StructBuilder) IfExists(options ...option) *StructBuilder {
	options = append(options, WithTypeBoolPtr())
	options = append(options, WithSQL("IF EXISTS"))
	return sb.Keyword("IfExists", options...)
}

func (sb *StructBuilder) Transient(options ...option) *StructBuilder {
	options = append(options, WithTypeBoolPtr())
	return sb.Keyword("Transient", options...)
}

// Identifier
func (sb *StructBuilder) Identifier(fieldName string, options ...option) *StructBuilder {
	return sb.Field("identifier", fieldName, options...)
}

func (sb *StructBuilder) AccountObjectIdentifier() *StructBuilder {
	return sb.Identifier("name", WithType("AccountObjectIdentifier"))
}

func (sb *StructBuilder) SchemaIdentifier() *StructBuilder {
	return sb.Identifier("name", WithType("SchemaIdentifier"))
}

// Variants

func (sb *StructBuilder) OneOf(fieldBuilders ...*FieldBuilder) *StructBuilder {

	return sb
}

// Builder

func (sb *StructBuilder) Build() *Struct {
	fields := make([]Field, len(sb.fields))
	for i, f := range sb.fields {
		tagsBuilder := strings.Builder{}
		tagsBuilder.WriteRune('`')
		tags := make([]string, 0)
		for k, v := range f.Tags {
			tagValues := strings.Join(v, ",")
			tags = append(tags, fmt.Sprintf("%s:\"%s\"", k, tagValues))
		}
		tagsBuilder.WriteString(strings.Join(tags, ", "))
		tagsBuilder.WriteRune('`')

		fields[i] = Field{
			Name: f.Name,
			Kind: f.Kind,
			Tags: tagsBuilder.String(),
		}
	}
	return &Struct{
		Name:   sb.name,
		Fields: fields,
	}
}
