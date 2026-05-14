// Package transform provides regex-based line transformation for logslice.
//
// A Transformer compiles a regular expression and a replacement template at
// construction time and then applies them to individual log lines via Apply.
// Named capture groups captured by the pattern can be referenced in the
// replacement string using the $name or ${name} syntax understood by
// regexp.ReplaceAllString.
//
// Multiple Transformers can be composed into a sequential pipeline with
// Chain, which feeds the output of each step as the input to the next.
//
// Typical usage:
//
//	tr, err := transform.New(`(?P<level>INFO|WARN|ERROR)`, "[${level}]")
//	if err != nil {
//		log.Fatal(err)
//	}
//	output := tr.Apply(line)
//
// For a multi-step pipeline:
//
//	output := transform.Chain(transformers, line)
package transform
