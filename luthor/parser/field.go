package parser

import (
	. "github.com/getnoops/terraform-gen/luthor/ast"
	"github.com/getnoops/terraform-gen/luthor/lexer"
)

func ParseType(source *Source) (*Type, error) {
	p := parser{
		lexer: lexer.New(source),
	}
	return p.parseType(false), p.err
}

func (p *parser) parseType(optional bool) *Type {
	if peek := p.peek(); peek.Kind == lexer.Name && peek.Value == "optional" {
		p.skip(lexer.Name)
		p.expect(lexer.ParenL)
		t := p.parseType(true)
		p.expect(lexer.ParenR)
		return t
	}

	token := p.peek()
	var kind TypeKind
	switch token.Value {
	case "string":
		kind = StringValue
	case "bool":
		kind = BooleanValue
	case "number":
		kind = NumberValue
	case "null":
		kind = NumberValue
	case "any":
		kind = AnyValue
	case "map":
		return p.parseMap(optional)
	case "list":
		return p.parseList(optional)
	case "object":
		return p.parseObject(optional)
	default:
		p.unexpectedError()
		return nil
	}

	p.next()

	return &Type{Position: &token.Pos, Name: token.Value, Kind: kind, Optional: optional}
}

func (p *parser) parseMap(optional bool) *Type {
	pos := p.peekPos()
	name := p.parseName()

	p.expect(lexer.ParenL)
	t := p.parseType(false)
	p.expect(lexer.ParenR)

	return &Type{Name: name, Kind: MapValue, Elem: t, Optional: optional, Position: pos}
}

func (p *parser) parseList(optional bool) *Type {
	pos := p.peekPos()
	name := p.parseName()

	p.expect(lexer.ParenL)
	t := p.parseType(false)
	p.expect(lexer.ParenR)

	return &Type{Name: name, Kind: ListValue, Elem: t, Optional: optional, Position: pos}
}

func (p *parser) parseObject(optional bool) *Type {
	pos := p.peekPos()
	name := p.parseName()
	var fields FieldList

	p.expect(lexer.ParenL)
	p.many(lexer.BraceL, lexer.BraceR, func() {
		fields = append(fields, p.parseObjectField())
	})
	p.expect(lexer.ParenR)

	return &Type{Name: name, Kind: ObjectValue, Optional: optional, Fields: fields, Position: pos}
}

func (p *parser) parseObjectField() *Field {
	pos := p.peekPos()
	name := p.parseName()
	p.expect(lexer.Equals)
	t := p.parseType(false)

	return &Field{Name: name, Type: t, Position: pos}
}

func (p *parser) parseName() string {
	token := p.expect(lexer.Name)

	return token.Value
}
