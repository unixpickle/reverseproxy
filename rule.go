package reverseproxy

import (
	"net/http"
	"net/url"
	pathlib "path"
	"strings"
)

type Rule struct {
	SourceHost string `json:source_host`
	SourcePath string `json:source_path`
	DestHost   string `json:dest_host`
	DestPath   string `json:dest_path`
	DestScheme string `json:dest_scheme`
	
	CaseSensitivePath bool `json:case_sensitive_path`
	CaseSensitiveHost bool `json:case_sensitive_host`
}

func (r Rule) MatchesRequest(req *http.Request) bool {
	if !r.CaseSensitiveHost {
		if strings.ToLower(r.SourceHost) != strings.ToLower(req.Host) {
			return false
		}
	} else if r.SourceHost != req.Host {
		return false
	}
	// If the source path is not absolute, anything should match
	if !pathlib.IsAbs(r.SourcePath) {
		return true
	}
	return PathContains(r.SourcePath, req.URL.Path, r.CaseSensitivePath)
}

func (r Rule) DestinationURL(req *http.Request) url.URL {
	newURL := *req.URL
	newURL.Scheme = r.DestScheme
	newURL.Host = r.DestHost
	
	// Compute the new path
	if pathlib.IsAbs(r.SourcePath) {
		rel := RelativePath(r.SourcePath, req.URL.Path, r.CaseSensitivePath)
		newURL.Path = pathlib.Join(r.DestPath, rel)
	}
	
	return newURL
}
