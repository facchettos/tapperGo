package gotapper

import (
	"bytes"
	"errors"
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
	switch order.WorkerConfig.Method {
	case http.MethodGet:
		executeGet(order)
	case http.MethodPut:
		executeRequest(order)
	case http.MethodPost:
		executePost(order)
	case http.MethodPatch:
		executePost(order)
	case http.MethodDelete:
		executeRequest(order)
	}

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
