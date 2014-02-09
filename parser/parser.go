package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"regexp"
	"strings"
	"text/template"

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

func (x *xeroxVisistor) Visit(nodex ast.Node) ast.Visitor {
	switch node := nodex.(type) {
	case *ast.CommentGroup:
		log.Printf("--> CommentGroup %+v\n", node)
		x.currentSpec = retrieveDescription(node.Text())
	case *ast.TypeSpec:
		// log.Printf("--> %+v\n", node.Doc)
		log.Printf("--> TypeSpec %+v\n", node)
	// case *ast.Ident:
	// 	log.Printf("--> Ident %+v\n", node)
	// 	if node.Obj != nil {
	// 		log.Printf("--> %+v\n", node.Obj)
	// 	}
	case *ast.Field:
		log.Printf("--> Field %+v %+v\n", node.Names, node.Type)
		if sfield, ok := node.Type.(*ast.Ident); ok {
			log.Printf("--> StructField %+v\n", sfield.Obj)
		}
	case *ast.StructType:
		log.Printf("--> StructType %+v\n", node.Fields.List)
		ast.Walk(x, node.Fields)
	// log.Printf("--> %+v\n", i.Fields)
	// fset := tokeny.NewFileSet()
	// cmap := ast.NewCommentMap(fset, i, []*ast.CommentGroup{x.currentCommentGroup})
	// for k, v := range cmap {
	// 	log.Printf("--> %+v\n", k)
	// 	log.Printf("--> %+v\n", v[0].List[0])
	// }
	// for _, field := range node.Fields.List {
	// 	log.Printf("--> %+v\n", "-----------")
	// 	// debug.PrintStack()
	// 	log.Printf("--> %+v\n", field.Type)
	// 	// log.Printf("--> %s\n", field.Type.(*ast.Ident).Name)
	// 	// log.Printf("--> %+v\n", reflect.TypeOf(field.Type).Name())
	// 	for _, name := range field.Names {
	// 		log.Printf("--> %+v\n", name.Name)
	// 	}
	// }
	case *ast.MapType:
		log.Printf("--> MapType %+v\n", node)
	case *ast.ArrayType:
		log.Printf("--> ArrayType %+v\n", node)
	case *ast.ChanType:
		log.Printf("--> ChanType %+v\n", node)
	case *ast.FuncType:
		log.Printf("--> FuncType %+v\n", node)
	case *ast.InterfaceType:
		log.Printf("--> InterfaceType %+v\n", node)
	}

	return x
}

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

var codeTmpl *template.Template

func runTmpl(name string, data interface{}) (codes string, err error) {
	buffer := bytes.NewBuffer([]byte{})
	err = codeTmpl.ExecuteTemplate(buffer, name, data)
	if err != nil {
		err = errutil.Wrap(err)
		return
	}
	codes = buffer.String()

	return
}

func init() {
	codeTmpl = template.Must(template.New("").Parse(`
		{{define "commonFunc"}}
			func Xerox{{.Name}}(sample {{.Type}}) {{.Type}} {
				copied := {{.Type}}{}
				{{.Body}}
				return copied
			}
		{{end}}

		{{define "starExpr"}}
			if {{.OPrefix}}.{{.Name}} != nil {
				val := *{{.OPrefix}}.{{.Name}}
				{{.CPrefix}}.{{.Name}} = &val
			}
		{{end}}

		{{define "structStarExpr"}}
			if {{.OPrefix}}.{{.Name}} != nil {
				{{.CPrefix}}.{{.Name}} = &{{.Type}}{}
				{{.Fields}}
			}
		{{end}}

		{{define "map"}}
			for key, value := range {{.OPrefix}}.{{.Name}} {
				{{.CPrefix}}.{{.Name}}[key] = value
			}
		{{end}}

		{{define "structMap"}}
			for key, value := range {{.OPrefix}}.{{.Name}} {
				newValue := {{.Type}}{}
				{{.Body}}
				{{.CPrefix}}.{{.Name}}[key] = newValue
			}
		{{end}}

		{{define "structPtrMap"}}
			for key, value := range {{.OPrefix}}.{{.Name}} {
				var newValue *{{.Type}}
				if value != nil {
					{{.Body}}
				}
				{{.CPrefix}}.{{.Name}}[key] = newValue
			}
		{{end}}

		{{define "array"}}
			for _, elt := range {{.OPrefix}}.{{.Name}} {
				{{if .IsEltPrimitive}}
					{{.CPrefix}}.{{.Name}} = append({{.CPrefix}}.{{.Name}}, elt)
				{{end}}
				{{if .IsEltPrimitivePtr}}
				{{end}}
				{{if .IsEltStruct}}
					newElt := {{.Type}}{}
					{{.Body}}
					{{.CPrefix}}.{{.Name}} = newElt
				{{end}}
				{{if .IsEltStructPtr}}
					newElt := &{{.Type}}{}
					if elt != nil {
						{{.Body}}
					}
					{{.CPrefix}}.{{.Name}} = newElt
				{{end}}
			}
		{{end}}
	`))
}

type (
	funcTmplData struct {
		Name string
		Type string
		Body string
	}

	starExprTmplData struct {
		OPrefix, CPrefix string
		Name             string
	}

	structStarExprTmplData struct {
		OPrefix, CPrefix string
		Name             string
		Type             string
		Fields           string
	}

	commonTmplData struct {
		OPrefix, CPrefix string
		Name             string
		Type             string
		Body             string
	}

	arrayTmplData struct {
		commonTmplData
		IsEltPrimitive    bool
		IsEltPrimitivePtr bool
		IsEltStruct       bool
		IsEltStructPtr    bool
	}
)

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

			for _, specx := range decl.Specs {
				spec, ok := specx.(*ast.TypeSpec)
				if !ok {
					continue
				}
				var specCodes string
				specCodes, err = genStructCodes(spec, "sample", "copied")
				if err != nil {
					err = errutil.Wrap(err)
					return
				}

				var funcCodes string
				funcCodes, err = runTmpl("commonFunc", funcTmplData{
					Name: spec.Name.Name,
					Type: spec.Name.Name,
					Body: specCodes,
				})
				if err != nil {
					err = errutil.Wrap(err)
					return
				}

				codes += funcCodes + "\n"
			}
		}
	}

	codes = "package main\n" + codes
	formatedCodes, err := format.Source([]byte(codes))
	if err != nil {
		err = errutil.Wrap(err)
		return
	}
	codes = string(formatedCodes)

	return
}

func genStructCodes(spec *ast.TypeSpec, oprefix, cprefix string) (codes string, err error) {
	structType, ok := spec.Type.(*ast.StructType)
	if !ok {
		return
	}

	for _, field := range structType.Fields.List {
		var fieldCodes string
		fieldCodes, err = genFieldCodes(field, oprefix, cprefix)
		if err != nil {
			err = errutil.Wrap(err)
			return
		}
		codes += fieldCodes
	}

	return
}

func genFieldCodes(field *ast.Field, oprefix, cprefix string) (codes string, err error) {
	for _, name := range field.Names {
		switch ftype := field.Type.(type) {
		case *ast.StructType:
		case *ast.Ident:
			if ftype.Obj != nil && ftype.Obj.Decl != nil {
				if decl, ok := ftype.Obj.Decl.(*ast.TypeSpec); ok {
					var typeCodes string
					typeCodes, err = genStructCodes(decl, oprefix+"."+name.Name, cprefix+"."+name.Name)
					if err != nil {
						err = errutil.Wrap(err)
						return
					}
					codes += typeCodes
				}
			} else {
				codes += fmt.Sprintf("%s.%s = %s.%s\n", cprefix, name, oprefix, name)
			}
		case *ast.StarExpr:
			if isPrimitiveType(ftype.X) {
				var starCodes string
				starCodes, err = runTmpl("starExpr", starExprTmplData{
					CPrefix: cprefix,
					OPrefix: oprefix,
					Name:    name.Name,
				})
				if err != nil {
					err = errutil.Wrap(err)
					return
				}
				codes += starCodes
			} else {
				switch xType := ftype.X.(type) {
				case *ast.Ident:
					if xType.Obj != nil && xType.Obj.Decl != nil {
						if decl, ok := xType.Obj.Decl.(*ast.TypeSpec); ok {
							var bodyCodes string
							bodyCodes, err = genStructCodes(decl, oprefix+"."+name.Name, cprefix+"."+name.Name)
							if err != nil {
								err = errutil.Wrap(err)
								return
							}
							var ptrCodes string
							ptrCodes, err = runTmpl("structStarExpr", structStarExprTmplData{
								OPrefix: oprefix,
								CPrefix: cprefix,
								Name:    name.Name,
								Type:    xType.Name,
								Fields:  bodyCodes,
							})
							if err != nil {
								err = errutil.Wrap(err)
								return
							}
							codes += ptrCodes
						}
					}
				}
			}
		case *ast.MapType:
			// _, ok := ftype.Key.(*ast.Ident)
			// if !ok {
			// 	continue
			// }
			if isPrimitiveType(ftype.Value) {
				val, ok := ftype.Value.(*ast.Ident)
				_ = val
				if !ok {
					continue
				}

				var mapCodes string
				mapCodes, err = runTmpl("map", commonTmplData{
					CPrefix: cprefix,
					OPrefix: oprefix,
					Name:    name.Name,
				})
				if err != nil {
					err = errutil.Wrap(err)
					return
				}
				codes += mapCodes
			} else {
				switch valType := ftype.Value.(type) {
				case *ast.Ident:
					var valCodes string
					decl := valType.Obj.Decl.(*ast.TypeSpec)
					valCodes, err = genStructCodes(decl, "value", "newValue")
					if err != nil {
						err = errutil.Wrap(err)
						return
					}
					var mapCodes string
					mapCodes, err = runTmpl("structMap", commonTmplData{
						CPrefix: cprefix,
						OPrefix: oprefix,
						Name:    name.Name,
						Type:    valType.Name,
						Body:    valCodes,
					})
					if err != nil {
						err = errutil.Wrap(err)
						return
					}
					codes += mapCodes
				case *ast.StarExpr:
					var valCodes string
					xType := valType.X.(*ast.Ident)
					decl := xType.Obj.Decl.(*ast.TypeSpec)
					valCodes, err = genStructCodes(decl, "value", "newValue")
					if err != nil {
						err = errutil.Wrap(err)
						return
					}
					var mapCodes string
					mapCodes, err = runTmpl("structPtrMap", commonTmplData{
						CPrefix: cprefix,
						OPrefix: oprefix,
						Name:    name.Name,
						Type:    xType.Name,
						Body:    valCodes,
					})
					if err != nil {
						err = errutil.Wrap(err)
						return
					}
					codes += mapCodes
				}
			}
		case *ast.ArrayType:
			switch eltType := ftype.Elt.(type) {
			case *ast.Ident:
				if eltType.Obj == nil {
					var eltCodes string
					eltCodes, err = runTmpl("array", arrayTmplData{
						commonTmplData: commonTmplData{
							CPrefix: cprefix,
							OPrefix: oprefix,
							Name:    name.Name,
						},
						IsEltPrimitive: true,
					})
					if err != nil {
						err = errutil.Wrap(err)
						return
					}

					codes += eltCodes
				} else {

				}
			case *ast.StarExpr:
			}
		case *ast.ChanType:
			log.Printf("--> ChanType %+v\n", ftype)
		case *ast.FuncType:
			log.Printf("--> FuncType %+v\n", ftype)
		case *ast.InterfaceType:
			log.Printf("--> InterfaceType %+v\n", ftype)
		}
	}

	return
}

func isPrimitiveType(exprx ast.Expr) bool {
	// if ident, ok := exprx.(*ast.Ident); ok {
	// 	return ident.Obj == nil
	// }
	switch expr := exprx.(type) {
	case *ast.Ident:
		return expr.Obj == nil
	case *ast.StarExpr:
		return false
	default:
		return true
	}

	return true
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
