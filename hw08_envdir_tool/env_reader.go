package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

func ReadFile(path string) (env EnvValue) {
	f, err := os.Open(path)
	if err != nil {
		return EnvValue{"", true}
	}

	defer f.Close()

	reader := bufio.NewReader(f)

	line, _, err := reader.ReadLine()
	if err != nil {
		return EnvValue{"", true}
	}

	line = bytes.ReplaceAll(line, []byte{0x00}, []byte("\n"))
	t := string(line)
	t = strings.TrimRight(t, "	 ")

	return EnvValue{t, false}
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	values := Environment{}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		s := f.Name()

		if strings.Contains(s, "=") {
			continue
		}

		path := filepath.Join(dir, s)
		envValue := ReadFile(path)

		values[s] = envValue
	}

	return values, nil
}
