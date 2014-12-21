package reverseproxy

import (
	"errors"
	"io"
	"net"
	"net/http"
	"sync"
)

// ProxyWebsocket proxies a websocket via a given rule.
func ProxyWebsocket(w http.ResponseWriter, r *http.Request, rule *Rule) error {
	// Make sure the rule is applicable.
	if !rule.MatchesRequest(r) {
		return errors.New("Request does not match rule.")
	}
	
	// Make sure we can hijack the ResponseWriter.
	hj, ok := w.(http.Hijacker)
	if !ok {
		return errors.New("Could not hijack connection")
	}
	
	// Open a raw connection to the destination host
	destURL := rule.DestinationURL(r)
	conn, err := net.Dial("tcp", destURL.Host)
	if err != nil {
		return err
	}
	
	// Update the headers and send the request to the target
	r.Header = RequestHeaders(r, true)
	r.Host = destURL.Host
	if err := r.Write(conn); err != nil {
		conn.Close()
		return err
	}
	
	// Hijack the response and proxy data.
	hjConn, hjStream, err := hj.Hijack()
	if err != nil {
		conn.Close()
		return err
	}
	
	BidirectionalPipe(hjStream, conn, func() {
		hjStream.Flush()
		conn.Close()
		hjConn.Close()
	})
	
	return nil
}

func BidirectionalPipe(a io.ReadWriter, b io.ReadWriter, closeBoth func()) {
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		io.Copy(b, a)
		closeBoth()
		wg.Done()
	}()
	go func() {
		io.Copy(a, b)
		closeBoth()
		wg.Done()
	}()
	wg.Wait()
}
