package reverseproxy

import "net/http"

// RequestHeaders adds "x-forwarded-*" headers to a request while removing
// hop-by-hop headers.
func RequestHeaders(r *http.Request, ws bool) http.Header {
	result := http.Header{}

	// Copy all regular headers
	for header, values := range r.Header {
		if !IsHopByHop(header, ws) && !IsForwardedHeader(header) {
			result[header] = values
		}
	}

	// Get the specific forwarded headers for this hop
	forwarded := map[string]string{"For": r.RemoteAddr, "Host": r.Host,
		"Proto": "http"}
	if r.URL.Scheme != "" {
		forwarded["Proto"] = r.URL.Scheme
	}
	
	// Old values remain in a comma-separated
	for key, value := range forwarded {
		header := "X-Forwarded-" + key
		if existing, ok := r.Header[header]; ok {
			value = existing[0] + ", " + value
		}
		result[header] = []string{value}
	}
	
	return result
}

// ResponseHeaders rewrites the response headers from an HTTP proxy target after
// removing hop-by-hop headers.
func ResponseHeaders(headers http.Header, ws bool) http.Header {
	result := http.Header{}

	// Copy all regular headers
	for header, values := range headers {
		if !IsHopByHop(header, ws) {
			result[header] = values
		}
	}

	return result
}

// IsHopByHop checks if a header is a hop-by-hop header that should be removed
// from the proxied request.
// If websocket is true, the Connection and Upgrade headers won't be considered
// as hop-by-hop headers.
func IsHopByHop(header string, ws bool) bool {
	hopByHop := []string{"Keep-Alive", "Proxy-Authenticate",
		"Proxy-Authorization", "Te", "Transfer-Encoding"}
	if !ws {
		hopByHop = append(hopByHop, "Upgrade", "Connection")
	}
	for _, val := range hopByHop {
		if val == header {
			return true
		}
	}
	return false
}

// IsForwardedHeader returns true if and only if the passed header is
// "X-Forwarded-" + suffix where suffix is "For", "Host", or "Proto".
func IsForwardedHeader(header string) bool {
	return header == "X-Forwarded-For" || header == "X-Forwarded-Proto" ||
		header == "X-Forwarded-Host"
}