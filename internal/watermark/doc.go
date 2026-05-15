// Package watermark provides a Watermark type that scans a stream of
// structured JSON log lines and tracks the earliest and latest values of a
// nominated timestamp field.
//
// Usage:
//
//	w, err := watermark.New("time")
//	if err != nil { ... }
//
//	for _, line := range lines {
//	    w.Record(line)
//	}
//
//	min, _ := w.Min()
//	max, _ := w.Max()
//	fmt.Println("span:", w.Span())
//
// Lines that are not valid JSON, do not contain the target field, or whose
// field value cannot be parsed as RFC 3339 are silently ignored.
package watermark
