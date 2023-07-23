package ast

type Field struct {
	Name     string
	Type     *Type
	Position *Position `dump:"-"`
}

type TypeKind int

const (
	StringValue TypeKind = iota
	BooleanValue
	NumberValue
	NullValue
	ListValue
	MapValue
	ObjectValue
	AnyValue
)

type Type struct {
	Name     string
	Kind     TypeKind
	Elem     *Type
	Optional bool
	Fields   FieldList
	Position *Position `dump:"-"`
}

type FieldList []*Field

func (v FieldList) ForName(name string) *Field {
	for _, f := range v {
		if f.Name == name {
			return f
		}
	}
	return nil
}
