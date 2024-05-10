package templates

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/queeck/cli/internal/services/templates/texts"
)

type Renderer interface {
	RenderCommonQuit() string
	RenderCommonSelectCommands(commands, description, equivalent string) string
}

const (
	interpolatorName      = `services/templates/interpolate`
	optionErrOnMissingKey = `missingkey=error`
)

type rendering struct {
	interpolator *template.Template
}

func New() Renderer {
	interpolator := template.New(interpolatorName).Option(optionErrOnMissingKey).Funcs(functionsMap)

	return &rendering{
		interpolator: interpolator,
	}
}

func (t *rendering) mustRender(source string, data any) string {
	clone := *t.interpolator
	parsed, err := (&clone).Parse(source)
	if err != nil {
		panic(err)
	}
	buf := &bytes.Buffer{}
	if err = parsed.Execute(buf, data); err != nil {
		panic(err)
	}
	return buf.String()
}

func (t *rendering) RenderCommonQuit() string {
	return texts.CommonQuit
}

func (t *rendering) RenderCommonSelectCommands(commands, description, equivalent string) string {
	return t.mustRender(strings.TrimLeft(texts.CommonSelectCommand, "\n"), map[string]string{
		"commands":    commands,
		"description": description,
		"equivalent":  equivalent,
	})
}
