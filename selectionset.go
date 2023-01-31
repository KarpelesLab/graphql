package graphql

import (
	"errors"
	"fmt"
	"strings"
)

type SelectionSet []Selection

func (s SelectionSet) String() string {
	if s == nil {
		return ""
	}

	var t []string
	for _, sub := range s {
		t = append(t, sub.String())
	}

	return "{" + strings.Join(t, " ") + "}"
}

type Selection interface {
	// TODO
	String() string
}

func (p *Parser) parseSelectionSet() (SelectionSet, error) {
	if p.cur() != '{' {
		return nil, fmt.Errorf("unexpected char %c while looking for SelectionSet start", p.cur())
	}
	if err := p.nextNotSpace(); err != nil {
		return nil, unexpected(err)
	}

	var res SelectionSet

	for {
		// can be either a field, or "..." followed by a named type (fragment spread) or "..." followed by optionally "on X" then "{"
		if p.is("...") {
			// that's a fragment thing
			if err := p.skip(3); err != nil {
				return nil, unexpected(err)
			}
			if !p.isName() {
				if p.cur() == '{' {
					// that's an inline fragment
					sl, err := p.parseSelectionSet() // yay for recursion
					if err != nil {
						return nil, err
					}
					res = append(res, &InlineFragment{SelectionSet: sl})
					continue
				}
				return nil, errors.New("... in a selection set must be followed by a name or a {")
			}
			frag := p.readName()
			if frag == "on" {
				cond, err := p.readTypeCondition() // we already have the "on"
				if err != nil {
					return nil, err
				}
				// InlineFragment
				sl, err := p.parseSelectionSet() // yay for recursion
				if err != nil {
					return nil, err
				}
				res = append(res, &InlineFragment{TypeCondition: cond, SelectionSet: sl})
			}

			// FragmentSpread
			dir, err := p.parseDirectives()
			if err != nil {
				return nil, err
			}
			res = append(res, &FragmentSpread{Name: frag, Directives: dir})
			continue
		}

		if p.cur() == '}' {
			// final
			p.nextNotSpace()
			return res, nil
		}

		f, err := p.parseField()
		if err != nil {
			return nil, err
		}

		res = append(res, f)
	}
}
