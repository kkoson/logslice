// Package labelmap extracts named capture groups from log lines and maps
// them to key=value label pairs suitable for structured logging or metric
// tagging pipelines.
//
// Usage:
//
//	m, err := labelmap.New(`(?P<level>\w+)\s+(?P<msg>.+)`, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	labels, ok := m.Map(line)
//	if ok {
//		fmt.Println(labels["level"], labels["msg"])
//	}
package labelmap
