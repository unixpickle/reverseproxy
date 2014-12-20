package reverseproxy

import (
	"errors"
	"io"
	"net/http"
)

// ProxyHTTP proxies an HTTP request via a given rule.
func ProxyHTTP(w http.ResponseWriter, r *http.Request, rule *Rule,
	client *http.Client) error {
	// Make sure the rule is applicable
	if !rule.MatchesRequest(r) {
		return errors.New("Request does not match rule.")
	}
	// Generate the request
	targetURL := rule.DestinationURL(r)
	req, err := http.NewRequest(r.Method, targetURL.String(), r.Body)
	if err != nil {
		return err
	}
	req.Header = RequestHeaders(r, false)
	
	// Send the request
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	
	// Write the response
	for header, value := range ResponseHeaders(res.Header, false) {
		w.Header()[header] = value
	}
	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)
	res.Body.Close()
	
	// w is automatically closed by the server
	return nil
}
