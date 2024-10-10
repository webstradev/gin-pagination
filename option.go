package pagination

type options struct {
	PageText    string
	SizeText    string
	Page        string
	PageSize    string
	MinPageSize int
	MaxPageSize int
}

type CustomOption func(opts *options)

var defaultOptions = options{
	PageText:    "page",
	SizeText:    "size",
	Page:        "1",
	PageSize:    "10",
	MinPageSize: 10,
	MaxPageSize: 100,
}
