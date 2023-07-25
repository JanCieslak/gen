package main

import (
	b "github.com/sf/go-gen/builder"
	"github.com/sf/go-gen/generator"
)

func main() {
	setRole := b.QueryStruct("SetRole").
		OneOfGroupped(
			b.FieldGroup().
				Static("hey").
				Static("wassup"),
			b.FieldGroup().
				Static("bye").
				Static("im_other_field"),
		).
		Static("normal")

	unsetRole := b.QueryStruct("UnsetRole")

	// Generate from builders
	generator.GenerateAll(
		b.QueryStruct("AlterRoleOptions").
			Alter().
			Static("role").
			IfExists(b.SingleQuoted()).
			AccountObjectIdentifier().
			OneOf(
				setRole.AsFieldBuilder("Set", b.WithTag("dll", "keyword")),
				unsetRole.AsFieldBuilder("Unset", b.WithTag("dll", "keyword")),
			).
			Build(),
		setRole.Build(),
		unsetRole.Build(),
	)
}
