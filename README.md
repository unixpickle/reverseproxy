# reverseproxy

**reverseproxy** is a lightweight HTTP reverse-proxy implementation.

# Overview

## What is a reverse proxy?

Reverse HTTP proxies are used for everything from load balancing to shared hosting. A reverse proxy transparently forwards connections from one server to another.

Take an example. You might connect to my server through the URL [http://aqnichol.com](http://aqnichol.com), but my server uses a reverse proxy to forward that connection to http://localhost:8080. Even though an outside client cannot connect to that port on my server, the reverse proxy can. This allows me to run different server applications on different ports on my server and still serve all of them through port 80.

## The API

With this API, you can forward an HTTP connection directly through a proxy. As can be seen in [demo/proxy_apple.go](demo/proxy_apple.go), the actual forwarding is one line of code:

    func handler(w http.ResponseWriter, r *http.Request) {
    	reverseproxy.Proxy(w, r, rule, false)
    }

The third argument, `rule`, is a `reverseproxy.Rule` object. This object has the following fields and can be easily created:

    type Rule struct {
    	SourceHost string `json:source_host`
    	SourcePath string `json:source_path`
    	DestHost   string `json:dest_host`
    	DestPath   string `json:dest_path`
    	DestScheme string `json:dest_scheme`

    	CaseSensitiveHost bool `json:case_sensitive_host`
    	CaseSensitivePath bool `json:case_sensitive_path`
    	CleanRequestPath  bool `json:clean_request_path`
    }

Let it be noted that the `Rule` type includes annotations so that it can be serialized directly to JSON. While the reverseproxy library does not serialize `Rule` objects itself, a user may want to.

# Extra features

Most HTTP reverse proxies do not allow the target URL to have a path. For example, you could forward "aqnichol.com" to "localhost:1337", but you could not forward "aqnichol.com/foo" to "localhost:1337/bar". The reverseproxy API makes this is possible. However, this does have some caveats. Your HTML, JavaScript, and CSS code must understand that it is being proxied in an unusual way. For example, if you proxy "aqnichol.com/foo/bar" to "localhost:1337", the absolute path "/" will have a very different meaning on "localhost:1337" than it has on "aqnichol.com/foo/bar".

# Demonstrations

The [demo](demo) folder contains a few programs which use this library. The folder also contains a [README](demo/README.md) with more detailed information on each demonstration.

# License

**reverseproxy** is licensed under the BSD 2-clause license. See [LICENSE](LICENSE).

```
Copyright (c) 2014-2015, Alex Nichol.
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer. 
2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
```