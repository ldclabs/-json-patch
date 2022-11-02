package jsonpatch

import (
	"strconv"
	"strings"
)

// Diff two JSON documents and generate a JSON Patch.
func Diff(src, dst []byte, opts *DiffOptions) (Patch, error) {
	return NewNode(src).Diff(NewNode(dst), opts)
}

// DiffOptions is used to customize the behavior of the Diff function.
type DiffOptions struct {
	// IDKey is the name of the key to use as the unique identifier for JSON object
	IDKey string
}

type collector struct {
	path  string
	patch Patch
}

func (c *collector) withPathToken(token string) string {
	if token == "" {
		return c.path
	}
	return c.path + "/" + token
}

func (c *collector) pushPathToken(token string) {
	c.path = c.withPathToken(token)
}

func (c *collector) popPathToken() {
	if i := strings.LastIndex(c.path, "/"); i >= 0 {
		c.path = c.path[:i]
	}
}

func (c *collector) replaceOp(token string, node *Node) error {
	raw, err := node.MarshalJSON()
	if err == nil {
		c.patch = append(c.patch, Operation{Op: "replace", Path: c.withPathToken(token), Value: raw})
	}
	return err
}

func (c *collector) addOp(token string, node *Node) error {
	raw, err := node.MarshalJSON()
	if err == nil {
		c.patch = append(c.patch, Operation{Op: "add", Path: c.withPathToken(token), Value: raw})
	}
	return err
}

func (c *collector) removeOp(token string) {
	c.patch = append(c.patch, Operation{Op: "remove", Path: c.withPathToken(token)})
}

// Diff two JSON nodes and generate a JSON Patch.
func (n *Node) Diff(target *Node, opts *DiffOptions) (Patch, error) {
	c := &collector{patch: make(Patch, 0)}
	if err := n.diff(target, c, opts); err != nil {
		return nil, err
	}
	return c.patch, nil
}

func (n *Node) diff(target *Node, c *collector, opts *DiffOptions) error {
	if n == nil || target == nil {
		return c.replaceOp("", target)
	}

	if n.Equal(target) {
		return nil
	}

	if target.which != n.which || target.which == eOther {
		return c.replaceOp("", target)
	}

	if n.which == eDoc {
		if opts != nil && opts.IDKey != "" {
			if v := n.doc.obj[opts.IDKey]; !v.isNull() && !v.Equal(target.doc.obj[opts.IDKey]) {
				return c.replaceOp("", target)
			}
		}

		for _, key := range n.doc.keys {
			if _, ok := target.doc.obj[key]; !ok {
				c.removeOp(encodePatchKey(key))
			}
		}

		for _, key := range target.doc.keys {
			node, ok := n.doc.obj[key]
			switch {
			case ok:
				c.pushPathToken(encodePatchKey(key))
				if err := node.diff(target.doc.obj[key], c, opts); err != nil {
					return err
				}
				c.popPathToken()

			default:
				if err := c.addOp(encodePatchKey(key), target.doc.obj[key]); err != nil {
					return err
				}
			}
		}

		return nil
	}

	nl := len(n.ary)
	for i, node := range target.ary {
		switch {
		case i < nl:
			c.pushPathToken(strconv.Itoa(i))
			if err := n.ary[i].diff(node, c, opts); err != nil {
				return err
			}
			c.popPathToken()

		default:
			if err := c.addOp(strconv.Itoa(i), node); err != nil {
				return err
			}
		}
	}

	for i := len(target.ary); i < nl; i++ {
		c.removeOp(strconv.Itoa(i))
	}

	return nil
}
