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
		path := filepath.Join(dir, s)
		file, err := os.Open(path)
		if err != nil {
			continue
		}
		defer file.Close()

		reader := bufio.NewReader(file)

		line, _, err := reader.ReadLine()
		if err != nil {
			values[s] = EnvValue{"", true}

			continue
		}

		line = bytes.ReplaceAll(line, []byte{0x00}, []byte("\n"))
		t := string(line)
		t = strings.TrimRight(t, "	 ")

		values[s] = EnvValue{t, false}
	}

	return values, nil
}
