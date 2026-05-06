// Package ratelimit implements a token-bucket rate limiter for controlling
// the throughput of log lines written to output.
//
// Usage:
//
//	l, err := ratelimit.New(500) // allow up to 500 lines/sec
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, line := range lines {
//		if l.Allow() {
//			fmt.Println(line)
//		}
//	}
//
// For pipelines that should never be throttled, use Unlimited():
//
//	l := ratelimit.Unlimited()
package ratelimit
