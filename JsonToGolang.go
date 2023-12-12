package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

func JsonToGolang(json, topLevel, packageName string) (string, error) {
	cmd := exec.Command(
		"quicktype",
		"--just-types",
		"--lang",
		"go",
		"--src-lang",
		"json",
		"--top-level",
		topLevel,
	)

	cmd.Stdin = bytes.NewReader([]byte(json))

	var stdoutBuf, stdErrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stdErrBuf

	err := cmd.Run()

	if err != nil {
		fmt.Println(stdoutBuf.String())
		return "", err
	}

	return `package ` + packageName + "\n\n" + stdoutBuf.String(), nil
}
