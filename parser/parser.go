package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"regexp"
	"strings"

	"github.com/bom-d-van/goutils/errutils"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

type xeroxVisistor struct {
	currentSpec    string
	currentAstNode struct {
		TypeSpec *ast.TypeSpec
	}
}

var xeroxMaker = regexp.MustCompile("@xerox( .*)?(\n)?")

func retrieveDescription(format string) string {
	matches := xeroxMaker.FindSubmatch([]byte(format))
	if len(matches) == 0 {
		return ""
	}

	return strings.Trim(string(matches[1]), " \n")
}

func GenCodes(path string) (codes string, err error) {
	visitor := &xeroxVisistor{}
	_ = visitor

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
	if err != nil {
		err = errutil.Wrap(err)
		return
	}

	var pkg *ast.Package
	for _, p := range pkgs {
		pkg = p
	}

	for _, file := range pkg.Files {
		for _, declx := range file.Decls {
			decl, ok := declx.(*ast.GenDecl)
			if !ok {
				continue
			}

			if !xeroxMaker.MatchString(decl.Doc.Text()) {
				continue
			}

			for _, specx := range decl.Specs {
				spec, ok := specx.(*ast.TypeSpec)
				if !ok {
					continue
				}
				var specCodes string
				specCodes, err = genStructCodes(spec, "sample", "copied", 0)
				if err != nil {
					err = errutil.Wrap(err)
					return
				}

				tname := spec.Name.Name
				codes += fmt.Sprintf(`
					func Xerox%s(sample %s) %s {
						copied := %s{}
						%s

						return copied
					}`, tname, tname, tname, tname, cleanNewLine(specCodes))
			}
		}
	}

	codes = "package " + pkg.Name + "\n" + codes
	formatedCodes, err := format.Source([]byte(codes))
	if err != nil {
		err = errutil.Wrap(err)
		return
	}
	codes = string(formatedCodes)

	return
}

func genStructCodes(spec *ast.TypeSpec, oprefix, cprefix string, level int) (codes string, err error) {
	structType, ok := spec.Type.(*ast.StructType)
	if !ok {
		return
	}

	for _, field := range structType.Fields.List {
		var fieldCodes string
		fieldCodes, err = genFieldCodes(field, oprefix, cprefix, level)
		if err != nil {
			err = errutil.Wrap(err)
			return
		}
		codes += fieldCodes
	}

	return
}

func genFieldCodes(field *ast.Field, oprefix, cprefix string, level int) (codes string, err error) {
	for _, name := range field.Names {
		var fieldCodes string
		fieldCodes, err = genExprCodes(field.Type, name, oprefix, cprefix, level)
		if err != nil {
			err = errutil.Wrap(err)
			return
		}

		codes += fieldCodes
	}

	return
}

type identHandler interface {
	primitive(ftype *ast.Ident, name *ast.Ident, old, copy string) (string, error)
	nonprimitive(ftype *ast.Ident, name *ast.Ident, old, copy string) (string, error)
}

func handle(handler identHandler, ftype *ast.Ident, name *ast.Ident, old, copy string) (string, error) {
	if isPrimitiveType(ftype) {
		return handler.primitive(ftype, name, old, copy)
	} else {
		return handler.nonprimitive(ftype, name, old, copy)
	}
}

type ident struct {
	level int
}

func (b *ident) primitive(ident *ast.Ident, name *ast.Ident, old, copy string) (string, error) {
	return fmt.Sprintf("\n%s.%s = %s.%s", copy, name, old, name), nil
}

func (b *ident) nonprimitive(ident *ast.Ident, name *ast.Ident, old, copy string) (codes string, err error) {
	decl := ident.Obj.Decl.(*ast.TypeSpec)
	codes, err = genStructCodes(decl, old+"."+name.Name, copy+"."+name.Name, b.level)
	if err != nil {
		err = errutil.Wrap(err)
		return
	}
	return
}

type starExpr struct {
	level int
}

func (s *starExpr) primitive(ident *ast.Ident, name *ast.Ident, old, copy string) (codes string, err error) {
	nameVal := name.Name
	codes = fmt.Sprintf(`
		if %s.%s != nil {
			val := *%s.%s
			%s.%s = &val
		}`, old, nameVal, old, nameVal, copy, nameVal)
	return
}

func (s *starExpr) nonprimitive(ident *ast.Ident, name *ast.Ident, old, copy string) (codes string, err error) {
	decl := ident.Obj.Decl.(*ast.TypeSpec)
	var bodyCodes string
	bodyCodes, err = genStructCodes(decl, old+"."+name.Name, copy+"."+name.Name, s.level)
	if err != nil {
		err = errutil.Wrap(err)
		return
	}
	codes += fmt.Sprintf(`
		if %s.%s != nil {
			%s.%s = new(%s)
			%s
		}`, old, name, old, name, ident.Name, cleanNewLine(bodyCodes))
	return
}

type mapIdent struct {
	level int
}

func (m *mapIdent) primitive(ident *ast.Ident, name *ast.Ident, old, copy string) (codes string, err error) {
	codes = fmt.Sprintf(`
		for %s, %s := range %s.%s {
			%s.%s[%[1]s] = %s
		}`, levelize("key", m.level), levelize("val", m.level), old, name, copy, name)
	return
}

func (m *mapIdent) nonprimitive(ident *ast.Ident, name *ast.Ident, old, copy string) (codes string, err error) {
	value, newValue := levelize("val", m.level), levelize("newVal", m.level)
	valCodes, err := genStructCodes(ident.Obj.Decl.(*ast.TypeSpec), value, newValue, m.level+1)
	if err != nil {
		err = errutil.Wrap(err)
		return
	}
	codes = fmt.Sprintf(`
		for %s, %s := range %[4]s.%s {
			%[3]s := %[8]s{}
			%s
			%[6]s.%s[%[1]s] = %[3]s
		}`,
		levelize("key", m.level), value, newValue,
		old, name, copy, name, ident.Name, cleanNewLine(valCodes))
	return
}

type mapStarExpr struct {
	level int
}

func (m *mapStarExpr) primitive(ftype *ast.Ident, name *ast.Ident, old, copy string) (codes string, err error) {
	codes = fmt.Sprintf(`
		for %s, %s := range %[4]s.%s {
			if %[2]s != nil {
				%s := *%[2]s
				%[6]s.%s[%[1]s] = &%[3]s
			} else {
				%[6]s.%s[%[1]s] = nil
			}
		}`, levelize("key", m.level), levelize("val", m.level), levelize("newVal", m.level), old, name, copy, name)
	return
}

func (m *mapStarExpr) nonprimitive(ident *ast.Ident, name *ast.Ident, old, copy string) (codes string, err error) {
	value, newValue := levelize("val", m.level), levelize("newVal", m.level)
	valCodes, err := genStructCodes(ident.Obj.Decl.(*ast.TypeSpec), value, newValue, m.level+1)
	if err != nil {
		err = errutil.Wrap(err)
		return
	}
	codes = fmt.Sprintf(`
		for %s, %s := range %[4]s.%s {
			if %[2]s != nil {
				%s := new(%[8]s)
				%s
				%[6]s.%s[%[1]s] = %[3]s
			} else {
				%[6]s.%s[%[1]s] = nil
			}
		}`,
		levelize("key", m.level), value, newValue,
		old, name, copy, name, ident.Name, cleanNewLine(valCodes))
	return
}

type arrayIdent struct {
	level int
}

func (a *arrayIdent) primitive(ident *ast.Ident, name *ast.Ident, old, copy string) (codes string, err error) {
	codes = fmt.Sprintf(`
		for _, %s := range %s.%s {
			%s.%s = append(%[4]s.%s, %[1]s)
		}`, levelize("elt", a.level), old, name, copy, name)
	return
}

func (a *arrayIdent) nonprimitive(ident *ast.Ident, name *ast.Ident, old, copy string) (codes string, err error) {
	elt, newElt := levelize("elt", a.level), levelize("newElt", a.level)
	valCodes, err := genStructCodes(ident.Obj.Decl.(*ast.TypeSpec), elt, newElt, a.level+1)
	if err != nil {
		err = errutil.Wrap(err)
		return
	}
	codes = fmt.Sprintf(`
		for _, %s := range %s.%s {
			%s := %s{}
			%s
			%s.%s = append(%[7]s.%s, %[4]s)
		}`, elt, old, name, newElt, ident.Name, cleanNewLine(valCodes), copy, name)
	return
}

type arrayStarExpr struct {
	level int
}

func (a *arrayStarExpr) primitive(ident *ast.Ident, name *ast.Ident, old, copy string) (codes string, err error) {
	codes = fmt.Sprintf(`
		for _, %[2]s := range %[3]s.%s {
			if %[2]s != nil {
				%[1]s := *%s
				%[5]s.%s = append(%[5]s.%s, &%[1]s)
			} else {
				%[5]s.%s = append(%[5]s.%s, nil)
			}
		}`, levelize("newElt", a.level), levelize("elt", a.level), old, name, copy, name)
	return
}

func (a *arrayStarExpr) nonprimitive(ident *ast.Ident, name *ast.Ident, old, copy string) (codes string, err error) {
	valCodes, err := genStructCodes(ident.Obj.Decl.(*ast.TypeSpec), levelize("elt", a.level), levelize("newElt", a.level), a.level+1)
	if err != nil {
		err = errutil.Wrap(err)
		return
	}
	codes = fmt.Sprintf(`
		for _, %s := range %[3]s.%s {
			if %[1]s != nil {
				%s := new(%[5]s)
				%s
				%[7]s.%s = append(%[7]s.%s, %[2]s)
			} else {
				%[7]s.%s = append(%[7]s.%s, nil)
			}
		}`, levelize("elt", a.level), levelize("newElt", a.level), old, name, ident.Name, cleanNewLine(valCodes), copy, name)
	return
}

type mapArray struct {
	level int
}

func (m *mapArray) primitive(ftype *ast.Ident, name *ast.Ident, old, copy string) (codes string, err error) {
	codes = fmt.Sprintf(`
		for key, val := range sample.mapArray {
			if val != nil {
				newVal := []string{}
				for _, elt1 := range itmes {
					newVal = apppend(newVal, elt1)
				}
				copied.mapArray[key] = newVal
			} else {
				copied.mapArray[key] = nil
			}
		}`)
	return
}

func (m *mapArray) nonprimitive(ftype *ast.Ident, name *ast.Ident, old, copy string) (codes string, err error) {
	return
}

// type mapIdent struct {
// 	level int
// }

// func (m *mapIdent) primitive(ftype *ast.Ident, name *ast.Ident, old, copy string) (codes string, err error) {
// 	return
// }

// func (m *mapIdent) nonprimitive(ftype *ast.Ident, name *ast.Ident, old, copy string) (codes string, err error) {
// 	return
// }

func genExprCodes(exprx ast.Expr, name *ast.Ident, old, copy string, level int) (codes string, err error) {
	switch expr := exprx.(type) {
	case *ast.Ident:
		return handle(&ident{level}, expr, name, old, copy)
	case *ast.StarExpr:
		return handle(&starExpr{level}, expr.X.(*ast.Ident), name, old, copy)
	case *ast.MapType:
		switch valType := expr.Value.(type) {
		case *ast.Ident:
			return handle(&mapIdent{level}, valType, name, old, copy)
		case *ast.StarExpr:
			return handle(&mapStarExpr{level}, valType.X.(*ast.Ident), name, old, copy)
		case *ast.ArrayType:
			switch elt := valType.Elt.(type) {
			case *ast.Ident:
			case *ast.StarExpr:
			case *ast.ArrayType:
			case *ast.MapType:
			}
			return handle(&mapStarExpr{level}, valType.X.(*ast.Ident), name, old, copy)
		case *ast.MapType:
		}
	case *ast.ArrayType:
		switch valType := expr.Elt.(type) {
		case *ast.Ident:
			return handle(&arrayIdent{level}, valType, name, old, copy)
		case *ast.StarExpr:
			return handle(&arrayStarExpr{level}, valType.X.(*ast.Ident), name, old, copy)
		}
		// switch eltType := expr.Elt.(type) {
		// case *ast.Ident:
		// 	if eltType.Obj == nil {
		// 		var eltCodes string
		// 		eltCodes, err = runTmpl("array", arrayTmplData{
		// 			commonTmplData: commonTmplData{
		// 				CPrefix: copy,
		// 				OPrefix: old,
		// 				Name:    name.Name,
		// 			},
		// 			IsEltPrimitive: true,
		// 		})
		// 		if err != nil {
		// 			err = errutil.Wrap(err)
		// 			return
		// 		}

		// 		codes += eltCodes
		// 	} else {

		// 	}
		// case *ast.StarExpr:
		// 	xtype := eltType.X.(*ast.Ident)
		// 	if xtype.Obj == nil {
		// 		var eltCodes string
		// 		eltCodes, err = runTmpl("array", arrayTmplData{
		// 			commonTmplData: commonTmplData{
		// 				CPrefix: copy,
		// 				OPrefix: old,
		// 				Name:    name.Name,
		// 			},
		// 			IsEltPrimitivePtr: true,
		// 		})
		// 		if err != nil {
		// 			err = errutil.Wrap(err)
		// 			return
		// 		}

		// 		codes += eltCodes
		// 	} else {

		// 	}
		// }
	case *ast.ChanType:
		log.Printf("--> ChanType %+v\n", expr)
	case *ast.FuncType:
		log.Printf("--> FuncType %+v\n", expr)
	case *ast.InterfaceType:
		log.Printf("--> InterfaceType %+v\n", expr)
	}

	return
}

func levelize(str string, level int) string {
	var levelStr string
	if level > 0 {
		levelStr = fmt.Sprintf("%d", level)
	}

	return str + levelStr
}

func cleanNewLine(codes string) string {
	if codes != "" && codes[0] == '\n' {
		codes = codes[1:]
	}
	return codes
}

func isPrimitiveType(expr *ast.Ident) bool {
	return expr.Obj == nil
}

// func walkGenDeclSpecs(specs []ast.Spec) {
// 	for _, spec := range specs {
// 		log.Printf("--> %+v\n", spec)
// 		switch node := spec.(type) {
// 		case *ast.TypeSpec:
// 			log.Printf("--> %+v\n", node.Type)
// 			// switch node.Type

// 		}
// 	}
// }

// func walkTypes(ntype *ast.Node) {
// 	switch node := ntype.(type) {
// 	case *ast.StructType:
// 		log.Printf("--> %+v\n", node)
// 		for _, field := range node.Fields.List {
// 			field.Type
// 		}
// 	case *ast.MapType:
// 	case *ast.ArrayType:
// 	case *ast.ChanType:
// 	case *ast.FuncType:
// 	case *ast.FuncType:
// 	case *ast.InterfaceType:
// 	case *ast.Ident:
// 	default:
// 		log.Printf("--> %+v\n", node)
// 	}
// }
