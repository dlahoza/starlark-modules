package convert

import (
	"sort"
	"testing"

	"github.com/DLag/starlark-modules/builtin"

	"github.com/stretchr/testify/assert"

	"github.com/starlight-go/starlight/convert"
	"go.starlark.net/starlark"
)

type testWrapper struct {
	testWrapperEmbeded
	Var1 string `starlark:"var1"`
	Var2 int64  `starlark:"var2"`
}

type testWrapperEmbeded struct {
	EmbVar1 string `starlark:"embvar1"`
}

func (t *testWrapper) func1(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, _ []starlark.Tuple) (starlark.Value, error) {
	t.Var1 += "func1text_"
	return starlark.None, nil
}

func (t *testWrapper) func2(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, _ []starlark.Tuple) (starlark.Value, error) {
	t.Var2++
	return starlark.None, nil
}

func (t *testWrapper) Methods() map[string]builtin.Function {
	return map[string]builtin.Function{
		"func1": t.func1,
		"func2": t.func2,
	}
}

func TestStarlarkStructWrapper(t *testing.T) {
	a := assert.New(t)
	thread := &starlark.Thread{}
	t.Run("HappyPath", func(t *testing.T) {
		v := new(testWrapper)
		v.Var1 = "start_"
		v.Var2 = 5
		predeclared := starlark.StringDict{
			"test": NewStruct(v),
		}
		script := `
k = dir(test)
i1 = test.var1
i2 = test.var2
i3 = test.embvar1
test.var1+="scripttext_"
test.func1()
test.var2-=3
test.func2()
a1 = test.var1
a2 = test.var2
`
		vars, err := starlark.ExecFile(thread, "script", script, predeclared)
		a.NotNil(vars)
		a.NoError(err)

		m := convert.FromStringDict(vars)
		a.NotNil(m)
		a.Equal("start_scripttext_func1text_", v.Var1)
		a.Equal(int64(3), v.Var2)
		a.Equal("start_", m["i1"])
		a.Equal(int64(5), m["i2"])
		a.Equal("start_scripttext_func1text_", m["a1"])
		a.Equal(int64(3), m["a2"])
		l := make([]string, 0)
		for _, val := range m["k"].([]interface{}) {
			l = append(l, val.(string))
		}
		expected := []string{"embvar1", "var1", "var2", "func1", "func2"}
		sort.Strings(l)
		sort.Strings(expected)
		a.Equal(expected, l)
	})
	t.Run("Unknown field", func(t *testing.T) {
		v := new(testWrapper)
		v.Var1 = "start_"
		v.Var2 = 5
		predeclared := starlark.StringDict{
			"test": NewStruct(v),
		}
		script := `
i1 = test.var1
i2 = test.var2
test.var1+="scripttext_"
test.func1()
test.var2-=3
test.func2()
test.var3 = 'abc'
a1 = test.var1
a2 = test.var2
`
		_, mod, err := starlark.SourceProgram("script", script, predeclared.Has)
		a.NoError(err)
		a.NotNil(mod)

		vars, err := mod.Init(thread, predeclared)
		vars.Freeze()
		a.NotNil(vars)
		a.Error(err)

		m := convert.FromStringDict(vars)
		a.NotNil(m)
		a.Equal("start_scripttext_func1text_", v.Var1)
		a.Equal(int64(3), v.Var2)
		a.Equal("start_", m["i1"])
		a.Equal(int64(5), m["i2"])
		a.NotEqual("start_scripttext_func1text_", m["a1"])
		a.NotEqual(int64(3), m["a2"])
	})
}
