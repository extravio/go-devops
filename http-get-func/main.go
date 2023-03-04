package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Response interface {
	GetResponse() string
}

type Page struct {
	Name string `json:"page"`
}

// {"page":"words","input":"word1","words":["word1"]}
type Words struct {
	Page  string   `json:"page"`
	Input string   `json:"input"`
	Words []string `json:"words"`
}

func (w Words) GetResponse() string {
	return fmt.Sprintf("%s", strings.Join(w.Words, ", "))
}

// {"page":"occurrence","words":{"word1":1,"word2":2}
type Occurrence struct {
	Page  string         `json:"page"`
	Words map[string]int `json:"words"`
}

func (o Occurrence) GetResponse() string {
	// check if an element exists
	// if val, ok := o.Words["word1"]; ok {
	// 	fmt.Printf("Found word1: %d\n", val)
	// }
	out := []string{}
	for word, occurrence := range o.Words {
		out = append(out, fmt.Sprintf("%s: %d", word, occurrence))
	}
	return fmt.Sprintf("%s", strings.Join(out, ", "))
}

func main() {
	// var args []string
	args := os.Args

	if len(args) < 2 {
		fmt.Printf("Usage: ./http-get <url>\n")
		os.Exit(1)
	}

	res, err := doRequest(args[1])
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	if res == nil {
		fmt.Printf("No Reponse\n")
		os.Exit(1)
	}

	fmt.Printf("Response: %s\n", res.GetResponse())

}

func doRequest(requestURL string) (Response, error) {

	if _, err := url.ParseRequestURI(requestURL); err != nil {
		return nil, fmt.Errorf("Validation error: URL is not valid: %s", err)
	}

	response, err := http.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("HTTP Get error: %s", err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, fmt.Errorf("ReadAll error: %s", err)
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Invalid output (HTTP Status Code: %d): %s", response.StatusCode, body)
	}

	var page Page
	err = json.Unmarshal(body, &page)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal error: %s", err)
	}

	switch page.Name {
	case "words":
		var words Words
		err = json.Unmarshal(body, &words)
		if err != nil {
			return nil, fmt.Errorf("Unmarshal error: %s", err)
		}
		return words, nil

	case "occurrence":
		var occurrence Occurrence
		err = json.Unmarshal(body, &occurrence)
		if err != nil {
			return nil, fmt.Errorf("Unmarshal error: %s", err)
		}
		return occurrence, nil
	}

	return nil, nil
}
