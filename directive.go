package graphql

import "strings"

type Directive struct {
	Directive string
	Arguments Arguments
}

func (d *Directive) String() string {
	return "@" + d.Directive + " " + d.Arguments.String()
}

type Directives []*Directive

func (ds Directives) String() string {
	var r []string

	for _, d := range ds {
		r = append(r, d.String())
	}

	return strings.Join(r, " ")
}

func (p *Parser) parseDirectives() (Directives, error) {
	var res Directives

	for {
		if p.cur() != '@' {
			return res, nil
		}
		if err := p.next(); err != nil {
			return nil, unexpected(err)
		}
		name := p.readName()
		args, err := p.parseArguments()
		if err != nil {
			return nil, err
		}
		res = append(res, &Directive{Directive: name, Arguments: args})
	}
}
