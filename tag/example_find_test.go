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
			<a href="https://notinteresting.com/1">Not interesting 1</a>
			<a href="https://notinteresting.com/1">Not interesting 1</a>
			<a href="https://notinteresting.com/1">Not interesting 1</a>
		</div>
	</body>
</html>
`

func ExampleFind() {
	// Find interesting content.
	interesting :=
		tag.Find(
			doc,
			"div",
			[]tag.Check{
				tag.Equal("id", "interesting"),
			}).Content()

	position := 0
	// Loop over all a tags in the interesting content.
	for position < len(interesting)-1 {
		a := tag.Find(interesting[position:], "a", []tag.Check{tag.Has("href")})
		if a == nil {
			break
		}

		fmt.Println(a.Content(), "->", a.Attr["href"])

		position += a.AfterClosureIndex
	}

	// Output:
	// Link 1 -> https://example.com/1
	// Link 2 -> https://example.com/2
	// Link 3 -> https://example.com/3
	// Link 4 -> https://example.com/4
}