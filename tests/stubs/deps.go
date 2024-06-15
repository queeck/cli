package stubs

import (
	"io"

	servicesConfig "github.com/queeck/cli/internal/services/config"
)

//go:generate moq -skip-ensure -out deps_moq.go . config input output

type config interface { //nolint:unused // used for mocks generations
	servicesConfig.Config
}

type input interface { //nolint:unused // used for mocks generations
	io.Reader
}

type output interface { //nolint:unused // used for mocks generations
	io.Writer
}
