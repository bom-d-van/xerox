package arraytest

// @xerox
type Array struct {
	infos      []float
	infoptrs   []*int
	structs    []AnotherData
	structptrs []*AnotherData
}

type AnotherData struct {
	TEST string
}

/*
package arraytest

func XeroxArray(sample Array) Array {
	copied := Array{}
	for _, elt := range sample.infos {
		copied.infos = append(copied.infos, elt)
	}
	for _, elt := range sample.infoptrs {
		if elt != nil {
			newElt := *elt
			copied.infoptrs = append(copied.infoptrs, &newElt)
		} else {
			copied.infoptrs = append(copied.infoptrs, nil)
		}
	}
	for _, elt := range sample.structs {
		newElt := AnotherData{}
		newElt.TEST = elt.TEST
		copied.structs = append(copied.structs, newElt)
	}
	for _, elt := range sample.structptrs {
		if elt != nil {
			newElt := new(AnotherData)
			newElt.TEST = elt.TEST
			copied.structptrs = append(copied.structptrs, newElt)
		} else {
			copied.structptrs = append(copied.structptrs, nil)
		}
	}

	return copied
}
*/
