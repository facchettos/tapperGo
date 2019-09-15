package gotapper

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

func selectField(jsonMap map[string]interface{}, fieldSelector FieldSelector) interface{} {
	fieldSlice := strings.Split(fieldSelector.Field, fieldSelector.Separator)
	tempMap := jsonMap
	for _, v := range fieldSlice {
		nodeType := reflect.ValueOf(jsonMap[v])
		if nodeType.Kind() != reflect.Map {
			break
		}
		tempMap = tempMap[v].(map[string]interface{})
	}

	return tempMap[fieldSlice[len(fieldSlice)-1]]
}

func checkSuccess(conditions TestDefinition, response *http.Response) (bool, error) {
	if res, err := checkStatus(conditions, response); !res {
		return false, err
	}

	return checkBody(conditions, response)
}

func checkStatus(conditions TestDefinition, response *http.Response) (bool, error) {
	if conditions.ExpectedStatus != response.StatusCode {
		return false, errors.New("wrong status code: expected " +
			string(conditions.ExpectedStatus) +
			" but found :" + string(response.StatusCode))
	}
	return true, nil

}

func checkBody(conditions TestDefinition, response *http.Response) (bool, error) {
	bodyAsString, err := readBody(response)
	if err != nil {
		return false, err
	}

	dynamic := make(map[string]interface{})
	json.Unmarshal([]byte(bodyAsString), &dynamic)
	field := selectField(dynamic, conditions.FieldSelector)

	return performCheck(conditions, field)
}

func checkType(a interface{}, b reflect.Type) (bool, error) {
	if reflect.TypeOf(a) == b {
		return true, nil
	} else {
		return false, errors.New("type mismatch : object is :" +
			reflect.TypeOf(a).String() + " but expected : " + b.String())
	}
}

func checkString(actual, expected string) (bool, error) {
	if actual != expected {
		return false, errors.New("Wrong result: found : " +
			actual +
			" while expecting :" + expected)
	}
	return true, nil
}

func checkInt(actual, expected int) (bool, error) {
	if actual != expected {
		return false, errors.New("Wrong result: found : " +
			strconv.FormatInt(int64(actual), 10) +
			" while expecting :" + strconv.FormatInt(int64(expected), 10))
	}
	return true, nil
}

func checkNumber(actual, expected float64) (bool, error) {
	if actual != expected {
		return false, errors.New("Wrong result: found : " +
			strconv.FormatFloat(actual, 'f', -1, 64) +
			" while expecting :" + strconv.FormatFloat(expected, 'f', -1, 64))
	}
	return true, nil
}

func checkBool(actual, expected bool) (bool, error) {

	if actual != expected {
		return false, errors.New("Wrong result: found : " +
			strconv.FormatBool(actual) +
			" while expecting :" + strconv.FormatBool(expected))
	}
	return true, nil
}

func checkLength(slice []interface{}, expectedSize int) (bool, error) {
	if len(slice) != expectedSize {
		return false, errors.New("Wrong size :" +
			string(len(slice)) +
			" while expecting :" + string(expectedSize))
	}
	return true, nil
}

func performCheck(condition TestDefinition, object interface{}) (bool, error) {
	switch condition.ExpectedType {
	case "int":
		return checkInt(object.(int), condition.ExpectedInt)
	case "bool":
		return checkBool(object.(bool), condition.ExpectedBool)
	case "number":
		return checkNumber(object.(float64), condition.ExpectedNumber)
	case "type":
		reflectedType, err := typeFromString(condition.ExpectedType)
		if err != nil {
			return false, err
		}
		return checkType(object, reflectedType)

	}
	return false, errors.New("unknow check operation")
}

func typeFromString(typeAsString string) (reflect.Type, error) {
	switch typeAsString {
	case "int":
		return reflect.TypeOf(1), nil
	case "number":
		return reflect.TypeOf(0.0), nil
	case "bool":
		return reflect.TypeOf(true), nil
	case "string":
		return reflect.TypeOf("string"), nil
	}
	return nil, errors.New("type not defined")
}
