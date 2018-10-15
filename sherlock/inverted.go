package sherlock

import (
	"errors"

	"github.com/armon/go-radix"
)

var (
	errTermMissing = errors.New("term or prefix not found in index")
)

// inverted index holds and retrieves postingLists
type inverted interface {
	insert(string, posting) error
	get(string) (*postingList, error)
	getByPrefix(string) ([]*postingList, error)
}

// radixIndex holds postingLists in a radix-tree for efficient
// term & prefix lookups
type radixInvertedIndex struct {
	tree *radix.Tree
}

func newRadixInvertedIndex() inverted {
	return &radixInvertedIndex{
		tree: radix.New(),
	}
}

func (r *radixInvertedIndex) insert(key string, p posting) error {
	if list, ok := r.tree.Get(key); ok {
		postingList := list.(*postingList)
		postingList.insert(p)
		return nil
	}

	postingList := postingList{
		term:     key,
		postings: []posting{p},
	}
	r.tree.Insert(key, &postingList)
	return nil
}

func (r *radixInvertedIndex) get(key string) (*postingList, error) {
	if v, ok := r.tree.Get(key); ok {
		pl := v.(*postingList)
		return pl, nil
	}

	return nil, errTermMissing
}

func (r *radixInvertedIndex) getByPrefix(key string) ([]*postingList, error) {

	collect := []*postingList{}
	r.tree.WalkPrefix(key, func(s string, v interface{}) bool {
		collect = append(collect, v.(*postingList))
		return false
	})

	if len(collect) > 0 {
		return collect, nil
	}

	return nil, errTermMissing
}
