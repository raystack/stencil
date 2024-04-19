package test_helper

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os/exec"
	"path/filepath"
)

func RunProtoc(
	rootDir string,
	includeImports bool,
	descSetOut string,
	protoFilesNames []string,
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
	if protoFilesNames == nil {
		protoFiles, _ := filepath.Glob(filepath.Join(rootDir, "./**/*.proto"))
		rootFiles, _ := filepath.Glob(filepath.Join(rootDir, "./*.proto"))
		args = append(args, rootFiles...)
		args = append(args, protoFiles...)
	} else {
		args = append(args, protoFilesNames...)
	}
	stderr := bytes.NewBuffer(nil)
	cmd := exec.Command(protocBinPath, args...)
	cmd.Stdout = stderr
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s returned error: %v %v", protocBinPath, err, stderr.String())
	}
	return nil
}

func GetRandomName() string {
	b := make([]byte, 10)
	rand.Read(b)
	return hex.EncodeToString(b)
}
