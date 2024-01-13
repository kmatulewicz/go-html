package tag

import (
	"errors"
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

// parse is a struct representing the current state of a parsing process
type parse struct {
	i     int               //a current position in the string
	r     rune              // a current rune
	name  string            // a current name
	value string            // a current value
	attr  map[string]string // a map of saved attributes
	state state             // a current state
}

// parseAttribute parses the s string into a map of attributes. The attribute name is always changed to lowercase.
//
// Loosely inspired by the algorithm described here:
// https://html.spec.whatwg.org/multipage/parsing.html#before-attribute-name-state
func parseAttribute(s string) map[string]string {

	// start with an empty map and the starting state is: Before attribute name state
	p := parse{
		attr:  map[string]string{},
		state: bn,
	}

	// read all the characters one by one
	for ; p.i < len(s); p.i++ {

		// read a rune
		p.r = rune(s[p.i])

		// choose a state
		switch p.state {
		case bn:
			// Before attribute name state
			if beforeName(&p) != nil {
				return map[string]string{}
			}
		case n:
			// Attribute name state
			if name(&p) != nil {
				return map[string]string{}
			}
		case an:
			// After attribute name state
			afterName(&p)
		case bv:
			// Before attribute value state
			beforeValue(&p)
		case vdq:
			// Attribute value (double-quoted) state
			valueDQ(&p)
		case vsq:
			// Attribute value (single-quoted) state
			valueSQ(&p)
		case v:
			// Attribute value (unquoted) state
			if value(&p) != nil {
				return map[string]string{}
			}
		case avq:
			// After attribute value (quoted) state
			if afterValueQ(&p) != nil {
				return map[string]string{}
			}
		}
	}

	// append the last attribute
	if p.name != "" {
		p.attr[p.name] = p.value
	}

	return p.attr
}

func beforeName(p *parse) error {
	// Before attribute name state
	switch {
	case unicode.IsSpace(p.r):
		// ignore the character
	case p.r == '=':
		// unexpected char
		return errors.New("unexpected char")
	default:
		// reconsume in the attribute name state
		p.state = n
		p.i--
	}

	return nil
}

func name(p *parse) error {
	// Attribute name state
	switch {
	case unicode.IsSpace(p.r):
		// switch to the after attribute name state
		p.state = an
	case p.r == '=':
		// switch to the before attribute value state
		p.state = bv
	case p.r == 0 || p.r == '"' || p.r == '\'' || p.r == '<':
		// unexpected char
		return errors.New("unexpected char")
	default:
		p.name += string(unicode.ToLower(p.r))
	}

	return nil
}

func afterName(p *parse) {
	// After attribute name state
	switch {
	case unicode.IsSpace(p.r):
		// ignore the character
	case p.r == '=':
		// switch to the before attribute value state
		p.state = bv
	default:
		// emit attribute without value; reconsume in the attribute name state
		p.state = n
		p.attr[p.name] = ""
		p.name = ""
		p.i--
	}
}

func beforeValue(p *parse) {
	// Before attribute value state
	switch {
	case unicode.IsSpace(p.r):
		// ignore the character
	case p.r == '"':
		// switch to the attribute value (double-quoted) state
		p.state = vdq
	case p.r == '\'':
		// switch to the attribute value (single-quoted) state
		p.state = vsq
	default:
		// reconsume in the attribute value (unquoted) state
		p.state = v
		p.i--
	}
}

func valueDQ(p *parse) {
	// Attribute value (double-quoted) state
	switch {
	case p.r == '"':
		// emit attribute without value; switch to the after attribute value (quoted) state
		p.attr[p.name] = p.value
		p.name = ""
		p.value = ""
		p.state = avq
	default:
		// append the current input character to the current attribute's value
		p.value += string(p.r)
	}
}

func valueSQ(p *parse) {
	// Attribute value (single-quoted) state
	switch {
	case p.r == '\'':
		// emit attribute without value; switch to the after attribute value (quoted) state
		p.attr[p.name] = p.value
		p.name = ""
		p.value = ""
		p.state = avq
	default:
		// append the current input character to the current attribute's value
		p.value += string(p.r)
	}
}

func value(p *parse) error {
	// Attribute value (unquoted) state
	switch {
	case unicode.IsSpace(p.r):
		// emit attribute without value; switch to the before attribute name state
		p.attr[p.name] = p.value
		p.name = ""
		p.value = ""
		p.state = bn
	case strings.ContainsAny(string(p.r), "\"'<=`"):
		// unexpected char
		return errors.New("unexpected char")
	default:
		// append the current input character to the current attribute's value
		p.value += string(p.r)
	}

	return nil
}

func afterValueQ(p *parse) error {
	// After attribute value (quoted) state
	switch {
	case unicode.IsSpace(p.r):
		// switch to the before attribute name state
		p.state = bn
	case p.r == '/':
		// reconsume in the attribute value (unquoted) state
		p.state = n
		p.i--
	default:
		// unexpected char
		return errors.New("unexpected char")
	}

	return nil
}
