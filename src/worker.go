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
	action       int
	workerConfig Config
}

func worker(orderChannel chan order) {
	configFromChannel := <-orderChannel
	executeOrder(configFromChannel)
	for {
		select {
		case order := <-orderChannel:
			if order.action == Stop {
				break
			} else if order.action == 2 {
				configFromChannel = order
			}
		case <-time.After(time.Duration(configFromChannel.workerConfig.Tick) * time.Second):
			executeOrder(configFromChannel)
		}
	}
}

func executeOrder(order order) (bool, error) {
	errorString := "Method unknown"
	switch order.workerConfig.Method {
	case http.MethodGet:
		fmt.Println("get")
	case http.MethodPut:
		fmt.Println("put")
	case http.MethodPost:
		fmt.Println("post")
	case http.MethodPatch:
		fmt.Println("patch")
	case http.MethodDelete:
		fmt.Println("delete")
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
