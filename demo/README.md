# Demos

Hopefully, these demonstrations will help you get your feet off the ground with this API.

## Testing WebSockets

If you run [websocket_demo.go](websocket_demo.go) in one command-line window, you will notice that you can connect to "http://localhost:1337/" and see how many people are actively viewing that page. You can use this to test that WebSockets forward correctly by running [general_proxy.go](general_proxy.go) in another command-line window:

    go run general_proxy.go 1338 "/" http localhost:1337 "/"

Now the URL "http://localhost:1338/" should be equivalent to "http://localhost:1337/".

**NOTE**: **websocket_demo.go** needs [Gorilla's websocket API](https://github.com/gorilla/websocket) to run correctly. You can download it with the following command:

    go get github.com/gorilla/websocket

## Proxying Apple's website!

My first test was to proxy "http://apple.com" to "http://localhost:1337"&mdash;it actually works surprisingly well.

You can run [proxy_apple.go](proxy_apple.go) to do this, or run **general_proxy.go**:

    go run general_proxy.go 1337 "/" http apple.com "/"
