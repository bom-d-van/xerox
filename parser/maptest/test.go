package maptest

// @xerox
type Map struct {
	info       map[int]string
	infoptr    map[int]*string
	structm    map[string]AnotherData
	structmptr map[string]*AnotherData
	// infoPtr *map[int]string
	mapArray  map[int][]string
	mapMap    map[int]map[int]int
	mapMapMap map[int]map[int]map[int]int
}

type AnotherData struct {
	TEST          string
	nesteds       map[string]*DeepNested
	complexNested map[init][]UltraDeep
}

type DeepNested struct {
	info           *string
	infos          map[float64]string
	nesteds        map[string]*UltraDeep
	nestedArrayVal map[init][]UltraDeep
}

type UltraDeep struct {
	level string
}

/*
package maptest

func XeroxMap(sample Map) Map {
	copied := Map{}
	if sample.info != nil {
		copied.info = map[int]string{}
		for key, val := range sample.info {
			copied.info[key] = val
		}
	}
	for key, val := range sample.infoptr {
		if val != nil {
			newVal := *val
			copied.infoptr[key] = &newVal
		} else {
			copied.infoptr[key] = nil
		}
	}
	for key, val := range sample.mapArray {
		copied.mapArray[key] = append(copied.mapArray[key], val)
	}
	for key, val := range sample.mapMap {
		if val != nil {
			newVal := map[int]int{}
			for key1, val1 := range val {
				newVal[key1] = val1
			}
			copied.mapMap[key] = newVal
		} else {
			copied.mapMap[key] = nil
		}
	}
	for key, val := range sample.mapMapMap {
		if val != nil {
			newVal := map[int]map[int]int{}
			for key1, val1 := range val {
				if val1 != nil {
					newVal1 := map[int]int{}
					for key2, val2 := range val1 {
						newVal1[key2] = vale
					}
					newVal[key1] = newVal1
				} else {
					newVal[key1] = nil
				}
			}
			copied.mapMapMap[key] = newVal
		} else {
			copied.mapMapMap[key] = nil
		}
	}
	for key, val := range sample.structm {
		newVal := AnotherData{}
		newVal.TEST = val.TEST
		for key1, val1 := range val.nesteds {
			if val1 != nil {
				newVal1 := new(DeepNested)
				if val1.info != nil {
					val := *val1.info
					newVal1.info = &val
				}
				for key2, val2 := range val1.infos {
					newVal1.infos[key2] = val2
				}
				for key2, val2 := range val1.nesteds {
					if val2 != nil {
						newVal2 := new(UltraDeep)
						newVal2.level = val2.level
						newVal1.nesteds[key2] = newVal2
					} else {
						newVal1.nesteds[key2] = nil
					}
				}
				newVal.nesteds[key1] = newVal1
			} else {
				newVal.nesteds[key1] = nil
			}
		}
		for key1, val1 := range val.complexNested {
			for _, val2 := range val1 {
				newVal2 := UltraDeep{}
				newVal2.level = val2.level
				newVal.complexNested[key1] = append(newVal.complexNested[key1], newVal2)
			}
		}
		copied.structm[key] = newVal
	}
	for key, val := range sample.structmptr {
		if val != nil {
			newVal := new(AnotherData)
			newVal.TEST = val.TEST
			for key1, val1 := range val.nesteds {
				if val1 != nil {
					newVal1 := new(DeepNested)
					if val1.info != nil {
						val := *val1.info
						newVal1.info = &val
					}
					for key2, val2 := range val1.infos {
						newVal1.infos[key2] = val2
					}
					for key2, val2 := range val1.nesteds {
						if val2 != nil {
							newVal2 := new(UltraDeep)
							newVal2.level = val2.level
							newVal1.nesteds[key2] = newVal2
						} else {
							newVal1.nesteds[key2] = nil
						}
					}
					newVal.nesteds[key1] = newVal1
				} else {
					newVal.nesteds[key1] = nil
				}
			}
			copied.structmptr[key] = newVal
		} else {
			copied.structmptr[key] = nil
		}
	}

	return copied
}*/
