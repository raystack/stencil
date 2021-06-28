package proto_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/odpf/stencil/server/proto"
	"github.com/stretchr/testify/assert"
)

func runProtoc(
	rootDir string,
	includeImports bool,
	descSetOut string,
) error {
	protocBinPath, err := exec.LookPath("protoc")
	if err != nil {
		return err
	}
	protocBinPath, err = filepath.EvalSymlinks(protocBinPath)
	if err != nil {
		return err
	}
	protocBinPath, err = filepath.Abs(protocBinPath)
	if err != nil {
		return err
	}
	protocIncludePath, err := filepath.Abs(filepath.Join(filepath.Dir(protocBinPath), "..", "include"))
	if err != nil {
		return err
	}
	args := []string{"-I", rootDir, "-I", protocIncludePath}
	args = append(args, fmt.Sprintf("--descriptor_set_out=%s", descSetOut))
	if includeImports {
		args = append(args, "--include_imports")
	}
	protoFiles, _ := filepath.Glob(filepath.Join(rootDir, "./**/*.proto"))
	rootFiles, _ := filepath.Glob(filepath.Join(rootDir, "./*.proto"))
	args = append(args, rootFiles...)
	args = append(args, protoFiles...)
	stderr := bytes.NewBuffer(nil)
	cmd := exec.Command(protocBinPath, args...)
	cmd.Stdout = stderr
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s returned error: %v %v", protocBinPath, err, stderr.String())
	}
	return nil
}

func getDescriptorData(t *testing.T, name string) ([]byte, []byte) {
	rule := strings.ToLower(name)
	root, _ := filepath.Abs(fmt.Sprintf("./testdata/%s/", rule))
	currentFileName := filepath.Join(t.TempDir(), "current.desc")
	prevFileName := filepath.Join(t.TempDir(), "prev.desc")
	err := runProtoc(filepath.Join(root, "current"), true, currentFileName)
	assert.NoError(t, err)
	err = runProtoc(filepath.Join(root, "previous"), true, prevFileName)
	assert.NoError(t, err)
	current, _ := ioutil.ReadFile(currentFileName)
	prev, _ := ioutil.ReadFile(prevFileName)
	return current, prev
}

func filter(strs []string, test func(string) bool) []string {
	var result []string
	for _, s := range strs {
		if test(s) {
			result = append(result, s)
		}
	}
	return result
}

func TestCompare(t *testing.T) {
	for _, test := range []struct {
		rule        string
		expectedErr []string
	}{
		{"FILE_NO_BREAKING_CHANGE", []string{
			`b/notfound.proto: file has been deleted`,
			`syntax.proto: syntax changed from "proto2" to "proto3"`,
			`package.proto: package changed from "filebreakingchange" to "filebreaking"`,
			`options/4.proto: File option "java package" changed from "com.stenciltest" to "com.stenciltest.change"`,
			`options/4.proto: File option "java outer classname" changed from "Teststencil" to "Teststencil.change"`,
			`options/5.proto: File option "java outer classname" changed from "Teststencil.valid" to ""`,
			`options/3.proto: all file options have been removed in current version`,
		}},
		{"ENUM_NO_BREAKING_CHANGE", []string{
			`1.proto: enum "a.Two" has been removed`,
			`1.proto: enum "a.Three.Four.Five" has been removed`,
			`1.proto: enum "a.Three.Eight" has been removed`,
			`1.proto: enumValue "a.Three.ONE" number changed from "0" to "1"`,
			`1.proto: enumValue "a.Three.TWO" number changed from "1" to "0"`,
			`1.proto: enumValue "SEVEN_SPECIFIED" deleted from enum "a.Three.Seven"`,
			`2.proto: enumValue "TEN_1" deleted from enum "a.Move"`,
			`2.proto: enumValue "TEN_2" deleted from enum "a.Move2"`,
			`2.proto: enum "a.MoveType" has been removed`,
		}},
		{"MESSAGE_NO_DELETE", []string{
			`1.proto: "a.Two" message has been removed`,
			`1.proto: "a.Three.Four.Five" message has been removed`,
			`1.proto: "a.Three.Seven" message has been removed`,
			`1.proto: "a.Nine" message has been removed`,
		}},
		{"FIELD_NO_BREAKING_CHANGE", []string{
			`1.proto: field "a.Two.three" is removed`,
			`1.proto: type has changed for "a.Three.three" from "a.Two" to "a.Three.Seven"`,
			`1.proto: type has changed for "a.Three.four" from "a.GroupEnums" to "a.GroupEnums2"`,
			`1.proto: number changed for "a.Three.Four.Six.four" from "4" to "5"`,
			`1.proto: type has changed for "a.Three.Four.Six.five" from "int32" to "string"`,
			`1.proto: label changed for "a.Three.Four.Six.six" from "repeated" to "optional"`,
			`1.proto: json name changed for "a.Three.Four.Six.eigth" from "foo" to "baz"`,
			`1.proto: field "a.Three.Eight.two" is removed`,
			`2.proto: field "a.One2.three" is removed`,
			`2.proto: field "a.One2.Two2" is removed`,
		}},
	} {
		t.Run(test.rule, func(t *testing.T) {
			allRules := []string{"FILE_NO_BREAKING_CHANGE", "ENUM_NO_BREAKING_CHANGE", "MESSAGE_NO_DELETE", "FIELD_NO_BREAKING_CHANGE"}
			current, prev := getDescriptorData(t, test.rule)
			skipRules := filter(allRules, func(s string) bool { return s != test.rule })
			err := proto.Compare(current, prev, skipRules)
			if err == nil {
				assert.Fail(t, "%s should return error", test.rule)
				return
			}
			errMsgs := strings.Split(err.Error(), "; ")

			assert.ElementsMatch(t, test.expectedErr, errMsgs)
		})
	}

	t.Run("should be able to skip rules", func(t *testing.T) {
		rule := "ENUM_NO_BREAKING_CHANGE"
		current, prev := getDescriptorData(t, rule)
		err := proto.Compare(current, prev, []string{rule})
		assert.Nil(t, err)
	})

	t.Run("should return nil if passed descriptors are backward compatibile", func(t *testing.T) {
		name := "valid"
		current, prev := getDescriptorData(t, name)
		err := proto.Compare(current, prev, []string{name})
		assert.Nil(t, err)
	})

	t.Run("should return err if passed data is not valid file descriptor set", func(t *testing.T) {
		err := proto.Compare([]byte("invalid bytes"), []byte("invalid bytes"), []string{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "descriptor set file is not valid.")
	})

	t.Run("should return err if passed file descriptor set is not fully contained", func(t *testing.T) {
		root, _ := filepath.Abs(fmt.Sprintf("./testdata/%s/", "valid"))
		fileName := filepath.Join(t.TempDir(), "current.desc")
		withImports := filepath.Join(t.TempDir(), "imports.desc")
		err := runProtoc(filepath.Join(root, "current"), false, fileName)
		assert.NoError(t, err)
		err = runProtoc(filepath.Join(root, "current"), true, withImports)
		assert.NoError(t, err)
		data, _ := ioutil.ReadFile(fileName)
		withImportsData, _ := ioutil.ReadFile(withImports)

		err = proto.Compare(data, withImportsData, []string{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "file is not fully contained descriptor file")

		err = proto.Compare(withImportsData, data, []string{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "file is not fully contained descriptor file")
	})
}
