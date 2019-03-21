package builtin

import (
	"go.starlark.net/starlark"
)

type Function func(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error)
