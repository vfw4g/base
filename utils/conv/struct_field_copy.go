package conv

import (
	"github.com/vfw4g/base/errors"
	"reflect"
	"strings"
)

//copy field from source to target.
//the target must be a point or a struct kind,or it will not be attended to.
//support the `fieldcopy` tag in target field.
func StructFieldCopy(src interface{}, target interface{}) (err error) {
	if src == nil || target == nil {
		return errors.New("source or target must not be nil.")
	}
	//source
	var st = reflect.TypeOf(src)
	var sv = reflect.ValueOf(src)

	if st.Kind() == reflect.Ptr {
		st = st.Elem()
		sv = sv.Elem()
	}
	if st.Kind() != reflect.Struct {
		return
	}
	ssr := &structRef{
		st, sv,
	}
	//target
	dt := reflect.TypeOf(target)
	if dt.Kind() != reflect.Ptr {
		return
	}
	dt = dt.Elem()
	if dt.Kind() != reflect.Struct {
		return
	}
	dv := reflect.ValueOf(target).Elem()
	for i := 0; i < dv.NumField(); i++ {
		df := dt.Field(i)
		dfv := dv.Field(i)
		if tv, ok := df.Tag.Lookup("fieldcopy"); ok {
			if strings.Contains(tv, "-") {
				continue
			}
			stv := strings.Split(tv, ",")
			switch len(stv) {
			case 0:
				//`fieldcopy:`
				if sfv, ok := resolveSourceValueFromField(df.Name, ssr, df.Type.Kind(), false); ok {
					dfv.Set(sfv)
				}
			case 1:
				if strings.TrimSpace(stv[0]) == "omitempty" {
					//`fieldcopy:"omitempty"`
					if sfv, ok := resolveSourceValueFromField(df.Name, ssr, df.Type.Kind(), true); ok {
						dfv.Set(sfv)
					}
				} else {
					//`fieldcopy:"nameFlag"`
					if sfv, ok := resolveSourceValueFromTag(strings.TrimSpace(stv[0]), ssr, df.Type.Kind(), false); ok {
						dfv.Set(sfv)
					}
				}
			case 2:
				if strings.TrimSpace(stv[0]) == "omitempty" {
					if strings.TrimSpace(stv[1]) != "" {
						//`fieldcopy:"omitempty,nameFlag"`
						if sfv, ok := resolveSourceValueFromTag(strings.TrimSpace(stv[1]), ssr, df.Type.Kind(), true); ok {
							dfv.Set(sfv)
						}
					} else if strings.TrimSpace(stv[1]) == "" {
						//`fieldcopy:"omitempty,"`
						//same as case 1
						if sfv, ok := resolveSourceValueFromField(df.Name, ssr, df.Type.Kind(), true); ok {
							dfv.Set(sfv)
						}
					}
				} else if strings.TrimSpace(stv[1]) == "omitempty" {
					if strings.TrimSpace(stv[0]) != "" {
						//`fieldcopy:"nameFlag,omitempty"`
						if sfv, ok := resolveSourceValueFromTag(strings.TrimSpace(stv[0]), ssr, df.Type.Kind(), true); ok {
							dfv.Set(sfv)
						}
					} else if strings.TrimSpace(stv[0]) == "" {
						//`fieldcopy:",omitempty"`
						//same as case 1
						if sfv, ok := resolveSourceValueFromField(df.Name, ssr, df.Type.Kind(), true); ok {
							dfv.Set(sfv)
						}
					}
				} else if strings.TrimSpace(stv[0]) == "" && strings.TrimSpace(stv[1]) == "" {
					//`fieldcopy:","`
					//same as case 0
					if sfv, ok := resolveSourceValueFromField(df.Name, ssr, df.Type.Kind(), false); ok {
						dfv.Set(sfv)
					}
				}
			default:
				//unSupport tag format
				//do nothing
			}
		} else {
			if sfv, ok := resolveSourceValueFromField(df.Name, ssr, df.Type.Kind(), false); ok {
				dfv.Set(sfv)
			}
		}
	}
	return
}

type structRef struct {
	st reflect.Type
	sv reflect.Value
}

func resolveSourceValueFromField(name string, ssr *structRef, kind reflect.Kind, omitempty bool) (value reflect.Value, ok bool) {
	if stf, exist := ssr.st.FieldByName(name); exist && kind == stf.Type.Kind() {
		value = ssr.sv.FieldByName(name)
		if omitempty && value.IsZero() {
			//ignore zero value
			return value, false
		}
		return value, true
	}
	return value, false
}

func resolveSourceValueFromTag(name string, ssr *structRef, kind reflect.Kind, omitempty bool) (value reflect.Value, ok bool) {
	//have tag,will return field value of the tag
	for i := 0; i < ssr.st.NumField(); i++ {
		ft := ssr.st.Field(i)
		if v, exist := ft.Tag.Lookup("fieldcopy"); exist {
			lv := strings.Split(v, ",")
			for _, item := range lv {
				if strings.TrimSpace(item) == name && kind == ft.Type.Kind() {
					value = ssr.sv.Field(i)
					if omitempty && value.IsZero() {
						//ignore zero value
						return value, false
					}
					return value, true
				}
			}
		}
	}
	//or no tag,will return zero value and false
	return value, false
}
