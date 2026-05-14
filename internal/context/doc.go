// Package context enriches log lines with pipeline metadata as they flow
// through logslice's processing stages.
//
// # Metadata fields
//
// Three private fields are injected into every JSON object:
//
//   - _line        – 1-based line number within the source stream.
//   - _source      – optional filename or URI of the originating file.
//   - _ingested_at – RFC 3339 timestamp recorded at ingestion time.
//
// These fields are prefixed with an underscore to minimise collision with
// real application fields. They can be stripped before final output with
// Strip so that downstream consumers receive clean records.
//
// # Usage
//
// Typical pipeline usage injects metadata at ingestion time and strips it
// just before writing output:
//
//	record, err := context.Inject(raw, context.Meta{
//		Line:       lineNum,
//		Source:     filename,
//		IngestedAt: time.Now(),
//	})
//	if err != nil {
//		return err
//	}
//	// ... process record through pipeline stages ...
//	clean, err := context.Strip(record)
package context
