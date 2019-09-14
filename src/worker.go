package gotapper

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	Start  int = 0
	Stop   int = 1
	Reload int = 2
)

type order struct {
	Action       int
	WorkerConfig Config
}

func worker(orderChannel chan order) {
	configFromChannel := <-orderChannel
	executeOrder(configFromChannel)
	for {
		select {
		case order := <-orderChannel:
			if order.Action == Stop {
				break
			} else if order.Action == 2 {
				configFromChannel = order
			}
		case <-time.After(time.Duration(configFromChannel.WorkerConfig.Tick) * time.Second):
			executeOrder(configFromChannel)
		}
	}
}

func executeOrder(order order) (bool, error) {
	errorString := "Method unknown"
	var statusCode int
	var respBody []byte
	var err error
	switch order.WorkerConfig.Method {
	case http.MethodGet:
		statusCode, respBody, err = executeGet(order)
	case http.MethodPut:
		statusCode, respBody, err = executeRequest(order)
	case http.MethodPost:
		statusCode, respBody, err = executePost(order)
	case http.MethodPatch:
		statusCode, respBody, err = executePost(order)
	case http.MethodDelete:
		statusCode, respBody, err = executeRequest(order)
	}
	fmt.Println(statusCode, respBody, err)
	return false, errors.New(errorString)
}

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

func executeCallBacks(requests []RequestDef) []RequestResult {
	resultChan := make(chan RequestResult)
	defer close(resultChan)
	resultSlice := make([]RequestResult, len(requests))
	for _, v := range requests {
		go retryPostUntil200(v, resultChan)
	}

	for i := 0; i < len(requests); i++ {
		resultSlice[i] = <-resultChan
	}

	return resultSlice
}

func checkSuccess(conditions ConditionDef, response *http.Response) bool {
	if !checkStatus(conditions, response) {
		return false
	}
	//TODO finish the function
	return false
}

func checkStatus(conditions ConditionDef, response *http.Response) bool {
	if conditions.ExpectedStatus != response.StatusCode {
		return false
	}
	return true

}

func checkBody(conditions ConditionDef, response *http.Response, expected interface{}) bool {
	switch conditions.Condition {
	case "in":
		fmt.Println("in")
	case "and":
		fmt.Println("and")
	case "equals":
		fmt.Println("equals")
	case "not":
		fmt.Println("not")
	}
	return true
}
