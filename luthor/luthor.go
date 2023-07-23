package luthor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/iancoleman/strcase"

	"github.com/getnoops/terraform-gen/luthor/ast"
	"github.com/getnoops/terraform-gen/luthor/parser"
)

var (
	upperTags = []string{"json", "yaml", "hcl"}
	lowerTags = []string{"json", "yaml", "cty"}
)

type Data struct {
	Modules []*ModuleData
	Structs []*StructType
}

type VariableData struct {
	Name        string
	Description string
	Type        *ast.Type
	Required    bool
}

type ModuleData struct {
	Name      string
	TypeName  string
	Module    *tfconfig.Module
	Variables []*VariableData
}

type StructType struct {
	Name   string
	Type   string
	Fields []*FieldType
}

type FieldType struct {
	Description string
	Name        string
	Type        string
	Tag         string
}

func LoadData(paths []string, ignoreList []string) (*Data, error) {
	var modules []*ModuleData
	for _, dir := range paths {
		matches, err := filepath.Glob(dir)
		if err != nil {
			return nil, fmt.Errorf("globbing module directories: %s", err)
		}
		for _, match := range matches {
			info, err := os.Stat(match)
			if err != nil {
				return nil, fmt.Errorf("statting module: %s", err)
			}
			if !info.IsDir() {
				continue
			}

			module, diags := tfconfig.LoadModule(match)
			if diags.HasErrors() {
				return nil, fmt.Errorf("loading module: %s", diags.Err())
			}

			name := filepath.Base(match)
			typeName := strcase.ToCamel(name)

			// variables.
			var variables []*VariableData
			for _, v := range module.Variables {
				t, err := parser.ParseType(&ast.Source{
					Name:  v.Name,
					Input: v.Type,
				})
				if err != nil {
					fmt.Println(v.Type)
					return nil, fmt.Errorf("parsing type: %s", err)
				}

				variables = append(variables, &VariableData{
					Name:        v.Name,
					Description: v.Description,
					Type:        t,
					Required:    v.Required,
				})
			}

			modules = append(modules, &ModuleData{
				Name:      name,
				TypeName:  typeName,
				Module:    module,
				Variables: variables,
			})
		}
	}

	builder := &structBuilder{
		ignoreList: ignoreList,
	}
	for _, m := range modules {
		builder.addModule(m)
	}

	data := &Data{
		Modules: modules,
		Structs: builder.structs,
	}
	return data, nil
}

type structBuilder struct {
	structs []*StructType
	err     error

	ignoreList []string
}

func (b *structBuilder) shouldIgnore(name string) bool {
	for _, ignore := range b.ignoreList {
		if strings.EqualFold(ignore, name) {
			return true
		}
	}
	return false
}

func (b *structBuilder) addModule(m *ModuleData) {
	if b.err != nil {
		return
	}

	s := &StructType{
		Name: m.TypeName,
		Type: "struct",
	}
	b.structs = append(b.structs, s)

	for _, v := range m.Variables {
		if b.shouldIgnore(v.Name) {
			continue
		}

		b.addField(s, m.TypeName, v.Description, v.Name, v.Required, v.Type, upperTags)
	}
}

func (b *structBuilder) addField(s *StructType, prefix string, description string, name string, required bool, t *ast.Type, tags []string) {
	if b.err != nil {
		return
	}

	proper := strcase.ToCamel(name)
	fieldType, err := b.fieldType(prefix+proper, t)
	if err != nil {
		b.err = err
		return
	}

	tagList := make([]string, len(tags))
	for i, t := range tags {
		tagList[i] = fmt.Sprintf("%s:\"%s\"", t, name)
	}
	if !required && !strings.HasPrefix(fieldType, "*") {
		fieldType = "*" + fieldType
	}

	s.Fields = append(s.Fields, &FieldType{
		Description: description,
		Name:        proper,
		Type:        fieldType,
		Tag:         strings.Join(tagList, " "),
	})
}

func (b *structBuilder) fieldType(name string, t *ast.Type) (string, error) {
	p := ""
	if t.Optional {
		p = "*"
	}

	if t.Kind == ast.StringValue {
		return p + "string", nil
	}
	if t.Kind == ast.BooleanValue {
		return p + "bool", nil
	}
	if t.Kind == ast.NumberValue {
		return p + "string", nil
	}
	if t.Kind == ast.NullValue {
		return p + "interface{}", nil
	}
	if t.Kind == ast.AnyValue {
		return p + "interface{}", nil
	}
	if t.Kind == ast.MapValue {
		innerT, err := b.fieldType(name, t.Elem)
		if err != nil {
			return "", err
		}
		return p + "map[string]" + innerT, nil
	}
	if t.Kind == ast.ListValue {
		innerT, err := b.fieldType(name, t.Elem)
		if err != nil {
			return "", err
		}
		return p + "[]" + innerT, nil
	}
	if t.Kind == ast.ObjectValue {
		// we need to create a new struct.
		nt, err := b.structType(name, t)
		if err != nil {
			return "", err
		}
		return p + nt, nil
	}

	return "", fmt.Errorf("unsupported type: %s", t.Name)
}

func (b *structBuilder) structType(name string, t *ast.Type) (string, error) {
	s := &StructType{
		Name: name,
		Type: "struct",
	}
	b.structs = append(b.structs, s)

	for _, f := range t.Fields {
		b.addField(s, name, "", f.Name, !f.Type.Optional, f.Type, lowerTags)
	}

	return name, nil
}
