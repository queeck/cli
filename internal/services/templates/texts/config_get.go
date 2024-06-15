package texts

const ConfigGetScreen = `
Enter key:

	{{ .inputKey }}

{{ .help }}

`

const ConfigGetKeyIsEmpty = `
{{ print "(key is empty)" | subtle }}
`

const ConfigGetKeyNotFound = `
{{ print "(key " .key " not found)" | subtle }}
`

const ConfigGetValueWithKey = `
{{ print .key " =" | subtle }} {{ .value | highlight }}
`

const ConfigGetValue = `
{{ .value }}
`
