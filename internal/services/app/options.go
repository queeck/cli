package app

type appOptions struct {
	fps int
}

type Option func(options *appOptions)

func WithFPS(fps int) Option {
	return func(options *appOptions) {
		options.fps = fps
	}
}
