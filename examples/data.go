package main

// Something else

// @xerox single array ptrs
type Data struct {
	Info       string // a info field
	info       string
	SubData    SubData
	subDataPtr *SubData
}

type AnotherData struct {
	TEST string
}

/*
Something Else
@xerox
*/
type SubData struct {
	Info string
	Int  int
}
