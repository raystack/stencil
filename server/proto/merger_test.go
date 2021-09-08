package proto_test

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"path/filepath"
	"testing"

	stencilProto "github.com/odpf/stencil/server/proto"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestMerge(t *testing.T) {

	t.Run("Run all cases", func(t *testing.T) {
		for _, test := range []struct {
			number      int
			description string
		}{
			{1, "should able to merge on simple new fields"},
			{2, "should able to merge on simple rename of fields"},
			{3, "should able to merge new field in nested message"},
			{4, "should able to merge on new field of nested message type"},
			{5, "should able to merge on new field of message type"},
			{6, "should able to merge on new field with dependency on well known types"},
			{7, "should able to merge on new field with dependency on other file"},
			{8, "should able to merge on addition of package name, options, imports"},
		} {
			t.Run(test.description, func(t *testing.T) {
				runTest(t, test.number)
			})
		}
	})
}

func runTest(t *testing.T, testNumber int) {
	previous, current, expected := getTestData(t, testNumber)
	expectedFDS := &descriptorpb.FileDescriptorSet{}
	err := proto.Unmarshal(expected, expectedFDS)
	assert.Nil(t, err)
	expectedFiles, err := protodesc.NewFiles(expectedFDS)
	assert.Nil(t, err)
	actual, err := stencilProto.Merge(current, previous, []string{})
	assert.Nil(t, err)
	actualFDS := &descriptorpb.FileDescriptorSet{}
	err = proto.Unmarshal(actual, actualFDS)
	assert.Nil(t, err)
	actualFiles, err := protodesc.NewFiles(actualFDS)
	assert.Nil(t, err)
	assert.Equal(t, expectedFiles.NumFiles(), actualFiles.NumFiles())

	actualFiles.RangeFiles(func(actualFD protoreflect.FileDescriptor) bool {
		expectedFiles.RangeFiles(func(expectedFD protoreflect.FileDescriptor) bool {
			if actualFD.Path() == expectedFD.Path() {
				assert.Equal(t, expectedFD.FullName(), actualFD.FullName())
				assert.Equal(t, expectedFD.Package(), actualFD.Package())
				assert.Equal(t, expectedFD.Options(), actualFD.Options())
				for i := 0; i<expectedFD.Imports().Len(); i++ {
					missingImport := true
					for j := 0; j< actualFD.Imports().Len(); j++ {
						if expectedFD.Imports().Get(i) == actualFD.Imports().Get(j) {
							missingImport = false
						}
					}
					assert.True(t, missingImport)
				}

				actualMDs := actualFD.Messages()
				expectedMDs := expectedFD.Messages()
				for i := 0; i < expectedMDs.Len(); i++ {
					expectedMD := expectedMDs.Get(i)
					actualMD := actualMDs.ByName(expectedMD.Name())
					assert.NotNil(t, actualMD)

					assertDescriptors(t, expectedMD, actualMD)
				}
			}
			return true
		})
		return true
	})
}

func assertDescriptors(t *testing.T, expected, actual protoreflect.MessageDescriptor) {
	if expected.Messages() != nil {
		expectedNestedMessages := expected.Messages()
		actualNestedMessages := actual.Messages()

		for i := 0; i < expectedNestedMessages.Len(); i++ {
			expectedNestedMsg := expectedNestedMessages.Get(i)
			actualNestedMsg := actualNestedMessages.ByName(expectedNestedMsg.Name())
			assert.NotNil(t, actualNestedMsg)

			assertDescriptors(t, expectedNestedMsg, actualNestedMsg)
		}
	}

	assert.Equal(t, expected.Name(), actual.Name())
	assert.Equal(t, expected.FullName(), actual.FullName())
	assert.Equal(t, expected.Options(), actual.Options())
	expectedFieldDescs := expected.Fields()
	actualFieldDescs := actual.Fields()
	assert.Equal(t, expectedFieldDescs.Len(), actualFieldDescs.Len())

	for j := 0; j < expectedFieldDescs.Len(); j++ {
		expectedFieldDesc := expectedFieldDescs.Get(j)
		actualFieldDesc := actualFieldDescs.Get(j)
		assert.Equal(t, expectedFieldDesc.FullName(), actualFieldDesc.FullName())
		assert.Equal(t, expectedFieldDesc.Number(), actualFieldDesc.Number())
		assert.Equal(t, expectedFieldDesc.Cardinality(), actualFieldDesc.Cardinality())
		assert.Equal(t, expectedFieldDesc.Kind(), actualFieldDesc.Kind())
		assert.Equal(t, expectedFieldDesc.JSONName(), actualFieldDesc.JSONName())
		assert.Equal(t, expectedFieldDesc.Options(), actualFieldDesc.Options())
	}
}

func getTestData(t *testing.T, testNumber int) ([]byte, []byte, []byte) {
	existingRoot, _ := filepath.Abs(fmt.Sprintf("./testdata/merge/existing/%d", testNumber))
	newRoot, _ := filepath.Abs(fmt.Sprintf("./testdata/merge/new/%d", testNumber))
	expectedRoot, _ := filepath.Abs(fmt.Sprintf("./testdata/merge/expected/%d", testNumber))
	currentFileName := filepath.Join(t.TempDir(), "existing.desc")
	prevFileName := filepath.Join(t.TempDir(), "new.desc")
	expectedFileName := filepath.Join(t.TempDir(), "expected.desc")
	err := runProtoc(existingRoot, true, currentFileName)
	assert.NoError(t, err)
	err = runProtoc(newRoot, true, prevFileName)
	assert.NoError(t, err)
	err = runProtoc(expectedRoot, true, expectedFileName)
	assert.NoError(t, err)
	current, _ := ioutil.ReadFile(currentFileName)
	prev, _ := ioutil.ReadFile(prevFileName)
	expected, _ := ioutil.ReadFile(expectedFileName)
	return current, prev, expected
}
