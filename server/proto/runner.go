package proto

import (
	"go.uber.org/multierr"
	"google.golang.org/protobuf/reflect/protoregistry"
)

func contains(list []string, elem string) bool {
	for _, v := range list {
		if v == elem {
			return true
		}
	}
	return false
}

func Compare(current, prev []byte, rulesToSkip []string) error {
	var err error
	var currentRegistry, previousRegistry *protoregistry.Files
	if currentRegistry, err = getRegistry(current); err != nil {
		return err
	}
	if previousRegistry, err = getRegistry(prev); err != nil {
		return err
	}
	var filteredRules []Rule
	for _, rule := range Rules {
		if !contains(rulesToSkip, rule.ID()) {
			filteredRules = append(filteredRules, rule)
		}
	}
	c := make(chan error, len(filteredRules))
	for _, rule := range filteredRules {
		r := rule
		go func() {
			ruleErr := r.Check(currentRegistry, previousRegistry)
			c <- ruleErr
		}()
	}
	for i := 0; i < len(filteredRules); i++ {
		err = multierr.Combine(err, <-c)
	}
	return err
}
