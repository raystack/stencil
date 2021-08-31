package api

import (
	"github.com/blang/semver/v4"
	v "github.com/go-playground/validator/v10"
)

var validate *v.Validate

func init() {
	validate = v.New()
	validate.SetTagName("binding")
	validate.RegisterValidation("version", func(fl v.FieldLevel) bool {
		version, ok := fl.Field().Interface().(string)
		if ok {
			if _, err := semver.Parse(version); err == nil {
				return true
			}
			return false
		}
		return false
	})
}
