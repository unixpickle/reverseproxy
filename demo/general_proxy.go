package main

import (
	"github.com/unixpickle/reverseproxy"
	"log"
	"net/http"
	"os"
	"strconv"
)

type MyHandler struct {
	rule reverseproxy.Rule
}

func main() {
	if len(os.Args) != 6 {
		log.Fatal("Usage: " + os.Args[0] + " <local port> <local path>" +
			" <remote protocol> <remote host> <remote path>")
	}
	localPort := os.Args[1]
	localPath := os.Args[2]
	remoteProtocol := os.Args[3]
	remoteHost := os.Args[4]
	remotePath := os.Args[5]
	if n, err := strconv.Atoi(localPort); err != nil || n < 0 || n > 65535 {
		log.Fatal("Invalid port number: " + localPort)
	}
	log.Print("Go to http://localhost:" + localPort + localPath)
	handler := new(MyHandler)
	handler.rule = reverseproxy.Rule{"localhost:" + localPort, localPath,
		remoteHost, remotePath, remoteProtocol, false, true, false}
	http.ListenAndServe(":"+localPort, handler)
}

func (self MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !self.rule.MatchesRequest(r) {
		log.Print("Request does not match: " + r.URL.String())
		w.WriteHeader(404)
		w.Write([]byte("404 not found."))
		return
	}
	log.Print("Proxying request: " + r.URL.String())
	if err := reverseproxy.Proxy(w, r, &self.rule, true); err != nil {
		log.Print("Error in " + r.URL.String() + ": " + err.Error())
	}
}
