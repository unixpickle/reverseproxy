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
	if n, err := strconv.Atoi(os.Args[1]); err != nil || n < 0 || n > 65535 {
		log.Fatal("Invalid port number: " + os.Args[1])
	}
	log.Print("Go to http://localhost:" + os.Args[1] + os.Args[2])
	handler := new(MyHandler)
	handler.rule = reverseproxy.Rule{"localhost:" + os.Args[1], os.Args[2],
		os.Args[4], os.Args[5], os.Args[3], false, true, false}
	http.ListenAndServe(":"+os.Args[1], handler)
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
