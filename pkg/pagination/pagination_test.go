package pagination_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/webstradev/gin-pagination/v2/pkg/pagination"
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
			pagination.New(),
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
			pagination.New(),
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
			pagination.New(),
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
			pagination.New(),
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
			pagination.New(),
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
			pagination.New(),
			url.Values{},
			1,
			10,
			"",
			"",
		},
		{
			"The first 100 results",
			pagination.New(),
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
			pagination.New(),
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
			pagination.New(
				pagination.WithPageText("pages"),
				pagination.WithSizeText("items"),
				pagination.WithDefaultPage(0),
				pagination.WithDefaultPageSize(5),
				pagination.WithMinPageSize(1),
				pagination.WithMaxPageSize(25),
			),
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

func TestPaginationHeaders(t *testing.T) {
	tests := []struct {
		name            string
		middleware      gin.HandlerFunc
		queryParams     url.Values
		expectedHeaders map[string]string
	}{
		{
			name:       "Default headers are set correctly",
			middleware: pagination.New(),
			queryParams: url.Values{
				"page": {"2"},
				"size": {"20"},
			},
			expectedHeaders: map[string]string{
				"X-Page": "2",
				"X-Size": "20",
			},
		},
		{
			name: "Custom text headers are set correctly",
			middleware: pagination.New(
				pagination.WithPageText("offset"),
				pagination.WithSizeText("limit"),
			),
			queryParams: url.Values{
				"offset": {"3"},
				"limit":  {"15"},
			},
			expectedHeaders: map[string]string{
				"X-Offset": "3",
				"X-Limit":  "15",
			},
		},
		{
			name:            "No headers on invalid input",
			middleware:      pagination.New(),
			queryParams:     url.Values{"page": {"invalid"}},
			expectedHeaders: map[string]string{},
		},
		{
			name:        "Default values are set in headers when no query params",
			middleware:  pagination.New(),
			queryParams: url.Values{},
			expectedHeaders: map[string]string{
				"X-Page": "1",
				"X-Size": "10",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			ctx.Request = &http.Request{
				URL: &url.URL{
					RawQuery: tt.queryParams.Encode(),
				},
			}

			tt.middleware(ctx)

			// Check for expected headers
			for headerKey, expectedValue := range tt.expectedHeaders {
				if gotValue := recorder.Header().Get(headerKey); gotValue != expectedValue {
					t.Errorf("Expected header %s to be %s, got %s", headerKey, expectedValue, gotValue)
				}
			}
		})
	}
}

func TestConstructHeader(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Single letter",
			input:    "a",
			expected: "X-A",
		},
		{
			name:     "Multiple letters",
			input:    "page",
			expected: "X-Page",
		},
		{
			name:     "Already capitalized",
			input:    "Page",
			expected: "X-Page",
		},
		{
			name:     "Mixed case",
			input:    "pageSize",
			expected: "X-PageSize",
		},
		{
			name:     "With numbers",
			input:    "page2",
			expected: "X-Page2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pagination.ConstructHeader(tt.input)
			if got != tt.expected {
				t.Errorf("ConstructHeader(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
