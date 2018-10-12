package sherlock

import (
	"github.com/alediaferia/prefixmap"
)

// Object is the basic interface for an object that you'd like to index
type Object interface {
	// ID returns a unique id for this object
	ID() string
}

// Index search backed by a prefix map
type Index struct {
	prefixMap *prefixmap.PrefixMap

	schema *Schema
}

//
func (i *Index) initWithSchema(schema Schema) {
	i.schema = &schema
	i.prefixMap = prefixmap.New()
}

// Index takes in a struct, processes it's struct tags, and indexes it's terms
func (i *Index) Index(v interface{}) error {
	if i.schema == nil {
		schema, err := NewSchemaFromStruct(v)
		if err != nil {
			return err
		}
		i.initWithSchema(*schema)
	}

	return nil
}

// Query takes a string and prefix searches it
func (i *Index) Query(q string) error {
	return nil
}
