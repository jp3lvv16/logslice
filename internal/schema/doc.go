// Package schema provides field-level validation for structured JSON log lines.
//
// A Validator is constructed from a slice of spec strings. Each spec has the
// form:
//
//	"fieldname"          – field must be present (any JSON type)
//	"fieldname:type"     – field must be present and match the given type
//
// Supported types: string, number, boolean, any.
//
// Example:
//
//	v, err := schema.New([]string{"level:string", "ts:number", "msg"})
//	if err != nil { ... }
//	if err := v.Validate(rawMsg); err != nil {
//	    // line does not conform to schema
//	}
package schema
