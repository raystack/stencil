package service

import (
	"io"
	"io/ioutil"
	"mime/multipart"

	"github.com/odpf/stencil/server/models"
)

func readDataFromReader(reader io.ReadCloser) ([]byte, error) {
	data, err := ioutil.ReadAll(reader)
	defer func() {
		reader.Close()
	}()
	return data, err
}

func readDataFromMultiPartFile(file *multipart.FileHeader) ([]byte, error) {
	fileReader, err := file.Open()
	if err != nil {
		return nil, err
	}
	return readDataFromReader(fileReader)
}

func isNotFoundErr(err error) bool {
	val, ok := err.(models.APIError)
	if ok && val.Code() == 404 {
		return true
	}
	return false
}
