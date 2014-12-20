package reverseproxy

import "net/http"

// Proxy a request through a given rule.
// If ws is true, WebSockets will be supported.
func Proxy(w http.ResponseWriter, r *http.Request, rule *Rule, ws bool) error {
	if r.Header.Get("Upgrade") == "websocket" && ws {
		return ProxyWebsocket(w, r, rule)
	} else {
		return ProxyHTTP(w, r, rule, &http.Client{})
	}
}
