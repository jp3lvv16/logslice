// Package dropout provides probabilistic message dropping for log streams.
//
// A Dropper is constructed with a drop rate in the range [0.0, 1.0).
// For each log line, Keep returns true if the line should be forwarded
// downstream and false if it should be silently discarded.
//
// This is useful for load-shedding in high-volume pipelines where
// processing every message is not required for statistical validity.
//
// Example:
//
//	d, err := dropout.New(0.5, time.Now().UnixNano())
//	if err != nil {
//		log.Fatal(err)
//	}
//	if d.Keep(line) {
//		// process line
//	}
package dropout
