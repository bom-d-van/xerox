package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/bom-d-van/goutil/errutil"
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

func filter(os.FileInfo) bool {
	return true
}

func GenCodes(path string) (codes string, err error) {
	visitor := &xeroxVisistor{}
	_ = visitor

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path, filter, parser.ParseComments)
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
				specCodes, err = genStructCodes(spec, "sample", "copied", "", 0)
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

func genStructCodes(spec *ast.TypeSpec, oprefix, cprefix, suffix string, level int) (codes string, err error) {
	structType, ok := spec.Type.(*ast.StructType)
	if !ok {
		return
	}

	for _, field := range structType.Fields.List {
		var fieldCodes string
		fieldCodes, err = genFieldCodes(field, oprefix, cprefix, suffix, level)
		if err != nil {
			err = errutil.Wrap(err)
			return
		}
		codes += fieldCodes
	}

	return
}

func genFieldCodes(field *ast.Field, oprefix, cprefix, suffix string, level int) (codes string, err error) {
	for _, name := range field.Names {
		var fieldCodes string
		fieldCodes, err = genExprCodes(field.Type, name, oprefix, cprefix, suffix, level)
		if err != nil {
			err = errutil.Wrap(err)
			return
		}

		codes += fieldCodes
	}

	return
}

func genExprCodes(ftypex ast.Expr, name *ast.Ident, old, copy, suffix string, level int) (codes string, err error) {
	switch ftype := ftypex.(type) {
	case *ast.Ident:
		if isPrimitiveType(ftype) {
			codes += fmt.Sprintf("\n%s%s.%s = %s.%s", copy, suffix, name, old, name)
		} else {
			if decl, ok := ftype.Obj.Decl.(*ast.TypeSpec); ok {
				var typeCodes string
				typeCodes, err = genStructCodes(decl, old+"."+name.Name, copy+"."+name.Name, suffix, level)
				if err != nil {
					err = errutil.Wrap(err)
					return
				}
				codes += typeCodes
			}
		}
	case *ast.StarExpr:
		xType := ftype.X.(*ast.Ident)
		if isPrimitiveType(xType) {
			nameVal := name.Name
			codes += fmt.Sprintf(`
				if %s.%s != nil {
					val := *%s.%s
					%s.%s = &val
				}`, old, nameVal, old, nameVal, copy, nameVal)
		} else {
			decl := xType.Obj.Decl.(*ast.TypeSpec)
			var bodyCodes string
			bodyCodes, err = genStructCodes(decl, old+"."+name.Name, copy+"."+name.Name, suffix, level)
			if err != nil {
				err = errutil.Wrap(err)
				return
			}
			codes += fmt.Sprintf(`
				if %s.%s != nil {
					%s.%s = new(%s)
					%s
				}`, old, name, old, name, xType.Name, cleanNewLine(bodyCodes))
		}
	case *ast.MapType:
		key, value, newValue := levelize("key", level), levelize("val", level), levelize("newVal", level)
		switch valType := ftype.Value.(type) {
		case *ast.Ident:
			if isPrimitiveType(valType) {
				codes += fmt.Sprintf(`
					for %s, %s := range %s.%s {
						%s.%s[%[1]s] = %s
					}`, key, value, old, name, copy, name)
			} else {
				var valCodes string
				valCodes, err = genStructCodes(valType.Obj.Decl.(*ast.TypeSpec), value, newValue, "["+key+"]", level+1)
				if err != nil {
					err = errutil.Wrap(err)
					return
				}
				codes += fmt.Sprintf(`
					for %s, %s := range %[4]s.%s {
						%[3]s := %[8]s{}
						%s
						%[6]s.%s[%[1]s] = %s
					}`, key, value, newValue, old, name, copy, name, valType.Name, cleanNewLine(valCodes))
			}
		case *ast.StarExpr:
			xtype := valType.X.(*ast.Ident)
			if isPrimitiveType(xtype) {
				codes += fmt.Sprintf(`
					for %s, %s := range %[4]s.%s {
						if %[2]s != nil {
							%s := *%[2]s
							%[6]s.%s[%[1]s] = &%[3]s
						} else {
							%[6]s.%s[%[1]s] = nil
						}
					}`, key, value, newValue, old, name, copy, name)
			} else {
				var valCodes string
				valCodes, err = genStructCodes(xtype.Obj.Decl.(*ast.TypeSpec), value, newValue, "["+key+"]", level+1)
				if err != nil {
					err = errutil.Wrap(err)
					return
				}
				codes += fmt.Sprintf(`
					for %s, %s := range %[4]s.%s {
						if %[2]s != nil {
							%s := new(%[8]s)
							%s
							%[6]s.%s[%[1]s] = %[3]s
						} else {
							%[6]s.%s[%[1]s] = nil
						}
					}`, key, value, newValue, old, name, copy, name, xtype.Name, cleanNewLine(valCodes))
			}
		}
	case *ast.ArrayType:
		elt, newElt := levelize("elt", level), levelize("newElt", level)
		switch valType := ftype.Elt.(type) {
		case *ast.Ident:
			if isPrimitiveType(valType) {
				codes += fmt.Sprintf(`
					for _, %s := range %s.%s {
						%s.%s = append(%[4]s.%s, %[1]s)
					}`, elt, old, name, copy, name)
			} else {
				var valCodes string
				valCodes, err = genStructCodes(valType.Obj.Decl.(*ast.TypeSpec), elt, newElt, "", level+1)
				if err != nil {
					err = errutil.Wrap(err)
					return
				}
				codes += fmt.Sprintf(`
					for _, %s := range %s.%s {
						%s := %s{}
						%s
						%s.%s = append(%[7]s.%s, %[4]s)
					}`, elt, old, name, newElt, valType.Name, cleanNewLine(valCodes), copy, name)
			}
		case *ast.StarExpr:
			xtype := valType.X.(*ast.Ident)
			if isPrimitiveType(xtype) {
				codes += fmt.Sprintf(`
					for _, %[2]s := range %[3]s.%s {
						if %[2]s != nil {
							%[1]s := *%s
							%[5]s.%s = append(%[5]s.%s, &%[1]s)
						} else {
							%[5]s.%s = append(%[5]s.%s, nil)
						}
					}`, newElt, elt, old, name, copy, name)
			} else {
				var valCodes string
				valCodes, err = genStructCodes(xtype.Obj.Decl.(*ast.TypeSpec), elt, newElt, "", level+1)
				if err != nil {
					err = errutil.Wrap(err)
					return
				}
				codes += fmt.Sprintf(`
					for _, %s := range %[3]s.%s {
						if %[1]s != nil {
							%s := new(%[5]s)
							%s
							%[7]s.%s = append(%[7]s.%s, %[2]s)
						} else {
							%[7]s.%s = append(%[7]s.%s, nil)
						}
					}`, elt, newElt, old, name, xtype.Name, cleanNewLine(valCodes), copy, name)
			}
		}
		// switch eltType := ftype.Elt.(type) {
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
		log.Printf("--> ChanType %+v\n", ftype)
	case *ast.FuncType:
		log.Printf("--> FuncType %+v\n", ftype)
	case *ast.InterfaceType:
		log.Printf("--> InterfaceType %+v\n", ftype)
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
