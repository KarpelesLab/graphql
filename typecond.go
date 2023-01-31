package graphql

import "errors"

type TypeCondition struct {
	NamedType string // "on X"
}

func (t *TypeCondition) String() string {
	return "on " + t.NamedType
}

func (p *Parser) parseTypeCondition() (*TypeCondition, error) {
	// expect "on" followed by TypeCondition name
	if !p.isName() {
		return nil, errors.New("type condition expected, need \"on\"")
	}
	on := p.readName()
	if on != "on" {
		return nil, errors.New("type condition expected, need \"on\"")
	}
	return p.readTypeCondition()
}

func (p *Parser) readTypeCondition() (*TypeCondition, error) {
	if !p.isName() {
		return nil, errors.New("type condition \"on\" must be followed by a NamedType")
	}
	t := p.readName()

	return &TypeCondition{NamedType: t}, nil
}
