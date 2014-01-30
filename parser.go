package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"regexp"

	"github.com/bom-d-van/goutil/errutil"
)

type xeroxVisistor struct {
	currentSpec string
}

var xeroxMaker = regexp.MustCompile("@xerox( .*)?(\n)?")

func (x *xeroxVisistor) Visit(node ast.Node) ast.Visitor {
	switch i := node.(type) {
	case *ast.CommentGroup:
		logger.Println(i.Text())
		// x.currentCommentGroup = i
	case *ast.TypeSpec:
		// logger.Printf("--> %+v\n", i.Doc)
	case *ast.Ident:
	case *ast.StructType:
		// logger.Printf("--> %+v\n", i.Fields)
		// fset := token.NewFileSet()
		// cmap := ast.NewCommentMap(fset, i, []*ast.CommentGroup{x.currentCommentGroup})
		// for k, v := range cmap {
		// 	logger.Printf("--> %+v\n", k)
		// 	logger.Printf("--> %+v\n", v[0].List[0])
		// }
	}
	return x
}

// type xeroxNode struct {
// 	pos, end token.Pos
// }

// func (n *xeroxNode) Pos() token.Pos {
// 	return n.pos
// }

// func (n *xeroxNode) End() token.Pos {
// 	return n.end
// }

func filter(os.FileInfo) bool {
	return true
}

func walk(path string) (err error) {
	visitor := &xeroxVisistor{}
	// node := &xeroxNode{}
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path, filter, parser.ParseComments)
	if err != nil {
		return errutil.Wrap(err)
	}
	for _, file := range pkgs["main"].Files {
		ast.Walk(visitor, file)
		// fset := token.NewFileSet()
		// cmap := ast.NewCommentMap(fset, file, file.Comments)
		// for key, value := range cmap {
		// 	logger.Printf("--> %+v\n", key)
		// 	for _, comment := range key.(*ast.File).Comments {
		// 		logger.Printf("--> %+v\n", comment.Text())
		// 	}
		// 	logger.Printf("--> %T %+v\n", value, value)
		// }
		// println("================")
	}
	return
}
