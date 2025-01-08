package templates

import (
	"bytes"
	"strings"
	"text/template"
)

type Renderer interface {
	Render(template string, keyValues ...string) string
}

const (
	interpolatorName        = `services/templates/interpolate`
	optionMissingKeyDefault = `missingkey=default`
)

type rendering struct {
	interpolator *template.Template
}

func New() Renderer {
	interpolator := template.New(interpolatorName).Option(optionMissingKeyDefault).Funcs(functionsMap)

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

func (r *rendering) Render(template string, keyValues ...string) string {
	data := keyValuesToData(keyValues...)
	template = strings.TrimPrefix(template, "\n")
	return r.mustRender(template, data)
}

func keyValuesToData(keyValues ...string) map[string]string {
	if len(keyValues) == 0 {
		return nil
	}
	length := len(keyValues) / 2 //nolint:mnd // is not a magic number
	if len(keyValues)%2 == 1 {
		length++
	}
	data := make(map[string]string, length)
	var key, value string
	for i, item := range keyValues {
		if i%2 == 0 {
			key = item
			value = ""
		} else {
			value = item
		}
		data[key] = value
	}
	return data
}
