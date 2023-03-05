package gomapper

type ApplyOptions func(options *configOptions)

type configOptions struct {
	CaseInsensitive bool
}

func defaultOptions() configOptions {
	return configOptions{
		CaseInsensitive: false,
	}
}

func WithCaseInsensitive(value bool) ApplyOptions {
	return func(options *configOptions) {
		options.CaseInsensitive = value
	}
}
