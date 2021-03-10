package models

import "fmt"

var (
	ErrMissingFormData = APIError{
		Code:    400,
		Message: "Missing fields in input data",
	}
	ErrUploadFailed = APIError{
		Code:    500,
		Message: "Upload failed",
	}
	ErrDownloadFailed = APIError{
		Code:    500,
		Message: "Download failed",
	}
)

type APIError struct {
	Code    int
	Message string
	error
}

func (a APIError) Error() string {
	return fmt.Sprintf("%d %s %s", a.Code, a.Message, a.error.Error())
}
