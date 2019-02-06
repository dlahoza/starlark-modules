package json

import (
	"errors"
	"testing"

	"github.com/DLag/starlark-modules/builtin"

	"github.com/starlight-go/starlight/convert"
	"github.com/stretchr/testify/assert"
	"go.starlark.net/resolve"
	"go.starlark.net/starlark"
)

func TestJson(t *testing.T) {
	a := assert.New(t)
	resolve.AllowFloat = true
	Marshal = func(v interface{}) ([]byte, error) {
		a.IsType((map[string]interface{})(nil), v)
		m := v.(map[string]interface{})
		a.Equal("b", m["a"])
		a.Equal(int64(123), m["c"])
		return []byte(`{"a":"b","c":123}`), nil
	}
	Unmarshal = func(data []byte, v interface{}) error {
		a.IsType((*map[string]interface{})(nil), v)
		a.NotNil(data)
		a.Equal([]byte(`{"a":"b","c":123}`), data)
		m := v.(*map[string]interface{})
		*m = map[string]interface{}{}
		(*m)["a"] = "b"
		(*m)["c"] = int64(123)
		return nil
	}

	predeclared := starlark.StringDict{
		"json": New(),
	}
	thread := &starlark.Thread{}
	t.Run("HappyPath", func(t *testing.T) {
		script := `
a1 = json.dump({'a': 'b', 'c': 123})
a2 = json.parse(a1)
`
		vars, err := starlark.ExecFile(thread, "script", script, predeclared)
		a.NoError(err)
		a.NotNil(vars)

		m := convert.FromStringDict(vars)
		a.NotNil(m)
		a.Equal(`{"a":"b","c":123}`, m["a1"])
		expected := map[string]interface{}{
			"a": "b",
			"c": int64(123),
		}
		a.EqualValues(expected, builtin.ConvertToStringMap(m).(map[string]interface{})["a2"])
	})
	t.Run("dumpWrongParams1", func(t *testing.T) {
		script := `
a1 = json.dump()
`
		vars, err := starlark.ExecFile(thread, "script", script, predeclared)
		a.Error(err)
		a.NotNil(vars)
	})
	t.Run("dumpWrongParams2", func(t *testing.T) {
		script := `
a1 = json.dump(123)
`
		vars, err := starlark.ExecFile(thread, "script", script, predeclared)
		a.Error(err)
		a.NotNil(vars)
	})
	t.Run("dumpMarshalError", func(t *testing.T) {
		Marshal = func(v interface{}) ([]byte, error) {
			return nil, errors.New("some err")
		}
		script := `
a1 = json.dump({'a': 'b', 'c': 123})
`
		vars, err := starlark.ExecFile(thread, "script", script, predeclared)
		a.Error(err)
		a.NotNil(vars)
	})
	t.Run("parseWrongType", func(t *testing.T) {
		script := `
a1 = json.parse({"a": "b"})
`
		vars, err := starlark.ExecFile(thread, "script", script, predeclared)
		a.Error(err)
		a.NotNil(vars)
	})
	t.Run("parseUnmarshalError", func(t *testing.T) {
		Unmarshal = func(data []byte, v interface{}) error {
			return errors.New("some err")
		}
		script := `
a1 = json.parse('{"a": "b"}')
`
		vars, err := starlark.ExecFile(thread, "script", script, predeclared)
		a.Error(err)
		a.NotNil(vars)
	})
}
