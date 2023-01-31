package graphql

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Fragment struct {
	Name          string         `json:"name"`
	TypeCondition *TypeCondition `json:"type_condition,omitempty"`
	Directives    Directives     `json:"directives,omitempty"`
	SelectionSet  SelectionSet   `json:"selection_set"`
}

func (f *Fragment) String() string {
	vals := []string{
		"fragment",
		f.Name,
		f.TypeCondition.String(),
		f.Directives.String(),
		f.SelectionSet.String(),
	}
	return strings.Join(vals, " ")
}

type InlineFragment struct {
	TypeCondition *TypeCondition `json:"type_condition,omitempty"`
	Directives    Directives     `json:"directives,omitempty"`
	SelectionSet  SelectionSet   `json:"selection_set"`
}

func (f *InlineFragment) MarshalJSON() ([]byte, error) {
	res := map[string]any{
		"type":          "inline_fragment",
		"selection_set": f.SelectionSet,
	}
	if f.TypeCondition != nil {
		res["type_condition"] = f.TypeCondition
	}
	if f.Directives != nil {
		res["directives"] = f.Directives
	}

	return json.Marshal(res)
}

func (i *InlineFragment) String() string {
	return "..." + i.TypeCondition.String() + " " + i.Directives.String() + " " + i.SelectionSet.String()
}

type FragmentSpread struct {
	Name       string     `json:"selection_set"`
	Directives Directives `json:"directives,omitempty"`
}

func (f *FragmentSpread) String() string {
	return "..." + f.Name + " " + f.Directives.String()
}

func (f *FragmentSpread) MarshalJSON() ([]byte, error) {
	res := map[string]any{
		"type": "fragment_spread",
		"name": f.Name,
	}
	if f.Directives != nil {
		res["directives"] = f.Directives
	}
	return json.Marshal(res)
}

func (p *Parser) readFragment() error {
	// at this point we already read "fragment"
	// fragmentFragmentNameTypeConditionDirectivesoptSelectionSet
	f := &Fragment{}
	if !p.isName() {
		return fmt.Errorf("fragment must be followed by a name")
	}
	f.Name = p.readName()
	if f.Name == "on" {
		// Name but not "on"
		return fmt.Errorf("'on' is not a valid fragment name")
	}

	cond, err := p.parseTypeCondition()
	if err != nil {
		return err
	}
	f.TypeCondition = cond

	sl, err := p.parseSelectionSet()
	if err != nil {
		return err
	}
	f.SelectionSet = sl

	// add fragment to p.doc
	if _, ok := p.doc.Fragments[f.Name]; ok {
		return fmt.Errorf("duplicate fragment name %s", f.Name)
	}
	p.doc.Fragments[f.Name] = f
	return nil
}
