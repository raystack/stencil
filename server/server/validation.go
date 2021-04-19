package server

import (
	"github.com/blang/semver/v4"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

//ValidateVersion validates if version is semantic version compatible or not
func ValidateVersion(fl validator.FieldLevel) bool {
	version, ok := fl.Field().Interface().(string)
	if ok {
		if _, err := semver.Parse(version); err == nil {
			return true
		}
		return false
	}
	return false
}

//ValidateVersionWithLatest validates if value is equal to latest or valid semantic version
func ValidateVersionWithLatest(fl validator.FieldLevel) bool {
	version, ok := fl.Field().Interface().(string)
	if ok {
		if _, err := semver.Parse(version); err == nil {
			return true
		}
		return version == "latest"
	}
	return false
}

func registerCustomValidations(e *gin.Engine) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("version", ValidateVersion)
		v.RegisterValidation("versionWithLatest", ValidateVersionWithLatest)
	}
}
