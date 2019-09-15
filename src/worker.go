package gotapper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
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
	var response *http.Response
	var err error
	switch order.WorkerConfig.Method {
	case http.MethodGet:
		response, err = executeGet(order)
	case http.MethodPut:
		response, err = executeRequest(order)
	case http.MethodPost:
		response, err = executePost(order)
	case http.MethodPatch:
		response, err = executePost(order)
	case http.MethodDelete:
		response, err = executeRequest(order)
	}
	fmt.Println(response, err)
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
