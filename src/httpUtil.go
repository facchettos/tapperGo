package gotapper

import (
	"bytes"
	"net/http"
)

func createUrlFromOrder(order order) string {
	return order.WorkerConfig.URL
}

func executeGet(order order) (*http.Response, error) {
	effectiveUrl := createUrlFromOrder(order)

	resp, err := http.Get(effectiveUrl)
	return resp, err
}

func executePost(order order) (*http.Response, error) {
	effectiveUrl := createUrlFromOrder(order)

	postBody := bytes.NewBuffer([]byte(order.WorkerConfig.Body))
	resp, err := http.Post(effectiveUrl, order.WorkerConfig.ContentType, postBody)
	return resp, err
}

func executeRequest(order order) (*http.Response, error) {
	client := &http.Client{}
	effectiveUrl := createUrlFromOrder(order)
	method := order.WorkerConfig.Method
	var reqBody *bytes.Buffer
	if method == http.MethodPut || method == http.MethodPost {
		reqBody = bytes.NewBuffer([]byte(order.WorkerConfig.Body))

	}
	req, err := http.NewRequest(order.WorkerConfig.Method, effectiveUrl, reqBody)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	return resp, nil
}
