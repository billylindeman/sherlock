package sherlock

import "github.com/alediaferia/prefixmap"

// Radix represents a radix tree
// the radix tree is used to store our posting lists
// This enables prefix search on real documents
type Radix interface {
}

type PostingList struct {
}

type Posting struct {
}

// PrefixMapRadix  Proxying calls to a prefixMap package for testing
type PrefixMapRadix struct {
	m *prefixmap.PrefixMap
}

func NewPrefixMapRadix() Radix {
	return &PrefixMapRadix{
		m: prefixmap.New(),
	}
}
