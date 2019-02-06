package builtin

import (
	"fmt"
	"reflect"

	"github.com/DLag/starlight/convert"
	"go.starlark.net/starlark"
)

func ConvertToStringMap(v interface{}) interface{} {
	switch i := v.(type) {
	case *convert.GoMap:
		return ConvertToStringMap(i.Value().Interface())
	case starlark.StringDict:
		m := convert.FromStringDict(i)
		for key := range m {
			m[key] = ConvertToStringMap(m[key])
		}
		return m
	case map[string]interface{}:
		for key := range i {
			i[key] = ConvertToStringMap(i[key])
		}
		return i
	case map[interface{}]interface{}:
		converted := make(map[string]interface{})
		for key, value := range i {
			strKey := fmt.Sprintf("%v", key)
			converted[strKey] = ConvertToStringMap(value)
		}
		return converted
	case []interface{}:
		for key, value := range i {
			i[key] = ConvertToStringMap(value)
		}
		return i
	case *convert.GoInterface:
		return i.Value().Interface()
	}
	return v
}

func ToValue(v interface{}) (starlark.Value, error) {
	if val, ok := v.(starlark.Value); ok {
		return val, nil
	}
	return convertValue(reflect.ValueOf(v))
}

func convertMapToDict(v reflect.Value) starlark.Value {
	d := starlark.Dict{}
	for _, k := range v.MapKeys() {
		kv, err := ToValue(k.Interface())
		if err != nil {
			continue
		}
		vv, err := ToValue(v.MapIndex(k).Interface())
		if err != nil {
			continue
		}
		err = d.SetKey(kv, vv)
		if err != nil {
			continue
		}
	}
	return &d
}

func convertSliceToList(v reflect.Value) starlark.Value {
	l := starlark.List{}
	for i := 0; i < v.Len(); i++ {
		vv, err := ToValue(v.Index(i).Interface())
		if err != nil {
			continue
		}
		err = l.Append(vv)
		if err != nil {
			continue
		}
	}
	return &l
}

func convertValue(val reflect.Value) (starlark.Value, error) {
	kind := val.Kind()
	if kind == reflect.Ptr {
		kind = val.Elem().Kind()
	}
	switch kind {
	case reflect.Bool:
		return starlark.Bool(val.Bool()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return starlark.MakeInt64(val.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return starlark.MakeUint64(val.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return starlark.Float(val.Float()), nil
	case reflect.Func:
		return convert.ToValue(val.Interface())
	case reflect.Map:
		return convertMapToDict(val), nil
	case reflect.String:
		return starlark.String(val.String()), nil
	case reflect.Slice, reflect.Array:
		return convertSliceToList(val), nil
	case reflect.Struct:
		return convert.ToValue(val.Interface())
	}

	return nil, fmt.Errorf("type %T is not a supported starlark type", val.Interface())
}
