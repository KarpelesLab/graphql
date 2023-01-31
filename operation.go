package graphql

import (
	"errors"
	"fmt"
	"strings"
)

type OperationType int

const (
	Query        OperationType = iota // query – a read‐only fetch.
	Mutation                          // mutation – a write followed by a fetch.
	Subscription                      // subscription – a long‐lived request that fetches data in response to source events.
)

func (o OperationType) String() string {
	switch o {
	case Query:
		return "query"
	case Mutation:
		return "mutation"
	case Subscription:
		return "subscription"
	default:
		panic("invalid operation")
	}
}

func (o OperationType) MarshalJSON() ([]byte, error) {
	switch o {
	case Query:
		return []byte(`"query"`), nil
	case Mutation:
		return []byte(`"mutation"`), nil
	case Subscription:
		return []byte(`"subscription"`), nil
	default:
		return nil, errors.New("invalid operation")
	}
}

type Operation struct {
	OperationType       OperationType       `json:"type"`
	Name                string              `json:"name"` // optional
	VariableDefinitions VariableDefinitions `json:"variable_definitions,omitempty"`
	Directives          Directives          `json:"directives,omitempty"`
	SelectionSet        SelectionSet        `json:"selection_set"`
}

func (op *Operation) String() string {
	vals := []string{
		op.OperationType.String(),
		op.Name,
		op.VariableDefinitions.String(),
		op.Directives.String(),
		op.SelectionSet.String(),
	}
	return strings.Join(vals, " ")
}

func (p *Parser) parseOperation() error {
	var err error

	// parse an operation
	if err = p.skipSpaces(); err != nil {
		return err
	}

	// default value
	op := &Operation{OperationType: Query}
	if p.isName() {
		opName := p.readName()
		switch strings.ToLower(opName) {
		case "query":
			op.OperationType = Query
		case "mutation":
			op.OperationType = Mutation
		case "subscription":
			op.OperationType = Subscription
		case "fragment":
			// go to fragment reading
			return p.readFragment()
		default:
			return fmt.Errorf("invalid operation %s", opName)
		}
		if p.isName() {
			// read actual operation name
			op.Name = p.readName()
		}
	}

	if p.cur() == '(' {
		op.VariableDefinitions, err = p.parseVariableDefinitions()
		if err != nil {
			return err
		}
	}

	op.Directives, err = p.parseDirectives()
	if err != nil {
		return err
	}

	if p.cur() != '{' {
		return fmt.Errorf("invalid operation, unexpected token %c", p.cur())
	}

	sl, err := p.parseSelectionSet()
	if err != nil {
		return fmt.Errorf("in %s %s: %w", op.OperationType, op.Name, err)
	}
	op.SelectionSet = sl

	// add operation to p.doc
	if _, ok := p.doc.Operations[op.Name]; ok {
		return fmt.Errorf("duplicate operation name %s", op.Name)
	}
	p.doc.Operations[op.Name] = op
	return nil
}
