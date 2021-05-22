package proto

import (
	"fmt"
	"log"
	"runtime/debug"

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
			// Only the panic in the current calling process can be captured in golang. If it is another goroutine, it cannot catch exceptions.
			// So gin.Recovery() won't be able to catch panics generated from rule checks.
			defer func() {
				if e := recover(); e != nil {
					log.Println(e, debug.Stack())
					c <- fmt.Errorf("internal error: %s validation rule failed", r.ID())
				}
			}()
			ruleErr := r.Check(currentRegistry, previousRegistry)
			c <- ruleErr
		}()
	}
	for i := 0; i < len(filteredRules); i++ {
		err = multierr.Combine(err, <-c)
	}
	return err
}
