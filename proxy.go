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

func proxyRequest(w http.ResponseWriter, r *http.Request, hosts []string,
	indices []int) {
	if r.Header.Get("Upgrade") == "websocket" {
		proxyWebSocket(w, r, hosts, indices)
	} else {
		proxyHTTP(w, r, hosts, indices)
	}
}
