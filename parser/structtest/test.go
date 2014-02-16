package structtest

// @xerox
type Struct struct {
	SubData, subData       SubData `something:test`
	SubDataPtr, subDataPtr *SubData
}

type SubData struct {
	Info        string
	Int         int
	AnotherData *AnotherData
}

type AnotherData struct {
	TEST string
}

/*
package structtest

func XeroxStruct(sample Struct) Struct {
	copied := Struct{}
	copied.SubData.Info = sample.SubData.Info
	copied.SubData.Int = sample.SubData.Int
	if sample.SubData.AnotherData != nil {
		sample.SubData.AnotherData = new(AnotherData)
		copied.SubData.AnotherData.TEST = sample.SubData.AnotherData.TEST
	}
	copied.subData.Info = sample.subData.Info
	copied.subData.Int = sample.subData.Int
	if sample.subData.AnotherData != nil {
		sample.subData.AnotherData = new(AnotherData)
		copied.subData.AnotherData.TEST = sample.subData.AnotherData.TEST
	}
	if sample.SubDataPtr != nil {
		sample.SubDataPtr = new(SubData)
		copied.SubDataPtr.Info = sample.SubDataPtr.Info
		copied.SubDataPtr.Int = sample.SubDataPtr.Int
		if sample.SubDataPtr.AnotherData != nil {
			sample.SubDataPtr.AnotherData = new(AnotherData)
			copied.SubDataPtr.AnotherData.TEST = sample.SubDataPtr.AnotherData.TEST
		}
	}
	if sample.subDataPtr != nil {
		sample.subDataPtr = new(SubData)
		copied.subDataPtr.Info = sample.subDataPtr.Info
		copied.subDataPtr.Int = sample.subDataPtr.Int
		if sample.subDataPtr.AnotherData != nil {
			sample.subDataPtr.AnotherData = new(AnotherData)
			copied.subDataPtr.AnotherData.TEST = sample.subDataPtr.AnotherData.TEST
		}
	}

	return copied
}
*/
