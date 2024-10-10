// Package pagination provides a middleware for Gin web framework to handle
// pagination. It allows for the usage of url parameters like `?page=1&size=25`
// to paginate data on your API. The values will be propagated throughout the
// request context.
package pagination

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// New returns a new pagination middleware with custom values.
func New(customOptions ...CustomOption) gin.HandlerFunc {
	opts := applyCustomOptionsToDefault(customOptions...)

	return func(c *gin.Context) {
		// Extract the page from the query string and convert it to an integer.
		pageStr := c.DefaultQuery(opts.PageText, opts.DefaultPage)
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{
					"error": "page number must be an integer",
				},
			)
			return
		}

		// Validate for positive page number.
		if page < 0 {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{
					"error": "page number must be positive",
				},
			)
			return
		}

		// Extract the size from the query string and convert it to an integer.
		sizeStr := c.DefaultQuery(opts.SizeText, opts.DefaultPageSize)
		size, err := strconv.Atoi(sizeStr)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{
					"error": "page size must be an integer",
				},
			)
			return
		}

		// Validate for min and max page size.
		if size < opts.MinPageSize || size > opts.MaxPageSize {
			errorMessage := fmt.Sprintf(
				"page size must be between %d and %d", opts.MinPageSize, opts.MaxPageSize,
			)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errorMessage})
			return
		}

		// Set the page and size in the gin context.
		c.Set(opts.PageText, page)
		c.Set(opts.SizeText, size)

		c.Next()
	}
}
