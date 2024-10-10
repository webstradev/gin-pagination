// Package pagination provides a middleware for Gin web framework to handle
// pagination. It allows for the usage of url parameters like `?page=1&size=25`
// to paginate data on your API.
package pagination

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// New returns a new pagination middleware with custom values.
func New(customOptions ...CustomOption) gin.HandlerFunc {
	opts := defaultOptions
	for _, customOption := range customOptions {
		customOption(&opts)
	}

	return func(c *gin.Context) {
		// Extract the page from the query string and convert it to an integer
		pageStr := c.DefaultQuery(opts.PageText, opts.Page)
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "page number must be an integer"})
			return
		}

		// Validate for positive page number
		if page < 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "page number must be positive"})
			return
		}

		// Extract the size from the query string and convert it to an integer
		sizeStr := c.DefaultQuery(opts.SizeText, opts.PageSize)
		size, err := strconv.Atoi(sizeStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "page size must be an integer"})
			return
		}

		// Validate for min and max page size
		if size < opts.MinPageSize || size > opts.MaxPageSize {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "page size must be between " + strconv.Itoa(opts.MinPageSize) + " and " + strconv.Itoa(opts.MaxPageSize)})
			return
		}

		// Set the page and size in the gin context
		c.Set(opts.PageText, page)
		c.Set(opts.SizeText, size)

		c.Next()
	}
}
