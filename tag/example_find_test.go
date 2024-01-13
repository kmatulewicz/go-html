package tag_test

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

func ExampleFind() {
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

	// Output:
	// Link 1 -> https://example.com/1
	// Link 2 -> https://example.com/2
	// Link 3 -> https://example.com/3
	// Link 4 -> https://example.com/4
}
