package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/valyala/fasthttp"
)

type HTTPClient interface {
	Do(req *Request) (*Response, error)
}

type Request struct {
	URL     string
	Method  string
	Body    []byte
	Headers map[string]string
}

type Response struct {
	StatusCode int
	Body       []byte
}

type StandardHTTPClient struct {
	client *http.Client
}

func NewStandardHTTPClient() *StandardHTTPClient {
	return &StandardHTTPClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 100,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}
}

func (c *StandardHTTPClient) Do(req *Request) (*Response, error) {
	httpReq, err := http.NewRequest(req.Method, req.URL, bytes.NewReader(req.Body))
	if err != nil {
		return nil, err
	}

	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		StatusCode: httpResp.StatusCode,
		Body:       body,
	}, nil
}

type FastHTTPClient struct {
	client *fasthttp.Client
}

func NewFastHTTPClient() *FastHTTPClient {
	return &FastHTTPClient{
		client: &fasthttp.Client{
			MaxConnsPerHost: 10000,
		},
	}
}

func (c *FastHTTPClient) Do(req *Request) (*Response, error) {
	fasthttpReq := fasthttp.AcquireRequest()
	fasthttpResp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(fasthttpReq)
	defer fasthttp.ReleaseResponse(fasthttpResp)

	fasthttpReq.SetRequestURI(req.URL)
	fasthttpReq.Header.SetMethod(req.Method)
	for k, v := range req.Headers {
		fasthttpReq.Header.Set(k, v)
	}
	fasthttpReq.SetBody(req.Body)

	err := c.client.DoTimeout(fasthttpReq, fasthttpResp, 10*time.Second)
	if err != nil {
		return nil, err
	}

	return &Response{
		StatusCode: fasthttpResp.StatusCode(),
		Body:       fasthttpResp.Body(),
	}, nil
}
