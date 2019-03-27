package easycc

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net/http"
)

type CCRequest struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    []byte
}

type CCResponse struct {
	Status       string // e.g. "200 OK"
	StatusCode   int    // e.g. 200
	Proto        string // e.g. "HTTP/1.0"
	ProtoMajor   int    // e.g. 1
	ProtoMinor   int    // e.g. 0
	Headers      http.Header
	Uncompressed bool
	Trailer      http.Header
	Body         []byte
	Request      *http.Request
	TLS          *tls.ConnectionState
	Err          error
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

	body, _ := ioutil.ReadAll(resp.Body)

	return &CCResponse{
		Status:       resp.Status,
		StatusCode:   resp.StatusCode,
		Proto:        resp.Proto,
		ProtoMajor:   resp.ProtoMajor,
		ProtoMinor:   resp.ProtoMinor,
		Headers:      resp.Header,
		Uncompressed: resp.Uncompressed,
		Trailer:      resp.Trailer,
		Body:         body,
		Request:      resp.Request,
		TLS:          resp.TLS,
		Err:          err,
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
