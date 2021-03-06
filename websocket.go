package reverseproxy

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"sync"
)

// ProxyWebSocket proxies a WebSocket request to a given host.
// This should only be used if the request had an "Upgrade: websocket" header.
func ProxyWebSocket(w http.ResponseWriter, r *http.Request, host string) {
	proxyWebSocket(w, r, []string{host}, []int{0})
}

func proxyWebSocket(w http.ResponseWriter, r *http.Request, hosts []string,
	indices []int) {
	// Make sure we can hijack the ResponseWriter.
	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "could not hijack connection", http.StatusBadGateway)
		return
	}

	// Open a raw connection to a destination host
	var conn net.Conn
	var err error
	var host string
	for _, i := range indices {
		host = hosts[i]
		conn, err = net.Dial("tcp", host)
		if err == nil {
			break
		}
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	// Update the headers and send the request to the target host
	r.Header = requestHeaders(r, host, true)
	r.Host = host
	if err := r.Write(conn); err != nil {
		conn.Close()
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	// Hijack the response and proxy the data.
	hjConn, hjStream, err := hj.Hijack()
	if err != nil {
		conn.Close()
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	bidirectionalPipe(&flushWriter{hjStream}, conn, func() {
		hjStream.Flush()
		conn.Close()
		hjConn.Close()
	})
}

// bidirectionalPipe pipes two io.ReadWriters into each other.
// When one io.ReadWriter is closed, closeBoth() is called.
// This method only returns once both streams have been closed.
func bidirectionalPipe(a io.ReadWriter, b io.ReadWriter, closeBoth func()) {
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

type flushWriter struct {
	b *bufio.ReadWriter
}

func (f *flushWriter) Read(buf []byte) (int, error) {
	return f.b.Read(buf)
}

func (f *flushWriter) Write(b []byte) (n int, err error) {
	n, err = f.b.Write(b)
	err1 := f.b.Flush()
	if err == nil {
		err = err1
	}
	return
}
