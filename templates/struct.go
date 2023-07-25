type {{ .Name }} struct {
	{{- range .Fields }}
	{{ .Name }} {{ .Kind }} {{ .Tags }}
	{{- end }}
}

