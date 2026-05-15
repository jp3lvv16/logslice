// Package sample implements reservoir sampling (Algorithm R) over a stream
// of structured JSON log messages.
//
// Reservoir sampling guarantees that every message in the stream has an equal
// probability of appearing in the final output, regardless of stream length.
// This is useful when you want a statistically representative subset of a
// very large log file without reading it twice.
//
// Usage:
//
//	s, err := sample.New(1000, time.Now().UnixNano())
//	if err != nil { ... }
//
//	for _, raw := range stream {
//		s.Add(raw)
//	}
//
//	for _, msg := range s.Results() {
//		fmt.Println(string(msg))
//	}
package sample
