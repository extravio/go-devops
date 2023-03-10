package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

type MockRoundTripper struct {
	RoundTripperOutput *http.Response
}

func (m MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Header.Get("Authorization") != "Bearer abc" {
		return nil, fmt.Errorf("Wrong authorizartion header: %s", req.Header.Get("Authorization"))
	}
	return m.RoundTripperOutput, nil
}

func TestRoundTrip(t *testing.T) {
	loginResponse := LoginResponse{
		Token: "abc",
	}
	loginResponseBytes, err := json.Marshal(loginResponse)
	if err != nil {
		t.Errorf("Marshal error: %s", err)
	}
	myJWTTransport := MyJWTTransport{
		transport: MockRoundTripper{
			RoundTripperOutput: &http.Response{
				StatusCode: 200,
			},
		},
		HTTPClient: MockClient{
			PostResponseOutput: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(loginResponseBytes)),
			},
		},
		password: "xyz",
	}

	req := &http.Request{
		Header: make(http.Header),
	}
	res, err := myJWTTransport.RoundTrip(req)
	if err != nil {
		t.Fatalf("Rountrip error: %s", err)
	}
	if res.StatusCode != 200 {
		t.Errorf("StatusCode is not 200, got:%d", res.StatusCode)
	}

	// 	apiInstance := API{
	// 		Options: Options{},
	// 		Client: MockClient{
	// 			ResponseOutput: &http.Response{
	// 				StatusCode: 200,
	// 				Body:       io.NopCloser(bytes.NewReader(wordsBytes)),
	// 			},
	// 		},
	// 	}
	// 	response, err := apiInstance.DoGetRequest("http://localhost/words")
	// 	if err != nil {
	// 		t.Errorf("DoGetRequest error: %s", err)
	// 	}
	// 	if response == nil {
	// 		t.Fatalf("reponse is empty")
	// 	}
	// 	if response.GetResponse() != strings.Join([]string{"a", "b"}, ", ") {
	// 		t.Errorf("Unexpected reponse: %s", response.GetResponse())
	// 	}

}
