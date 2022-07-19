package sql2go

const tpl = `type {{.TableName | ToCamelCase}} struct {
{{range .Column}}
	{{$length := len .Comment}}
	
	{{if gt $length 0 }}// {{.Comment}}
	{{else}}//{{.Name}} 
	{{end}}
	{{$typeLength := len .Type}}
	
	{{if gt $typeLength 0 }}// {{.Name || ToCamelCase}}  {{.Type}} {{.Tag}}
	{{else}}//{{.Name}} 
	{{end}}
{{end}}
}	
`

//
//func (model {{.TableName | ToCamelCase}}) TableName() string  {
//	return "{{.TableName}}"
//}

type StructTemplateDB struct {
	TableName string
	Columns   []*StructColumn
}

type StructTemplate struct {
	tpl string
}

func NewStructTemplate() *StructTemplate {
	return &StructTemplate{tpl: tpl}
}

type StructColumn struct {
	Name    string
	Type    string
	Tag     string
	Comment string
}
