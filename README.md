# JSON-Patch
[![ci](https://github.com/ldclabs/json-patch/actions/workflows/ci.yml/badge.svg)](https://github.com/ldclabs/json-patch/actions?query=workflow%3Aci)
[![codecov](https://codecov.io/gh/ldclabs/json-patch/branch/main/graph/badge.svg?token=2G1SE83FY5)](https://codecov.io/gh/ldclabs/json-patch)

`jsonpatch` is a library which provides functionality for applying
[RFC6902 JSON patches](https://datatracker.ietf.org/doc/html/rfc6902) on JSON.

## Documentation

[Go-Documentation](https://pkg.go.dev/github.com/ldclabs/json-patch)

## Import

```go
// package jsonpatch
import "github.com/ldclabs/json-patch"
```

## Examples

### Create and apply a JSON Patch

```go
package main

import (
	"fmt"

	jsonpatch "github.com/ldclabs/json-patch"
)

func main() {
	original := []byte(`{"name": "John", "age": 24, "height": 3.21}`)
	patchDoc := []byte(`[
		{"op": "replace", "path": "/name", "value": "Jane"},
		{"op": "remove", "path": "/height"}
	]`)

	patch, err := jsonpatch.NewPatch(patchDoc)
	if err != nil {
		panic(err)
	}
	modified, err := patch.Apply(original)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", modified)
	// {"name":"Jane","age":24}
}
```

### Create a JSON Patch from Diff

```go
package main

import (
	"fmt"

	jsonpatch "github.com/ldclabs/json-patch"
)

func main() {
	original := []byte(`{"name": "John", "age": 24, "height": 3.21}`)
	target := []byte(`{"name":"Jane","age":24}`)

	patch, err := Diff(original, target, nil)
	if err != nil {
		panic(err)
	}
	patchDoc, err := json.Marshal(patch)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", patchDoc)
	// [{"op":"remove","path":"/height"},{"op":"replace","path":"/name","value":"Jane"}]

	patch, err = Diff(original, target, &jsonpatch.DiffOptions{IDKey: "name"})
	if err != nil {
		panic(err)
	}
	patchDoc, err = json.Marshal(patch)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", patchDoc)
	// [{"op":"replace","path":"","value":{"name":"Jane","age":24}}]
}
```

### Create a Node and apply more Patchs

```go
package main

import (
	"fmt"

	jsonpatch "github.com/ldclabs/json-patch"
)

func main() {
	original := []byte(`{"name": "John", "age": 24, "height": 3.21}`)
	node := jsonpatch.NewNode(original)

	patch, err := jsonpatch.NewPatch([]byte(`[
		{"op": "replace", "path": "/name", "value": "Jane"},
		{"op": "remove", "path": "/height"}
	]`))
	if err != nil {
		panic(err)
	}
	err = node.Patch(patch, nil)
	if err != nil {
		panic(err)
	}

	modified, err := node.MarshalJSON()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", modified)
	// {"name":"Jane","age":24}

	patch, err = jsonpatch.NewPatch([]byte(`[
		{"op": "replace", "path": "/age", "value": 25}
	]`))
	if err != nil {
		panic(err)
	}
	err = node.Patch(patch, nil)
	if err != nil {
		panic(err)
	}

	modified, err = node.MarshalJSON()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", modified)
	// {"name":"Jane","age":25}
}
```

### Get value by path

```go
package main

import (
	"fmt"

	jsonpatch "github.com/ldclabs/json-patch"
)

func main() {
	doc := []byte(`{
		"baz": "qux",
		"foo": [ "a", 2, "c" ]
	}`)
	node := jsonpatch.NewNode(doc)
	value, err := node.GetValue("/foo/0", nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", string(value))
	// "a"
}
```

### Find children by test operations

```go
package main

import (
	"fmt"

	jsonpatch "github.com/ldclabs/json-patch"
)

func main() {
	doc := []byte(`["root", ["p",
		["span", {"data-type": "text"},
			["span", {"data-type": "leaf"}, "Hello 1"],
			["span", {"data-type": "leaf"}, "Hello 2"],
			["span", {"data-type": "leaf"}, "Hello 3"],
			["span", {"data-type": null}, "Hello 4"]
		]
	]]`)

	node := jsonpatch.NewNode(doc)
	tests := jsonpatch.PVs{
		{"/0", []byte(`"span"`)},
		{"/1/data-type", []byte(`"leaf"`)},
	}

	result, err := node.FindChildren(tests, nil)
	if err != nil {
		panic(err)
	}
	for _, r := range result {
		fmt.Printf("Path: \"%s\", Value: %s\n", r.Path, string(r.Value))
	}
	// Path: "/1/1/2", Value: ["span", {"data-type": "leaf"}, "Hello 1"]
	// Path: "/1/1/3", Value: ["span", {"data-type": "leaf"}, "Hello 2"]
	// Path: "/1/1/4", Value: ["span", {"data-type": "leaf"}, "Hello 3"]
}
```