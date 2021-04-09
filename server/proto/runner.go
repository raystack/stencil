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
	c := make(chan error, len(Rules))
	for _, rule := range Rules {
		r := rule
		go func() {
			ruleErr := r.Check(currentRegistry, previousRegistry)
			c <- ruleErr
		}()
	}
	for i := 0; i < len(Rules); i++ {
		err = multierr.Combine(err, <-c)
	}
	return err
}
