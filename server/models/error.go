package models

import (
	"errors"
	"fmt"
)

var (
	ErrMissingFormData = &apiErr{
		code:    400,
		message: "Missing fields in input data",
	}
	ErrUploadFailed = &apiErr{
		code:    500,
		message: "Upload failed",
	}
	ErrUploadInvalidFile = &apiErr{
		code:    400,
		message: "Unable to read uploaded file",
	}
	ErrDownloadFailed = &apiErr{
		code:    500,
		message: "Download failed",
	}
	ErrMetadataUpdateFailed = &apiErr{
		code:    500,
		message: "Metadata update failed",
	}
	ErrGetMetadataFailed = &apiErr{
		code:    500,
		message: "Unable to get metadata information",
	}
	ErrNotFound = &apiErr{
		code:    404,
		message: "Not found",
	}
	ErrConflict = &apiErr{
		code:    409,
		message: "Resource already exist",
	}
	ErrCancel = &apiErr{
		code:    500,
		message: "Operation was cancelled",
	}
	ErrTimeout = &apiErr{
		code:    408,
		message: "Operation was timedout",
	}
	ErrStoreInternal = &apiErr{
		code:    500,
		message: "Internal backend store error",
	}
	ErrUnknown = &apiErr{
		code:    500,
		message: "Internal server error",
	}
)

//APIError returns API response with provided code and error message
type APIError interface {
	Code() int
	Message() string
	error
}

type apiErr struct {
	code    int
	message string
	error
}

func (a *apiErr) Code() int {
	return a.code
}

func (a *apiErr) Message() string {
	if a.message == "" && a.error != nil {
		return a.error.Error()
	}
	return a.message
}

func (a *apiErr) Error() string {
	var err error
	err = a.error
	if err == nil {
		err = errors.New("")
	}
	return fmt.Sprintf("Err: %s, Message: %s", err.Error(), a.Message())
}

//NewAPIError helper function to contruct API error
func NewAPIError(code int, message string, err error) APIError { return &apiErr{code, message, err} }

func WrapAPIError(err APIError, rootErr error) APIError {
	return &apiErr{err.Code(), err.Message(), rootErr}
}
