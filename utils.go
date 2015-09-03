package reggy

import (
	"regexp"
	"strings"
)

var (
	specs      = regexp.MustCompile(`\W+`)
	allSlashes = regexp.MustCompile(`/+`)
	paramd     = regexp.MustCompile(`^{[\w\W]+}$`)
	picker     = regexp.MustCompile(`^:[\w\W]+$`)
	special    = regexp.MustCompile(`{\w+:[\w\W]+}`)
	anyvalue   = `[\w\W]+`
)

// HasParam returns true/false if the special pattern {:[..]} exists in the string
func HasParam(p string) bool {
	return special.MatchString(p)
}

// HasPick matches string of type :id,:name
func HasPick(p string) bool {
	return picker.MatchString(p)
}

//YankSpecial provides a means of extracting parts of form `{id:[\d+]}`
func YankSpecial(val string) (string, string, bool) {
	if HasPick(val) {
		cls := strings.TrimPrefix(val, ":")
		return cls, anyvalue, true
	}

	if !paramd.MatchString(val) {
		cls := specs.ReplaceAllString(val, "")
		return cls, cls, false
	}

	part := strings.Split(removeCurly(val), ":")
	// mrk := removeBracket(part[1])
	return part[0], removeBracket(part[1]), true
}

func removeCurly(s string) string {
	return strings.TrimPrefix(strings.TrimSuffix(s, "}"), "{")
}

func removeBracket(s string) string {
	return strings.TrimPrefix(strings.TrimSuffix(s, "]"), "[")
}

func splitPattern(c string) []string {
	return strings.Split(c, "/")
}
