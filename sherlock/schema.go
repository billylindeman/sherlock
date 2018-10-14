package sherlock

import (
	"errors"
	"reflect"
	"strings"
)

var (
	// ErrTypeNotStruct returned when someone tries to index a non-struct type
	ErrTypeNotStruct = errors.New("Indexer requires a Struct type")
)

const tagName = "sherlock"

// Schema represents the indexing criteria for a given object
type Schema struct {
	Fields []FieldRule
}

func (s *Schema) analyze(v interface{}) (analysis, error) {
	a := analysis{
		tokens: []token{},
	}

	doc := reflect.ValueOf(v)
	for _, f := range s.Fields {
		// log.Println("analyzing field: ", f.fieldName)
		text := doc.FieldByName(f.fieldName)
		// fmt.Println("found value: ", text.String())
		norm := normalize(text.String())
		pos := 0
		for i, word := range strings.Split(norm, " ") {
			tok := token{
				value:    strings.ToLower(word),
				field:    f.fieldName,
				position: i + 1,
				weight:   f.weight,
			}
			a.tokens = append(a.tokens, tok)
			pos += len(word)
		}
	}

	return a, nil

}

// FieldRule contains rule information for each individual field being indexed on a given object
type FieldRule struct {
	Omit bool

	weight    int
	field     reflect.StructField
	tag       string
	fieldName string
}

// NewSchemaFromStruct builds a document schema by reflecting over the passed in struct
func NewSchemaFromStruct(v interface{}) (*Schema, error) {
	out := Schema{}
	s := reflect.ValueOf(v)

	for i := 0; i < s.NumField(); i++ {
		// Get the field & tag value
		f := s.Type().Field(i)
		tag := f.Tag.Get(tagName)

		// for _, t := range strings.Split(tag, ",") {
		// 	// log.Println("Found option: ", t)
		// }

		r := FieldRule{
			fieldName: f.Name,
			tag:       tag,
		}

		// Skip if tag is not defined or ignored
		if tag == "" || tag == "-" {
			r.Omit = true
		}

		// fmt.Printf("Found Tag: %v \n build schema: %+v ", tag, r)

		out.Fields = append(out.Fields, r)
	}

	return &out, nil
}
