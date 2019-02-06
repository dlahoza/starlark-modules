package builtin

import (
	"errors"

	"go.starlark.net/starlark"
)

type Function func(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error)

type builtin struct {
	functions map[string]Function
}

func New(builders map[string]Function) starlark.Value {
	return &builtin{
		functions: builders,
	}
}

// Attr returns a starlark value that wraps the method or field with the given
// name.
func (m *builtin) Attr(name string) (starlark.Value, error) {
	if v, ok := m.functions[name]; ok {
		return starlark.NewBuiltin(name, v), nil
	}
	return nil, nil
}

// AttrNames returns the list of all fields and methods on this struct.
func (m *builtin) AttrNames() (list []string) {
	for k := range m.functions {
		list = append(list, k)
	}
	return
}

// String returns the string representation of the value.
// Starlark string values are quoted as if by Python's repr.
func (m *builtin) String() string {
	return ""
}

// Type returns a short string describing the value's type.
func (m *builtin) Type() string {
	return "starlark_module"
}

// Freeze causes the value, and all values transitively
// reachable from it through collections and closures, to be
// marked as frozen.  All subsequent mutations to the data
// structure through this API will fail dynamically, making the
// data structure immutable and safe for publishing to other
// Starlark interpreters running concurrently.
func (m *builtin) Freeze() {}

// Truth returns the truth value of an object.
func (m *builtin) Truth() starlark.Bool {
	return true
}

// Hash returns a function of x such that Equals(x, y) => Hash(x) == Hash(y).
// Hash may fail if the value's type is not hashable, or if the value
// contains a non-hashable value.
func (m *builtin) Hash() (uint32, error) {
	return 0, errors.New("starlark_module is not hashable")
}
