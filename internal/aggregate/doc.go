// Package aggregate provides a simple field-frequency counter for structured
// (JSON) log lines.
//
// Usage:
//
//	c := aggregate.New("level")
//	for _, raw := range logLines {
//		c.Add(raw)
//	}
//	c.Print(os.Stdout)
//
// Output example:
//
//	   142  info
//	    37  error
//	    12  warn
//	     1  <missing>
//
// Values that are not JSON strings are represented as their raw JSON
// encoding (e.g. numbers, booleans).  Lines that are not valid JSON or
// that do not contain the requested field are counted under the synthetic
// key "<missing>".
package aggregate
