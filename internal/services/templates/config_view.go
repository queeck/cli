package templates

import (
	"fmt"
	"strings"

	"github.com/queeck/cli/internal/services/templates/texts"
)

type RendererConfigView interface {
	RenderConfigViewHeaderTitle(path string) string
	RenderConfigViewFooterInfo(scrollPercent float64) string
	RenderConfigViewScreen(header, viewport, footer, help string) string
}

func (r *rendering) RenderConfigViewHeaderTitle(path string) string {
	return r.mustRender(texts.ConfigViewHeaderTitle, map[string]string{
		"path": path,
	})
}

func (r *rendering) RenderConfigViewFooterInfo(scrollPercent float64) string {
	return r.mustRender(texts.ConfigViewFooterInfo, map[string]string{
		"scrollPercentage": fmt.Sprintf("%3.f", scrollPercent*100), //nolint:mnd // not a magic number
	})
}

func (r *rendering) RenderConfigViewScreen(header, viewport, footer, help string) string {
	return r.mustRender(strings.TrimLeft(texts.ConfigViewScreen, "\n"), map[string]string{
		"header":   header,
		"viewport": viewport,
		"footer":   footer,
		"help":     help,
	})
}
