package graphql

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Field struct {
	Alias        string       `json:"alias,omitempty"`
	Name         string       `json:"name"`
	Arguments    Arguments    `json:"arguments,omitempty"`
	Directives   Directives   `json:"directives,omitempty"`
	SelectionSet SelectionSet `json:"selection_set,omitempty"`
}

func (f *Field) String() string {
	var t []string

	if f.Alias == "" {
		t = append(t, f.Name)
	} else {
		t = append(t, f.Alias+":", f.Name)
	}

	if v := f.Arguments.String(); v != "" {
		t = append(t, v)
	}
	if v := f.Directives.String(); v != "" {
		t = append(t, v)
	}
	if v := f.SelectionSet.String(); v != "" {
		t = append(t, v)
	}
	return strings.Join(t, " ")
}

func (f *Field) MarshalJSON() ([]byte, error) {
	res := map[string]any{
		"type": "field",
		"name": f.Name,
	}
	if f.Alias != "" {
		res["alias"] = f.Alias
	}
	if f.Arguments != nil {
		res["arguments"] = f.Arguments
	}
	if f.Directives != nil {
		res["directives"] = f.Directives
	}
	if f.SelectionSet != nil {
		res["selection_set"] = f.SelectionSet
	}
	return json.Marshal(res)
}

func (p *Parser) parseField() (*Field, error) {
	// parse a field: Alias Name Arguments Directives SelectionSet
	// only Name is required
	if !p.isName() {
		return nil, fmt.Errorf("expected to find a field name but got a %c", p.cur())
	}
	f := &Field{}
	f.Name = p.readName()

	if p.cur() == ':' {
		// that was actually an alias
		f.Alias = f.Name

		if err := p.nextNotSpace(); err != nil {
			return nil, unexpected(err)
		}
		if !p.isName() {
			return nil, fmt.Errorf("expected to find a field name after alias but got a %c", p.cur())
		}
		f.Name = p.readName()
	}

	if p.cur() == '(' {
		// arguments
		args, err := p.parseArguments()
		if err != nil {
			return nil, err
		}
		f.Arguments = args
	}
	if p.cur() == '@' {
		// directives
		dir, err := p.parseDirectives()
		if err != nil {
			return nil, err
		}
		f.Directives = dir
	}
	if p.cur() == '{' {
		// selection set
		sl, err := p.parseSelectionSet()
		if err != nil {
			return nil, err
		}
		f.SelectionSet = sl
	}

	return f, nil
}
