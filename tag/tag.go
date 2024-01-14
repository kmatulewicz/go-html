/*
The go-html/tag package provides a convenient and flexible method to search for an HTML tag with a specific name and attributes.
Useful for web crawlers to quickly extract data from websites.
It does not implement the full HTML specification, so there might be cases where it will not work correctly.
*/
package tag

import (
	"strings"
)

// Tag is a representation of an HTML Tag found in doc.
type Tag struct {
	Name              string            // The name of the tag.
	Attr              map[string]string // The map of attributes map[attr_name]attr_val. Attribute names are always lowercase.
	ContentIndex      int               // The index points to the next character after the opening tag's closure in doc (it might be outside the doc range).
	AfterClosureIndex int               // The index points to the next character after the closing tag's closure in doc (it might be outside the doc range).
	doc               string            // A String where the tag was found.
	checks            []Check           // A slice of check functions used to find the tag
}

// Content returns a string between the starting tag and the closing tag of t.
func (t *Tag) Content() string {
	// return empty content if t is nil
	if t == nil {
		return ""
	}

	if t.AfterClosureIndex < 1 {
		// there is no closing tag; return empty content
		return ""
	}

	// find the index of the closing tag's beginning
	closureIndex := strings.LastIndex(t.doc[t.ContentIndex:t.AfterClosureIndex], "</")

	// return content between those tags
	return t.doc[t.ContentIndex : t.ContentIndex+closureIndex]
}

// Return the next *Tag with the same name and check functions
func (t *Tag) Next() *Tag {
	if t == nil {
		return nil
	}

	newT := Find(t.doc[t.ContentIndex:], t.Name, t.checks)
	if newT == nil {
		return nil
	}

	// change the starting point of the doc
	newT.doc = t.doc
	newT.ContentIndex += t.ContentIndex
	newT.AfterClosureIndex += t.ContentIndex

	return newT
}

// Find returns a *Tag struct representing a tag found in the s string, which has the n name and satisfies all f functions.
func Find(s string, n string, f []Check) *Tag {
	// starting at the beginning of s
	pos := 0

loop:
	// as far as the end of s is not reached
	for pos < len(s)-1 {

		//search for the start and end positions
		start, end := findStartTag(s[pos:], n)
		if start == -1 {
			// no such tag
			return nil
		}

		// set the positions in relation to s
		start += pos
		end += pos

		// create a tag for f checks
		t := &Tag{
			Name:              n,
			Attr:              parseAttribute(s[start+len(n)+1 : end-1]),
			ContentIndex:      end,
			AfterClosureIndex: getAfterClosureIndex(s, n, end),
			doc:               s,
			checks:            f,
		}

		// check if t will pass all f
		if !passChecks(f, t) {
			// continue after the current tag if checks failed
			pos = end
			continue loop
		}

		// return a found tag
		return t
	}

	// no such tag
	return nil
}

// passChecks returns true if t pass all checks; returns false if not
func passChecks(checks []Check, t *Tag) bool {
	// loop over all checks
	for _, c := range checks {
		if !c(t) {
			// the tag does not satisfy the function
			return false
		}
	}

	// the tag does satisfy the function
	return true
}
