package reverseproxy

import "testing"

func TestPathComponents(t *testing.T) {
	testPathComponentsCase(t, "/foo/bar/baz", []string{"foo", "bar", "baz"})
	testPathComponentsCase(t, "", []string{})
	testPathComponentsCase(t, "/", []string{})
	testPathComponentsCase(t, "//", []string{})
	testPathComponentsCase(t, "///", []string{})
	testPathComponentsCase(t, "foo///bar/baz/", []string{"foo", "bar", "baz"})
}

func testPathComponentsCase(t *testing.T, str string, comps []string) {
	res := PathComponents(str)
	if len(res) != len(comps) {
		t.Error("Invalid path components for '" + str + "'")
		return
	}
	for i, x := range comps {
		if x != res[i] {
			t.Error("Invalid path components for '" + str + "'")
			return
		}
	}
}

func TestJoinComponents(t *testing.T) {
	testJoinComponentsCase(t, "", false)
	testJoinComponentsCase(t, "/", true)
	testJoinComponentsCase(t, "foo", false, "foo")
	testJoinComponentsCase(t, "foo/bar/baz", false, "foo", "bar", "baz")
	testJoinComponentsCase(t, "/foo/bar/baz", true, "foo", "bar", "baz")
}

func testJoinComponentsCase(t *testing.T, answer string, abs bool,
	comps ...string) {
	if JoinComponents(comps, abs) != answer {
		t.Error("Invalid joined components for", comps)
	}
}

func TestRelativePath(t *testing.T) {
	testRelativePathCase(t, true, "/foo", "/bar", "../bar")
	testRelativePathCase(t, true, "/foo", "/bar", "../bar")
	testRelativePathCase(t, true, "/foo/bar/bill", "/foo/bob/bar",
		"../../bob/bar")
	testRelativePathCase(t, true, "/foo/Bar/bill", "/foo/bar/bill/gates",
		"../../bar/bill/gates")
	testRelativePathCase(t, false, "/foo/Bar/bill", "/foo/bar/bill/gates",
		"gates")
	testRelativePathCase(t, true, "/foo", "/foo", "")
	testRelativePathCase(t, true, "/foo", "/foo/", "")
	testRelativePathCase(t, true, "foo", "/bar", "/bar")
	testRelativePathCase(t, true, "/foo", "bar", "bar")
}

func testRelativePathCase(t *testing.T, cs bool, from, to, answer string) {
	if RelativePath(from, to, cs) != answer {
		t.Error("Invalid relative path from '" + from + "' to '" + to + "'")
	}
}

func TestPathContains(t *testing.T) {
	testPathContainsCase(t, true, "/foo/bar", "/foo/bar", true)
	testPathContainsCase(t, true, "/foo/bar", "/foo/bar/", true)
	testPathContainsCase(t, true, "/foo/bar/", "/foo/bar", true)
	testPathContainsCase(t, true, "foo/bar", "foo/bar", true)
	testPathContainsCase(t, true, "foo/bar", "foo/bar/", true)
	testPathContainsCase(t, true, "foo/bar/", "foo/bar", true)
	testPathContainsCase(t, true, "/test", "/foo/bar", false)
	testPathContainsCase(t, true, "test", "foo/bar", false)
	testPathContainsCase(t, true, "abc", "/abc", false)
	testPathContainsCase(t, true, "/abc", "anything", true)
	testPathContainsCase(t, true, "/test", "/", false)
	testPathContainsCase(t, true, "/foo/Bar/baz", "/foo/bar/baz/joe", false)
	testPathContainsCase(t, false, "/foo/Bar/baz", "/foo/bar/baz/joe", true)
}

func testPathContainsCase(t *testing.T, cs bool, super, sub string, res bool) {
	if PathContains(super, sub, cs) != res {
		if res {
			t.Error("Path '" + super + "' should contain '" + sub + "'")
		} else {
			t.Error("Path '" + super + "' should not contain '" + sub + "'")
		}
	}
}
