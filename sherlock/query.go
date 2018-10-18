//
// query.go
// billy lindeman <billy@lnd.mn>
//
// this is a collection of functions that process posting lists
// after they've been retrieved for a given query
//
package sherlock

import (
	"fmt"
	"sort"
)

// the query pipeline should look something like this
//   query, load, search, score, collect -> matches

type match interface {
	term() string
	termCount() int

	postings() []posting
}

type searcher interface {
	search(index inverted) []match
}

// termSearcher loads the postingList for a single term
type termSearcher struct {
	term string
}

func (t *termSearcher) search(i inverted) []match {
	fmt.Println("[termSearcher] retrieving ", t.term)

	if val, err := i.get(t.term); err == nil {
		return []match{
			termMatch{
				p: val,
			},
		}
	}
	return []match{}
}

// prefixSearcher loads the posting lists for any terms that match the prefix
type prefixSearcher struct {
	prefix string
}

func (p *prefixSearcher) search(i inverted) []match {
	matches := []match{}

	fmt.Println("[prefixSearcher] retrieving prefix: ", p.prefix)
	if val, err := i.getByPrefix(p.prefix); err == nil {
		for _, pl := range val {
			matches = append(matches, termMatch{p: pl})
		}
	}
	fmt.Printf("prefixSearcher] found %v terms\n", len(matches))
	return matches
}

// termMatch represents a postingList hit in the inverted index (when executed by a termSearcher)
type termMatch struct {
	p *postingList
}

func (m termMatch) termCount() int {
	return m.p.termFreq
}

func (m termMatch) term() string {
	return m.p.term
}

func (m termMatch) postings() []posting {
	return m.p.postings
}

func (m termMatch) String() string {
	return fmt.Sprintf("%v(%v)", m.term(), m.termCount())
}

// unionSearcher collects the results from child searchers and merges them
type unionSearcher struct {
	searchers []searcher
}

func (u *unionSearcher) search(i inverted) []match {
	fmt.Println("[unionSearcher] running")
	matches := []match{}
	for _, s := range u.searchers {
		if s != nil {
			matches = append(matches, s.search(i)...)
		}
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].termCount() < matches[j].termCount()
	})
	return matches
}

// intersectionSearcher performs an intersection on the child termSearchers (and)
type intersectionSearcher struct {
	searcher searcher
}

type intersectMatch struct {
	docID uint64

	postingList []posting
}

func (m intersectMatch) termCount() int {
	// return m.p.termFreq
	return -1
}

func (m intersectMatch) term() string {
	return ""
}

func (m intersectMatch) postings() []posting {
	return m.postingList
}

func (m intersectMatch) String() string {
	return fmt.Sprintf("%v(%v)", m.docID, len(m.postingList))
}

// todo
func (s *intersectionSearcher) search(i inverted) []match {

	fmt.Println("[intersectionSearcher] running")
	matches := s.searcher.search(i)

	fmt.Println("[intersectionSearcher] processing termMatches: ", matches)
	// if we only have 1 term hit we return the result as is
	// (this operation basically becomes an identity op)
	if len(matches) < 1 {
		return []match{}
	}

	// seed the intersectionMatch with the first term
	intermediate := []*intersectMatch{}
	for _, p := range matches[0].postings() {
		m := &intersectMatch{
			docID:       p.docID,
			postingList: []posting{p},
		}
		intermediate = append(intermediate, m)
	}

	// loop through the remaining terms and intersect it with the intermediary result
	for i := 1; i < len(matches); i++ {
		merging := matches[i].postings()
		merged := s.matchIntersect(intermediate, merging)
		intermediate = merged
	}

	// cast back to match{}
	out := []match{}
	for _, m := range intermediate {
		out = append(out, m)
	}
	return out
}

// matchIntersect performs a two postingList set intersection in O(len(p1)+len(p2)) time
func (s intersectionSearcher) matchIntersect(curMatch []*intersectMatch, p2 []posting) []*intersectMatch {
	matches := []*intersectMatch{}

	p1 := []posting{}
	for _, pp := range curMatch {
		p1 = append(p1, pp.postingList...)
	}

	p1idx := 0
	p2idx := 0

	for p1idx < len(p1) && p2idx < len(p2) {
		if p1[p1idx].docID == p2[p2idx].docID {
			m := &intersectMatch{
				docID:       p1[p1idx].docID,
				postingList: []posting{p1[p1idx], p2[p2idx]},
			}
			matches = append(matches, m)

			p1idx++
			p2idx++
		} else if p1[p1idx].docID < p2[p2idx].docID {
			p1idx++
		} else {
			p2idx++
		}

	}

	return matches
}

// merges posting lists
// based on algorithm 2.12 from into to ir (manning)

// // phrase proximity calculation
// // (during positional intersection of posting lists)
// // based on algorithm 2.12 from into to ir (manning)
// const withinKWords = 1
// type phraseMatch struct {
// 	p1term string
// 	p2term string
// 	p1     hit
// 	p2     hit
// }
// answers := make(map[uint64][]phraseMatch)

// var p1, p2 *posting
// if len(merged) > 1 {
// 	p1 = &merged[0]
// 	p2 = &merged[1]
// 	idx := 1

// 	for p1 != nil && p2 != nil {
// 		// terms in the same document
// 		if p1.docID == p2.docID {
// 			fmt.Printf("p1: %#v -- p2: %#v\n", p1.term, p2.term)

// 			l := []hit{}
// 			for _, pp1 := range p1.positions {
// 				fmt.Printf("pp1: %#v\n", pp1)
// 				for _, pp2 := range p2.positions {
// 					if abs(pp1.position-pp2.position) <= withinKWords {
// 						l = append(l, pp2)
// 					} else if pp2.position > pp1.position {
// 						break
// 					}
// 				}

// 				for len(l) > 0 && abs(l[0].position-pp1.position) > withinKWords {
// 					fmt.Println("purgging step?")
// 					l = append(l[:0], l[1:]...) // remove item 0 from slice
// 				}

// 				// fmt.Println(l)
// 				for _, h := range l {
// 					// make sure match object is in order (makes phrase evalution easier)
// 					if pp1.position < h.position {
// 						m := phraseMatch{
// 							p1term: p1.term,
// 							p2term: p2.term,
// 							p1:     pp1,
// 							p2:     h,
// 						}
// 						answers[p1.docID] = append(answers[p1.docID], m)
// 					} else {
// 						m := phraseMatch{
// 							p1term: p2.term,
// 							p2term: p1.term,
// 							p1:     h,
// 							p2:     pp1,
// 						}
// 						answers[p1.docID] = append(answers[p1.docID], m)

// 					}
// 				}
// 			}

// 			idx++
// 			if idx < len(merged) {
// 				p1 = p2
// 				p2 = &merged[idx]
// 				continue
// 			} else {
// 				break
// 			}
// 		} else if p1.docID < p2.docID {
// 			idx++
// 			p1 = nil
// 			if idx < len(merged) {
// 				p1 = &merged[idx]
// 				continue
// 			}
// 		} else {
// 			idx++
// 			p2 = nil
// 			if idx < len(merged) {
// 				p2 = &merged[idx]
// 			}
// 		}
// 	}
