package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/billylindeman/sherlock/sherlock"
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

	for _, doc := range corpus {
		s.Index(doc)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to sherlock!")
	for {
		fmt.Print("search:> ")
		text, _ := reader.ReadString('\n')
		s.Query(text)
	}
}
