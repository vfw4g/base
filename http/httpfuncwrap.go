package vfwhttp

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/joomcode/errorx"
	"gopkg.in/go-playground/validator.v8"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	vgin = binding.Validator
)

const (
	IN_PATH  = "inpath"
	IN_QUERY = "inquery"
	IN_BODY  = "inbody"
)

func HttpFuncWrap(f HttpApiFunc) gin.HandlerFunc {
	return func(cxt *gin.Context) {

		bindStruct := func(in interface{}) (err error) {
			t := reflect.TypeOf(in)
			if t.Kind() != reflect.Ptr {
				return
			}
			v := reflect.ValueOf(in).Elem()
			return structEvaluate(v, cxt)
		}

		req := &ApiRequest{
			BindStruct: bindStruct,
			Context:    cxt,
		}

		rsp := f(req)

		//respon logic
		if rsp.e != nil {
			cxt.JSON(500, rsp.e)
		} else if rsp.o != nil {
			cxt.JSON(200, rsp.o)
			// cxt.SecureJSON(200,rsp.o)
		} else if rsp.m != nil {
			cxt.JSON(200, rsp.m)
			// cxt.SecureJSON(200,rsp.m)
		} else {
			cxt.JSON(200, map[string]string{
				"rsp_msg": "ok",
			})
		}
	}
}

func structEvaluate(v reflect.Value, cxt *gin.Context) (err error) {

	t := v.Type()
	tnf := t.NumField()
	for i := 0; i < tnf; i++ {
		//迭代每一个field
		tf := t.Field(i)
		if err = structFieldEvaluate(v, tf, cxt); err != nil {
			return
		}
	}
	return
}

///////////////////////////////////////////////////////////////////////////////

type FieldEvaluaterFunc func(reflect.Value, reflect.StructField, *gin.Context) error

func structFiledEvaluate(v reflect.Value, sf reflect.StructField, cxt *gin.Context) (err error) {
	if sf.Type.Kind() != reflect.Struct {
		return errorx.InitializationFailed.New(`the "%s\t" tag support struct filed only. `, "struct")
	}
	tv := sf.Tag.Get("struct")

	if strings.TrimSpace(tv) == "+" {
		n := reflect.New(sf.Type).Elem()
		structEvaluate(n, cxt)
		v.FieldByName(sf.Name).Set(n)
		return
	} else if tv == "json" {
		o := reflect.New(sf.Type).Interface()
		if err = cxt.BindJSON(o); err != nil {
			return errorx.InitializationFailed.WrapWithNoMessage(err)
		}
		tfv := reflect.ValueOf(o).Elem()
		v.FieldByName(sf.Name).Set(tfv)
	}
	return err
}

func queryStructFiledEvaluate(v reflect.Value, sf reflect.StructField, cxt *gin.Context) (err error) {
	if sf.Type.Kind() == reflect.Struct {
		return errorx.InitializationFailed.New(`the "%s\t" tag don't support struct filed. `, "query")
	}
	tv := sf.Tag.Get("query")

	tagElems := strings.Split(tv, ",")
	key := tagElems[0]
	if key == "" {
		key = sf.Name
	}
	kv := cxt.Query(key)

	if err = validateField(sf, kv); err != nil {
		return errorx.InitializationFailed.WrapWithNoMessage(err)
	}
	tfv, err := fieldSetVal(sf.Type.Kind(), kv)
	if err != nil {
		return errorx.InitializationFailed.WrapWithNoMessage(err)
	}
	v.FieldByName(sf.Name).Set(tfv)
	return
}

func pathStructFiledEvaluate(v reflect.Value, sf reflect.StructField, cxt *gin.Context) (err error) {
	if sf.Type.Kind() == reflect.Struct {
		return errorx.InitializationFailed.New(`the "%s\t" tag don't support struct filed. `, "path")
	}
	tv := sf.Tag.Get("path")

	tagElems := strings.Split(tv, ",")
	key := tagElems[0]
	if key == "" {
		key = sf.Name
	}

	kv := cxt.Param(key)

	if err = validateField(sf, kv); err != nil {
		return errorx.InitializationFailed.WrapWithNoMessage(err)
	}
	tfv, err := fieldSetVal(sf.Type.Kind(), kv)
	if err != nil {
		return errorx.InitializationFailed.WrapWithNoMessage(err)
	}
	v.FieldByName(sf.Name).Set(tfv)
	return err
}

func fieldEvaluaterResolve(tf reflect.StructField) FieldEvaluaterFunc {

	if ok, _ := parseTagName(tf, "query"); ok {
		return queryStructFiledEvaluate
	} else if ok, _ := parseTagName(tf, "path"); ok {
		return pathStructFiledEvaluate
	} else if ok, _ := parseTagName(tf, "struct"); ok {
		return structFiledEvaluate
	}
	return nil
}

///////////////////////////////////////////////////////////////////////////////

func structFieldEvaluate(v reflect.Value, tf reflect.StructField, cxt *gin.Context) error {

	f := fieldEvaluaterResolve(tf)

	return f(v, tf, cxt)
}

func validateField(tf reflect.StructField, kv string) error {
	v8 := vgin.Engine().(*validator.Validate)

	if ok, tagItems := parseTagName(tf, "binding"); ok {
		if err := v8.Field(kv, tagItems); err != nil {
			return errorx.InitializationFailed.New("expected:%s" + tagItems)
		}
	}
	return nil
}

func fieldSetVal(tftk reflect.Kind, kv string) (reflect.Value, error) {
	switch tftk {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
		res, err := toInt(kv)
		return reflect.ValueOf(res), err
	case reflect.Bool:
		res, err := toBoolean(kv)
		return reflect.ValueOf(res), err
	case reflect.Float32, reflect.Float64:
		res, err := toFloat(kv)
		return reflect.ValueOf(res), err
	case reflect.String:
		return reflect.ValueOf(kv), nil
	}
	return reflect.ValueOf(nil),
		errorx.InitializationFailed.New("don't support the type of filed:%s", tftk.String())
}

func parseTagName(sf reflect.StructField, tagName string) (bool, string) {
	ftag := sf.Tag
	tagValue, ok := ftag.Lookup(tagName)
	if !ok {
		return false, ""
	}
	return true, tagValue
}

func hasDefaultValue(s []string) (bool, int, string) {
	for i := 0; i < 2; i++ {
		defaultList := strings.SplitN(s[i], "=", 2)
		if defaultList[0] == "default" {
			return true, i, defaultList[1]
		}
	}
	return false, 0, ""
}

// ToFloat convert the input string to a float, or 0.0 if the input is not a float.
func toFloat(str string) (float64, error) {
	res, err := strconv.ParseFloat(str, 64)
	if err != nil {
		res = 0.0
	}
	return res, err
}

// ToInt convert the input string or any int type to an integer type 64, or 0 if the input is not an integer.
func toInt(value interface{}) (res int64, err error) {
	val := reflect.ValueOf(value)

	switch value.(type) {
	case int, int8, int16, int32, int64:
		res = val.Int()
	case uint, uint8, uint16, uint32, uint64:
		res = int64(val.Uint())
	case string:
		if isInt(val.String()) {
			res, err = strconv.ParseInt(val.String(), 0, 64)
			if err != nil {
				res = 0
			}
		} else {
			err = fmt.Errorf("math: square root of negative number %g", value)
			res = 0
		}
	default:
		err = fmt.Errorf("math: square root of negative number %g", value)
		res = 0
	}

	return
}

// ToBoolean convert the input string to a boolean.
func toBoolean(str string) (bool, error) {
	return strconv.ParseBool(str)
}

// IsInt check if the string is an integer. Empty string is valid.
func isInt(str string) bool {
	if isNull(str) {
		return true
	}
	intPattern := "^(?:[-+]?(?:0|[1-9][0-9]*))$"
	return regexp.MustCompile(intPattern).MatchString(str)
}

// IsNull check if the string is null.
func isNull(str string) bool {
	return len(str) == 0
}

///////////////////////////////////////////////////////////////////////////////

type HttpApiFunc func(*ApiRequest) ApiResponse

type ApiRequest struct {
	BindStruct func(interface{}) error
	*gin.Context
}

type ApiResponse struct {
	e error
	m map[string]interface{}
	o interface{}
}

func (rsp ApiResponse) BindError(err error) ApiResponse {
	rsp.e = err
	return rsp
}

func (rsp ApiResponse) AddAttribute(k string, v interface{}) ApiResponse {
	if rsp.m == nil {
		rsp.m = make(map[string]interface{})
	}
	rsp.m[k] = v
	return rsp
}

func (rsp ApiResponse) AddObject(o interface{}) ApiResponse {
	rsp.o = o
	return rsp
}
