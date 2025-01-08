package view

const templateInitializing = `
	Initializing...
`

const templateHeaderTitle = `Config: {{ .path | highlight }}`

const templateFooterInfo = `{{ print .scrollPercentage " %" | special }}`

const templateScreen = `
{{ .header }}
{{ .viewport }}
{{ .footer }}
{{ .help }}`
