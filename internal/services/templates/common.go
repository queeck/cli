package templates

import (
	"strings"

	"github.com/queeck/cli/internal/services/templates/texts"
)

type RendererCommon interface {
	RenderCommonQuit() string
	RenderCommonSelectCommands(commands, description, equivalent string) string
}

func (r *rendering) RenderCommonQuit() string {
	return texts.CommonQuit
}

func (r *rendering) RenderCommonSelectCommands(commands, description, equivalent string) string {
	return r.mustRender(strings.TrimLeft(texts.CommonSelectCommand, "\n"), map[string]string{
		"commands":    commands,
		"description": description,
		"equivalent":  equivalent,
	})
}
