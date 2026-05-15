// Package grok provides named-pattern log parsing inspired by the Logstash
// grok filter.
//
// # Overview
//
// A Parser is created from a pattern string containing %{NAME} tokens. Each
// token is expanded to a named capturing group using the built-in pattern
// library (IP, TIMESTAMP, LOGLEVEL, WORD, NUMBER, GREEDYDATA) or caller-
// supplied extras. The resulting regexp is compiled once and reused across
// calls to Parse.
//
// # Usage
//
//	p, err := grok.New("%{TIMESTAMP} %{LOGLEVEL} %{GREEDYDATA}", nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fields := p.Parse(line) // nil when line does not match
//
// # Custom patterns
//
// Pass a map[string]string of name → raw-regexp pairs as the second argument
// to New. Custom patterns may use named capturing groups; the group name
// becomes the field key returned by Parse.
package grok
