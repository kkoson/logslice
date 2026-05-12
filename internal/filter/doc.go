// Package filter provides regex-based log line filtering with support for
// multi-stage pipelines.
//
// A Pipeline is an ordered sequence of Filter instances. Each Filter wraps a
// compiled regular expression and an optional invert flag:
//
//   - When invert is false the filter passes lines that MATCH the pattern.
//   - When invert is true  the filter passes lines that do NOT match the pattern.
//
// A log line is emitted only when every filter in the pipeline returns true,
// making it straightforward to compose include and exclude rules.
//
// Example usage:
//
//	pipeline := filter.NewPipeline()
//
//	include, err := filter.NewFilter(`ERROR|WARN`, false)
//	if err != nil { log.Fatal(err) }
//	pipeline.Add(include)
//
//	exclude, err := filter.NewFilter(`healthcheck`, true)
//	if err != nil { log.Fatal(err) }
//	pipeline.Add(exclude)
//
//	if pipeline.Match(line) {
//		fmt.Println(line)
//	}
package filter
