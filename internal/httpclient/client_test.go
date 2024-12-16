package httpclient

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchHTML(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse string
		headers        map[string]string
		statusCode     int
		contentType    string
		expectError    bool
	}{
		{
			name:           "Valid HTML Response",
			serverResponse: "<html><body>Hello World</body></html>",
			headers:        nil,
			statusCode:     http.StatusOK,
			contentType:    "text/html",
			expectError:    false,
		},
		{
			name:           "Non-HTML Content Type",
			serverResponse: "Not HTML content",
			headers:        nil,
			statusCode:     http.StatusOK,
			contentType:    "application/json",
			expectError:    true,
		},
		{
			name:           "Non-200 Status Code",
			serverResponse: "Error",
			headers:        nil,
			statusCode:     http.StatusInternalServerError,
			contentType:    "text/html",
			expectError:    true,
		},
		{
			name:           "Custom Header",
			serverResponse: "<html><body>Custom Header</body></html>",
			headers:        map[string]string{"X-Test-Header": "TestValue"},
			statusCode:     http.StatusOK,
			contentType:    "text/html",
			expectError:    false,
		},
		{
			name:           "Empty Response",
			serverResponse: "",
			headers:        nil,
			statusCode:     http.StatusOK,
			contentType:    "text/html",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up a test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Check for custom headers
				for key, value := range tt.headers {
					if r.Header.Get(key) != value {
						t.Errorf("Expected header %s to be %s, got %s", key, value, r.Header.Get(key))
					}
				}
				w.Header().Set("Content-Type", tt.contentType)
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			// Call FetchHTML with the test server URL and headers
			body, err := FetchHTML(server.URL, tt.headers)

			// Validate the result
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected an error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error but got: %v", err)
				}
				if body != tt.serverResponse {
					t.Errorf("Expected body to be %q, got %q", tt.serverResponse, body)
				}
			}
		})
	}
}
