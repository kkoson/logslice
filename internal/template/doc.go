// Package template reformats log lines by combining a named-group regex with
// a Go text/template string.
//
// # Overview
//
// A [Renderer] is created with a regex pattern and a template string. When
// [Renderer.Apply] is called on a log line:
//
//  1. The regex is matched against the line.
//  2. Named capture groups are collected into a map[string]string.
//  3. The template is executed with that map as its data object.
//  4. The rendered string is returned.
//
// Lines that do not match the pattern are returned unchanged, making the
// renderer safe to use in a mixed-format log stream.
//
// # Example
//
//	pattern := `(?P<ts>[\d:T-]+)\s+(?P<level>\w+)\s+(?P<msg>.*)`
//	tmplStr := `{{.ts}} | {{.level}} | {{.msg}}`
//	r, err := template.New(pattern, tmplStr)
//	if err != nil { /* handle */ }
//	out, err := r.Apply("2024-01-02T15:04:05 INFO server started")
package template
