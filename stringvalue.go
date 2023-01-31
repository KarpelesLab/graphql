package graphql

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type StringValue string

func (s StringValue) String() string {
	// escape
	return "\"" + strconv.Quote(string(s)) + "\""
}

func (s StringValue) MarshalJSON() ([]byte, error) {
	res := map[string]any{
		"type":  "string_value",
		"value": string(s),
	}
	return json.Marshal(res)
}

func (p *Parser) parseStringValue() (Value, error) {
	// at this point.p.cur() == '"'

	// String value or blockstring.
	// Blockstring: """ [anything] """ (can contain \""" which becomes """)
	// string value cannot contain line terminator, but can contain many escapes including \\ \" \/ \b \f \n \r \t \u[0-9A-Fa-f]{4}

	if p.is("\"\"\"") {
		return nil, errors.New("blockstring not supported yet")
	}

	buf := &bytes.Buffer{}

	// read string char by char
	for {
		if err := p.next(); err != nil {
			return nil, unexpected(err)
		}
		c := p.cur()

		if c == '"' {
			// end of string, finally!
			p.nextNotSpace()
			return StringValue(buf.String()), nil
		}
		if c == '\n' {
			// error
			return nil, errors.New("string value cannot contain LineTerminator")
		}
		if c == '\\' {
			// escape value
			if err := p.next(); err != nil {
				return nil, unexpected(err)
			}
			c = p.cur()

			switch c {
			case '\\', '/', '"':
				// insert as is
				buf.WriteByte(c)
			case 'b':
				buf.WriteByte('\b')
			case 'f':
				buf.WriteByte('\f')
			case 'n':
				buf.WriteByte('\n')
			case 'r':
				buf.WriteByte('\r')
			case 't':
				buf.WriteByte('\t')
			case 'u':
				// read unicode
				// \u[0-9A-Fa-f]{4}
				uv, err := p.take(4)
				if err != nil {
					return nil, err
				}
				// need to parse hex value ([0-9a-fA-F])
				v, err := strconv.ParseUint(uv, 16, 32)
				if err != nil {
					return nil, err
				}
				buf.WriteRune(rune(v))
			default:
				return nil, fmt.Errorf("invalid escape sequence in StringValue: \\%c", c)
			}
		}

		buf.WriteByte(c)
	}
}
