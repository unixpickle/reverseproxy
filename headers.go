package reverseproxy

import (
	"net/http"
	"strings"
)

// requestHeaders adds "x-forwarded-*" headers to a request while removing
// hop-by-hop headers.
func requestHeaders(r *http.Request, host string, ws bool) http.Header {
	result := http.Header{}

	// Copy all regular headers.
	for header, values := range r.Header {
		if !isHopByHop(header, ws) && !isForwardedHeader(header) {
			result[header] = values
		}
	}

	// Get the specific forwarded headers for this hop.
	forwarded := map[string]string{"For": extractRemoteIP(r.RemoteAddr),
		"Host": r.Host, "Proto": "http"}
	if r.URL.Scheme != "" {
		forwarded["Proto"] = r.URL.Scheme
	}

	// Old values remain in a comma-separated list.
	for key, value := range forwarded {
		header := "X-Forwarded-" + key
		if existing, ok := r.Header[header]; ok {
			value = existing[0] + ", " + value
		}
		result[header] = []string{value}
	}

	return result
}

// responseHeaders rewrites the response headers from an HTTP proxy target after
// removing hop-by-hop headers.
func responseHeaders(r *http.Response, host, scheme string,
	ws bool) http.Header {
	result := http.Header{}
	headers := r.Header

	// Copy all regular headers
	for header, values := range headers {
		if !isHopByHop(header, ws) {
			result[header] = values
		}
	}

	// TODO: rewrite the Host header if possible

	return result
}

// extractRemoteIP takes out the port number from a RemoteAddr.
func extractRemoteIP(remoteAddr string) string {
	if strings.HasPrefix(remoteAddr, "[") {
		// Extract IPv6 addresses from "[addr]:port"
		return strings.Split(remoteAddr[1:], "]")[0]
	} else {
		// Extract IPv4 addresses from "addr:port"
		return strings.Split(remoteAddr, ":")[0]
	}
}

// isHopByHop checks if a header is a hop-by-hop header that should be removed
// from the proxied request.
// If websocket is true, the Connection and Upgrade headers won't be considered
// as hop-by-hop headers.
func isHopByHop(header string, ws bool) bool {
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

// isForwardedHeader returns true if and only if the passed header is
// "X-Forwarded-" + suffix where suffix is "For", "Host", or "Proto".
func isForwardedHeader(header string) bool {
	return header == "X-Forwarded-For" || header == "X-Forwarded-Proto" ||
		header == "X-Forwarded-Host"
}
