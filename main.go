package main

type Document struct {
	Title string `sherlock:"term,boost=2"`
	Body  string `sherlock:"term"`
}

func main() {
	s := sherlock.NewIndex()

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
