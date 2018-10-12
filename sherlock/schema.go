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
	Omit bool

	tag       string
	fieldName string
}

// NewSchemaFromStruct builds a document schema by reflecting over the passed in struct
func NewSchemaFromStruct(v interface{}) (*Schema, error) {
	s := reflect.ValueOf(v)

	for i := 0; i < s.NumField(); i++ {
		// Get the field & tag value
		f := s.Type().Field(i)
		tag := f.Tag.Get(tagName)

		r := SchemaRule{
			fieldName: f.Name,
			tag:       tag,
		}

		// Skip if tag is not defined or ignored
		if tag == "" || tag == "-" {
			r.Omit = true
		}

		fmt.Printf("Found Tag: %v \n build schema: %+v ", tag, r)
	}
	return nil, nil
}
