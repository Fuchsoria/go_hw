package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	// Place your code here
	t.Run("check wrong folder", func(t *testing.T) {
		_, err := ReadDir("testdata/nonexist")

		require.Error(t, err, "wrong folder should return error")
	})

	t.Run("check empty folder", func(t *testing.T) {
		dir, err := ioutil.TempDir("", "test")
		if err != nil {
			fmt.Println(err)
		}
		defer os.RemoveAll(dir)

		envs, err := ReadDir(dir)

		require.Equal(t, 0, len(envs), "envs length should be 0")
		require.NoError(t, err, "should be done without errors")
	})

	t.Run("check folder", func(t *testing.T) {
		dir, err := ioutil.TempDir("", "test")
		if err != nil {
			fmt.Println(err)
		}

		defer os.RemoveAll(dir)

		file1, err := ioutil.TempFile(dir, "")
		if err != nil {
			fmt.Println(err)
		}
		io.WriteString(file1, "first")
		fInfo1, err := file1.Stat()
		if err != nil {
			fmt.Println(err)
		}

		file2, err := ioutil.TempFile(dir, "")
		if err != nil {
			fmt.Println(err)
		}
		io.WriteString(file2, "second")
		fInfo2, err := file2.Stat()
		if err != nil {
			fmt.Println(err)
		}

		envs, err := ReadDir(dir)

		require.Equal(t, "first", envs[fInfo1.Name()].Value, "first test value should be equal to 'first'")
		require.Equal(t, "second", envs[fInfo2.Name()].Value, "first test value should be equal to 'second'")
		require.NoError(t, err, "should be done without errors")
	})
}
