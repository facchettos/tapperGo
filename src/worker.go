package gotapper

import (
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
	workerConfig config
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

func executeOrder(order order) (bool, string) {
	errorString := "Method unknown"
	switch order.workerConfig.Method {
	case http.MethodGet:
		fmt.Println("this is a get")
	case http.MethodPut:
		fmt.Println("toto")
	case http.MethodPost:
		fmt.Println("post")
	case http.MethodPatch:
		fmt.Println("patch")
	case http.MethodDelete:
		fmt.Println("delete")
	}

	return false, errorString
}
