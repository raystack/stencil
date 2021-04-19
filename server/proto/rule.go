package proto

import (
	"google.golang.org/protobuf/reflect/protoregistry"
)

type CheckFunc func(current, prev *protoregistry.Files) error

type Rule interface {
	ID() string
	Check(*protoregistry.Files, *protoregistry.Files) error
}

type rule struct {
	id          string
	Description string
	CheckFn     CheckFunc
}

func (r *rule) ID() string {
	return r.id
}

func (r *rule) Check(current, prev *protoregistry.Files) error {
	return r.CheckFn(current, prev)
}

func NewRule(id, description string, check CheckFunc) *rule {
	return &rule{id, description, check}
}
