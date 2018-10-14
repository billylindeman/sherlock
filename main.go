package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/billylindeman/sherlock/sherlock"
)

type Document struct {
	// Title string `json:"play_name" sherlock:"weight=10"`
	Body string `json:"text_entry" sherlock:"weight=5"`
	Line string `json:"line_number"`
}

var corpus = []Document{
	Document{
		Body: "hello there son",
	},
	Document{
		Body: "different entirely",
	},
	Document{
		Body: "super different entirely",
	},
	Document{
		Body: "it aint be like it is, but it do",
	},
}

func main() {
	s := sherlock.Index{}

	// corpus := []Document{}
	// b, _ := ioutil.ReadFile("./shakespeare.json")
	// for _, j := range strings.Split(string(b), "\n") {
	// 	d := Document{}
	// 	json.Unmarshal([]byte(j), &d)
	// 	corpus = append(corpus, d)
	// }

	fmt.Printf("\nLoaded %v docs\n", len(corpus))
	fmt.Print("Indexing.")
	t1 := time.Now()
	for i, doc := range corpus {
		s.Index(doc)

		if i%5000 == 0 {
			fmt.Print(".")
		}
	}
	t2 := time.Now()
	fmt.Println("Done")
	fmt.Printf("Indexed %v docs in %v\n", len(corpus), t2.Sub(t1))

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\nWelcome to sherlock!")
	for {
		fmt.Print("search:> ")
		text, _ := reader.ReadString('\n')

		t1 := time.Now()
		results, _ := s.Query(text)
		t2 := time.Now()

		fmt.Printf("Found %v results in %v\n", len(results), t2.Sub(t1))

		for i := 0; i < 10; i++ {
			if i < len(results) {
				fmt.Printf("\t[%v](score:%v) %#v\n", i, results[i].Score, results[i].Object)
			}
		}
	}
}
