package random

import (
	"testing"

	"github.com/DLag/starlight/convert"

	"go.starlark.net/resolve"
	"go.starlark.net/starlark"

	"github.com/stretchr/testify/assert"
)

func TestRandom(t *testing.T) {
	a := assert.New(t)
	resolve.AllowFloat = true
	Seed = func(seed int64) {
		a.NotZero(seed)
	}
	Float64 = func() float64 {
		return 0.234
	}
	Int63n = func(n int64) int64 {
		a.Equal(int64(35), n)
		return 21
	}
	predeclared := starlark.StringDict{
		"random": New(),
	}
	thread := &starlark.Thread{}
	t.Run("HappyPath", func(t *testing.T) {
		script := `
random.seed()
a1 = random.randint(5, 40)
a2 = random.randint(42, 7)
a3 = random.random()
a4 = random.uniform(2.0, 200.0)
a5 = random.uniform(201.0, 3.0)
`
		vars, err := starlark.ExecFile(thread, "script", script, predeclared)
		a.NotNil(vars)
		a.NoError(err)

		m := convert.FromStringDict(vars)
		a.NotNil(m)
		a.Equal(int64(26), m["a1"])
		a.Equal(int64(28), m["a2"])
		a.Equal(0.234, m["a3"])
		a.Equal(0.234*(200.0-2.0)+2.0, m["a4"])
		a.Equal(0.234*(201.0-3.0)+3.0, m["a5"])
	})
	t.Run("randint error", func(t *testing.T) {
		script := `
random.seed()
a1 = random.randint(5, 40.0)
`
		vars, err := starlark.ExecFile(thread, "script", script, predeclared)
		a.NotNil(vars)
		a.Error(err)
	})
	t.Run("uniform error", func(t *testing.T) {
		script := `
random.seed()
a1 = random.uniform(5, 40.0)
`
		vars, err := starlark.ExecFile(thread, "script", script, predeclared)
		a.NotNil(vars)
		a.Error(err)
	})
}
