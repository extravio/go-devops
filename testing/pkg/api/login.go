package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type LoginRequest struct {
	Password string `json:"password"`
}
type LoginResponse struct {
	Token string `json:"token"`
}

// the first string type can be omitted
func doLoginRequest(client ClientIface, requestURL, password string) (string, error) {
	LoginRequest := LoginRequest{
		Password: password,
	}
	body, err := json.Marshal(LoginRequest)
	if err != nil {
		return "", fmt.Errorf("Marshall Error: %s", err)
	}

	// use NewBuffer that implements a "Read" function
	response, err := client.Post(requestURL, "application/json", bytes.NewBuffer(body))

	if err != nil {
		return "", fmt.Errorf("HTTP Post error: %s", err)
	}

	defer response.Body.Close()

	resBody, err := io.ReadAll(response.Body)

	if err != nil {
		return "", fmt.Errorf("ReadAll error: %s", err)
	}

	if response.StatusCode != 200 {
		return "", fmt.Errorf("Invalid output (HTTP Status Code: %d): %s", response.StatusCode, resBody)
	}

	if !json.Valid(resBody) {
		return "", RequestError{
			HTTPCode: response.StatusCode,
			Body:     string(resBody),
			Err:      fmt.Sprintf("No valid JSON returned"),
		}
	}

	var loginResponse LoginResponse
	err = json.Unmarshal(resBody, &loginResponse)
	if err != nil {
		return "", RequestError{
			HTTPCode: response.StatusCode,
			Body:     string(resBody),
			Err:      fmt.Sprintf("Page Unmarshal error: %s", err),
		}
	}

	return loginResponse.Token, nil
}
