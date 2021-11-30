package schema

import "go.uber.org/multierr"

type ValidationStrategy func(ParsedSchema, ParsedSchema) error
type CompatibilityFn func(ParsedSchema, []ParsedSchema) error

func validateLatest(strategy ValidationStrategy) CompatibilityFn {
	return func(ps1 ParsedSchema, ps2 []ParsedSchema) error {
		for _, prev := range ps2 {
			return strategy(ps1, prev)
		}
		return nil
	}
}

func validateAll(strategy ValidationStrategy) CompatibilityFn {
	return func(ps1 ParsedSchema, ps2 []ParsedSchema) error {
		var err error
		for _, prev := range ps2 {
			e := strategy(ps1, prev)
			err = multierr.Combine(err, e)
		}
		return err
	}
}

func backwardStrategy(current, prev ParsedSchema) error {
	return current.IsBackwardCompatible(prev)
}

func forwardStrategy(current, prev ParsedSchema) error {
	return current.IsForwardCompatible(prev)
}

func fullStrategy(current, prev ParsedSchema) error {
	return current.IsFullCompatible(prev)
}

func defaultCompatibilityFn(current ParsedSchema, prevs []ParsedSchema) error {
	return nil
}

func getCompatibilityChecker(compatibility string) CompatibilityFn {
	switch compatibility {
	case "COMPATIBILITY_BACKWARD":
		return validateLatest(backwardStrategy)
	case "COMPATIBILITY_BACKWARD_TRANSITIVE":
		return validateAll(backwardStrategy)
	case "COMPATIBILITY_FORWARD":
		return validateLatest(forwardStrategy)
	case "COMPATIBILITY_FORWARD_TRANSITIVE":
		return validateAll(forwardStrategy)
	case "COMPATIBILITY_FULL":
		return validateLatest(fullStrategy)
	case "COMPATIBILITY_FULL_TRANSITIVE":
		return validateAll(fullStrategy)
	default:
		return defaultCompatibilityFn
	}
}
