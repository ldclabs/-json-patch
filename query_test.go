// (c) 2022-2022, LDC Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package jsonpatch

import (
	"strings"
	"testing"
)

type GetValueCase struct {
	doc, path string
	result    []byte
	err       string
}

var GetValueCases = []GetValueCase{
	{
		`{ "baz": "qux" }`,
		"/baz",
		[]byte(`"qux"`),
		"",
	},
	{
		`{
	    "baz": "qux",
	    "foo": [ "a", 2, "c" ]
	  }`,
		"/foo/0",
		[]byte(`"a"`),
		"",
	},
	{
		`{
	    "baz": "qux",
	    "foo": [ "a", 2, "c" ]
	  }`,
		"/foo/1",
		[]byte(`2`),
		"",
	},
	{
		`{
	    "baz": "qux",
	    "foo": [ "a", 2, "c", {"baz": null} ]
	  }`,
		"/foo/3/baz",
		[]byte(`null`),
		"",
	},
	{
		`{
	    "baz": "qux",
	    "foo": [ "a", 2, "c", {"baz": null}, null ]
	  }`,
		"/foo/4",
		[]byte(`null`),
		"",
	},
	{
		`{ "foo": {} }`,
		"/foo",
		[]byte(`{}`),
		"",
	},
	{
		`{ "foo": [ ] }`,
		"/foo",
		[]byte(`[]`),
		"",
	},
	{
		`{ "foo": null }`,
		"/foo",
		[]byte(`null`),
		"",
	},
	{
		`{ "baz/foo": "qux" }`,
		"/baz~1foo",
		[]byte(`"qux"`),
		"",
	},
	{
		`{
	    "baz": "qux",
	    "foo": [ "a", 2, "c" ]
	  }`,
		"/fooo",
		nil,
		`unable to get nonexistent key "fooo", missing value`,
	},
}

func TestGetValueByPath(t *testing.T) {
	for _, c := range GetValueCases {
		res, err := GetValueByPath([]byte(c.doc), c.path)
		if c.err != "" {
			if err == nil || !strings.Contains(err.Error(), c.err) {
				t.Errorf("Testing failed when it should have error for [%s]: expected [%s], got [%v]",
					string(c.doc), c.err, err)
			}
		} else if err != nil {
			t.Errorf("Testing failed when it should have passed for [%s]: %v", string(c.doc), err)
		} else {
			if string(res) != string(c.result) {
				t.Errorf("Testing failed for [%s]: expected [%s], got [%s]", string(c.doc), string(c.result), string(res))
			}
		}
	}
}

type FindChildrenCase struct {
	doc    []byte
	tests  []*PV
	result []*PV
}

var FindChildrenCases = []FindChildrenCase{
	{
		[]byte(`{ "baz": "qux" }`),
		[]*PV{{"/baz", []byte(`"qux"`)}},
		[]*PV{{"", []byte(`{"baz": "qux"}`)}},
	},
	{
		[]byte(`{
	    "baz": "qux",
	    "foo": [ "a", 2, "c" ]
	  }`),
		[]*PV{{"/foo/0", []byte(`"a"`)}},
		[]*PV{{"", []byte(`{
				"baz": "qux",
				"foo": [ "a", 2, "c" ]
			}`),
		}},
	},
	{
		[]byte(`{
	    "baz": "qux",
	    "foo": [ "a", 2, "c" ]
	  }`),
		[]*PV{{"/1", []byte(`2`)}},
		[]*PV{{"/foo", []byte(`[ "a", 2, "c" ]`)}},
	},
	{
		[]byte(`{
	    "baz": "qux",
	    "foo": [ "a", 2, "c" ]
	  }`),
		[]*PV{{"/fooo", nil}},
		[]*PV{},
	},
	{
		[]byte(`{ "foo": {} }`),
		[]*PV{{"/foo", []byte(`{}`)}},
		[]*PV{{"", []byte(`{ "foo": {} }`)}},
	},
	{
		[]byte(`{ "foo": [ ] }`),
		[]*PV{{"/foo", []byte(`[]`)}},
		[]*PV{{"", []byte(`{ "foo": [ ] }`)}},
	},
	{
		[]byte(`{ "foo": null }`),
		[]*PV{{"/foo", nil}},
		[]*PV{{"", []byte(`{ "foo": null }`)}},
	},
	{
		[]byte(`{ "foo": null }`),
		[]*PV{{"/foo", []byte("")}},
		[]*PV{{"", []byte(`{ "foo": null }`)}},
	},
	{
		[]byte(`{ "foo": null }`),
		[]*PV{{"/foo", []byte("null")}},
		[]*PV{{"", []byte(`{ "foo": null }`)}},
	},
	{
		[]byte(`{ "foo": "" }`),
		[]*PV{{"/foo", []byte(`""`)}},
		[]*PV{{"", []byte(`{ "foo": "" }`)}},
	},
	{
		[]byte(`{ "baz/foo": "qux" }`),
		[]*PV{{"/baz~1foo", []byte(`"qux"`)}},
		[]*PV{{"", []byte(`{ "baz/foo": "qux" }`)}},
	},
	{
		[]byte(`{ "baz/foo": [ "qux" ] }`),
		[]*PV{{"/0", []byte(`"qux"`)}},
		[]*PV{{"/baz~1foo", []byte(`["qux"]`)}},
	},
	{
		[]byte(`[
			"root",
			["object", { "id": "id1" }],
			["object", { "id": "id2" }]
		]`),
		[]*PV{{"/0", []byte(`"object"`)}},
		[]*PV{
			{"/1", []byte(`["object", { "id": "id1" }]`)},
			{"/2", []byte(`["object", { "id": "id2" }]`)},
		},
	},
	{
		[]byte(`[
			"root",
			["object", { "id": "id1" }],
			["object", { "id": "id2" }]
		]`),
		[]*PV{{"/1/id", []byte(`"id1"`)}},
		[]*PV{{"/1", []byte(`["object", { "id": "id1" }]`)}},
	},
	{
		[]byte(`[
			"root",
			["object", { "id": "id1" }],
			["object", { "id": "id2" }]
		]`),
		[]*PV{{"/1", []byte(`{ "id": "id1" }`)}},
		[]*PV{{"/1", []byte(`["object", { "id": "id1" }]`)}},
	},
	{
		[]byte(`[
			"root",
			["object", { "id": "" }],
			["object", { "id": null }]
		]`),
		[]*PV{{"/1/id", []byte(`""`)}},
		[]*PV{{"/1", []byte(`["object", { "id": "" }]`)}},
	},
	{
		[]byte(`[
			"root",
			["object", { "id": "" }],
			["object", { "id": null }]
		]`),
		[]*PV{{"/1/id", []byte(`null`)}},
		[]*PV{{"/2", []byte(`["object", { "id": null }]`)}},
	},
	{
		[]byte(`[
			"root",
			["object", { "id": "" }],
			["object", { "id": null }]
		]`),
		[]*PV{{"/1/id", []byte(`null`)}},
		[]*PV{{"/2", []byte(`["object", { "id": null }]`)}},
	},
	{
		[]byte(`[
			"root",
			["object", { "id": "" }],
			["object", { "id": null }]
		]`),
		[]*PV{{"/1/id", []byte(`""`)}},
		[]*PV{{"/1", []byte(`["object", { "id": "" }]`)}},
	},
	{
		[]byte(`[
			"root",
			["object1", { "id": "" }],
			["object2", { "id": null }]
		]`),
		[]*PV{
			{"/0", []byte(`"object2"`)},
			{"/1/id", []byte(`null`)},
		},
		[]*PV{{"/2", []byte(`["object2", { "id": null }]`)}},
	},
	{
		[]byte(`[
			"root",
			["object1", { "id": "" }],
			["object2", { "id": null }]
		]`),
		[]*PV{
			{"/0", []byte(`"root"`)},
			{"/1/0", []byte(`"object1"`)},
			{"/1/1/id", []byte(`""`)},
		},
		[]*PV{{"", []byte(`[
				"root",
				["object1", { "id": "" }],
				["object2", { "id": null }]
			]`)},
		},
	},
	{
		[]byte(`[
			"root",
			["object1", { "id": "" }],
			["object2", { "id": null }]
		]`),
		[]*PV{
			{"/0", []byte(`"root"`)},
			{"/1/0", []byte(`"object1"`)},
			{"/1/1/id", []byte(`""`)},
			{"/2", []byte(`["object2", { "id": null }]`)},
		},
		[]*PV{
			{"", []byte(`[
				"root",
				["object1", { "id": "" }],
				["object2", { "id": null }]
			]`)},
		},
	},
	{
		[]byte(`["root", ["p",
			["span", {"data-type": "text"},
				["span", {"data-type": "leaf"}, "Hello 1"],
				["span", {"data-type": "leaf"}, "Hello 2"],
				["span", {"data-type": "leaf"}, "Hello 3"],
				["span", {"data-type": null}, "Hello 4"]
			]
		]]`),
		[]*PV{{"/0", []byte(`"span"`)}, {"/1/data-type", []byte(`"leaf"`)}},
		[]*PV{
			{"/1/1/2", []byte(`["span", {"data-type": "leaf"}, "Hello 1"]`)},
			{"/1/1/3", []byte(`["span", {"data-type": "leaf"}, "Hello 2"]`)},
			{"/1/1/4", []byte(`["span", {"data-type": "leaf"}, "Hello 3"]`)},
		},
	},
	{
		[]byte(`["root", ["p",
			["span", {"data-type": "text"},
				["span", {"data-type": "leaf"}, "Hello 1"],
				["span", {"data-type": "leaf"}, "Hello 2"],
				["span", {"data-type": "leaf"}, "Hello 3"],
				["span", {"data-type": null}, "Hello 4"]
			]
		]]`),
		[]*PV{{"/0", []byte(`"span"`)}, {"/1/data-type", nil}},
		[]*PV{{"/1/1/5", []byte(`["span", {"data-type": null}, "Hello 4"]`)}},
	},
	{
		[]byte(`["root", ["p",
			["span", {"data-type": "text"},
				["span", {"data-type": "leaf"}, "Hello 1"],
				["span", {"data-type": "leaf"}, "Hello 2"],
				["span", {"data-type": "leaf"}, "Hello 3"],
				["span", {"data-type": null}, "Hello 4"]
			]
		]]`),
		[]*PV{{"/0", []byte(`"span"`)}},
		[]*PV{
			{"/1/1", []byte(`["span", {"data-type": "text"},
			["span", {"data-type": "leaf"}, "Hello 1"],
			["span", {"data-type": "leaf"}, "Hello 2"],
			["span", {"data-type": "leaf"}, "Hello 3"],
			["span", {"data-type": null}, "Hello 4"]]`)},
			{"/1/1/2", []byte(`["span", {"data-type": "leaf"}, "Hello 1"]`)},
			{"/1/1/3", []byte(`["span", {"data-type": "leaf"}, "Hello 2"]`)},
			{"/1/1/4", []byte(`["span", {"data-type": "leaf"}, "Hello 3"]`)},
			{"/1/1/5", []byte(`["span", {"data-type": null}, "Hello 4"]`)},
		},
	},
}

func TestFindChildren(t *testing.T) {
	for i, c := range FindChildrenCases {
		res, err := NewNode(c.doc).FindChildren(c.tests, nil)

		if err != nil {
			t.Errorf("Testing failed when case %d should have passed: %s", i, err)
		} else {
			if len(res) != len(c.result) {
				t.Errorf("Testing failed for case %d, %s: expected %d, got %d",
					i, string(c.doc), len(c.result), len(res))
			}
			for j := range res {
				if c.result[j].Path != res[j].Path {
					t.Errorf("Testing failed for case %d, %s: expected path [%s], got [%s]",
						i, string(c.doc), c.result[j].Path, res[j].Path)
				} else if !Equal(c.result[j].Value, res[j].Value) {
					t.Errorf("Testing failed for case %d, %v: expected [%s], got [%s]",
						i, string(c.doc), string(c.result[j].Value), string(res[j].Value))
				}
			}
		}
	}
}
