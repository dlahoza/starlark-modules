package random

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/DLag/starlark-modules/builtin"

	"github.com/starlight-go/starlight/convert"
	"go.starlark.net/starlark"
)

func New() starlark.Value {
	return builtin.New(map[string]builtin.Function{
		"seed":    randSeed,
		"randint": randInt,
		"random":  random,
		"uniform": uniform,
	})
}

func randSeed(_ *starlark.Thread, _ *starlark.Builtin, _ starlark.Tuple, _ []starlark.Tuple) (starlark.Value, error) {
	rand.Seed(time.Now().UnixNano())
	return starlark.None, nil
}

func random(_ *starlark.Thread, _ *starlark.Builtin, _ starlark.Tuple, _ []starlark.Tuple) (starlark.Value, error) {
	return convert.ToValue(rand.Float64())
}

func uniform(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, _ []starlark.Tuple) (starlark.Value, error) {
	fname := "random.uniform"
	if args.Len() != 2 || args.Index(0).Type() != "float" || args.Index(1).Type() != "float" {
		return starlark.None, fmt.Errorf("wrong args, should be %s(int, int)", fname)
	}
	var min, max starlark.Float
	a, _ := args.Index(0).(starlark.Float)
	b, _ := args.Index(1).(starlark.Float)
	if a < b {
		min = a
		max = b
	} else {
		min = b
		max = a
	}

	max = max - min
	return convert.ToValue(starlark.Float(rand.Float64())*max + min)
}

func randInt(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, _ []starlark.Tuple) (starlark.Value, error) {
	fname := "random.randint"
	if args.Len() != 2 || args.Index(0).Type() != "int" || args.Index(1).Type() != "int" {
		return starlark.None, fmt.Errorf("wrong args, should be %s(int, int)", fname)
	}
	var min, max int64
	a, _ := args.Index(0).(starlark.Int).Int64()
	b, _ := args.Index(1).(starlark.Int).Int64()
	if a <= b {
		min = a
		max = b
	} else {
		min = b
		max = a
	}

	max = max - min
	return convert.ToValue(rand.Int63n(max) + min)
}
