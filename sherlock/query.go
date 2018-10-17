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

	if val, err := i.getByPrefix(p.prefix); err == nil {
		for _, pl := range val {
			matches = append(matches, termMatch{p: val})
		}
	}
	return []match{}
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
	matches := []match{}
	for _, s := range u.searchers {
		matches = append(matches, s.search(i)...)
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

// first gen result sorting
//
// results := []QueryResult{}
// for docID, postingList := range grouped {
// 	r := QueryResult{
// 		Object: i.documents[docID],
// 	}

// 	r.Score += len(postingList) ^ 2

// 	if len(answers[docID]) > 0 {
// 		matches := answers[docID]

// 		sort.Slice(matches, func(i, j int) bool {
// 			return matches[i].p1.position < matches[j].p1.position
// 		})

// 		fmt.Printf("matches: %#v\n", matches)

// 		termIdx := 0
// 		matchIdx := 0

// 		curScore := 0
// 		pendingHit := false

// 		for termIdx < len(terms) && matchIdx < len(matches) {
// 			fmt.Printf("loop term: %v match %v \n", termIdx, matchIdx)
// 			distance := abs(matches[matchIdx].p2.position - matches[matchIdx].p1.position)

// 			if matches[matchIdx].p1term == terms[termIdx] && distance == 1 {
// 				pendingHit = true

// 				termIdx++
// 				fmt.Printf("loop term: %v match %v \n", termIdx, matchIdx)
// 				if termIdx == len(terms) {
// 					break
// 				}
// 			} else {
// 				matchIdx++
// 				continue
// 			}

// 			if matches[matchIdx].p2term == terms[termIdx] && pendingHit {
// 				// phrase bigram hit
// 				fmt.Printf("bigram hit: %v->%v \n", matches[matchIdx].p1term, matches[matchIdx].p2term)
// 				curScore += 50 / distance

// 				matchIdx++
// 			} else {
// 				termIdx++
// 			}

// 			pendingHit = false
// 		}

// 		r.Score += (curScore)
// 		fmt.Printf("phraseScore: %v, rScore: %v\n", curScore, r.Score)
// 		// r.Score = totalScore
// 	}

// 	results = append(results, r)
// }
