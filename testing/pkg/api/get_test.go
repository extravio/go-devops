package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

type MockClient struct {
	ResponseOutput *http.Response
}

func (m MockClient) Get(url string) (resp *http.Response, err error) {
	return m.ResponseOutput, nil
}

func TestDoGetRequest(t *testing.T) {
	words := WordsPage{
		Page: Page{"words"},
		Words: Words{
			Input: "abc",
			Words: []string{"a", "b"},
		},
	}

	wordsBytes, err := json.Marshal(words)
	if err != nil {
		t.Errorf("Marshal error: %s", err)
	}

	apiInstance := API{
		Options: Options{},
		Client: MockClient{
			ResponseOutput: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(wordsBytes)),
			},
		},
	}
	response, err := apiInstance.DoGetRequest("http://localhost/words")
	if err != nil {
		t.Errorf("DoGetRequest error: %s", err)
	}
	if response == nil {
		t.Fatalf("reponse is empty")
	}
	if response.GetResponse() != strings.Join([]string{"a", "b"}, ", ") {
		t.Errorf("Unexpected reponse: %s", response.GetResponse())
	}

}
