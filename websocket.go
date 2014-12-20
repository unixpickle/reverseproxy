package reverseproxy

import (
	"errors"
	"net/http"
)

// ProxyWebsocket proxies a websocket via a given rule.
func ProxyWebsocket(w http.ResponseWriter, r *http.Request, rule *Rule) error {
	return errors.New("Not yet implemented")
}
