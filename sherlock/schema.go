package sherlock

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	// ErrTypeNotStruct returned when someone tries to index a non-struct type
	ErrTypeNotStruct = errors.New("Indexer requires a Struct type")
)

const tagName = "sherlock"

// Schema represents the indexing criteria for a given object
type Schema struct {
	Rules []SchemaRule
}

// SchemaRule contains rule information for each individual field being indexed on a given object
type SchemaRule struct {
	Type reflect.Type
	Omit bool
}

// NewSchemaFromStruct builds a document schema by reflecting over the passed in struct
func NewSchemaFromStruct(v interface{}) (*Schema, error) {
	s := reflect.ValueOf(v)

	for i := 0; i < s.NumField(); i++ {
		// Get the field tag value
		tag := s.Type().Field(i).Tag.Get(tagName)

		// Skip if tag is not defined or ignored
		if tag == "" || tag == "-" {
			continue
		}

		fmt.Println("Found Tag: ", tag)
	}
	return nil, nil
}
