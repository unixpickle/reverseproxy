package reverseproxy

import (
	"errors"
	"net/http"
)

// Proxy a request through a given rule.
// If ws is true, WebSockets will be supported.
func Proxy(w http.ResponseWriter, r *http.Request, rule *Rule, ws bool) error {
	if !rule.MatchesRequest(r) {
		return errors.New("Request does not match rule.")
	}
	if r.Header.Get("Upgrade") == "websocket" && ws {
		return ProxyWebsocket(w, r, rule)
	} else {
		return ProxyHTTP(w, r, rule, &http.Client{})
	}
}
