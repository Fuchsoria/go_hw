package main

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("empty cmd length", func(t *testing.T) {
		strs := []string{}

		code := RunCmd(strs, Environment{})

		require.Equal(t, 1, code, "should return 1 code")
	})

	t.Run("test echoline return", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip()
		}

		strs := []string{"/bin/bash", "testdata/echoline.sh"}

		code := RunCmd(strs, Environment{})

		require.Equal(t, 0, code, "should return 0 code")
	})
}
