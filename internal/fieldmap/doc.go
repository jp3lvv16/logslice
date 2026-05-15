// Package fieldmap provides a Mapper that renames JSON object keys according
// to a set of caller-supplied "old:new" specs.
//
// # Usage
//
//	specs := []string{"level:severity", "msg:message"}
//	m, err := fieldmap.New(specs)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	for _, line := range lines {
//		fmt.Println(string(m.Apply(line)))
//	}
//
// # Spec format
//
// Each spec is a colon-separated pair:  old:new
//
// If the source key is absent in a given message the message is passed through
// unchanged.  If the same source key appears in multiple specs the last one
// wins.
package fieldmap
