package reverseproxy

import (
	"math/rand"
	"net/http"
	"sync"
)

// ProxyRequest proxies a request to a given host.
// This will handle WebSockets intelligently.
func ProxyRequest(w http.ResponseWriter, r *http.Request, host string) {
	if r.Header.Get("Upgrade") == "websocket" {
		ProxyWebSocket(w, r, host)
	} else {
		ProxyHTTP(w, r, host)
	}
}

func proxyRequest(w http.ResponseWriter, r *http.Request, hosts []string,
	indices []int) {
	if r.Header.Get("Upgrade") == "websocket" {
		proxyWebSocket(w, r, hosts, indices)
	} else {
		proxyHTTP(w, r, hosts, indices)
	}
}

// A RuleTable associates zero or more target hosts with a destination hosts.
type RuleTable map[string][]string

// Copy returns a deep copy of a RuleTable.
func (r RuleTable) Copy() RuleTable {
	res := RuleTable{}
	for key, val := range r {
		res[key] = make([]string, len(val))
		copy(res[key], val)
	}
	return res
}

// A Proxy handles HTTP requests and forwards them through a rule table.
type Proxy struct {
	lock  sync.RWMutex
	rules RuleTable
}

// NewProxy creates a Proxy with an initial RuleTable.
func NewProxy(rules RuleTable) *Proxy {
	return &Proxy{sync.RWMutex{}, rules.Copy()}
}

// RuleTable returns the Proxy's current rule table.
func (p *Proxy) RuleTable() RuleTable {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.rules.Copy()
}

// ServeHTTP routes a request to the HTTP server.
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Pick a host
	p.lock.RLock()
	rules, found := p.rules[r.URL.Host]
	if !found {
		// This is the "no forward rule" rule
		rules, found = p.rules["*"]
	}
	if !found {
		p.lock.RUnlock()
		w.Write([]byte("No forward rule found."))
		return
	}
	hosts := make([]string, len(rules))
	copy(hosts, rules)
	p.lock.RUnlock()

	// For now, load balancing will be random.
	proxyRequest(w, r, hosts, rand.Perm(len(hosts)))
}

// SetRuleTable updates the rule table used by the Proxy.
func (p *Proxy) SetRuleTable(t RuleTable) {
	p.lock.Lock()
	p.rules = t
	p.lock.Unlock()
}
