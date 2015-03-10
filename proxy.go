package reverseproxy

import "net/http"

// ProxyRequest proxies a request to a given host.
// This will handle WebSockets intelligently.
func ProxyRequest(w http.ResponseWriter, r *http.Request, host string) {
	if r.Header.Get("Upgrade") == "websocket" {
		ProxyWebSocket(w, r, host)
	} else {
		ProxyHTTP(w, r, host)
	}
}


