package protobuf_test

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"

	"github.com/goto/stencil/formats/protobuf"
	"github.com/goto/stencil/test_helper"

	"github.com/stretchr/testify/assert"
)

func getDescriptorData(t *testing.T, path string, includeImports bool) []byte {
	t.Helper()
	root, _ := filepath.Abs(path)
	log.Println(t.Name())
	targetFile := filepath.Join(t.TempDir(), test_helper.GetRandomName())
	err := test_helper.RunProtoc(root, includeImports, targetFile, nil)
	assert.NoError(t, err)
	data, err := ioutil.ReadFile(targetFile)
	assert.NoError(t, err)
	return data
}

func TestGetParsedSchema(t *testing.T) {
	t.Run("should return error if protobuf data is not valid", func(t *testing.T) {
		data := []byte("invalid data")
		_, err := protobuf.GetParsedSchema(data)
		assert.Error(t, err)
	})
	t.Run("should be able to parse valid protobuf data", func(t *testing.T) {
		data := getDescriptorData(t, "./testdata/valid", true)
		parsedSchema, err := protobuf.GetParsedSchema(data)
		assert.NoError(t, err)
		assert.NotNil(t, parsedSchema)
	})
	t.Run("should return error if protobuf data is not fully contained", func(t *testing.T) {
		data := getDescriptorData(t, "./testdata/valid", false)
		_, err := protobuf.GetParsedSchema(data)
		assert.Error(t, err)
	})
}
