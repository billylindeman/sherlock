//
// store.go
// billy lindeman <billy@lnd.mn>
//
// storage for the full documents that have been indexed
// currently just memory interfaces
// in the futre we'll use some sort of disk-backed k/v such as badger
//

package sherlock

import "errors"

var (
	errDocumentNotFound = errors.New("document not found in store")
)

// store key/value store to put the documents that have been indexed
type store interface {
	insert(uint64, interface{}) error
	get(uint64) (interface{}, error)
}

// memoryStore is an inmemory DocStore backed by a basic map
type memoryStore struct {
	documents map[uint64]interface{}
}

// Insert puts a doc into the map
func (s *memoryStore) insert(docID uint64, v interface{}) error {
	if s.documents == nil {
		s.documents = make(map[uint64]interface{})
	}

	s.documents[docID] = v
	return nil
}

// Get retrieves a doc from the map
func (s *memoryStore) get(docID uint64) (interface{}, error) {
	if s.documents == nil {
		return nil, errDocumentNotFound
	}
	if doc, ok := s.documents[docID]; ok {
		return doc, nil
	}
	return nil, errDocumentNotFound
}
