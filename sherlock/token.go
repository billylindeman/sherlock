package sherlock

// token holds information about a Token found in a document
type token struct {
	value    string
	position int
	field    string
	weight   int
}

type analysis struct {
	tokens []token
}
