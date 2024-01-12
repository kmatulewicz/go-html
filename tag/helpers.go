package tag

import (
	"strings"
	"unicode"
)

// getAfterClosureIndex returns the index of the next character after the closing tag's end
// for the tag named n of the document doc, starting at the index i.
// Returns -1 if there is no closing tag.
func getAfterClosureIndex(doc, n string, i int) int {
	// number of closing tags to be found
	count := 1
	// start at the beginning of a content
	pos := i

	// until some closing tags are left to be found
	for count > 0 {
		// check if the position is within the bounds of doc
		if pos >= len(doc)-1 {
			return -1
		}

		// find the next closure tag
		_, closure := findEndTag(doc[pos:], n)
		if closure == -1 {
			// there is no closure
			return -1
		}

		// find the next tag with the same name
		_, next := findStartTag(doc[pos:], n)

		// check if there is no next tag with the same name, or if it is after the closure tag
		if next == -1 || closure < next {
			// decrease the number of closing tags to be found
			count--
			// set the position after the beginning of the closing tag.
			pos += closure
		} else {
			// if the next tag with the same name is before the closing tag, increase the number of closing tags to be found
			count++
			// and set the position at the beginning of the next tag with the same name
			pos += next
		}
	}

	// return the index of the next character after the closing tag's end
	return pos
}

// findTag returns the start and end points of the tag with the name n (n includes the tag's opening sequence) in s
func findTag(s, n string) (start, end int) {
	// start at the beginning of s
	pos := 0
	// until reached the end of s
	for pos < len(s)-1 {
		// localize the beginning of a tag, starting from the position pos
		start = strings.Index(s[pos:], n)
		if start == -1 {
			// there is no tag with n name in s after pos
			return -1, -1
		}
		// set the start point in relation to s
		start = start + pos

		// localize the nearest closure of a found tag
		end = strings.Index(s[start:], ">")
		if end == -1 {
			// there are no closures in s after pos
			return -1, -1
		}
		// set the end point in relation to s
		end = start + end + 1

		// calculate the position of the next char after the found tag name
		endNameCharIndex := start + len(n)
		// check if the found tag name is not only the beginning of the proper name (after the name is some other character, then a closure or a white space)
		if !(unicode.IsSpace(rune(s[endNameCharIndex])) || s[endNameCharIndex] == '>' || s[endNameCharIndex] == '/') {
			// continue after the current tag
			pos = end
			continue
		}

		// return the positions
		return start, end
	}

	// no such tag
	return -1, -1
}

// findStartTag returns the start and end points of the opening tag with the name n in s
func findStartTag(s, n string) (start, end int) {
	return findTag(s, "<"+n)
}

// findStartTag returns the start and end points of the closing tag with the name n in s
func findEndTag(s, n string) (start, end int) {
	return findTag(s, "</"+n)
}
