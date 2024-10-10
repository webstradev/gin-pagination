// Package pagination provides a middleware for Gin web framework to handle
// pagination. It allows for the usage of url parameters like `?page=1&size=25`
// to paginate data on your API.
package pagination

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	defaultPageText    = "page"
	defaultSizeText    = "size"
	defaultPage        = "1"
	defaultPageSize    = "10"
	defaultMinPageSize = 10
	defaultMaxPageSize = 100
)

// Default returns a new pagination middleware with default values.
func Default() gin.HandlerFunc {
	return New(
		defaultPageText,
		defaultSizeText,
		defaultPage,
		defaultPageSize,
		defaultMinPageSize,
		defaultMaxPageSize,
	)
}

// New returns a new pagination middleware with custom values.
func New(pageText, sizeText, defaultPage, defaultPageSize string, minPageSize, maxPageSize int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the page from the query string and convert it to an integer
		pageStr := c.DefaultQuery(pageText, defaultPage)
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
		sizeStr := c.DefaultQuery(sizeText, defaultPageSize)
		size, err := strconv.Atoi(sizeStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "page size must be an integer"})
			return
		}

		// Validate for min and max page size
		if size < minPageSize || size > maxPageSize {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "page size must be between " + strconv.Itoa(minPageSize) + " and " + strconv.Itoa(maxPageSize)})
			return
		}

		// Set the page and size in the gin context
		c.Set(pageText, page)
		c.Set(sizeText, size)

		c.Next()
	}
}
