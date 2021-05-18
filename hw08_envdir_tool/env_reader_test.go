package main

import (
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
		require.Nilf(t, err, "err should be nil")

		defer os.RemoveAll(dir)

		envs, err := ReadDir(dir)

		require.Equal(t, 0, len(envs), "envs length should be 0")
		require.NoError(t, err, "should be done without errors")
	})

	t.Run("check folder", func(t *testing.T) {
		dir, err := ioutil.TempDir("", "test")
		require.Nilf(t, err, "err should be nil")

		defer os.RemoveAll(dir)

		file1, err := ioutil.TempFile(dir, "")
		require.Nilf(t, err, "err should be nil")
		io.WriteString(file1, "first")

		fInfo1, err := file1.Stat()
		require.Nilf(t, err, "err should be nil")

		file2, err := ioutil.TempFile(dir, "")
		require.Nilf(t, err, "err should be nil")
		io.WriteString(file2, "second")

		fInfo2, err := file2.Stat()
		require.Nilf(t, err, "err should be nil")

		envs, err := ReadDir(dir)

		require.Equal(t, "first", envs[fInfo1.Name()].Value, "first test value should be equal to 'first'")
		require.Equal(t, "second", envs[fInfo2.Name()].Value, "first test value should be equal to 'second'")
		require.NoError(t, err, "should be done without errors")
	})
}
