package main

import (
	"go/parser"
	"go/token"
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

var walkSamples = map[string]string{
	// "simpletest": "/Users/bom_d_van/Code/go/workspace/src/github.com/bom-d-van/xerox/parser",
	// "structtest": "/Users/bom_d_van/Code/go/workspace/src/github.com/bom-d-van/xerox/parser",
	"maptest":    "/Users/bom_d_van/Code/go/workspace/src/github.com/bom-d-van/xerox/parser",
	// "arraytest":  "/Users/bom_d_van/Code/go/workspace/src/github.com/bom-d-van/xerox/parser",
}

func (p *ParserSuite) TestWalk(c *C) {
	for pkg, path := range walkSamples {
		fset := token.NewFileSet()
		pkgfixture := path + "/" + pkg
		pkgs, err := parser.ParseDir(fset, pkgfixture, filter, parser.ParseComments)
		c.Check(err, Equals, nil)
		comments := pkgs[pkg].Files[pkgfixture+"/test.go"].Comments
		expectation := comments[len(comments)-1].Text()

		codes, err := GenCodes(pkgfixture)
		c.Check(err, Equals, nil)
		c.Check(codes, Equals, expectation)
		println(codes)
	}
}
