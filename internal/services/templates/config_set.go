package templates

import (
	"strings"

	"github.com/queeck/cli/internal/services/templates/texts"
)

type RendererConfigSet interface {
	RenderConfigSetKeyScreen(inputKey, help string) string
	RenderConfigSetValueScreen(inputValue, key, help string) string
	RenderConfigSetKeyNotFound(key string) string
	RenderConfigSetKeyHasComplexType(key string) string
}

func (r *rendering) RenderConfigSetKeyScreen(inputKey, help string) string {
	return r.mustRender(strings.TrimLeft(texts.ConfigSetKeyScreen, "\n"), map[string]string{
		"inputKey": inputKey,
		"help":     help,
	})
}

func (r *rendering) RenderConfigSetValueScreen(inputValue, key, help string) string {
	return r.mustRender(strings.TrimLeft(texts.ConfigSetValueScreen, "\n"), map[string]string{
		"inputValue": inputValue,
		"key":        key,
		"help":       help,
	})
}

func (r *rendering) RenderConfigSetKeyNotFound(key string) string {
	return r.mustRender(strings.TrimLeft(texts.ConfigSetKeyNotFound, "\n"), map[string]string{
		"key": key,
	})
}

func (r *rendering) RenderConfigSetKeyHasComplexType(key string) string {
	return r.mustRender(strings.TrimLeft(texts.ConfigSetKeyHasComplexType, "\n"), map[string]string{
		"key": key,
	})
}
