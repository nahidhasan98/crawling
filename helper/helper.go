package helper

import "net/http"

func GETRequest(apiURL string) *http.Response {
	// for panic() error Recovery
	defer ErrorRecovery()

	req, err := http.NewRequest("GET", apiURL, nil)
	ErrorCheck(err)

	client := &http.Client{}
	response, err := client.Do(req)
	ErrorCheck(err)

	return response
}
