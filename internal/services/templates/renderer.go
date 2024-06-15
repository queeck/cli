package templates

import (
	"bytes"
	"text/template"
)

type Renderer interface {
	RendererCommon
	RendererConfigGet
	RendererConfigSet
	RendererConfigView
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

func (r *rendering) mustRender(source string, data any) string {
	clone := *r.interpolator
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
