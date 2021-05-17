package main

import (
	"math"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	defer os.RemoveAll("/tmp/copy/")

	t.Run("random 'from' path", func(t *testing.T) {
		err := Copy("/tmp/copy/randomfoldername", "/tmp/copy/test1/output.txt", 0, 0)

		require.ErrorIs(t, err, ErrUnsupportedFile, "should throw error when can't open file")
	})

	t.Run("wrong offset", func(t *testing.T) {
		err := Copy("./testdata/input.txt", "/tmp/copy/test4/output.txt", math.MaxInt64, 0)

		require.ErrorIs(t, err, ErrOffsetExceedsFileSize, "should throw error when offset is greater than filesize")
	})

	t.Run("random 'to' path", func(t *testing.T) {
		os.RemoveAll("/tmp/copy/test2/")

		err := Copy("./testdata/input.txt", "/tmp/copy/test2/output.txt", 0, 0)

		require.NoError(t, err, "should create destination folder")
	})

	t.Run("wrong limit arg", func(t *testing.T) {
		err := Copy("./testdata/input.txt", "/tmp/copy/test3/output.txt", 0, math.MaxInt64)

		require.NoError(t, err, "should copy full file without errors")
	})
}
