package pagination_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
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
		expectError    bool
		errorContains  string
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
			true,
			"page parameter must be an integer",
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
			true,
			"size parameter must be an integer",
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
			true,
			"page number must be positive",
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
			true,
			"size must be between",
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
			true,
			"size must be between",
		},
		{
			"Default Handling",
			pagination.New(),
			url.Values{},
			1,
			10,
			"",
			"",
			false,
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
			false,
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
			false,
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
			false,
			"",
		},
		{
			"Invalid Custom Page Param - Bad Request",
			pagination.New(
				pagination.WithPageText("offset"),
			),
			url.Values{
				"offset": {"-1"},
			},
			0,
			0,
			"offset",
			"",
			true,
			"offset number must be positive",
		},
		{
			"Invalid Custom Size Param - Bad Request",
			pagination.New(
				pagination.WithSizeText("limit"),
				pagination.WithMinPageSize(1),
				pagination.WithMaxPageSize(25),
			),
			url.Values{
				"limit": {"-1"},
			},
			0,
			0,
			"",
			"limit",
			true,
			"limit must be between 1 and 25",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			// Create a test context
			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)

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
			gotSize := ctx.GetInt(tt.customSizeText)

			// Check if the page and pageSize are set correctly or if an error is expected
			if tt.expectError {
				if gotPage != 0 || gotSize != 0 {
					t.Errorf("Expected error, but got page %d and size %d", gotPage, gotSize)
				}
				if !strings.Contains(recorder.Body.String(), tt.errorContains) {
					t.Errorf("Expected error message to contain '%s', but got '%s'", tt.errorContains, recorder.Body.String())
				}
			} else {
				if gotPage != tt.expectedPage {
					t.Errorf("Expected page %d, got %d", tt.expectedPage, gotPage)
				}
				if gotSize != tt.expectedSize {
					t.Errorf("Expected size %d, got %d", tt.expectedSize, gotSize)
				}
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
			name:       "Default headers without prefix",
			middleware: pagination.New(pagination.WithHeaderPrefix("")),
			queryParams: url.Values{
				"page": {"2"},
				"size": {"20"},
			},
			expectedHeaders: map[string]string{
				"page": "2",
				"size": "20",
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
				"X-limit":  "15",
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
			middleware:  pagination.New(pagination.WithHeaderPrefix("X-")),
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
