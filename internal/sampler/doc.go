// Package sampler implements deterministic rate-based sampling for log lines.
//
// Usage:
//
//	s, err := sampler.New(10)  // keep every 10th line
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, line := range lines {
//		if s.Keep() {
//			// process line
//		}
//	}
//
// A rate of 1 disables sampling (every line is kept). Rates less than 1
// are rejected with an error.
//
// The sampler is not safe for concurrent use without external synchronisation.
package sampler
