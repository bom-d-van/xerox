package main

import (
	"testing"
	"github.com/bom-d-van/goutil/errutil"

	"log"
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
	// 	{`descriptions has nothing
	// to do with the markder.
	// yes, it is as it is`, ``},
}

func (s *ParserSuite) TestXeroxMarker(c *C) {
	for _, sample := range xeroxMarkderSamples {
		c.Check(sample[0], gocheckutil.RegexpMatches, xeroxMaker)
		c.Check(xeroxMaker.FindString(sample[0]), Equals, sample[1])
	}
}

func (p *ParserSuite) TestWalk(c *C) {
	codes, err := GenCodes("/Users/bom_d_van/Code/go/workspace/src/github.com/bom-d-van/xerox/examples")
	println(codes)
	if err != nil {
		log.Printf("--> %+v\n", err.(errutil.Err).Details())
	}
}
