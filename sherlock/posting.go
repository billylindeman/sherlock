//
// posting.go
// billy lindeman <billy@lnd.mn>
//
// data structures for the search index
//
package sherlock

// postingList is the root structure for every term in the radix tree
type postingList struct {
	term     string
	termFreq int
	postings []posting
}

func (l *postingList) insert(posting posting) {
	l.postings = append(l.postings, posting)
	l.termFreq += len(posting.hits)
}

// posting represents all occurences of a term within a document
type posting struct {
	docID uint64

	hits []hit
}

// hit represents an occurence of a term in a document
type hit struct {
	position int
	fieldID  int
}
