package sql2go

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
	"strings"
	"study/tools/internal/word"
)

const tpl = `package model

{{range $v := .Imports}}
import "{{$v}}" 
{{- end}}

type {{.TableName | ToCamelCase}} struct {
{{range .Columns}}
	{{$length := len .Comment}}
	{{- if gt $length 0 }}// {{.Comment}} {{- end}}
	{{$typeLength := len .Tag}}	
	{{- if gt $typeLength 0 }} {{.Name | ToCamelCase}}  {{.Type}} {{.Tag}} {{else}} {{.Name | ToCamelCase}} {{.Type}} {{- end}}
{{end}}
}	

func ({{.TableName}}  *{{.TableName | ToCamelCase}}) TableName() string  {
	return "{{.TableName}}"
}
`

var DBTypeToGoType = map[string]string{
	"int":              "int",
	"tinyint":          "int8",
	"float":            "float32",
	"bigint":           "int64",
	"bit":              "byte",
	"boolean":          "bool",
	"char":             "string",
	"date":             "time.Time",
	"datetime":         "time.Time",
	"double":           "int64",
	"double precision": "int64",
	"integer":          "int32",
	"longblob":         "int64",
	"longtext":         "string",
	"mediumint":        "int32",
	"mediumtext":       "string",
	"smallint":         "int8",
	"text":             "string",
	"time":             "time.Time",
	"timestamp":        "string",
	"tinytext":         "string",
	"varchar":          "string",
	"year data type":   "string",
}

type goStructTemplateDB struct {
	TableName string
	Columns   []*goStructColumn
	Imports   map[string]template.HTML
}

type goStructColumn struct {
	Name    string
	Type    string
	Tag     template.HTML
	Comment template.HTML
}

type goStructTemplate struct {
	tpl string
}

func NewStructTemplate() *goStructTemplate {
	return &goStructTemplate{tpl: tpl}
}

func (t *goStructTemplate) AssemblyColumns(columns []*TableColumn) ([]*goStructColumn, map[string]template.HTML) {
	sc := make([]*goStructColumn, 0, len(columns))
	imports := make(map[string]template.HTML, 0)
	for _, column := range columns {
		tag := "`json:" + `"` + column.ColumnName + `"` + "`"
		_type := DBTypeToGoType[strings.Split(column.ColumnType, "(")[0]]
		switch _type {
		case "time.Time", "time":
			imports[_type] = "time"
		}

		sc = append(sc, &goStructColumn{
			Name:    column.ColumnName,
			Type:    _type,
			Tag:     template.HTML(tag),
			Comment: template.HTML(column.ColumnComment),
		})
	}
	return sc, imports
}

func (t *goStructTemplate) Generate(dbName, tableName string, columns []*TableColumn) error {
	sc, imports := t.AssemblyColumns(columns)
	tpl := template.Must(template.New("sql2go").Funcs(template.FuncMap{
		"ToCamelCase": word.UnderscoreToUpperCamelCase,
		"html": func(s string) string {
			return s
		},
	}).Parse(t.tpl))

	tplDB := &goStructTemplateDB{
		TableName: tableName,
		Columns:   sc,
		Imports:   imports,
	}

	if err := tpl.Execute(os.Stdout, tplDB); err != nil {
		return err
	}

	dir := fmt.Sprintf("template/model/%s/", dbName)
	if err := os.MkdirAll(dir, os.ModeDir); err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s/%s.go", dir, tableName)
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}

	if err := tpl.Execute(file, tplDB); err != nil {
		return err
	}

	output, err := exec.Command("go", "fmt", fileName).Output()
	log.Println("output: ", string(output))
	if err != nil {
		log.Println("waning: ", err)
	}

	return nil
}
