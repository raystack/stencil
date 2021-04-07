package proto

import (
	"go.uber.org/multierr"
	"google.golang.org/protobuf/reflect/protoregistry"
)

func Compare(current, prev []byte) error {
	var err error
	var currentRegistry, previousRegistry *protoregistry.Files
	if currentRegistry, err = getRegistry(current); err != nil {
		return err
	}
	if previousRegistry, err = getRegistry(prev); err != nil {
		return err
	}
	for _, rule := range Rules {
		ruleErr := rule.Check(currentRegistry, previousRegistry)
		err = multierr.Combine(err, ruleErr)
	}
	return err
}
