package sherlock

import "errors"

var (
	ErrDocumentNotFound = errors.New("document not found in store")
)

// DocStore key/value store to put the documents that have been indexed
type DocStore interface {
	Insert(uint64, interface{}) error
	Get(uint64) (interface{}, error)
}

// MemoryStore is an inmemory DocStore backed by a basic map
type MemoryStore struct {
	documents map[uint64]interface{}
}

// Insert puts a doc into the map
func (s *MemoryStore) Insert(docID uint64, v interface{}) error {
	if s.documents == nil {
		s.documents = make(map[uint64]interface{})
	}

	s.documents[docID] = v
	return nil
}

// Get retrieves a doc from the map
func (s *MemoryStore) Get(docID uint64) (interface{}, error) {
	if s.documents == nil {
		return nil, ErrDocumentNotFound
	}
	if doc, ok := s.documents[docID]; ok {
		return doc, nil
	}
	return nil, ErrDocumentNotFound
}