package convert

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/DLag/starlark-modules/builtin"
	sconvert "github.com/DLag/starlight/convert"

	"go.starlark.net/starlark"
)

var StructTags = []string{"starlark"}

type MethodsWrapper interface {
	Methods() map[string]builtin.Function
}

type StarlarkStruct struct {
	methods map[string]starlark.Value
	fields  map[string]reflect.Value
}

func NewStruct(v interface{}) starlark.Value {
	st := &StarlarkStruct{
		methods: make(map[string]starlark.Value),
		fields:  make(map[string]reflect.Value),
	}
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return st
	}
	if fw, ok := v.(MethodsWrapper); ok {
		for k, m := range fw.Methods() {
			st.methods[k] = starlark.NewBuiltin(k, m)
		}
	} else {
		for i := 0; i < val.NumMethod(); i++ {
			name := val.Type().Method(i).Name
			method := val.Method(i)
			if method.Kind() == reflect.Invalid {
				continue
			}
			if m, ok := method.Interface().(func(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error)); ok {
				st.methods[name] = starlark.NewBuiltin(name, m)
			} else {
				if m, err := ToValue(method); err != nil && m.Type() == (&starlark.Builtin{}).Type() {
					st.methods[name] = m
				}
			}
		}
	}
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		if val.Field(i).Kind() == reflect.Invalid {
			continue
		}
		var name string
		for _, t := range StructTags {
			name = field.Tag.Get(t)
		}
		switch name {
		case "-":
			continue
		case "":
			name = field.Name
		}
		st.fields[name] = val.Field(i)
	}
	return st
}

// Attr returns a starlark value that wraps the method or field with the given
// name.
func (s *StarlarkStruct) Attr(name string) (starlark.Value, error) {
	if s.fields == nil || s.methods == nil {
		return nil, nil
	}
	if method, ok := s.methods[name]; ok {
		return method, nil
	}
	if field, ok := s.fields[name]; ok {
		return ToValue(field.Interface())
	}
	return nil, nil
}

// AttrNames returns the list of all fields and methods on this struct.
func (s *StarlarkStruct) AttrNames() []string {
	nameMap := make(map[string]struct{})
	for name := range s.methods {
		nameMap[name] = struct{}{}
	}
	for name := range s.fields {
		nameMap[name] = struct{}{}
	}
	names := make([]string, len(nameMap))
	i := 0
	for name := range nameMap {
		names[i] = name
		i++
	}
	return names
}

// SetField sets the struct field with the given name with the given value.
func (s *StarlarkStruct) SetField(name string, val starlark.Value) error {
	field := s.fields[name]
	if field.CanSet() {
		val := conv(val, field.Type())
		field.Set(val)
		return nil
	}
	return fmt.Errorf("%s is not a settable field", name)
}

// String returns the string representation of the value.
// Starlark string values are quoted as if by Python's repr.
func (s *StarlarkStruct) String() string {
	return fmt.Sprint(s.fields)
}

// Type returns a short string describing the value's type.
func (s *StarlarkStruct) Type() string {
	return "starlark_go_struct"
}

// Freeze causes the value, and all values transitively
// reachable from it through collections and closures, to be
// marked as frozen.  All subsequent mutations to the data
// structure through this API will fail dynamically, making the
// data structure immutable and safe for publishing to other
// Starlark interpreters running concurrently.
func (s *StarlarkStruct) Freeze() {}

// Truth returns the truth value of an object.
func (s *StarlarkStruct) Truth() starlark.Bool {
	return true
}

// Hash returns a function of x such that Equals(x, y) => Hash(x) == Hash(y).
// Hash may fail if the value's type is not hashable, or if the value
// contains a non-hashable value.
func (s *StarlarkStruct) Hash() (uint32, error) {
	return 0, errors.New("starlark_go_struct is not hashable")
}

// conv tries to convert v to t if v is not assignable to t.
func conv(v starlark.Value, t reflect.Type) reflect.Value {
	out := reflect.ValueOf(sconvert.FromValue(v))
	if !out.Type().AssignableTo(t) {
		return out.Convert(t)
	}
	return out
}
