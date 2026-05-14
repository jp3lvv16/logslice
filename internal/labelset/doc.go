// Package labelset attaches a fixed set of static key/value labels to
// every structured log message that passes through the pipeline.
//
// # Usage
//
//	l, err := labelset.New([]string{
//		"env=production",
//		"region=us-east-1",
//		"service=api",
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	outMsg, err := l.Apply(rawMsg)
//
// Labels are injected only when the key is absent from the incoming
// message, so application-level fields always take precedence over
// the configured defaults.
package labelset
