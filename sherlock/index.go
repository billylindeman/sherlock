package sherlock

import (
	"fmt"
	"strings"

	"github.com/alediaferia/prefixmap"
)

// Index search backed by a prefix map
type Index struct {
	numDocs   uint64
	prefixMap *prefixmap.PrefixMap
	schema    *Schema

	documents map[uint64]interface{}
}

type posting struct {
	docID     uint64
	positions []hit
}

// func (p *posting) score() int {
// 	score := 0

// 	// for _, hit := range p.positions {

// 	// }
// }

// hit represents a hit of a term in a document
type hit struct {
	position int
	weight   int
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

	analysis, err := i.schema.analyze(v)
	if err != nil {
		return err
	}

	postings := make(map[string]*posting)

	// Process tokens into posting list
	for _, tok := range analysis.tokens {
		if _, ok := postings[tok.value]; !ok {
			postings[tok.value] = &posting{
				docID: i.numDocs,
			}
		}
		list := postings[tok.value]
		h := hit{
			position: tok.position,
			weight:   tok.weight,
		}
		list.positions = append(list.positions, h)
	}

	fmt.Printf("index built posting list: %#v", postings)
	// Merge into inverted index
	for term, postingList := range postings {
		i.prefixMap.Insert(term, postingList)
	}

	i.numDocs++
	return nil
}

// Query takes a string and prefix searches it
func (i *Index) Query(q string) ([]interface{}, error) {
	norm := strings.ToLower(q)
	terms := strings.Split(norm, " ")

	merged := []posting{}
	for _, term := range terms {
		// grab postings for q
		postings := i.prefixMap.GetByPrefix(term)
		// sort postings based on scoring
		for _, p := range postings {
			p := p.(*posting)
			merged = append(merged, *p)

			fmt.Printf("hit: %#v\n", p)
		}
	}

	//	return results

	return nil, nil
}
