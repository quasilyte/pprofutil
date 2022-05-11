package pprofutil

import (
	"strings"

	"github.com/google/pprof/profile"
)

// WalkSamples is a convenience function for looping over time samples.
// For instance, it can be used to loop over time/ns data.
//
// It calls a given function for every non-empty location
// with all "lines" combined into a stack.
func WalkSamples(p *profile.Profile, visit func(s Sample)) {
	var stack []profile.Line
	for _, s := range p.Sample {
		if len(s.Location) == 0 || len(s.Location[0].Line) == 0 {
			continue
		}
		stack = stack[:0] // reuse memory for every stack
		for _, loc := range s.Location {
			stack = append(stack, loc.Line...)
		}
		visit(Sample{
			Value: s.Value[1],
			Stack: stack,
		})
	}
}

// Sample holds data of the current sample in a convenient way.
//
// See WalkSamples.
type Sample struct {
	Value int64

	// Stack contains line entries for this sample.
	// Stack[0] is a current function.
	// Stack[1:] is a slice of the callers.
	//
	// The stack is never empty if this Sample object is provided
	// by the WalkSamples function.
	Stack []profile.Line
}

// Symbol represents a parsed profile.proto function name.
type Symbol struct {
	PkgPath  string
	PkgName  string
	TypeName string
	FuncName string
}

// ParseFuncName parses a profile.proto Function.Name.
func ParseFuncName(s string) Symbol {
	return parseFuncName(s)
}

func parseFuncName(s string) Symbol {
	rest := s
	lastSlash := strings.LastIndexByte(s, '/')
	if lastSlash != -1 {
		rest = s[lastSlash+len("/"):]
	}

	i := strings.IndexByte(rest, '.')
	if i == -1 {
		return Symbol{FuncName: s}
	}

	var sym Symbol
	sym.PkgName = rest[:i]
	if lastSlash != -1 {
		sym.PkgPath = s[:lastSlash+i+len("/")]
	} else {
		sym.PkgPath = sym.PkgName
	}
	rest = rest[i+1:]

	// A simple case: we have () that surround the receiver.
	if strings.HasPrefix(rest, "(") {
		offset := 1
		if strings.HasPrefix(rest, "(*") {
			offset++
		}
		rparen := strings.IndexByte(rest, ')')
		if rparen == -1 {
			return Symbol{}
		}
		sym.TypeName = rest[offset:rparen]
		resultFuncName := rest[rparen+len(")."):]
		sym.FuncName = trimLambdaSuffix(resultFuncName)
		return sym
	}

	// Possible ambiguity: if symbol looks like `x.func1`, there are at least two
	// possible interpretations:
	// 1. `x` is a type name, `func1` is a method name
	// 2. `x` is a func name, `func1` means "first anonymous function inside x"
	// Since `func%d` is not a common method name, we try to resolve this ambiguity
	// in favor of (2).
	// See https://groups.google.com/g/golang-nuts/c/sAY9RDSfZX8
	rest = trimLambdaSuffix(rest)
	if dotPos := strings.IndexByte(rest, '.'); dotPos != -1 {
		sym.TypeName = rest[:dotPos]
		sym.FuncName = rest[dotPos+1:]
		return sym
	}

	sym.FuncName = rest
	return sym
}

func trimLambdaSuffix(s string) string {
	end := len(s) - 1
	for {
		i := end
		for s[i] >= '0' && s[i] <= '9' {
			i--
		}
		found := false
		if strings.HasSuffix(s[:i+1], ".func") {
			i -= len(".func")
			found = true
		} else if s[i] == '.' {
			i--
			found = true
		}
		if !found {
			break
		}
		end = i
	}
	return s[:end+1]
}
