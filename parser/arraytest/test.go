package arraytest

// @xerox
type Array struct {
	infos      []float
	infoptrs   []*int
	structs    []AnotherData
	structptrs []*AnotherData
}

type AnotherData struct {
	TEST    string
	nesteds []*DeepNested
}

type DeepNested struct {
	info    *string
	infos   []string
	nesteds []*UltraDeep
}

type UltraDeep struct {
	level string
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
		for _, elt1 := range elt.nesteds {
			newElt1 := new(DeepNested)
			if elt1.info != nil {
				val := *elt1.info
				newElt1.info = &val
			}
			for _, elt2 := range elt1.infos {
				newElt1.infos = append(newElt1.infos, elt2)
			}
			for _, elt2 := range elt1.nesteds {
				if elt2 != nil {
					newElt2 := new(UltraDeep)
					newElt2.level = elt2.level
					newElt1.nesteds = append(newElt1.nesteds, newElt2)
				} else {
					newElt1.nesteds = append(newElt1.nesteds, nil)
				}
			}
			newElt.nesteds = append(newElt.nesteds, newElt2)
		}
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
