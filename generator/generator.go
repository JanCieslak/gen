package generator

import (
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/sf/go-gen/builder"
)

func GenerateAll(structs ...*builder.Struct) {
	fileTemplateContent, _ := ioutil.ReadFile("templates/file.go")
	fileTemplate, _ := template.New("file").Parse(string(fileTemplateContent))

	structTemplateContent, _ := ioutil.ReadFile("templates/struct.go")
	structTemplate, _ := template.New("struct").Parse(string(structTemplateContent))

	structBuilder := strings.Builder{}
	for _, s := range structs {
		_ = structTemplate.Execute(&structBuilder, s)
	}

	_ = fileTemplate.Execute(os.Stdout, structBuilder.String())
}
