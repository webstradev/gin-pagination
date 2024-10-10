package pagination

type Options struct {
	PageText    string
	SizeText    string
	Page        string
	PageSize    string
	MinPageSize int
	MaxPageSize int
}

type CustomOption func(opts *Options)

var defaultOptions = Options{
	PageText:    "page",
	SizeText:    "size",
	Page:        "1",
	PageSize:    "10",
	MinPageSize: 10,
	MaxPageSize: 100,
}
