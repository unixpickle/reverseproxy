package reverseproxy

import (
	"io"
	"net/http"
)

// ProxyHTTP proxies an HTTP request to a given destination host.
// This will not handle WebSockets intelligently.
func ProxyHTTP(w http.ResponseWriter, r *http.Request, host string) {
	// Generate the request for the proxy destination.
	targetURL := *r.URL
	targetURL.Host = host
	targetURL.Scheme = "http"
	req, err := http.NewRequest(r.Method, targetURL.String(), r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	req.Header = requestHeaders(r, host, false)

	// Send the request
	res, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	// Write the response
	respHead := responseHeaders(res, r.URL.Host, r.URL.Scheme, false)
	for header, value := range respHead {
		w.Header()[header] = value
	}
	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)
	res.Body.Close()
}
