// nolint: lll // long line is ok for expected text
package functional

import (
	"context"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/require"

	"github.com/queeck/cli/tests/stubs"
)

const configGetOutput = `Enter key:

	config/<key.with.dots>                                                                                                          

← : move left • tab : autocomplete • esc / ctrl+c : quit • enter : select

`

const configGetNameOutput = `John Doe
`

func TestConfigGet(t *testing.T) {
	t.Run("expect that 'config get' prints output for config.get command", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
		defer cancel()

		cfg := stubs.ConfigInMemory(``)
		args := []string{"config", "get"}

		out := stubs.Writer()
		in := stubs.ReaderEmpty()

		program := makeApp(cfg, args, in, out)

		err := program.Run(ctx)
		require.Error(t, err)
		require.Equal(t, tea.ErrProgramKilled, err)

		frames := out.FramesWithLF()
		require.NotEmpty(t, frames)

		require.Equal(t, configGetOutput, frames[len(frames)-1])
	})

	t.Run("expect that 'config get name' returns value for key 'name'", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
		defer cancel()

		cfg := stubs.ConfigInMemory(`{"name":"John Doe"}`)
		args := []string{"config", "get", "name"}

		out := stubs.Writer()
		in := stubs.ReaderEmpty()

		program := makeApp(cfg, args, in, out)

		err := program.Run(ctx)
		require.NoError(t, err)

		frames := out.FramesWithLF()
		require.NotEmpty(t, frames)

		require.Equal(t, configGetNameOutput, frames[len(frames)-1])
	})
}
