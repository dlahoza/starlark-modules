# JSON Starlark module

## Usage

    predeclared := starlark.StringDict{
        "json": json.New()
        ...
    }
    starlark.ExecFile(thread, filename, nil, predeclared)

## Supported functions

    json.parse(string) dict

Parses JSON string to Starlark dict. Returns error and None value on parse error. Uses json.Unmarshal Go function.

    json.dump(dict) string

Serialize Starlark dict to JSON. Returns error and None value on parse error. Uses json.Marshal Go function.
