// (c) 2022-2022, LDC Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package jsonpatch

import (
	"encoding/json"
	"fmt"
)

func ExamplePatch_Apply() {
	original := []byte(`{"name": "John", "age": 24, "height": 3.21}`)
	patchDoc := []byte(`[
		{"op": "replace", "path": "/name", "value": "Jane"},
		{"op": "remove", "path": "/height"}
	]`)

	patch, err := NewPatch(patchDoc)
	if err != nil {
		panic(err)
	}
	modified, err := patch.Apply(original)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", modified)

	// Output:
	// {"name":"Jane","age":24}
}

func ExampleDiff() {
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

	patch, err = Diff(original, target, &DiffOptions{IDKey: "name"})
	if err != nil {
		panic(err)
	}
	patchDoc, err = json.Marshal(patch)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", patchDoc)
	// [{"op":"replace","path":"","value":{"name":"Jane","age":24}}]

	// Output:
	// [{"op":"remove","path":"/height"},{"op":"replace","path":"/name","value":"Jane"}]
	// [{"op":"replace","path":"","value":{"name":"Jane","age":24}}]
}

func ExampleNode_Patch() {
	original := []byte(`{"name": "John", "age": 24, "height": 3.21}`)
	patchDoc0 := []byte(`[
		{"op": "replace", "path": "/name", "value": "Jane"},
		{"op": "remove", "path": "/height"}
	]`)
	patchDoc1 := []byte(`[
		{"op": "replace", "path": "/age", "value": 25}
	]`)

	node := NewNode(original)
	patch, err := NewPatch(patchDoc0)
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

	patch, err = NewPatch(patchDoc1)
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

	// Output:
	// {"name":"Jane","age":24}
	// {"name":"Jane","age":25}
}

func ExampleNode_GetValue() {
	doc := []byte(`{
		"baz": "qux",
		"foo": [ "a", 2, "c" ]
	}`)
	node := NewNode(doc)

	value, err := node.GetValue("/foo/0", nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", string(value))

	// Output:
	// "a"
}

func ExampleNode_FindChildren() {
	doc := []byte(`["root", ["p",
		["span", {"data-type": "text"},
			["span", {"data-type": "leaf"}, "Hello 1"],
			["span", {"data-type": "leaf"}, "Hello 2"],
			["span", {"data-type": "leaf"}, "Hello 3"],
			["span", {"data-type": null}, "Hello 4"]
		]
	]]`)

	node := NewNode(doc)
	tests := PVs{
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

	// Output:
	// Path: "/1/1/2", Value: ["span", {"data-type": "leaf"}, "Hello 1"]
	// Path: "/1/1/3", Value: ["span", {"data-type": "leaf"}, "Hello 2"]
	// Path: "/1/1/4", Value: ["span", {"data-type": "leaf"}, "Hello 3"]
}
