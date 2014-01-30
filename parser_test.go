package main

import (
	"fmt"
	"testing"

	"github.com/bom-d-van/goutil/gocheckutil"
	. "launchpad.net/gocheck"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type ParserSuite struct{}

var _ = Suite(&ParserSuite{})

var xeroxMarkderSamples = [][]string{
	{`something else
@xerox single array ptrs
something else`, `@xerox single array ptrs
`},
	{`@xerox`, `@xerox`},
	{`something else
@xerox`, `@xerox`},
	{`something else
@xerox
something else`, `@xerox
`},
}

func (s *ParserSuite) TestXeroxMarker(c *C) {
	for _, sample := range xeroxMarkderSamples {
		c.Check(sample[0], gocheckutil.RegexpMatches, xeroxMaker)
		c.Check(xeroxMaker.FindString(sample[0]), Equals, sample[1])
		println("-------")
		for _, val := range xeroxMaker.FindSubmatch([]byte(sample[0])) {
			fmt.Printf("--> %+v\n", string(val))
		}
	}
}
