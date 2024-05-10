package cli

type Arguments interface {
	Commands() []string
	Flags() []string
	HasFlag(flag string) bool
	Options() map[string]string
	Option(option string) (value string, has bool)
	PositionalArguments() []string
	ToMap() map[string]any
}

type Args struct {
	commands            []string
	flagsMap            map[string]bool
	options             map[string]string
	positionalArguments []string
}

func (a *Args) Commands() []string {
	return a.commands
}

func (a *Args) Flags() []string {
	flags := make([]string, 0, len(a.flagsMap))
	for flag := range a.flagsMap {
		flags = append(flags, flag)
	}
	return flags
}

func (a *Args) HasFlag(flag string) bool {
	return a.flagsMap[flag]
}

func (a *Args) Option(option string) (value string, has bool) {
	value, has = a.options[option]
	return value, has
}

func (a *Args) Options() map[string]string {
	return a.options
}

func (a *Args) PositionalArguments() []string {
	return a.positionalArguments
}

func (a *Args) ToMap() map[string]any {
	return map[string]any{
		"commands":  a.commands,
		"flags":     a.Flags(),
		"options":   a.options,
		"arguments": a.positionalArguments,
	}
}
