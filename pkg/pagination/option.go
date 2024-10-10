package pagination

type options struct {
	PageText        string
	SizeText        string
	DefaultPage     string
	DefaultPageSize string
	MinPageSize     int
	MaxPageSize     int
}

var defaultOptions = options{
	PageText:        "page",
	SizeText:        "size",
	DefaultPage:     "1",
	DefaultPageSize: "10",
	MinPageSize:     10,
	MaxPageSize:     100,
}

type CustomOption func(opts *options)

func WithPageText(pageText string) CustomOption {
	return func(opts *options) {
		opts.PageText = pageText
	}
}

func WithSizeText(sizeText string) CustomOption {
	return func(opts *options) {
		opts.SizeText = sizeText
	}
}

func WithDefaultPage(page string) CustomOption {
	return func(opts *options) {
		opts.DefaultPage = page
	}
}

func WithDefaultPageSize(pageSize string) CustomOption {
	return func(opts *options) {
		opts.DefaultPageSize = pageSize
	}
}

func WithMinPageSize(minPageSize int) CustomOption {
	return func(opts *options) {
		opts.MinPageSize = minPageSize
	}
}

func WithMaxPageSize(maxPageSize int) CustomOption {
	return func(opts *options) {
		opts.MaxPageSize = maxPageSize
	}
}

func applyCustomOptionsToDefault(customOptions ...CustomOption) options {
	opts := defaultOptions
	for _, customOption := range customOptions {
		customOption(&opts)
	}
	return opts
}
