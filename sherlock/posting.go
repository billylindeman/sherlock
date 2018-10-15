//
// posting.go
// billy lindeman <billy@lnd.mn>
//
// data structures for the search index
//
package sherlock

// postingList is the root structure for every term in the radix tree
type postingList struct {
	term     string
	termFreq int
	postings []posting
}

func (l *postingList) insert(posting posting) {
	l.postings = append(l.postings, posting)
	l.termFreq += len(posting.hits)
}

// posting represents all occurences of a term within a document
type posting struct {
	docID uint64

	hits []hit
}

// hit represents an occurence of a term in a document
type hit struct {
	position int
	fieldID  int
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
