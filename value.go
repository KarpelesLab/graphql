package graphql

import (
	"encoding/json"
	"errors"
	"fmt"
)

// https://spec.graphql.org/June2018/#Value

type Value interface {
	// TODO
	String() string
}

type VariableValue struct {
	Var string
}

func (v *VariableValue) String() string {
	return "$" + v.Var
}

func (v *VariableValue) MarshalJSON() ([]byte, error) {
	res := map[string]any{
		"type":     "variable",
		"variable": v.Var,
	}
	return json.Marshal(res)
}

type EnumValue string

func (v EnumValue) String() string {
	return string(v)
}

func (v EnumValue) MarshalJSON() ([]byte, error) {
	res := map[string]any{
		"type":  "enum_value",
		"value": string(v),
	}
	return json.Marshal(res)
}

type BooleanValue bool

func (v BooleanValue) String() string {
	if v {
		return "true"
	} else {
		return "false"
	}
}

func (v BooleanValue) MarshalJSON() ([]byte, error) {
	res := map[string]any{
		"type":  "bool_value",
		"value": bool(v),
	}
	return json.Marshal(res)
}

type NullValue struct{}

func (v NullValue) String() string {
	return "null"
}

func (v NullValue) MarshalJSON() ([]byte, error) {
	res := map[string]any{
		"type": "null_value",
	}
	return json.Marshal(res)
}

func (p *Parser) parseValue() (Value, error) {
	// can be a number of things...
	switch p.cur() {
	case '$': // variable
		if err := p.next(); err != nil {
			return nil, unexpected(err)
		}
		if !p.isName() {
			return nil, errors.New("variable must be followed by a name")
		}
		return &VariableValue{Var: p.readName()}, nil
	case '"':
		return p.parseStringValue()
	default:
		if p.isName() {
			nam := p.readName()

			switch nam {
			case "true":
				return BooleanValue(true), nil
			case "false":
				return BooleanValue(false), nil
			case "null":
				return NullValue{}, nil
			}
			return EnumValue(nam), nil
		}
		return nil, fmt.Errorf("unsupported value character %c", p.cur())
	}
}
