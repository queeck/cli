// nolint: lll // long line is ok for expected text
package functional

import (
	"context"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/require"

	"github.com/queeck/cli/tests/stubs"
)

const configSetOutput = `Enter key:

	config/<key.with.dots>                                                                                                          

← : move left • tab : autocomplete • esc / ctrl+c : quit • enter : select

`

const configSetNameOutput = `Enter value for key name:

	 = <value>                                                                                                                  

← : move left • tab : autocomplete • esc / ctrl+c : quit • enter : select

`

const configSetNameJonathanOutput = `name = Jonathan
`

func TestConfigSet(t *testing.T) {
	t.Run("expect that 'config set' prints output for config.set command", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
		defer cancel()

		cfg := stubs.ConfigInMemory(``)
		args := []string{"config", "set"}

		out := stubs.Writer()
		in := stubs.ReaderEmpty()

		program := makeApp(cfg, args, in, out)

		err := program.Run(ctx)
		require.Error(t, err)
		require.ErrorIs(t, err, tea.ErrProgramKilled)

		frames := out.FramesWithLF()
		require.NotEmpty(t, frames)

		require.Equal(t, configSetOutput, frames[len(frames)-1])
	})

	t.Run("expect that 'config set name' asks for value for key 'name'", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
		defer cancel()

		cfg := stubs.ConfigInMemory(``)
		args := []string{"config", "set", "name"}

		out := stubs.Writer()
		in := stubs.ReaderEmpty()

		program := makeApp(cfg, args, in, out)

		err := program.Run(ctx)
		require.Error(t, err)
		require.ErrorIs(t, err, tea.ErrProgramKilled)

		frames := out.FramesWithLF()
		require.NotEmpty(t, frames)

		require.Equal(t, configSetNameOutput, frames[len(frames)-1])
	})

	t.Run("expect that 'config set name Jonathan' prints new value and updates config", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
		defer cancel()

		cfg := stubs.ConfigInMemory(`{"name":"John Doe"}`)
		require.Equal(t, "name = John Doe", cfg.View())

		args := []string{"config", "set", "name", "Jonathan"}
		out := stubs.Writer()
		in := stubs.ReaderEmpty()

		program := makeApp(cfg, args, in, out)

		err := program.Run(ctx)
		require.NoError(t, err)

		frames := out.FramesWithLF()
		require.NotEmpty(t, frames)

		require.Equal(t, configSetNameJonathanOutput, frames[len(frames)-1])

		require.Equal(t, "name = Jonathan", cfg.View())
	})
}
