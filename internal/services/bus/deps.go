package bus

import (
	"github.com/queeck/cli/internal/pkg/cli"
	serviceConfig "github.com/queeck/cli/internal/services/config"
	servicesState "github.com/queeck/cli/internal/services/state"
	servicesTemplates "github.com/queeck/cli/internal/services/templates"
)

//go:generate moq -skip-ensure -out deps_moq_test.go . arguments state templates config

type arguments interface {
	cli.Arguments
}

type state interface {
	servicesState.State
}

type templates interface {
	servicesTemplates.Renderer
}

type config interface {
	serviceConfig.Config
}
