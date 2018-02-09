package changelog

import (
	"html/template"
	"io"
)

// Render renders a changelog in Markdown to the given writer.
func Render(changelog Changelog, writer io.Writer) {
	template, err := template.New("test").Funcs(template.FuncMap{
		"deref": func(i *Release) Release { return *i },
	}).Parse(`# Changelog

{{.Description}}

## [Unreleased]
{{- if .Unreleased.Added }}

### Added
{{range .Unreleased.Added }}
- {{ .Description -}}
{{- end }}
{{- end }}
{{- if .Unreleased.Changed }}

### Changed
{{range .Unreleased.Changed }}
- {{ .Description -}}
{{- end }}
{{- end }}
{{- if .Unreleased.Deprecated }}

### Deprecated
{{range .Unreleased.Deprecated }}
- {{ .Description -}}
{{- end }}
{{- end }}
{{- if .Unreleased.Removed }}

### Removed
{{range .Unreleased.Removed }}
- {{ .Description -}}
{{- end }}
{{- end }}
{{- if .Unreleased.Fixed }}

### Fixed
{{range .Unreleased.Fixed }}
- {{ .Description -}}
{{- end }}
{{-  end }}
{{- if .Unreleased.Security }}

### Security
{{range .Unreleased.Security }}
- {{ .Description -}}
{{- end }}
{{- end }}
{{-  range .Releases }}

## [{{.Name}}] - {{.Date}}
{{- if .Added }}

### Added
{{range .Added }}
- {{ .Description -}}
{{- end }}
{{- end }}
{{- if .Changed }}

### Changed
{{range .Changed }}
- {{ .Description -}}
{{- end }}
{{- end }}
{{- if .Deprecated }}

### Deprecated
{{range .Deprecated }}
- {{ .Description -}}
{{- end }}
{{- end }}
{{- if .Removed }}

### Removed
{{range .Removed }}
- {{ .Description -}}
{{- end }}
{{- end }}
{{- if .Fixed }}

### Fixed
{{range .Fixed }}
- {{ .Description -}}
{{- end }}
{{- end }}
{{- if .Security }}

### Security
{{range .Security }}
- {{ .Description -}}
{{- end }}{{- end }}{{ end }}

[Unreleased]: {{.URL}}{{.LatestRelease.Name}}...HEAD
{{- range $i, $e := .Releases}}
{{- if $e.PreviousRelease}}
[{{$e.Name}}]: {{$.URL}}{{(deref $e.PreviousRelease).Name}}...{{$e.Name}}
{{- end}}
{{- end}}
`)
	if err != nil {
		panic(err)
	}

	err = template.Execute(writer, changelog)
	if err != nil {
		panic(err)
	}
}
