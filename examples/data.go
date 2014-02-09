package examples

// Something else

// @xerox single array ptrs
//
type Data struct {
	Info string // a info field
	info int

	InfoPtr *string
	// ptrPtr  **int

	SubData, subData       SubData `something:test`
	SubDataPtr, subDataPtr *SubData

	Map           map[int]string
	subDataMap    map[string]SubData
	subDataPtrMap map[string]*SubData
	// mapPtr        *map[int]string
	// ValPtrMap  map[int]*string

	infos       []float64
	subdatas    []SubData
	subDataPtrs []*SubData

	recv    <-chan int
	send    chan<- int
	channel chan int

	function func(input string) (name string)

	EmbeddedData
	AnonymousField struct {
		degrees float64
	}

	// TODO: Recursive Type
	// Parent *Data
	// Childs []Data
}

type EmbeddedData struct {
	EmbeddedInfo string
}

// is this a comment group

type AnotherData struct {
	TEST string
}

/*
Something Else
@xerox
*/
type SubData struct {
	Info        string
	Int         int
	AnotherData *AnotherData
}
