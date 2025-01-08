package executor

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	pkgOS "github.com/queeck/cli/internal/pkg/os"
)

func TestExecutor_Execute(t *testing.T) {
	currentDirectory, err := os.Getwd()
	require.NoError(t, err)

	t.Run("get current directory on windows", func(t *testing.T) {
		if !pkgOS.IsWindows() {
			t.Skip("skip — running not on windows")
		}
		var result string
		result, err = New().Execute(`cd`)
		require.NoError(t, err)
		require.Truef(t, strings.HasSuffix(result, currentDirectory+"\r\n"),
			"expected path to current directory, got: %s", result)
	})

	t.Run("get current directory on linux", func(t *testing.T) {
		if !pkgOS.IsLinux() {
			t.Skip("skip — running not on linux")
		}
		var result string
		result, err = New().Execute(`pwd`)
		require.NoError(t, err)
		require.Truef(t, strings.HasSuffix(result, currentDirectory+"\n"),
			"expected path to current directory, got: %s", result)
	})

	t.Run("get current directory on darwin", func(t *testing.T) {
		if !pkgOS.IsDarwin() {
			t.Skip("skip — running not on darwin")
		}
		var result string
		result, err = New().Execute(`pwd`)
		require.NoError(t, err)
		require.Truef(t, strings.HasSuffix(result, currentDirectory+"\n"),
			"expected path to current directory, got: %s", result)
	})
}
