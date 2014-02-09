package examples

// single
func XeroxData(data Data) Data {
	newData := Data{}
	newData.Info = data.Info
	newData.info = data.info
	newData.SubData.Info = data.subData.Info
	newData.SubData.Int = data.subData.Int
	if data.subData.AnotherData != nil {
		newData.SubData.AnotherData = &AnotherData{}
		newData.SubData.AnotherData.TEST = data.SubData.AnotherData.TEST
	}
	// newData.Something = data.Something
	if data.InfoPtr != nil {
		val := *data.InfoPtr
		newData.InfoPtr = &val
	}

	for _, subData := range data.subdatas {
		newSubData := SubData{}
		newSubData.Info = subData.Info
		newSubData.Int = subData.Int
		if subData.AnotherData != nil {
			newSubData.AnotherData = &AnotherData{}
			newSubData.AnotherData.TEST = subData.AnotherData.TEST
		}
		newData.subdatas = append(newData.subdatas, newSubData)
	}

	// newData.function = func() (names string) { return "test" }

	for key, value := range data.subDataMap {
		newVal := SubData{}
		newVal.Info = value.Info
		newVal.Int = value.Int
		if value.AnotherData != nil {
			newVal.AnotherData = &AnotherData{}
			newVal.AnotherData.TEST = value.AnotherData.TEST
		}
		newData.subDataMap[key] = newVal
	}

	if data.subDataPtr != nil {
		newData.subDataPtr = &SubData{}
		newData.subDataPtr.Info = data.subDataPtr.Info
		newData.subDataPtr.Int = data.subDataPtr.Int
		if data.subDataPtr.AnotherData != nil {
			newData.subDataPtr.AnotherData = &AnotherData{}
			newData.subDataPtr.AnotherData.TEST = data.subDataPtr.AnotherData.TEST
		}
	}
	return newData
}

// func XeroxData(sample Data) (copied Data) {
// 	copied.Info = sample.Info
// 	copied.info = sample.Info
// 	copied.SubData.Info = sample.SubData.Info
// 	copied.SubData.Int = sample.SubData.Int
// 	if sample.SubData.AnotherData != nil {
// 		anotherData := &AnotherData{}
// 		anotherData.TEST = sample.SubData.AnotherData.TEST
// 		copied.SubData.AnotherData = anotherData
// 	}
// }

// type Struct struct {
// 	Name   string
// 	Type   string
// 	Fields []*Type
// 	IsStar bool
// }

type Type struct {
	Name string
	Type string

	Expression string

	// Parent *Type
	Prefix string

	IsPrimitive bool // i.e. string, int, float, etc

	IsStruct bool

	IsStar     bool
	TmpVarName string

	IsMap bool
	// KeyType   *Type
	// ValueType *Type

	IsArray bool
	// EleType *Type
}

var primitiveCopy = `{{.CopiedPrefix}}.{{.Name}} = {{.SamplePrefix}}.{{.Name}}`

var arrayCopy = `
for _, ele := {{.SamplePrefix}}.{{.Name}} {
	{{if not .EleType.IsPrimitive}}

	{{else}}
		{{.CopiedPrefix}}.{{.Name}} = append({{.CopiedPrefix}}.{{.Name}}, ele)
	{{end}}
}
`

var mapCopy = `
for key, value := range sample.{{.Name}} {
	{{if not .ValueType.IsPrimitive}}

	{{else}}
		{{.CopiedPrefix}}.{{.Name}}[key] = value
	{{end}}
}
`

var structCopy = `
{{range .Fields}}
	{{if not .IsPrimitive}}

	{{else}}
		{{.CopiedPrefix}}.{{.Name}} = {{.SamplePrefix}}.{{.Name}}
	{{end}}
{{end}}
`

var tmpl = `
func Xerox{{.Name}}(sample {{.Type}}) {{.Type}} {
	copied := {{.Type}}
	{{range .Fields}}
		{{if .IsMap}}
			for key, value := range sample.{{.Name}} {
				copied.{{.Name}}[key] = value
				if {{.ValueType}} {

				}
			}
		{{else}}
			{{if .IsArray}}
			{{else}}
				{{if .IsStar}}
				{{else}}
					copied.{{.Name}} = sample.{{.Name}}
				{{end}}
			{{end}}
		{{end}}
	{{end}}
	return copied
}
`

// array
func XeroxDatas(datas []Data) []Data {
	return datas
}

// ptrs
func XeroxDataPtrs(datas []*Data) []*Data {
	return datas
}

// tomap
func XerodDataToMap(data Data) map[string]interface{} {
	return map[string]interface{}{}
}
