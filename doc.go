package graphqldoc

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

type TypeRef struct {
	Kind   string   `json:"kind"`
	Name   string   `json:"name"`
	OfType *TypeRef `json:"ofType"`
}

type InputValue struct {
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	DefaultValue interface{} `json:"defaultValue"`
	Type         TypeRef     `json:"type"`
}

type TypeField struct {
	Name              string       `json:"name"`
	Description       string       `json:"description"`
	Args              []InputValue `json:"args"`
	Type              TypeRef      `json:"type"`
	IsDeprecated      bool         `json:"isDeprecated"`
	DeprecationReason string       `json:"deprecationReason"`
}

type EnumValues struct {
	Name              string `json:"name"`
	Description       string `json:"description"`
	IsDeprecated      bool   `json:"isDeprecated"`
	DeprecationReason string `json:"deprecationReason"`
}

type FullType struct {
	Kind          string       `json:"kind"`
	Name          string       `json:"name"`
	Description   string       `json:"description"`
	Fields        []TypeField  `json:"fields"`
	InputFields   []InputValue `json:"inputFields"`
	Interfaces    []TypeRef    `json:"interfaces"`
	EnumValues    []EnumValues `json:"enumValues"`
	PossibleTypes []TypeRef    `json:"possibleTypes"`
}

type TypeDirective struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Args        []InputValue `json:"args"`
	OnOperation bool         `json:"onOperation"`
	onFragment  bool         `json:"onFragment"`
	onField     bool         `json:"onField"`
}

type Schema struct {
	QueryType        FullType        `json:"queryType"`
	MutationType     FullType        `json:"mutationType"`
	SubscriptionType FullType        `json:"subscriptionType"`
	Types            []FullType      `json:"types"`
	Directives       []TypeDirective `json:"directives"`
}

var (
	dir           = "doc/"
	queryFile     = dir + "query.md"
	objectFile    = dir + "object.md"
	mutationFile  = dir + "mutation.md"
	scalarFile    = dir + "scalar.md"
	enumFile      = dir + "enum.md"
	interfaceFile = dir + "interface.md"
)

func init() {
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		os.RemoveAll(dir)
		os.Remove(dir)
	}

	os.Mkdir(dir, 0755)
}

func generateDocs(schema Schema) {
	queries(schema.QueryType)
	mutations(schema.MutationType)

	var scalar []FullType
	var enum []FullType
	var inter []FullType
	var object []FullType
	for _, v := range schema.Types {
		if !strings.Contains(v.Name, "__") {
			switch v.Kind {
			case "SCALAR":
				scalar = append(scalar, v)
				break
			case "ENUM":
				enum = append(enum, v)
				break
			case "INTERFACE":
				inter = append(inter, v)
				break
			case "OBJECT":
				object = append(object, v)
				break
			}
		}
	}

	scalars(scalar)
	enums(enum)
	interfaces(inter)
	objects(object)

}

func queries(query FullType) {
	f, err := os.OpenFile(queryFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	checkError(err)
	data, err := Asset("template/schema.tmpl")
	checkError(err)

	t := template.Must(temp(string(data)))
	err = t.Execute(f, query)
	checkError(err)
}

func mutations(mutation FullType) {
	f, err := os.OpenFile(mutationFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	checkError(err)
	data, err := Asset("template/schema.tmpl")
	checkError(err)

	t := template.Must(temp(string(data)))
	err = t.Execute(f, mutation)
	checkError(err)
}

func scalars(scalars []FullType) {
	f, err := os.OpenFile(scalarFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	checkError(err)
	data, err := Asset("template/scalar.tmpl")
	checkError(err)

	t := template.Must(temp(string(data)))
	err = t.Execute(f, scalars)
	checkError(err)
}

func objects(scalars []FullType) {
	f, err := os.OpenFile(objectFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	checkError(err)
	data, err := Asset("template/object.tmpl")
	checkError(err)

	t := template.Must(temp(string(data)))
	err = t.Execute(f, scalars)
	checkError(err)
}

func enums(enums []FullType) {
	f, err := os.OpenFile(enumFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	checkError(err)
	data, err := Asset("template/enum.tmpl")
	checkError(err)

	t := template.Must(temp(string(data)))
	err = t.Execute(f, enums)
	checkError(err)
}

func interfaces(interfaces []FullType) {
	f, err := os.OpenFile(interfaceFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	checkError(err)
	data, err := Asset("template/interface.tmpl")
	checkError(err)

	t := template.Must(temp(string(data)))
	err = t.Execute(f, interfaces)
	checkError(err)
}

func temp(data string) (*template.Template, error) {
	p, err := template.New("MD").Funcs(template.FuncMap{
		"getType": func(t *TypeRef) interface{} {
			value := struct {
				Name string
				Type string
				Kind string
			}{Type: "%s"}
			for t.OfType != nil {
				if t.Kind == "NON_NULL" {
					value.Type = value.Type + "!"
				}
				if t.Kind == "LIST" {
					value.Type = "[" + value.Type + "]"
				}
				t = t.OfType
			}
			value.Name = t.Name
			value.Kind = t.Kind
			value.Type = fmt.Sprintf(value.Type, value.Name)
			if t.Kind == "SCALAR" {
				value.Name = dir + "scalar#" + value.Name
			}
			if t.Kind == "OBJECT" {
				value.Name = dir + "object#" + value.Name
			}
			value.Name = strings.Replace(strings.ToLower(value.Name), " ", "-", -1)
			return value

		},
	}).Parse(data)
	checkError(err)
	return p, err
}
