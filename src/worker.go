package gotapper

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	Start   int = 0
	Stop    int = 1
	Restart int = 2
)

type order struct {
	Action       int
	WorkerConfig Config
}

func worker(orderStruct Config, orderChannel chan int) {
	var err error
	for {
		err = nil
		select {
		case order := <-orderChannel:
			if order == Stop {
				break
			} else if order == Restart {
				err = executeOrders(orderStruct.Conditions)

			}
		case <-time.After(time.Duration(orderStruct.Tick) * time.Second):
			err = executeOrders(orderStruct.Conditions)
		}
		executeCallBacksFromStatus(orderStruct, err)
	}
}

func executeCallBacksFromStatus(order Config, err error) {
	if err != nil {
		executeCallBacks(order.CallBackUrlsFailure)
	} else {
		executeCallBacks(order.CallBackUrlsSuccess)
	}
}

func executeOrders(reqs []TestDefinition) error {
	for _, v := range reqs {
		time.Sleep(time.Duration(v.SleepBefore) * time.Second)
		_, err := executeOrder(v)
		if err != nil {
			return err
		}
		time.Sleep(time.Duration(v.SleepAfter) * time.Second)
	}
	return nil
}

func executeOrder(order TestDefinition) (bool, error) {
	errorString := "Method unknown"
	var response *http.Response
	var err error
	switch order.Method {
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
