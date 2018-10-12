
# Sherlock

Sherlock is a simple and embeddable search index for go.  It uses a prefix tree based inverted index so it can serve as an autocomplete search for full documents.


# Usage

Sherlock is designed to be simple to use and embeddable in any go application. 

```go
type Document struct {
	Title    string `sherlock:"term,boost=2"`
	Body     string `sherlock:"term"`
}

func main() {
    s := sherlock.NewIndex()

    corpus := []Document{
        Document {
        Title: "Example",
        Body: "This is an example",
        },
        Document {
        Title: "Real world",
        Body: "This is the real world",
    },
    }


    s.Index()

}
