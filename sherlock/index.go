//
// index.go
// billy lindeman <billy@lnd.mn>
//
// main interface into sherlock
//
package sherlock

import (
	"fmt"
	"sort"
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
	i.prefixMap = prefixmap.New()
	i.documents = make(map[uint64]interface{})
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

	// fmt.Printf("index built posting list: %#v", postings)
	// Merge into inverted index
	for term, postingList := range postings {
		fmt.Printf("indexing %v posting list: %#v\n", term, postingList)
		i.prefixMap.Insert(term, postingList)
	}

	i.documents[i.numDocs] = v

	i.numDocs++

	return nil
}

// Query takes a string and prefix searches it
func (i *Index) Query(q string) ([]QueryResult, error) {
	norm := normalize(q)
	terms := []string{}
	for _, t := range strings.Split(norm, " ") {
		terms = append(terms, strings.TrimSpace(t))
	}

	fmt.Printf("Got query terms: %v\n", terms)

	dedupe := map[string]bool{}

	for _, term := range terms {
		dedupe[term] = true
	}

	byTerm := make(map[string][]posting)
	j := 0
	for term := range dedupe {
		// grab postings for q
		var postings []interface{}
		postings = i.prefixMap.Get(term)

		for _, p := range postings {
			p := p.(*posting)
			byTerm[term] = append(byTerm[term], *p)
		}

		j++
	}

	// merge posting lists
	intermediate := []postings{}
	for term, list := range byTerm {

	}

	// phrase proximity calculation
	// (during positional intersection of posting lists)
	// based on algorithm 2.12 from into to ir (manning)
	const withinKWords = 1
	type phraseMatch struct {
		p1term string
		p2term string
		p1     hit
		p2     hit
	}
	answers := make(map[uint64][]phraseMatch)

		// fmt.Printf("Found phrase distances: %#v\n", answers)
	}

	grouped := make(map[uint64][]posting)
	//	return results
	for _, p := range merged {
		if _, ok := grouped[p.docID]; !ok {
			grouped[p.docID] = []posting{p}
			continue
		}
		grouped[p.docID] = append(grouped[p.docID], p)
	}

	results := []QueryResult{}
	for docID, postingList := range grouped {
		r := QueryResult{
			Object: i.documents[docID],
		}

		r.Score += len(postingList) ^ 2

		if len(answers[docID]) > 0 {
			matches := answers[docID]

			sort.Slice(matches, func(i, j int) bool {
				return matches[i].p1.position < matches[j].p1.position
			})

			fmt.Printf("matches: %#v\n", matches)

			termIdx := 0
			matchIdx := 0

			curScore := 0
			pendingHit := false

			for termIdx < len(terms) && matchIdx < len(matches) {
				fmt.Printf("loop term: %v match %v \n", termIdx, matchIdx)
				distance := abs(matches[matchIdx].p2.position - matches[matchIdx].p1.position)

				if matches[matchIdx].p1term == terms[termIdx] && distance == 1 {
					pendingHit = true

					termIdx++
					fmt.Printf("loop term: %v match %v \n", termIdx, matchIdx)
					if termIdx == len(terms) {
						break
					}
				} else {
					matchIdx++
					continue
				}

				if matches[matchIdx].p2term == terms[termIdx] && pendingHit {
					// phrase bigram hit
					fmt.Printf("bigram hit: %v->%v \n", matches[matchIdx].p1term, matches[matchIdx].p2term)
					curScore += 50 / distance

					matchIdx++
				} else {
					termIdx++
				}

				pendingHit = false
			}

			r.Score += (curScore)
			fmt.Printf("phraseScore: %v, rScore: %v\n", curScore, r.Score)
			// r.Score = totalScore
		}

		results = append(results, r)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	// return []QueryResult{}, nil
	return results, nil
}
