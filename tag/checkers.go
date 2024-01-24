package tag

import "strings"

// Check is a type of function that takes *Tag as an argument and returns a boolean value.
// []Check is used as an argument for the Find function.
type Check func(*Tag) bool

// Has determines if the attribute of the given name exists in the tag.
// attr is case-insensitive.
func Has(attr string) Check {
	return func(t *Tag) bool {
		_, ok := t.Attr[strings.ToLower(attr)]
		return ok
	}
}

// Contains determines if the value of the attr attribute contains the s string.
// Returns false if the attribute does not exist.
// attr is case-insensitive.
func Contains(attr, s string) Check {
	return func(t *Tag) bool {
		v, ok := t.Attr[strings.ToLower(attr)]
		if !ok {
			return false
		}

		return strings.Contains(v, s)
	}
}

// Equal determines if the value of the attr attribute is equal to the s string.
// Returns false if the attribute does not exist.
// attr is case-insensitive.
func Equal(attr, s string) Check {
	return func(t *Tag) bool {
		v, ok := t.Attr[strings.ToLower(attr)]
		if !ok {
			return false
		}

		return v == s
	}
}

// NotEmpty determines if the value of the attr attribute is not empty.
// Returns false if the attribute does not exist.
// attr is case-insensitive.
func NotEmpty(attr string) Check {
	return func(t *Tag) bool {
		v, ok := t.Attr[strings.ToLower(attr)]
		if !ok || len(v) == 0 {
			return false
		}

		return true
	}
}
