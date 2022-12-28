package pagination

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestPaginationMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		middleware     gin.HandlerFunc
		queryParams    url.Values
		expectedPage   int
		expectedSize   int
		customPageText string
		customSizeText string
	}{
		{
			"Non int Page Param - Bad Request",
			Default(),
			url.Values{
				"page": {"notanumber"},
			},
			0,
			0,
			"",
			"",
		},
		{
			"Non int Size Param - Bad Request",
			Default(),
			url.Values{
				"page": {"1"},
				"size": {"notanumber"},
			},
			0,
			0,
			"",
			"",
		},
		{
			"Negative Page Param - Bad Request",
			Default(),
			url.Values{
				"page": {"-1"},
			},
			0,
			0,
			"",
			"",
		},
		{
			"Size below min - Bad Request",
			Default(),
			url.Values{
				"page": {"1"},
				"size": {"0"},
			},
			0,
			0,
			"",
			"",
		},
		{
			"Size above max - Bad Request",
			Default(),
			url.Values{
				"page": {"1"},
				"size": {"101"},
			},
			0,
			0,
			"",
			"",
		},
		{
			"Default Handling",
			Default(),
			url.Values{},
			1,
			10,
			"",
			"",
		},
		{
			"The first 100 results",
			Default(),
			url.Values{
				"page": {"1"},
				"size": {"100"},
			},
			1,
			100,
			"",
			"",
		},
		{
			"The second 20 results",
			Default(),
			url.Values{
				"page": {"2"},
				"size": {"20"},
			},
			2,
			20,
			"",
			"",
		},
		{
			"Custom Handling",
			New("pages", "items", "0", "5", 1, 25),
			url.Values{},
			0,
			5,
			"pages",
			"items",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			// Create a test context
			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

			// Add the query parameters to the Request of the test context
			ctx.Request = &http.Request{
				URL: &url.URL{
					RawQuery: url.Values(tt.queryParams).Encode(),
				},
			}

			// Call middleware on the test context
			tt.middleware(ctx)

			// Handle custom page and size text
			if tt.customPageText == "" {
				tt.customPageText = "page"
			}

			if tt.customSizeText == "" {
				tt.customSizeText = "size"
			}

			gotPage := ctx.GetInt(tt.customPageText)
			// Check if the page and pageSize are set correctly
			if gotPage != tt.expectedPage {
				t.Errorf("Expected page %d, got %d", tt.expectedPage, gotPage)
			}

			gotSize := ctx.GetInt(tt.customSizeText)
			if gotSize != tt.expectedSize {
				t.Errorf("Expected size %d, got %d", tt.expectedSize, gotSize)
			}
		})
	}
}
