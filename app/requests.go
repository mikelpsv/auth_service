package app

import (
	"net/http"
	"strconv"
)

func GetSimpleValue(r *http.Request, paramName string) (value string, exists bool)  {
	val, exists := r.URL.Query()[paramName]
	return val[0], exists
}

func GetSimpleValueAsInt(r *http.Request, paramName string) (value int64, exists bool, err error)  {
	val_map, exists := r.URL.Query()[paramName]
	value, err = strconv.ParseInt(val_map[0], 10, 64)
	return value, exists, err
}