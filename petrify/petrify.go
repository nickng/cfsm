// Package petrify contains encoding for ordinary string into petrify-accepted
// format. The encoding is incomplete.
//
package petrify // import "github.com/nickng/cfsm/petrify"

import "strings"

const Tmpl = `
-- Machines #{{ .ID }}
-- {{ multiline .Comment }}
.outputs
.state graph
{{ range .Edges -}}{{ . }}{{- end }}
{{- if .Start -}}
.marking q{{ .ID }}{{ .Start.ID }}
{{- else -}}
-- Start state not set
{{- end }}
.end
`

// Encode string to petrify accepted format.
func Encode(s string) string {
	r := strings.NewReplacer("{", "LBRACE", "}", "RBRACE", ".", "DOT", "(", "LPAREN", ")", "RPAREN", "/", "SLASH")
	return r.Replace(s)
}

// Decode string from petrify-encoded format to normal text.
func Decode(s string) string {
	r := strings.NewReplacer("LBRACE", "{", "RBRACE", "}", "DOT", ".", "LPAREN", "(", "RPAREN", ")", "SLASH", "/")
	return r.Replace(s)
}
