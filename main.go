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

func CCTest(reqCC *CCRequest, concurrency uint) []*CCResponse {
	chs := make([]chan *CCResponse, concurrency) // 保存请求结果
	var index uint = 0
	for index = 0; index < concurrency; index++ {
		chs[index] = make(chan *CCResponse)
		go func(ch chan *CCResponse) {
			ch <- sendRequest(reqCC)
		}(chs[index])
	}
	respArr := make([]*CCResponse, concurrency)
	for index, ch := range chs {
		respArr[index] = <-ch
	}
	return respArr
}
