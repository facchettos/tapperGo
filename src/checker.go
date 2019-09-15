package gotapper

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
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

func checkSuccess(conditions ConditionDef, response *http.Response) bool {
	if !checkStatus(conditions, response) {
		return false
	}

	return checkBody(conditions, response)
}

func checkStatus(conditions ConditionDef, response *http.Response) bool {
	if conditions.ExpectedStatus != response.StatusCode {
		return false
	}
	return true

}

func checkBody(conditions ConditionDef, response *http.Response) bool {
	bodyAsString, err := readBody(response)
	if err != nil {
		return false
	}

	dynamic := make(map[string]interface{})
	json.Unmarshal([]byte(bodyAsString), &dynamic)
	field := selectField(dynamic, conditions.FieldSelector)

	return performCheck(conditions, field)
}

func checkType(a interface{}, b reflect.Type) bool {
	return reflect.TypeOf(a) == b
}

func checkString(actual, expected string) bool {
	return actual == expected
}

func checkInt(actual, expected int) bool {
	return actual == expected
}

func checkNumber(actual, expected float64) bool {
	return actual == expected
}

func checkBool(actual, expected bool) bool {
	return actual == expected
}

func checkLength(slice []interface{}, expectedSize int) bool {
	return len(slice) == expectedSize
}

func performCheck(condition ConditionDef, object interface{}) bool {
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
			return false
		}
		return checkType(object, reflectedType)

	}
	return false
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
