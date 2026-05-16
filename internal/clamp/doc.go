// Package clamp provides a Clamper that constrains a numeric JSON field to a
// closed interval [min, max].
//
// # Usage
//
//	c, err := clamp.New("latency_ms", 0, 30_000)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	out, err := c.Apply(msg)
//
// Values below min are raised to min; values above max are lowered to max.
// Fields that are absent or contain non-numeric JSON values are left unchanged.
package clamp
