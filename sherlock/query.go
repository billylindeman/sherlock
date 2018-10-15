//
// query.go
// billy lindeman <billy@lnd.mn>
//
// this is a collection of functions that process posting lists
// after they've been retrieved for a given query
//
package sherlock

type termMatch struct {
	term     string
	position int
}

type queryResult struct {
	termMatches termMatch
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
