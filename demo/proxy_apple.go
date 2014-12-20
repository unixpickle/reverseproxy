package main

import (
	"fmt"
	"github.com/unixpickle/reverseproxy"
	"net/http"
)

var rule *reverseproxy.Rule

func handler(w http.ResponseWriter, r *http.Request) {
	reverseproxy.Proxy(w, r, rule, false)
}

func main() {
	portStr := ":1337"
	rule = &reverseproxy.Rule{"localhost" + portStr, "", "www.apple.com", "",
		"http", false, false, false}
	fmt.Println("Check out http://localhost" + portStr)
	http.HandleFunc("/", handler)
	http.ListenAndServe(portStr, nil)
}
