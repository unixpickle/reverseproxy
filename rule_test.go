package reverseproxy

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestMatchesRequest(t *testing.T) {
	rule := &Rule{"google.com", "", "localhost:1337", "/google", "http",
		false, false, false}
	for _, rule.DestPath = range []string{"", "/"} {
		for _, rule.CaseSensitiveHost = range []bool{false, true} {
			ensureMatch(t, rule, "google.com/foo", true)
			ensureMatch(t, rule, "google.com", true)
			ensureMatch(t, rule, "google.com/", true)
			ensureMatch(t, rule, "Google.com/foo", !rule.CaseSensitiveHost)
			ensureMatch(t, rule, "Google.com", !rule.CaseSensitiveHost)
			ensureMatch(t, rule, "Google.com/", !rule.CaseSensitiveHost)
			ensureMatch(t, rule, "www.google.com", false)
		}
	}
	rule = &Rule{"apple.com", "/hacked", "localhost:1337", "/foobar/", "https",
		false, false, false}
	for _, rule.CaseSensitivePath = range []bool{false, true} {
		for _, rule.CleanRequestPath = range []bool{false, true} {
			ensureMatch(t, rule, "apple.com", false)
			ensureMatch(t, rule, "apple.com/", false)
			ensureMatch(t, rule, "apple.com/foo", false)
			ensureMatch(t, rule, "apple.com/bar/hacked", false)
			ensureMatch(t, rule, "apple.com/hacked", true)
			ensureMatch(t, rule, "apple.com/hacked/", true)
			ensureMatch(t, rule, "apple.com/hacked/bar", true)
			ensureMatch(t, rule, "apple.com/Hacked", !rule.CaseSensitivePath)
			ensureMatch(t, rule, "apple.com/Hacked/", !rule.CaseSensitivePath)
			ensureMatch(t, rule, "apple.com/hacked/../hacked/foo", true)
			ensureMatch(t, rule, "apple.com/Hacked/bar",
				!rule.CaseSensitivePath)
			ensureMatch(t, rule, "apple.com/foo/../hacked",
				rule.CleanRequestPath)
			ensureMatch(t, rule, "apple.com/foo/../hacked/../hacked",
				rule.CleanRequestPath)
			ensureMatch(t, rule, "apple.com/hacked/../etc/passwd",
				!rule.CleanRequestPath)
		}
	}
}

func dummyRequest(urlStr string) *http.Request {
	r := new(http.Request)
	u := new(url.URL)
	r.URL = u
	firstIdx := strings.Index(urlStr, "/")
	if firstIdx == -1 {
		r.Host = urlStr
		u.Path = ""
	} else {
		r.Host = urlStr[0:firstIdx]
		u.Path = urlStr[firstIdx:]
	}
	return r
}

func ensureMatch(t *testing.T, rule *Rule, x string, match bool) {
	r := dummyRequest(x)
	if rule.MatchesRequest(r) != match {
		if match {
			t.Error("Rule", *rule, "should match", x)
		} else {
			t.Error("Rule", *rule, "should not match", x)
		}
	}
}
