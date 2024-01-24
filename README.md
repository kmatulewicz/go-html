[![Build Status](https://github.com/kmatulewicz/go-html/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/kmatulewicz/go-html/actions/workflows/go.yml?query=branch%3Amain)
[![Go Report Card](https://goreportcard.com/badge/github.com/kmatulewicz/go-html)](https://goreportcard.com/report/github.com/kmatulewicz/go-html)
[![codecov](https://codecov.io/gh/kmatulewicz/go-html/graph/badge.svg?token=TWCZIJDDCB)](https://codecov.io/gh/kmatulewicz/go-html)

# go-html/tag

The go-html/tag package provides a convenient and flexible method to search for an HTML tag with a specific name and attributes. It is useful for web crawlers to quickly extract data from websites. It does not implement the full HTML specification, so there might be cases where it will not work correctly.

### Installation

```sh
go get github.com/kmatulewicz/go-html
```

### Usage

#### Find function

The function `Find(s string, n string, f []Check) *Tag` is used to find the specified HTML tag; it takes as arguments:

- **s** - a string containing HTML where the tag needs to be found
- **n** - the name of the tag you are looking for
- **f** - a slice of Check functions used to validate the tag, usually its attributes.

The function returns a pointer to the Tag structure or a nil pointer if there is no such tag in the provided string.

#### Check functions

Currently, in the tag module are available those Check functions:

- `func Has(attr string) Check` - it determines if the attribute of the given name exists in the tag,
- `func NotEmpty(attr string) Check` - it determines if the value of the attr attribute is not empty,
- `func Contains(attr, s string) Check` - it determines if the value of the attr attribute contains the s string,
- `func Equal(attr, s string) Check` - it determines if the value of the attr attribute is equal to the s string.

All attribute names are case-unsensitive.

You can use as many Check functions as you wish; a tag will be considered a result if all of them are satisfied.

You can write your own Check functions using closure, e.g.:
```go
// HasXClasses checks if tag has x classes
func HasXClasses(x int) Check {
	return func(t *Tag) bool {
		v, ok := t.Attr["class"]
		if !ok {
			return false
		}

		if len(strings.Split(v, " ")) == x {
			return true
		}

		return false
	}
}
```

#### *Tag structure

Find returns a pointer to a Tag structure, which has some exported methods:

- `func (t *Tag) Next() *Tag` - returns the next tag of the same name and satisfy the same Check functions, it is useful in loops,
- `func (t *Tag) Content() string` - returns a string that is between the opening and closing tags. If there is no closing tag or the tag is nil, it will return an empty string.

Tag structure also has some exported fields:

```
Name              string            // The name of the tag.
Attr              map[string]string // The map of attributes map[attr_name]attr_val. Attribute names are always lowercase.
ContentIndex      int               // The index points to the next character after the opening tag's closure in doc (it might be outside the doc range).
AfterClosureIndex int               // The index points to the next character after the closing tag's closure in doc (it might be outside the doc range).
```

### Example

```go
package main

import (
	"fmt"

	"github.com/kmatulewicz/go-html/tag"
)

const doc = `
<html>
	<body>
		<div id="interesting">
			<a href="https://example.com/1">Link 1</a>
			<a href="https://example.com/2">Link 2</a>
			<a href="https://example.com/3">Link 3</a>
			<a href="https://example.com/4">Link 4</a>
		</div>
		<div id="not">
			<a href="https://notinteresting.com/1">Not interesting 1</a>
		</div>
	</body>
</html>
`

func main() {
	// Find interesting content.
	interesting :=
		tag.Find(
			doc,   // the HTML document
			"div", // a name of the tag to be found
			[]tag.Check{ // a slice of Check functions to check if the tag is correct (all of them need to return true)
				tag.Equal("id", "interesting"), // it returns true if the tag has an id equal to interesting
			},
		).Content() // returns the content between the opening and closing tags

	// Loop over all a tags in the interesting content if they have a href attribute.
	a := tag.Find(interesting, "a", []tag.Check{tag.Has("href")})
	for ; a != nil; a = a.Next() {
		fmt.Println(a.Content(), "->", a.Attr["href"])
	}
}
```
Output:
```sh
Link 1 -> https://example.com/1
Link 2 -> https://example.com/2
Link 3 -> https://example.com/3
Link 4 -> https://example.com/4
```
