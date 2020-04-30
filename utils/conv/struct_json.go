package conv

import (
	"github.com/json-iterator/go"
	"github.com/vfw4g/base/errors"
)

func Struct2Json(v interface{}) (j string, err error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if b, err := json.Marshal(&v); err != nil {
		return "", errors.WrapMark(err, "json Marshal error")
	} else {
		return string(b), nil
	}
}

func Json2Struct(data []byte, v interface{}) (err error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	return json.Unmarshal(data, &v)
}
