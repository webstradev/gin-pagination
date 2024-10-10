package pagination

type options struct {
	PageText        string
	SizeText        string
	DefaultPage     int
	DefaultPageSize int
	MinPageSize     int
	MaxPageSize     int
}

var defaultOptions = options{
	PageText:        "page",
	SizeText:        "size",
	DefaultPage:     1,
	DefaultPageSize: 10,
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

func WithDefaultPage(page int) CustomOption {
	return func(opts *options) {
		opts.DefaultPage = page
	}
}

func WithDefaultPageSize(pageSize int) CustomOption {
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
