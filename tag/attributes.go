package tag

import (
	"io"
	"strings"
	"unicode"
)

// state is an enum type
type state int

const (
	bn  state = iota // Before attribute name state
	n                // Attribute name state
	an               // After attribute name state
	bv               // Before attribute value state
	vdq              // Attribute value (double-quoted) state
	vsq              // Attribute value (single-quoted) state
	v                // Attribute value (unquoted) state
	avq              // After attribute value (quoted) state
)

// parseAttribute parses the s string into a map of attributes. The attribute name is always changed to lowercase.
//
// Loosely inspired by the algorithm described here:
// https://html.spec.whatwg.org/multipage/parsing.html#before-attribute-name-state
func parseAttribute(s string) map[string]string {

	// start with an empty map
	attr := map[string]string{}

	// create a reader from the s string
	reader := strings.NewReader(s)

	// the starting state is: Before attribute name state
	state := bn

	name := ""
	value := ""

	// read all the characters one by one
	for r, _, err := reader.ReadRune(); err == nil; r, _, err = reader.ReadRune() {

		// choose a state
		switch state {
		case bn:
			// Before attribute name state
			switch {
			case unicode.IsSpace(r):
				// ignore the character
			case r == '=':
				// unexpected char
				return map[string]string{}
			default:
				// reconsume in the attribute name state
				state = n
				reader.Seek(-1, io.SeekCurrent)
			}
		case n:
			// Attribute name state
			switch {
			case unicode.IsSpace(r):
				// switch to the after attribute name state
				state = an
			case r == '=':
				// switch to the before attribute value state
				state = bv
			case r == 0 || r == '"' || r == '\'' || r == '<':
				// unexpected char
				return map[string]string{}
			default:
				name += string(unicode.ToLower(r))
			}
		case an:
			// After attribute name state
			switch {
			case unicode.IsSpace(r):
				// ignore the character
			case r == '=':
				// switch to the before attribute value state
				state = bv
			default:
				// emit attribute without value; reconsume in the attribute name state
				state = n
				attr[name] = ""
				name = ""
				reader.Seek(-1, io.SeekCurrent)
			}
		case bv:
			// Before attribute value state
			switch {
			case unicode.IsSpace(r):
				// ignore the character
			case r == '"':
				// switch to the attribute value (double-quoted) state
				state = vdq
			case r == '\'':
				// switch to the attribute value (single-quoted) state
				state = vsq
			default:
				// reconsume in the attribute value (unquoted) state
				state = v
				reader.Seek(-1, io.SeekCurrent)
			}
		case vdq:
			// Attribute value (double-quoted) state
			switch {
			case r == '"':
				// emit attribute without value; switch to the after attribute value (quoted) state
				attr[name] = value
				name = ""
				value = ""
				state = avq
			default:
				// append the current input character to the current attribute's value
				value += string(r)
			}
		case vsq:
			// Attribute value (single-quoted) state
			switch {
			case r == '\'':
				// emit attribute without value; switch to the after attribute value (quoted) state
				attr[name] = value
				name = ""
				value = ""
				state = avq
			default:
				// append the current input character to the current attribute's value
				value += string(r)
			}
		case v:
			// Attribute value (unquoted) state
			switch {
			case unicode.IsSpace(r):
				// emit attribute without value; switch to the before attribute name state
				attr[name] = value
				name = ""
				value = ""
				state = bn
			case strings.ContainsAny(string(r), "\"'<=`"):
				// unexpected char
				return map[string]string{}
			default:
				// append the current input character to the current attribute's value
				value += string(r)
			}

		case avq:
			// After attribute value (quoted) state
			switch {
			case unicode.IsSpace(r):
				// switch to the before attribute name state
				state = bn
			case r == '/':
				// reconsume in the attribute value (unquoted) state
				state = n
				reader.Seek(-1, io.SeekCurrent)
			default:
				// unexpected char
				return map[string]string{}
			}
		}
	}

	// append the last attribute
	if name != "" {
		attr[name] = value
	}

	return attr
}
