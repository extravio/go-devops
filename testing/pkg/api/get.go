package api

import (
	"encoding/json"
	"fmt"
	"io"
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
	// Page  string   `json:"page"`
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

// declared for testing (to mock the response from the API)
type WordsPage struct {
	Page
	Words
}

func (a API) DoGetRequest(requestURL string) (Response, error) {

	response, err := a.Client.Get(requestURL)
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

	if !json.Valid(body) {
		return nil, RequestError{
			HTTPCode: response.StatusCode,
			Body:     string(body),
			Err:      fmt.Sprintf("No valid JSON returned"),
		}
	}

	var page Page
	err = json.Unmarshal(body, &page)
	if err != nil {
		return nil, RequestError{
			HTTPCode: response.StatusCode,
			Body:     string(body),
			Err:      fmt.Sprintf("Page Unmarshal error: %s", err),
		}
	}

	switch page.Name {
	case "words":
		var words Words
		err = json.Unmarshal(body, &words)
		if err != nil {
			return nil, RequestError{
				HTTPCode: response.StatusCode,
				Body:     string(body),
				Err:      fmt.Sprintf("Words Unmarshal error: %s", err),
			}
		}
		return words, nil

	case "occurrence":
		var occurrence Occurrence
		err = json.Unmarshal(body, &occurrence)
		if err != nil {
			return nil, RequestError{
				HTTPCode: response.StatusCode,
				Body:     string(body),
				Err:      fmt.Sprintf("Occurrence Unmarshal error: %s", err),
			}
		}
		return occurrence, nil
	}

	return nil, nil
}
