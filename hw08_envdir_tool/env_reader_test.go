package main

import (
	"io"
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
		os.MkdirAll("/tmp/checkempty", 0777)
		defer os.RemoveAll("/tmp/checkempty")

		envs, err := ReadDir("/tmp/checkempty")

		require.Equal(t, 0, len(envs), "envs length should be 0")
		require.NoError(t, err, "should be done without errors")
	})

	t.Run("check folder", func(t *testing.T) {
		os.MkdirAll("/tmp/checkfolder", 0777)
		defer os.RemoveAll("/tmp/checkfolder")

		file1, _ := os.Create("/tmp/checkfolder/TESTONE")
		io.WriteString(file1, "first")

		file2, _ := os.Create("/tmp/checkfolder/TESTTWO")
		io.WriteString(file2, "second")

		file3, _ := os.Create("/tmp/checkfolder/TESTTHREE")
		io.WriteString(file3, "third")

		envs, err := ReadDir("/tmp/checkfolder")

		require.Equal(t, "first", envs["TESTONE"].Value, "first test value should be equal to 'first'")
		require.Equal(t, "second", envs["TESTTWO"].Value, "first test value should be equal to 'second'")
		require.Equal(t, "third", envs["TESTTHREE"].Value, "first test value should be equal to 'third'")
		require.NoError(t, err, "should be done without errors")
	})
}
