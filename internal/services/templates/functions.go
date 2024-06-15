package templates

import (
	"strings"
	"text/template"

	"github.com/queeck/cli/internal/pkg/styles"
)

var functionsMap = template.FuncMap{
	"subtle": func(args ...string) string {
		return styles.ColorForegroundSubtle(strings.Join(args, ""))
	},
	"special": func(args ...string) string {
		return styles.ColorForegroundSpecial(strings.Join(args, ""))
	},
	"highlight": func(args ...string) string {
		return styles.ColorForegroundHighlight(strings.Join(args, ""))
	},
}
