package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Page struct {
	Name string `json:"page"`
}

// {"page":"words","input":"word1","words":["word1"]}
type Words struct {
	Page  string   `json:"page"`
	Input string   `json:"input"`
	Words []string `json:"words"`
}

// {"page":"occurrence","words":{"word1":1,"word2":2}
type Occurrence struct {
	Page  string         `json:"page"`
	Words map[string]int `json:"words"`
}

func main() {
	// var args []string
	args := os.Args

	if len(args) < 2 {
		fmt.Printf("Usage: ./http-get <url>\n")
		os.Exit(1)
	}

	if _, err := url.ParseRequestURI(args[1]); err != nil {
		fmt.Printf("URL is in invalid format: %s\n", err)
		os.Exit(1)
	}

	response, err := http.Get(args[1])
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	if response.StatusCode != 200 {
		fmt.Printf("Invalid output (HTTP Status Code: %d): %s\n", response.StatusCode, body)
		os.Exit(1)
	}

	var page Page
	err = json.Unmarshal(body, &page)
	if err != nil {
		log.Fatal(err)
	}

	switch page.Name {
	case "words":
		var words Words
		err = json.Unmarshal(body, &words)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("JSON parsed\nPage: %s\nWords: %s\n", page.Name, strings.Join(words.Words, ", "))
	case "occurrence":
		var occurence Occurrence
		err = json.Unmarshal(body, &occurence)
		if err != nil {
			log.Fatal(err)
		}

		// check if an element exists
		if val, ok := occurence.Words["word1"]; ok {
			fmt.Printf("Found word1: %d\n", val)
		}

		for word, occurrence := range occurence.Words {
			fmt.Printf("%s: %d\n", word, occurrence)
		}
	default:
		fmt.Printf("Page not found\n")
	}

}
