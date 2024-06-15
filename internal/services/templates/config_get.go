package templates

import (
	"strings"

	"github.com/queeck/cli/internal/services/templates/texts"
)

type RendererConfigGet interface {
	RenderConfigGetScreen(inputKey, help string) string
	RenderConfigGetKeyNotFound(key string) string
	RenderConfigGetKeyIsEmpty() string
	RenderConfigGetValueWithKey(key, value string) string
	RenderConfigGetValue(value string) string
}

func (r *rendering) RenderConfigGetScreen(inputKey, help string) string {
	return r.mustRender(strings.TrimLeft(texts.ConfigGetScreen, "\n"), map[string]string{
		"inputKey": inputKey,
		"help":     help,
	})
}

func (r *rendering) RenderConfigGetKeyNotFound(key string) string {
	return r.mustRender(strings.TrimLeft(texts.ConfigGetKeyNotFound, "\n"), map[string]string{
		"key": key,
	})
}

func (r *rendering) RenderConfigGetKeyIsEmpty() string {
	return r.mustRender(strings.TrimLeft(texts.ConfigGetKeyIsEmpty, "\n"), map[string]string{})
}

func (r *rendering) RenderConfigGetValueWithKey(key, value string) string {
	return r.mustRender(strings.TrimLeft(texts.ConfigGetValueWithKey, "\n"), map[string]string{
		"key":   key,
		"value": value,
	})
}

func (r *rendering) RenderConfigGetValue(value string) string {
	return r.mustRender(strings.TrimLeft(texts.ConfigGetValue, "\n"), map[string]string{
		"value": value,
	})
}
