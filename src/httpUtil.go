package gotapper

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func createUrlFromOrder(order order) string {
	return order.WorkerConfig.URL
}

func executeGet(order order) (int, []byte, error) {
	effectiveUrl := createUrlFromOrder(order)

	resp, err := http.Get(effectiveUrl)
	if err != nil {
		return 0, nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}
	return resp.StatusCode, body, nil
}

func executePost(order order) (int, []byte, error) {
	effectiveUrl := createUrlFromOrder(order)

	postBody := bytes.NewBuffer([]byte(order.WorkerConfig.Body))
	resp, err := http.Post(effectiveUrl, order.WorkerConfig.ContentType, postBody)

	if err != nil {
		return 0, nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	return resp.StatusCode, body, nil
}

func executeRequest(order order) (int, []byte, error) {
	client := &http.Client{}
	effectiveUrl := createUrlFromOrder(order)
	method := order.WorkerConfig.Method
	var reqBody *bytes.Buffer
	if method == http.MethodPut || method == http.MethodPost {
		reqBody = bytes.NewBuffer([]byte(order.WorkerConfig.Body))

	}
	req, err := http.NewRequest(order.WorkerConfig.Method, effectiveUrl, reqBody)
	if err != nil {
		return 0, nil, err
	}
	resp, err := client.Do(req)

	if err != nil {
		return 0, nil, err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	return resp.StatusCode, body, nil
}
