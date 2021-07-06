package reverseproxy

import (
	"io"
	"net/http"
)

var transport *http.Transport = &http.Transport{DisableKeepAlives: true}

// ProxyHTTP proxies an HTTP request to a given destination host.
// This will not handle WebSockets intelligently.
func ProxyHTTP(w http.ResponseWriter, r *http.Request, host string) {
	proxyHTTP(w, r, []string{host}, []int{0})
}

func proxyHTTP(w http.ResponseWriter, r *http.Request, hosts []string,
	indices []int) {
	var res *http.Response
	var err error
	for _, i := range indices {
		host := hosts[i]

		// Generate the request for the proxy destination.
		targetURL := *r.URL
		targetURL.Host = host
		targetURL.Scheme = "http"
		var req *http.Request
		req, err = http.NewRequest(r.Method, targetURL.String(), r.Body)
		if err != nil {
			continue
		}
		req.Header = requestHeaders(r, host, false)
		req.Host = r.Host

		// Send the request
		res, err = transport.RoundTrip(req)
		// If an error occurs the request's body may have been read and trying a
		// new host would be pointless. However, I do not currently check for
		// this.
		if err != nil {
			continue
		}
		break
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	// Write the response
	respHead := responseHeaders(res, r.Host, r.URL.Scheme, false)
	for header, value := range respHead {
		w.Header()[header] = value
	}
	w.WriteHeader(res.StatusCode)
	if respHead.Get("X-Content-Type-Options") == "nosniff" {
		// This header tends to be set for chunked/streaming responses.
		io.Copy(&httpFlushWriter{w: w, f: w.(http.Flusher)}, res.Body)
	} else {
		io.Copy(w, res.Body)
	}
	res.Body.Close()
}

type httpFlushWriter struct {
	w io.Writer
	f http.Flusher
}

func (h *httpFlushWriter) Write(buf []byte) (int, error) {
	n, err := h.w.Write(buf)
	h.f.Flush()
	return n, err
}
