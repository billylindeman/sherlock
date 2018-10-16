//
// index.go
// billy lindeman <billy@lnd.mn>
//
// main interface into sherlock
//
package sherlock

import (
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

// Index search backed by a prefix map
type Index struct {
	numDocs uint64

	schema *Schema

	inverted inverted
	store    store
}

// QueryResult gather results
type QueryResult struct {
	Object interface{}
	Score  int

	docID        uint64
	postingsList []posting
}

//
func (i *Index) initWithSchema(schema Schema) {
	i.schema = &schema
	i.inverted = newRadixInvertedIndex()
	i.store = &memoryStore{}
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

	postingMap := make(map[string]*posting)
	// Process tokens into posting list
	for _, tok := range analysis.tokens {
		if _, ok := postingMap[tok.value]; !ok {
			postingMap[tok.value] = &posting{
				docID: i.numDocs,
				hits:  []hit{},
			}
		}

		list := postingMap[tok.value]
		h := hit{
			position: tok.position,
			fieldID:  tok.fieldID,
		}
		list.hits = append(list.hits, h)
	}

	// fmt.Printf("index built posting list: %#v", postingMap)
	// Merge into inverted index
	for term, post := range postingMap {
		// fmt.Printf("indexing %v posting list: %#v\n", term, post)
		i.inverted.insert(term, *post)
	}

	i.store.insert(i.numDocs, v)
	i.numDocs++
	return nil
}

// Query takes a string and prefix searches it
func (i *Index) Query(q string) ([]QueryResult, error) {
	norm := normalize(q)
	split := strings.Split(norm, " ")

	termSearchers := []searcher{}
	for _, t := range split {
		trim := strings.TrimSpace(t)
		if t != "\n" {
			searcher := &termSearcher{
				term: trim,
			}
			termSearchers = append(termSearchers, searcher)
		}
	}

	var plan searcher

	plan = &intersectionSearcher{
		searcher: &unionSearcher{
			searchers: termSearchers,
		},
	}
	// plan = &unionSearcher{
	// 	searchers: termSearchers,
	// }

	fmt.Printf("built query plan: ")
	spew.Dump(plan)

	matches := plan.search(i.inverted)
	// fmt.Printf("found matches: ", matches)
	// spew.Dump(matches)

	results := []QueryResult{}
	for _, m := range matches {
		if im, ok := m.(*intersectMatch); ok {
			doc, _ := i.store.get(im.docID)
			// fmt.Printf("[%x](%v) -> %#v\n", im.docID, len(im.postings()), doc)

			qr := QueryResult{
				Object: doc,
				docID:  im.docID,
			}

			results = append(results, qr)
		}
	}

	// fmt.Println("found matches: ", matches)

	// Grab posting lists from the inverted index
	// lists := []*postingList{}
	// for _, term := range terms {
	// 	val, _ := i.inverted.get(term)
	// 	lists = append(lists, val)
	// }

	// // Sort them from least common word to most common
	// sort.Slice(lists, func(i, j int) bool {
	// 	return lists[i].termFreq < lists[j].termFreq
	// })

	// perform a positional intersection of the postingLists
	// we walk each posting list in order (by docID) and along the way
	// note positional differences between our search terms
	// this will produce a set of phrase matches that we can use later to sort results

	// return []QueryResult{}, nil
	return results, nil
}
