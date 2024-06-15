package functional

import (
	"context"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/require"

	"github.com/queeck/cli/tests/stubs"
)

const configViewOutput = `╭────────────────────────────╮                                                  
│ Config: tests/stubs/config ├──────────────────────────────────────────────────
╰────────────────────────────╯                                                  
name = John Doe                                                                 
                                                                       ╭───────╮
───────────────────────────────────────────────────────────────────────┤ 100 % │
                                                                       ╰───────╯
h : toggle help • q / esc : quit`

func TestConfigView(t *testing.T) {
	t.Run("expect that 'config view' prints config", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
		defer cancel()

		cfg := stubs.ConfigInMemory(`{"name":"John Doe"}`)
		args := []string{"config", "view"}

		out := stubs.Writer()
		in := stubs.ReaderEmpty()

		program := makeApp(cfg, args, in, out)

		program.CommandBus().CommandConfigView().Update(tea.WindowSizeMsg{
			Width:  80,
			Height: 8,
		}) // Prepare sizing for viewport

		err := program.Run(ctx)
		require.Error(t, err)
		require.Equal(t, tea.ErrProgramKilled, err)

		frames := out.FramesWithLF()
		require.NotEmpty(t, frames)

		require.Equal(t, configViewOutput, frames[len(frames)-1])
	})
}
