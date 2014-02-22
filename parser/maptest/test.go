package maptest

// @xerox
type Map struct {
	info       map[int]string
	infoptr    map[int]*string
	structm    map[string]AnotherData
	structmptr map[string]*AnotherData
	// infoPtr *map[int]string
}

type AnotherData struct {
	TEST    string
	nesteds map[string]*DeepNested
}

type DeepNested struct {
	info    *string
	infos   map[float64]string
	nesteds map[string]*UltraDeep

	// nestedArrayVal map[init]*[]UltraDeep
}

type UltraDeep struct {
	level string
}

/*
package maptest

func XeroxMap(sample Map) Map {
	copied := Map{}
	for key, value := range sample.info {
		copied.info[key] = value
	}
	for key, value := range sample.infoptr {
		if value != nil {
			newValue := *value
			copied.infoptr[key] = &newValue
		} else {
			copied.infoptr[key] = nil
		}
	}
	for key, value := range sample.structm {
		copied.structm[key] = AnotherData{}
		copied.structm[key].TEST = value.TEST
		for key1, value1 := range sample.nesteds {

		}
	}
	for key, value := range sample.structmptr {
		if value != nil {
			copied.structmptr[key] = new(AnotherData)
			copied.structmptr[key].TEST = value.TEST
		} else {
			copied.structmptr[key] = nil
		}
	}

	return copied
}
*/
