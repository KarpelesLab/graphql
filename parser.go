package graphql

import (
	"bytes"
	"io"
	"strings"
	"unicode"
)

// see: https://spec.graphql.org/June2018/

type Parser struct {
	str string
	pos int

	// parser state
	doc *Document
}

func Parse(v string) (*Document, error) {
	p := &Parser{str: v, doc: newDocument()}
	err := p.parse()
	if err != nil {
		return nil, err
	}
	return p.doc, nil
}

func (p *Parser) parse() error {
	// seek to first token
	if err := p.skipSpaces(); err != nil {
		if err == io.EOF {
			// nothing in this query
			return nil
		}
		return err
	}

	// main parser loop
	for {
		err := p.parseOperation()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

// cur returns the current byte
func (p *Parser) cur() byte {
	if p.eof() {
		return 0
	}
	return p.str[p.pos]
}

// peek returns the next character
func (p *Parser) peek() byte {
	if p.pos+1 >= len(p.str) {
		return 0
	}
	return p.str[p.pos+1]
}

func (p *Parser) eof() bool {
	return p.pos == len(p.str)
}

// buf returns a []byte of what remains in the buffer
func (p *Parser) buf() string {
	return p.str[p.pos:]
}

// take returns the first N bytes after moving the read point
func (p *Parser) take(ln int) (string, error) {
	b := p.buf()
	if len(b) < ln {
		return "", io.ErrUnexpectedEOF
	}
	v := b[:ln]
	p.pos += ln

	return v, nil
}

func (p *Parser) is(pfx string) bool {
	return strings.HasPrefix(p.str[p.pos:], pfx)
}

// skip to next character that is not a space
func (p *Parser) nextNotSpace() error {
	return p.skip(1)
}

// skip n chars and then drops any spaces
func (p *Parser) skip(n int) error {
	if p.pos+n > len(p.str) {
		p.pos = len(p.str)
		return io.EOF
	}
	p.pos += n
	return p.skipSpaces()
}

// next advances the parser to the next character
func (p *Parser) next() error {
	if p.pos >= len(p.str) {
		return io.EOF
	}
	p.pos += 1
	return nil
}

// skipSpaces advances the parser until the next non-space character, or does
// nothing if already as a non-space character
func (p *Parser) skipSpaces() error {
	for {
		c := p.cur()
		if unicode.IsSpace(rune(c)) || c == 0 || c == ',' {
			if err := p.next(); err != nil {
				return err
			}
			continue
		}
		if c == '#' {
			// Comment :: # [any char except line terminator]
			// skip until we find '\r' or '\n'
			for {
				if err := p.next(); err != nil {
					return err
				}
				c = p.cur()
				if (c == '\r') || (c == '\n') {
					break
				}
			}
		}
		// not a comment, not a space, so we found something
		return nil
	}
}

func (p *Parser) readName() string {
	// Name :: /[_A-Za-z][_0-9A-Za-z]*/
	buf := &bytes.Buffer{}
	first := true
	for {
		c := p.cur()
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			if first {
				// invalid as first char
				return ""
			}
			fallthrough
		case '_':
			fallthrough
		case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z':
			fallthrough
		case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
			buf.WriteByte(c)
			if err := p.next(); err == nil {
				first = false
				break
			}
			fallthrough
		default:
			p.skipSpaces()
			return buf.String()
		}
	}
}

// isName returns true if the current char is a name char
func (p *Parser) isName() bool {
	// Name :: /[_A-Za-z][_0-9A-Za-z]*/
	switch p.cur() {
	case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z':
		fallthrough
	case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
		fallthrough
	case '_':
		return true
	default:
		return false
	}
}
