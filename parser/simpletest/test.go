package simpletest

// @xerox
type Simple struct {
	info    string
	infoPtr *string
}

/*
package simpletest

func XeroxSimple(sample Simple) Simple {
	copied := Simple{}
	copied.info = sample.info
	if sample.infoPtr != nil {
		val := *sample.infoPtr
		copied.infoPtr = &val
	}

	return copied
}
*/
