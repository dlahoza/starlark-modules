package json

import (
	"encoding/json"
	"fmt"

	"github.com/DLag/starlark-modules/builtin"

	"github.com/DLag/starlight/convert"
	"go.starlark.net/starlark"
)

func New() starlark.Value {
	return builtin.New(map[string]builtin.Function{
		"parse": jsonParse,
		"dump":  jsonDump,
	})
}

var (
	Marshal   = json.Marshal
	Unmarshal = json.Unmarshal
)

func jsonParse(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, _ []starlark.Tuple) (starlark.Value, error) {
	fname := "json.parse"
	if args.Len() != 1 || args.Index(0).Type() != "string" {
		return starlark.None, fmt.Errorf("wrong args, should be %s(string)", fname)
	}
	buf := args.Index(0).(starlark.String).GoString()
	var v map[string]interface{}
	err := Unmarshal([]byte(buf), &v)
	if err != nil {
		return starlark.None, err
	}
	return convert.MakeDict(v)
}

func jsonDump(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, _ []starlark.Tuple) (starlark.Value, error) {
	fname := "json.dump"
	if args.Len() != 1 || args.Index(0).Type() != "dict" {
		return starlark.None, fmt.Errorf("wrong args, should be %s(dict)", fname)
	}
	d := builtin.ConvertToStringMap(convert.FromDict(args.Index(0).(*starlark.Dict)))
	buf, err := Marshal(d)
	if err != nil {
		return starlark.None, err
	}
	return convert.ToValue(string(buf))
}
