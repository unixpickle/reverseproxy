package reverseproxy

import "net/http"

// Proxy proxies a request through a given rule.
// If ws is true, WebSockets will be supported.
// If any error is returned, it is the caller's responsibility to close the
// response.
// In the case of a regular HTTP request, this method will not close the
// response writer even if no error occurs.
func Proxy(w http.ResponseWriter, r *http.Request, rule *Rule, ws bool) error {
	if r.Header.Get("Upgrade") == "websocket" && ws {
		return ProxyWebsocket(w, r, rule)
	} else {
		return ProxyHTTP(w, r, rule, http.DefaultTransport)
	}
}
