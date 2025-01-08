package models

const TemplateQuit = `

		ʕ •ɷ•ʔ ฅ		Bye!

`

const TemplateSelectCommand = `
Let’s make it queeck 🐸
Choose command:
{{ .commands }}

{{ "Description" | subtle }}: {{ .description }}
{{ "Equivalent" | subtle }}:  {{ .equivalent }}
`
