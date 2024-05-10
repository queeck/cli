package cli

import (
	"strings"
)

func Parse(arguments []string, options ...Option) Arguments {
	settings := &Settings{}
	for _, option := range options {
		option(settings)
	}

	return newParser(settings).Parse(arguments)
}

type parser struct {
	isCommandsExtracted bool
	commandsCount       int
	settings            *Settings
	args                *Args
}

func newParser(settings *Settings) *parser {
	return &parser{
		isCommandsExtracted: false,
		commandsCount:       0,
		settings:            settings,
		args: &Args{
			commands:            make([]string, 0),
			flagsMap:            make(map[string]bool),
			options:             make(map[string]string),
			positionalArguments: make([]string, 0),
		},
	}
}

func (p *parser) Parse(arguments []string) *Args {
	for i := 0; i < len(arguments); i++ {
		current := arguments[i]
		next := ""
		if i+1 < len(arguments) {
			next = arguments[i+1]
		}

		if current == "" {
			continue
		}

		if isNextUsed := p.extract(current, next); isNextUsed {
			i++
		}
	}

	return p.args
}

func (p *parser) extractCommand(current string) (done bool) {
	if p.isCommandsExtracted {
		return false
	}
	if current[0] == '-' {
		p.isCommandsExtracted = true
		return false
	}
	p.args.commands = append(p.args.commands, current)
	p.commandsCount++
	if p.commandsCount == p.settings.commandsCount {
		p.isCommandsExtracted = true
	}
	return true
}

func (p *parser) extractShortFlags(option string) (done bool) {
	if len(p.settings.shortFlagsMap) != 0 && len(option) <= len(p.settings.shortFlagsMap) {
		optionFlagsMap := make(map[string]bool)
		match := true
		for n := range option {
			letter := string(option[n])
			optionFlagsMap[letter] = true
			if !p.settings.shortFlagsMap[letter] {
				match = false
				break
			}
		}
		if match {
			for k := range optionFlagsMap {
				if p.args.flagsMap == nil {
					p.args.flagsMap = make(map[string]bool, len(optionFlagsMap))
				}
				p.args.flagsMap[k] = true
			}
			return true
		}
	}
	return false
}

func (p *parser) extractOptionOrFlag(option, next string) (isNextUsed bool) {
	value := ""
	if strings.Contains(option, `=`) {
		chunks := strings.Split(option, `=`)
		option = chunks[0]
		value = chunks[1]
	}
	if !p.settings.flags[option] && next != "" && next[0] != '-' {
		value = next
		isNextUsed = true
	}

	if value == "" {
		p.args.flagsMap[option] = true
	} else {
		p.args.options[option] = value
	}

	return isNextUsed
}

func (p *parser) extract(current, next string) (isNextUsed bool) {
	if p.extractCommand(current) {
		return false
	}

	if current[0] != '-' {
		p.args.positionalArguments = append(p.args.positionalArguments, current)
		return false
	}

	option := strings.TrimLeft(current, `-`)

	if p.extractShortFlags(option) {
		return false
	}

	return p.extractOptionOrFlag(option, next)
}
