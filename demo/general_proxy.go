package main

import (
	"github.com/unixpickle/reverseproxy"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatal("Usage: " + os.Args[0] + " <local port>" +
			" <remote host>")
	}
	localPort := os.Args[1]
	remoteHost := os.Args[2]
	if n, err := strconv.Atoi(localPort); err != nil || n < 0 || n > 65535 {
		log.Fatal("Invalid port number: " + localPort)
	}
	log.Print("Go to http://localhost:" + localPort)
	handler := reverseproxy.NewProxy(map[string][]string{
		"*": []string{remoteHost},
	})
	http.ListenAndServe(":"+localPort, handler)
}

