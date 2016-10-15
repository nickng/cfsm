// Package petrify contains encoding for ordinary string into petrify-accepted
// format. The encoding is incomplete.
//
package petrify // import "github.com/nickng/cfsm/petrify"

import "strings"

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
