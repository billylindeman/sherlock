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

type posting struct {
	docID     uint64
	positions []hit
}

func (p *posting) score() int {
	return len(p.positions)
}

type QueryResult struct {
	Object interface{}
	Score  int
}

// hit represents a hit of a term in a document
type hit struct {
	position int
	weight   int
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
		i.prefixMap.Insert(term, postingList)
	}

	i.documents[i.numDocs] = v

	i.numDocs++

	return nil
}

// Query takes a string and prefix searches it
func (i *Index) Query(q string) ([]QueryResult, error) {
	norm := strings.ToLower(q)
	terms := strings.Split(norm, " ")

	dedupe := map[string]bool{}
	fmt.Printf("Got query terms: %v", terms)
	for _, term := range terms {
		t := strings.TrimSpace(term)
		dedupe[t] = true
	}

	merged := []posting{}
	j := 0
	for term := range dedupe {
		// grab postings for q
		var postings []interface{}
		if j == len(terms)-1 {
			// if we're on the last term lets do a prefix search
			postings = i.prefixMap.GetByPrefix(term)
		} else {
			postings = i.prefixMap.Get(term)
		}

		// sort postings based on scoring
		for _, p := range postings {
			p := p.(*posting)
			merged = append(merged, *p)
		}

		j++
	}

	// phrase proximity calculation
	// based on algorithm 2.12 from into to ir (manning)
	const withinKWords = 2

	type phraseMatch struct {
		p1 hit
		p2 hit
	}
	answers := make(map[uint64][]phraseMatch)

	if len(terms) > 1 {
		for i := 0; i < len(merged)-1; i++ {
			for j := 1; j < len(merged); j++ {
				p1 := merged[i]
				p2 := merged[j]

				// terms in the same document
				if p1.docID == p2.docID {
					l := []hit{}
					for _, pp1 := range p1.positions {
						for _, pp2 := range p2.positions {
							if abs(pp1.position-pp2.position) <= withinKWords {
								l = append(l, pp2)
							} else if pp2.position > pp1.position {
								break
							}
						}

						for len(l) > 0 && abs(l[0].position-pp1.position) > withinKWords {
							l = append(l[:0], l[1:]...) // remove item 0 from slice
						}

						for _, h := range l {
							m := phraseMatch{
								p1: pp1,
								p2: h,
							}
							answers[p1.docID] = append(answers[p1.docID], m)
						}
					}
				}
			}
		}

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

		for _, p := range postingList {
			r.Score += (2 * p.score())
		}

		if len(answers[docID]) > 0 {
			totalScore := 100
			// fmt.Printf("doc phrase hit, score:  %#v", answers)
			for _, pm := range answers[docID] {
				totalScore += abs(pm.p2.position - pm.p1.position)
			}
			r.Score = totalScore
		}

		results = append(results, r)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score < results[j].Score
	})

	// return []QueryResult{}, nil
	return results, nil
}

func abs(n int) int {
	y := n >> 63
	return (n ^ y) - y
}
