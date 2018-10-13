
# Sherlock

Sherlock is a simple and embeddable search index for go.  It uses a prefix tree based inverted index so it can serve as an autocomplete search for full documents.


# Usage

Sherlock is designed to be simple to use and embeddable in any go application. 

```go
type Document struct {
	Title    string `sherlock:"weight=10"`
	Body     string `sherlock:"weight=5"`
	Tags []string `sherlock:"weight=10,facet=term"` // We want facets for these results
}

func main() {
	s := sherlock.Index{}

	corpus := []Document{
		Document{
			Title: "Example",
			Body:  "This is an example",
		},
		Document{
			Title: "Real world",
			Body:  "This is the real world",
		},
	}

	for _, doc := range corpus {
		s.Index(doc)
	}

}
```
