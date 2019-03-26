package easycc

import (
	"bytes"
	"net/http"
)

type CCRequest struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    []byte
}

type CCResponse struct {
	Response *http.Response
	Err      error
}

func sendRequest(reqCC *CCRequest) *CCResponse {
	client := &http.Client{}

	req, err := http.NewRequest(
		reqCC.Method,
		reqCC.URL,
		bytes.NewReader(reqCC.Body))
	if err != nil {
		return &CCResponse{
			Err: err,
		}
	}
	for key, value := range reqCC.Headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	defer resp.Body.Close()

	return &CCResponse{
		Response: resp,
		Err:      err,
	}
}
