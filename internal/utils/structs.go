package utils

import (
	"fmt"
	"reflect"
)

// IsStructEmpty checks if any fields of a struct are empty
func IsStructEmpty(s interface{}) bool {
	v := reflect.ValueOf(s)
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Interface() != reflect.Zero(v.Field(i).Type()).Interface() {
			fmt.Println(s)
			return false
		}
	}
	return true
}