package texts

const ConfigViewInitializing = `
	Initializing...
`

const ConfigViewHeaderTitle = `Config: {{ .path | highlight }}`

const ConfigViewFooterInfo = `{{ print .scrollPercentage " %" | special }}`

const ConfigViewScreen = `
{{ .header }}
{{ .viewport }}
{{ .footer }}
{{ .help }}`
