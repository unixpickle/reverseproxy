package main

import (
	"fmt"
	"github.com/unixpickle/reverseproxy"
	"net/http"
)

func main() {
	portStr := ":1338"
	table := map[string][]string{"*": []string{"www.apple.com"}}
	fmt.Println("Check out http://localhost" + portStr)
	http.Handle("/", reverseproxy.NewProxy(table))
	http.ListenAndServe(portStr, nil)
}
