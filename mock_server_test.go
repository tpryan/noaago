package noaago

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// setupMockServer creates a mock HTTP server for testing
func setupMockServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *Client) {
	server := httptest.NewServer(handler)

	client := NewClient()
	client.dataHost = server.URL[7:] // remove http://
	client.metadataHost = server.URL[7:]

	// Override the URL builders to use http instead of https for the mock server
	// Since the client constructs https URLs by default, we need a way to force http for tests
	// Or we can just use the mock server's client which is standard, but the client implementation
	// hardcodes "https" scheme in dataUrl and metadataUrl.
	// To test this properly without modifying the main code too much to support http (which is only for tests),
	// we can make the scheme configurable or just let the test fail if it tries to connect to https://127.0.0.1...

	// Actually, the client uses `c.dataHost` in `dataUrl`.
	// If we set `c.dataHost` to include the scheme, `url.URL` might get confused if we also set Scheme="https".
	// Let's modify `client.go` slightly to support custom schemes or just handle it here.

	// A better approach for testing:
	// We can replace the `HTTPClient` with one that uses a custom Transport that intercepts requests
	// and routes them to our mock server, REGARDLESS of the URL scheme/host.

	return server, client
}

// mockTransport intercepts traffic and directs it to the mock server
type mockTransport struct {
	serverURL string
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Rewrite the request to go to the mock server
	// We keep the path and query, but change scheme and host
	u, _ := req.URL.Parse(m.serverURL)
	req.URL.Scheme = u.Scheme
	req.URL.Host = u.Host
	return http.DefaultTransport.RoundTrip(req)
}

func setupMockClient(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *Client) {
	server := httptest.NewServer(handler)

	client := NewClient()
	client.HTTPClient = &http.Client{
		Transport: &mockTransport{serverURL: server.URL},
	}

	return server, client
}
