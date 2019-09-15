package gotapper

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func createClient(config TestDefinition) http.Client {
	return http.Client{
		Timeout: time.Duration(config.ExpectedDuration) * time.Second,
	}
}

//takes the request definition and builds the url with the parameters
func createUrlFromOrder(order TestDefinition) string {
	var builder strings.Builder
	builder.WriteString(order.URL)
	builder.WriteString("?")
	counter := 0
	length := len(order.Arguments)
	for k, v := range order.Arguments {
		builder.WriteString(k + "=" + v)
		if counter < length {
			builder.WriteString("&")
		}
		counter++
	}
	return builder.String()
}

func executeGet(order TestDefinition) (*http.Response, error) {
	effectiveUrl := createUrlFromOrder(order)
	client := createClient(order)
	resp, err := client.Get(effectiveUrl)
	return resp, err
}

func executePost(order TestDefinition) (*http.Response, error) {
	effectiveUrl := createUrlFromOrder(order)

	client := createClient(order)
	postBody := bytes.NewBuffer([]byte(order.Body))
	resp, err := client.Post(effectiveUrl, order.ContentType, postBody)
	return resp, err
}

func executeRequest(order TestDefinition) (*http.Response, error) {
	client := createClient(order)
	effectiveUrl := createUrlFromOrder(order)
	method := order.Method
	var reqBody *bytes.Buffer
	if method == http.MethodPut || method == http.MethodPost {
		reqBody = bytes.NewBuffer([]byte(order.Body))

	}
	req, err := http.NewRequest(order.Method, effectiveUrl, reqBody)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	return resp, nil
}

func readBody(resp *http.Response) (string, error) {
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

//this is for the callbacks
func retryPostUntil200(reqDef RequestDef, resChan chan RequestResult) {
	for i := 1; i <= reqDef.Retries; i++ {
		resp, err := http.Post(reqDef.Url, reqDef.ContentType, bytes.NewBuffer([]byte(reqDef.Body)))
		if err != nil {
			if i == reqDef.Retries {
				resChan <- RequestResult{StatusCode: resp.StatusCode, Error: err, Name: reqDef.Name}
			}
			continue
		}
		resChan <- RequestResult{StatusCode: resp.StatusCode, Error: nil, Name: reqDef.Name}
	}
}
