// Package pagination provides a middleware for Gin web framework to handle
// pagination. It allows for the usage of url parameters like `?page=1&size=25`
// to paginate data on your API. The values will be propagated throughout the
// request context.
package pagination

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// New returns a new pagination middleware with custom values.
func New(customOptions ...CustomOption) gin.HandlerFunc {

	return func(c *gin.Context) {
		p := &paginator{
			opts: applyCustomOptionsToDefault(customOptions...),
			c:    c,
		}

		page, err := p.getPageFromRequest()
		if err != nil {
			p.abortWithBadRequest(err)
			return
		}

		if err := p.validatePage(page); err != nil {
			p.abortWithBadRequest(err)
			return
		}

		pageSize, err := p.getPageSizeFromRequest()
		if err != nil {
			p.abortWithBadRequest(err)
			return
		}

		if err := p.validatePageSize(pageSize); err != nil {
			p.abortWithBadRequest(err)
			return
		}

		p.setPageAndPageSize(page, pageSize)

		p.Next()
	}
}

type paginator struct {
	opts options
	c    *gin.Context
}

func (p *paginator) abortWithBadRequest(err error) {
	p.c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}

func (p *paginator) getPageFromRequest() (int, error) {
	return p.getIntValueWithDefault(p.opts.PageText, strconv.Itoa(p.opts.DefaultPage))
}

func (p *paginator) getPageSizeFromRequest() (int, error) {
	return p.getIntValueWithDefault(p.opts.SizeText, strconv.Itoa(p.opts.DefaultPageSize))
}

func (p *paginator) getIntValueWithDefault(key string, defaultValue string) (int, error) {
	valueStr := p.c.DefaultQuery(key, defaultValue)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("%s parameter must be an integer", key)
	}

	return value, nil
}

func (p *paginator) validatePage(page int) error {
	if page < 0 {
		return errors.New("page number must be positive")
	}

	return nil
}

func (p *paginator) validatePageSize(size int) error {
	if size < p.opts.MinPageSize || size > p.opts.MaxPageSize {
		return errors.New("page size must be between %d and %d")
	}

	return nil
}

func (p *paginator) setPageAndPageSize(page int, size int) {
	p.c.Set(p.opts.PageText, page)
	p.c.Header(p.opts.PageText, strconv.Itoa(page))
	p.c.Set(p.opts.SizeText, size)
	p.c.Header(p.opts.SizeText, strconv.Itoa(size))
}

func (p *paginator) Next() {
	p.c.Next()
}
