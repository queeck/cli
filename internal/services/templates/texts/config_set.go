package texts

const ConfigSetKeyScreen = `
Enter key:

	{{ .inputKey }}

{{ .help }}

`

const ConfigSetValueScreen = `
Enter value for key {{ .key | highlight }}:

	{{ .inputValue }}

{{ .help }}

`

const ConfigSetKeyNotFound = `
Key {{ .key | highlight  }} not found.
`

const ConfigSetKeyHasComplexType = `
Key {{ .key | highlight  }} has complex type, try specifying the key deeper.
`
