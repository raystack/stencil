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

func TestCompare(t *testing.T) {
	for _, test := range []struct {
		rule           string
		isErrNil       bool
		skipRules      []string
		errContains    []string
		errNotContains []string
	}{
		{"FILE_NO_BREAKING_CHANGE", false, []string{}, []string{"all file options have been removed in options/3.proto current version", "java package for options/4.proto changed from com.stenciltest to com.stenciltest.change",
			"java outer classname for options/4.proto changed from Teststencil to Teststencil.change", "go package for options/4.proto changed from com.stenciltest to com.stenciltest.change",
			"package for package.proto changed from filebreakingchange to filebreaking", "syntax for syntax.proto changed from proto2 to proto3", "\"b/notfound.proto\" file has been deleted in current version"}, []string{"com.stenciltest.valid"}},
		{"ENUM_NO_BREAKING_CHANGE", false, []string{}, []string{"a.Two enum has been removed from current version", "enumValue a.Three.SEVEN_SPECIFIED deleted from current version", "a.Three.Eight enum has been removed from current version",
			"enumValue a.Three.ONE number changed from 0 to 1", "enumValue a.Three.TWO number changed from 1 to 0", "a.Three.Four.Five enum has been removed from current version", "enumValue a.TEN_1 deleted from current version",
			"enumValue a.TEN_2 deleted from current version"}, []string{}},
		{"MESSAGE_NO_DELETE", false, []string{}, []string{"a.Two has been removed in current version", "a.Three.Four.Five has been removed in current version", "a.Three.Seven has been removed in current version"}, []string{}},
		{"FIELD_NO_BREAKING_CHANGE", false, []string{}, []string{"a.One has been removed in current version", "a.Nine has been removed in current version", "field a.One.one is removed in current version", "field a.Two.three is removed in current version",
			"number changed for a.Two.four from 4 to 5", "type has changed for a.Two.five from int32 to string", "label changed for a.Two.six from repeated to optional", "json name changed for a.Two.eigth from foo to baz", "field a.Two.nine is removed in current version",
			"field a.Three.three is removed in current version", "number changed for a.Three.four from 4 to 5", "type has changed for a.Three.five from int32 to string", "label changed for a.Three.six from repeated to optional", "json name changed for a.Three.eigth from foo to baz",
			"field a.Three.Four.Five.three is removed in current version", "field a.Three.Four.Six.three is removed in current version", "number changed for a.Three.Four.Six.four from 4 to 5", "type has changed for a.Three.Four.Six.five from int32 to string", "label changed for a.Three.Four.Six.six from repeated to optional",
			"json name changed for a.Three.Four.Six.eigth from foo to baz", "field a.Three.Seven.three is removed in current version", "field a.Three.Eight.two is removed in current version", "field a.Nine.one is removed in current version", "field a.Nine.two is removed in current version", "field a.Nine.three is removed in current version", "field a.One2.three is removed in current version"}, []string{}},
		{"MESSAGE_NO_DELETE", true, []string{"MESSAGE_NO_DELETE", "FIELD_NO_BREAKING_CHANGE", "ENUM_NO_BREAKING_CHANGE", "FILE_NO_BREAKING_CHANGE"}, []string{}, []string{}},
		{"valid", true, []string{}, []string{}, []string{}},
	} {
		t.Run(test.rule, func(t *testing.T) {
			rule := strings.ToLower(test.rule)
			root, _ := filepath.Abs(fmt.Sprintf("./testdata/%s/", rule))
			currentFileName := filepath.Join(t.TempDir(), "current.desc")
			prevFileName := filepath.Join(t.TempDir(), "prev.desc")
			err := runProtoc(filepath.Join(root, "current"), true, currentFileName)
			assert.NoError(t, err)
			err = runProtoc(filepath.Join(root, "previous"), true, prevFileName)
			assert.NoError(t, err)
			current, _ := ioutil.ReadFile(currentFileName)
			prev, _ := ioutil.ReadFile(prevFileName)
			err = proto.Compare(current, prev, test.skipRules)
			if test.isErrNil {
				assert.Nil(t, err)
				return
			}
			for _, str := range test.errContains {
				assert.Contains(t, err.Error(), str)
			}
			for _, str := range test.errNotContains {
				assert.NotContains(t, err.Error(), str)
			}
		})
	}

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
