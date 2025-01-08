package set

const templateKeyScreen = `
Enter key:

	{{ .inputKey }}

{{ .help }}

`

const templateValueScreen = `
Enter value for key {{ .key | highlight }}:

	{{ .inputValue }}

{{ .help }}

`

const templateKeyNotFound = `
Key {{ .key | highlight  }} not found.
`

const templateKeyHasComplexType = `
Key {{ .key | highlight  }} has complex type, try specifying the key deeper.
`
