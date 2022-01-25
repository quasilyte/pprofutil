package pprofutil

import (
	"testing"
)

func TestParseFuncName(t *testing.T) {
	tests := []struct {
		s        string
		pkgPath  string
		pkgName  string
		typeName string
		funcName string
	}{
		{"", "", "", "", ""},
		{"indexbyte", "", "", "", "indexbyte"},
		{"strings.SplitN", "strings", "strings", "", "SplitN"},
		{"testing.(*B).launch", "testing", "testing", "B", "launch"},
		{"a.(Example).b", "a", "a", "Example", "b"},
		{"a.(Example).b.func1.1", "a", "a", "Example", "b"},
		{"runtime.gcBgMarkWorker.func2", "runtime", "runtime", "", "gcBgMarkWorker"},
		{"runtime.gcMarkDone.func1.1", "runtime", "runtime", "", "gcMarkDone"},
		{"github.com/quasilyte/gogrep.(*matcher).matchNodeWithInst", "github.com/quasilyte/gogrep", "gogrep", "matcher", "matchNodeWithInst"},
		{"github.com/quasilyte/gogrep.(*matcher).matchNodeWithInst.func1", "github.com/quasilyte/gogrep", "gogrep", "matcher", "matchNodeWithInst"},
		{"github.com/quasilyte/gogrep.(*matcher).matchNodeWithInst.func1.1", "github.com/quasilyte/gogrep", "gogrep", "matcher", "matchNodeWithInst"},
		{"aaa/bbb.(CCC).fff.func1", "aaa/bbb", "bbb", "CCC", "fff"},
		{"aaa/bbb.(*CCC).fff.func1", "aaa/bbb", "bbb", "CCC", "fff"},
		{"/aaa/bbb.(*CCC).fff.func1", "/aaa/bbb", "bbb", "CCC", "fff"},
		{"/aaa/bbb.(*CCC).fff.func10.20.30", "/aaa/bbb", "bbb", "CCC", "fff"},
		{"/aaa/bbb.(*CCC).fff.func10.func20.30", "/aaa/bbb", "bbb", "CCC", "fff"},
		{"/aaa/bbb.(*CCC).fff.func10.func20.func30", "/aaa/bbb", "bbb", "CCC", "fff"},
		{"/aaa/bbb.(*CCC).func10", "/aaa/bbb", "bbb", "CCC", "func10"},
		{"/aaa/bbb.(*CCC).func10.1", "/aaa/bbb", "bbb", "CCC", "func10"},
		{"/aaa/bbb.(*CCC).func10.func1", "/aaa/bbb", "bbb", "CCC", "func10"},
		{"aaa.com/bbb.ccc/ddd.(EEE).fff", "aaa.com/bbb.ccc/ddd", "ddd", "EEE", "fff"},
		{"aaa.com/bbb.ccc/ddd.EEE.fff", "aaa.com/bbb.ccc/ddd", "ddd", "EEE", "fff"},
		{"aaa.com/bbb.ccc/ddd.EEE.fff.func10", "aaa.com/bbb.ccc/ddd", "ddd", "EEE", "fff"},
		{"aaa.com/bbb.ccc/ddd.EEE.fff.func10.1", "aaa.com/bbb.ccc/ddd", "ddd", "EEE", "fff"},
		{"aaa.com/bbb.ccc/ddd.EEE.fff.func10.func1", "aaa.com/bbb.ccc/ddd", "ddd", "EEE", "fff"},
		{"reflectlite.flag.kind", "reflectlite", "reflectlite", "flag", "kind"},
		{"internal/reflectlite.flag.kind", "internal/reflectlite", "reflectlite", "flag", "kind"},

		// Handling ambiguous cases.
		// See https://groups.google.com/g/golang-nuts/c/sAY9RDSfZX8
		{"pkg.sym.func1", "pkg", "pkg", "", "sym"},
		{"pkg.(sym).func1", "pkg", "pkg", "sym", "func1"},
		{"aaa.com/bbb/pkg.sym.func1", "aaa.com/bbb/pkg", "pkg", "", "sym"},
		{"aaa.com/bbb/pkg.(sym).func1", "aaa.com/bbb/pkg", "pkg", "sym", "func1"},
		{"pkg.sym.func1.2", "pkg", "pkg", "", "sym"},
		{"pkg.(sym).func1.2", "pkg", "pkg", "sym", "func1"},
		{"aaa.com/bbb/pkg.sym.func1.2", "aaa.com/bbb/pkg", "pkg", "", "sym"},
		{"aaa.com/bbb/pkg.(sym).func1.2", "aaa.com/bbb/pkg", "pkg", "sym", "func1"},
		{"pkg.sym.func1.func2", "pkg", "pkg", "", "sym"},
		{"pkg.(sym).func1.func2", "pkg", "pkg", "sym", "func1"},
		{"aaa.com/bbb/pkg.sym.func1.func2", "aaa.com/bbb/pkg", "pkg", "", "sym"},
		{"aaa.com/bbb/pkg.(sym).func1.func2", "aaa.com/bbb/pkg", "pkg", "sym", "func1"},
	}

	for _, test := range tests {
		have := parseFuncName(test.s)
		want := Symbol{
			PkgPath:  test.pkgPath,
			PkgName:  test.pkgName,
			TypeName: test.typeName,
			FuncName: test.funcName,
		}
		if have != want {
			t.Fatalf("parseFuncName(%q):\nhave: %#v\nwant: %#v", test.s, have, want)
		}
	}
}
