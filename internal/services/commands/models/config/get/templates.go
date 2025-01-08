package get

const templateScreen = `
Enter key:

	{{ .inputKey }}

{{ .help }}

`

const templateKeyIsEmpty = `
{{ print "(key is empty)" | subtle }}
`

const templateKeyNotFound = `
{{ print "(key " .key " not found)" | subtle }}
`

const templateValueWithKey = `
{{ print .key " =" | subtle }} {{ .value | highlight }}
`

const templateValue = `
{{ .value }}
`
