package json

import (
	"fmt"
	"strings"
)

type diffKind int

type diff struct {
	kind diffKind
	msg  string
}

type compatibilityErr struct {
	notAllowed []diffKind
	diffs      []diff
}

func (d diffKind) contains(others []diffKind) bool {
	for _, v := range others {
		if v == d {
			return true
		}
	}
	return false
}

func (c *compatibilityErr) add(kind diffKind, location string, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if kind.contains(c.notAllowed) && msg != "" {
		c.diffs = append(c.diffs, diff{kind: kind, msg: fmt.Sprintf("%s: %s", location, msg)})
	}
}

func (c *compatibilityErr) isEmpty() bool {
	return len(c.diffs) == 0
}

func (c *compatibilityErr) Error() string {
	var msgs []string
	for _, val := range c.diffs {
		msgs = append(msgs, val.msg)
	}
	return strings.Join(msgs, ";")
}
