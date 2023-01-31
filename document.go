package graphql

import (
	"bytes"
	"fmt"
)

type Document struct {
	Operations map[string]*Operation `json:"operations"`
	Fragments  map[string]*Fragment  `json:"fragments,omitempty"`
}

func newDocument() *Document {
	return &Document{
		Operations: make(map[string]*Operation),
		Fragments:  make(map[string]*Fragment),
	}
}

func (d *Document) String() string {
	buf := &bytes.Buffer{}

	for _, v := range d.Operations {
		fmt.Fprintf(buf, "%s\n\n", v.String())
	}
	for _, v := range d.Fragments {
		fmt.Fprintf(buf, "%s\n\n", v.String())
	}
	return buf.String()
}
