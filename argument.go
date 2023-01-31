package graphql

import (
	"fmt"
	"strings"
)

type Arguments map[string]Value

func (a Arguments) String() string {
	if a == nil {
		return ""
	}
	var t []string

	for k, v := range a {
		t = append(t, k+":"+v.String())
	}
	return "(" + strings.Join(t, " ") + ")"
}

func (p *Parser) parseArguments() (Arguments, error) {
	if p.cur() != '(' {
		return nil, fmt.Errorf("expected ( for beginning of arguments, got %c", p.cur())
	}
	if err := p.nextNotSpace(); err != nil {
		return nil, unexpected(err)
	}

	args := make(Arguments)

	for {
		if p.cur() == ')' {
			// end of args
			p.nextNotSpace()
			return args, nil
		}
		// expect: name: value
		if !p.isName() {
			return nil, fmt.Errorf("expected name in arguments but got a %c", p.cur())
		}
		name := p.readName()
		if p.cur() != ':' {
			return nil, fmt.Errorf("expected a colon after name, got a %c", p.cur())
		}
		if err := p.nextNotSpace(); err != nil {
			return nil, unexpected(err)
		}
		val, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		args[name] = val
	}
}
