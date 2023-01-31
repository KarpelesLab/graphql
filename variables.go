package graphql

import (
	"errors"
	"fmt"
)

type VariableDefinition struct {
	Variable     string
	Type         string
	DefaultValue Value // optional
}

type VariableDefinitions []*VariableDefinition

func (v VariableDefinitions) String() string {
	if v == nil {
		return ""
	}
	return "TODO VariableDefinitions"
}

func (p *Parser) parseVariableDefinitions() (VariableDefinitions, error) {
	if p.cur() != '(' {
		return nil, nil
	}
	if err := p.skip(1); err != nil {
		return nil, unexpected(err)
	}

	var res VariableDefinitions

	for {
		if p.cur() == ')' {
			return res, unexpected(p.skip(1))
		}

		if p.cur() != '$' {
			return nil, fmt.Errorf("variable definition value must start with a $")
		}
		if err := p.next(); err != nil {
			return nil, unexpected(err)
		}
		name := p.readName()
		if name == "" {
			return nil, fmt.Errorf("variable name must be a name")
		}

		return nil, errors.New("TODO VariableDefinitions")

	}
}
