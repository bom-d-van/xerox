package examples

// Something else

// @xerox single array ptrs
//
type Data struct {
	Info                   string // a info field
	InfoPtr                *string
	ptrPtr                 **int
	info                   int
	SubData, subData       SubData `something:test`
	SubDataPtr, subDataPtr *SubData
	mapPtr                 *map[int]string
	subdatas               []SubData
	subDataPtrs            []*SubData
	infos                  []float64
	subDataMap             map[string]SubData
	recv                   <-chan int
	send                   chan<- int
	channel                chan int
	function               func(input string) (name string)
	EmbeddedData
	AnonymousField struct {
		degrees float64
	}
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
