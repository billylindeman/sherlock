package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/billylindeman/sherlock/sherlock"
	"github.com/sbwhitecap/tqdm"
)

type Document struct {
	Title string `json:"play_name" sherlock:"weight=10"`
	Body  string `json:"text_entry" sherlock:"weight=5"`

	Line string `json:"line_number"`
}

func main() {
	s := sherlock.Index{}

	corpus := []Document{}

	b, _ := ioutil.ReadFile("./shakespeare.json")
	for _, j := range strings.Split(string(b), "\n") {
		d := Document{}
		json.Unmarshal([]byte(j), &d)
		corpus = append(corpus, d)
	}

	fmt.Printf("\nLoaded %v docs\n", len(corpus))
	tqdm.R(0, len(corpus), func(v interface{}) (brk bool) {
		idx := v.(int)
		doc := corpus[idx]
		s.Index(doc)
		return
	})

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\nWelcome to sherlock!")
	for {
		fmt.Print("search:> ")
		text, _ := reader.ReadString('\n')

		t1 := time.Now()
		results, _ := s.Query(text)
		t2 := time.Now()

		if len(results) > 0 {
			fmt.Printf("Found %v results in %v\n", len(results), t2.Sub(t1))
			fmt.Println("Top result: ", results[0])
		}

	}
}
