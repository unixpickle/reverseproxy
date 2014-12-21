package reverseproxy

import (
	pathlib "path"
	"strings"
)

// PathComponents splits a path into its components.
// None of the returned path components will include a "/".
// If the empty string or "/" is passed, an empty list will be returned.
// Redundant slashes are always ignored.
func PathComponents(path string) []string {
	comps := strings.Split(path, "/")
	// Remove any empty components (incase they do "//Users////alex/" or some
	// such trickery)
	for i := 0; i < len(comps); i++ {
		if len(comps[i]) == 0 {
			for j := i; j < len(comps)-1; j++ {
				comps[j] = comps[j+1]
			}
			comps = comps[0 : len(comps)-1]
			i--
		}
	}
	return comps
}

// JoinComponents joins path components by adding slashes between them.
// If absolute is true, the returned path will begin with a leading /.
// The result will not end with a "/" unless an empty list is passed and an
// absolute result is requested.
func JoinComponents(comps []string, absolute bool) string {
	res := strings.Join(comps, "/")
	if absolute {
		return "/" + res
	} else {
		return res
	}
}

// RelativePath computes the shortest way to get from a given path to a
// different path.
// This is most useful when both paths are absolute, for example going from
// "/home/bill/Downloads" to "/home/joe/Downloads" would yield
// "../../joe/Downloads".
// If both paths are absolute, the shortest relative path to get from the source
// to the destination is returned.
// If either one of the paths is not absolute, the destination path will be
// returned.
func RelativePath(from, to string, caseSensitive bool) string {
	fromAbs := pathlib.IsAbs(from)
	toAbs := pathlib.IsAbs(to)
	if !fromAbs || !toAbs {
		return to
	}
	fromComps := PathComponents(from)
	toComps := PathComponents(to)
	// See how many identical components there are
	common := 0
	for i, x := range fromComps {
		if i >= len(toComps) {
			break
		}
		if caseSensitive {
			if toComps[i] != x {
				break
			}
		} else if strings.ToLower(toComps[i]) != strings.ToLower(x) {
			break
		}
		common++
	}
	numDots := len(fromComps) - common
	numComps := len(toComps) - common
	res := make([]string, 0, numDots+numComps)
	for i := 0; i < numDots; i++ {
		res = append(res, "..")
	}
	for i := 0; i < numComps; i++ {
		res = append(res, toComps[common+i])
	}
	return JoinComponents(res, false)
}

// PathContains determines whether a path contains another path.
// If the container is relative while the containee is absolute, false is
// returned.
// If the container is absolute while the containee is relative, true is
// returned.
func PathContains(container, containee string, caseSensitive bool) bool {
	superAbs := pathlib.IsAbs(container)
	subAbs := pathlib.IsAbs(containee)
	if !superAbs && subAbs {
		return false
	} else if superAbs && !subAbs {
		return true
	}
	superComps := PathComponents(container)
	subComps := PathComponents(containee)
	if len(superComps) > len(subComps) {
		return false
	}
	for i, comp := range superComps {
		if caseSensitive {
			if comp != subComps[i] {
				return false
			}
		} else {
			lc := strings.ToLower(comp)
			if lc != strings.ToLower(subComps[i]) {
				return false
			}
		}
	}
	return true
}
