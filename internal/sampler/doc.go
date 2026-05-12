// Package sampler implements probabilistic log line sampling for logslice.
//
// Overview
//
// When processing very high-volume log streams it is often impractical to
// inspect every line. The sampler package provides a lightweight mechanism
// to stochastically retain only a configurable fraction of lines before
// they are passed downstream to the filter pipeline or aggregator.
//
// Usage
//
//	// Keep roughly 10 % of lines using a random source seeded from time.
//	src := rand.NewSource(time.Now().UnixNano())
//	s, err := sampler.New(0.10, src)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	for _, line := range lines {
//		if s.Keep(line) {
//			// process line
//		}
//	}
//
// A rate of 1.0 disables sampling and every line is kept.
package sampler
