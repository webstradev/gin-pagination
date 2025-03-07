package pagination

// CustomOption is a function that allows for customizing the pagination
// middleware.
type CustomOption func(opts *options)

// WithPageText allows for customizing the page parameter name.
func WithPageText(pageText string) CustomOption {
	return func(opts *options) {
		opts.PageText = pageText
	}
}

// WithSizeText allows for customizing the size parameter name.
func WithSizeText(sizeText string) CustomOption {
	return func(opts *options) {
		opts.SizeText = sizeText
	}
}

// WithDefaultPage allows for customizing the default page number.
func WithDefaultPage(page int) CustomOption {
	return func(opts *options) {
		opts.DefaultPage = page
	}
}

// WithDefaultPageSize allows for customizing the default page size.
func WithDefaultPageSize(pageSize int) CustomOption {
	return func(opts *options) {
		opts.DefaultPageSize = pageSize
	}
}

// WithMinPageSize allows for customizing the minimum page size.
func WithMinPageSize(minPageSize int) CustomOption {
	return func(opts *options) {
		opts.MinPageSize = minPageSize
	}
}

// WithMaxPageSize allows for customizing the maximum page size.
func WithMaxPageSize(maxPageSize int) CustomOption {
	return func(opts *options) {
		opts.MaxPageSize = maxPageSize
	}
}

// WithHeaderPrefix allows for customizing the header prefix.
// Empty strings are allowed and will result in no prefix.
func WithHeaderPrefix(headerPrefix string) CustomOption {
	return func(opts *options) {
		opts.HeaderPrefix = headerPrefix
	}
}

type options struct {
	PageText        string
	SizeText        string
	DefaultPage     int
	DefaultPageSize int
	MinPageSize     int
	MaxPageSize     int
	HeaderPrefix    string
}

var defaultOptions = options{
	PageText:        "page",
	SizeText:        "size",
	DefaultPage:     1,
	DefaultPageSize: 10,
	MinPageSize:     10,
	MaxPageSize:     100,
	HeaderPrefix:    "x-",
}

func applyCustomOptionsToDefault(customOptions ...CustomOption) options {
	opts := defaultOptions
	for _, customOption := range customOptions {
		customOption(&opts)
	}
	return opts
}
