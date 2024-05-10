package cli

import (
	"strings"
)

type Settings struct {
	commandsCount int
	flags         map[string]bool
	shortFlagsMap map[string]bool
}

type Option func(*Settings)

func WithCommandsCount(commandsCount int) Option {
	return func(s *Settings) {
		s.commandsCount = commandsCount
	}
}

func WithFlags(flags ...string) Option {
	return func(s *Settings) {
		for _, flag := range flags {
			if s.flags == nil {
				s.flags = make(map[string]bool, len(flags))
			}
			flag = strings.TrimLeft(flag, `-`)
			s.flags[flag] = true
			if len(flag) == 1 {
				if s.shortFlagsMap == nil {
					s.shortFlagsMap = make(map[string]bool, 1)
				}
				s.shortFlagsMap[flag] = true
			}
		}
	}
}
