// (c) 2022-2022, LDC Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package jsonpatch

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollector(t *testing.T) {
	assert := assert.New(t)

	c := &collector{patch: make(Patch, 0)}
	assert.Equal("/abc", c.withPathToken(encodePatchKey("abc")))
	assert.Equal("/a~0c", c.withPathToken(encodePatchKey("a~c")))
	assert.Equal("/a~1c", c.withPathToken(encodePatchKey("a/c")))
	assert.Equal("/0", c.withPathToken(strconv.Itoa(0)))
	assert.Equal("/99", c.withPathToken(strconv.Itoa(99)))

	c.pushPathToken(encodePatchKey("list"))
	assert.Equal("/list", c.path)
	c.pushPathToken(strconv.Itoa(1))
	assert.Equal("/list/1", c.path)

	assert.Equal("/list/1/abc", c.withPathToken(encodePatchKey("abc")))
	assert.Equal("/list/1/a~0c", c.withPathToken(encodePatchKey("a~c")))
	assert.Equal("/list/1/a~1c", c.withPathToken(encodePatchKey("a/c")))
	assert.Equal("/list/1/0", c.withPathToken(strconv.Itoa(0)))
	assert.Equal("/list/1/99", c.withPathToken(strconv.Itoa(99)))

	c.pushPathToken(encodePatchKey("a/c"))
	assert.Equal("/list/1/a~1c", c.path)
	c.popPathToken()
	assert.Equal("/list/1", c.path)
	c.popPathToken()
	assert.Equal("/list", c.path)
	c.popPathToken()
	assert.Equal("", c.path)

	c.replaceOp(NewNode([]byte(`{}`)))
	assert.Equal(1, len(c.patch))
	assert.Equal(Operation{Op: "replace", Path: "", Value: []byte(`{}`)}, c.patch[0])

	c.addOp(encodePatchKey("a/c"), NewNode([]byte(`"abc"`)))
	assert.Equal(2, len(c.patch))
	assert.Equal(Operation{Op: "add", Path: "/a~1c", Value: []byte(`"abc"`)}, c.patch[1])

	c.removeOp(encodePatchKey("a/c"))
	assert.Equal(3, len(c.patch))
	assert.Equal(Operation{Op: "remove", Path: "/a~1c"}, c.patch[2])
}

func TestDiff(t *testing.T) {
	assert := assert.New(t)

	for i, c := range Cases {

		patch, err := Diff([]byte(c.doc), []byte(c.result))
		if !assert.NoErrorf(err, "Failed to diff at case %d\nSrc: %s\nDst: %s\n",
			i, reformatJSON(c.doc), reformatJSON(c.result)) {
			continue
		}

		out, err := patch.Apply([]byte(c.doc))
		if !assert.NoErrorf(err, "Failed to apply patch at case %d\nSrc: %s\nDst: %s\nPatch:%s\n",
			i, reformatJSON(c.doc), reformatJSON(c.result), mustJSONString(patch)) {
			continue
		}

		assert.Truef(compareJSON(string(out), c.result), "Not equal at case %d\nSrc: %s\nDst: %s\nOut:%s\nPatch:%s\n",
			i, reformatJSON(c.doc), reformatJSON(c.result), reformatJSON(string(out)), mustJSONString(patch))
	}
}
