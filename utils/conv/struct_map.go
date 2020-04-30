package conv

import (
	"github.com/goinggo/mapstructure"
	"reflect"
)

func Struct2Map(obj interface{}) map[string]interface{} {
	to := reflect.TypeOf(obj)
	vo := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < to.NumField(); i++ {
		data[to.Field(i).Name] = vo.Field(i).Interface()
	}
	return data
}

func Map2Struct(m map[string]interface{}, ps interface{}) error {
	return mapstructure.Decode(m, ps)
}
