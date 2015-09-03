package reggy

import (
	"path"
	"regexp"
	"strings"
)

var (
	specs       = regexp.MustCompile(`\W+`)
	allSlashes  = regexp.MustCompile(`/+`)
	paramd      = regexp.MustCompile(`^{[\w\W]+}$`)
	picker      = regexp.MustCompile(`^:[\w\W]+$`)
	special     = regexp.MustCompile(`{\w+:[\w\W]+}|:[\w]+`)
	morespecial = regexp.MustCompile(`{\w+:[\w\W]+}`)
	anyvalue    = `[\w\W]+`

	//MoreSlashes this to check for more than one forward slahes
	MoreSlashes = regexp.MustCompile(`/+`)
)

//RemoveCurly removes '{' and '}' from any string
func RemoveCurly(s string) string {
	return strings.TrimPrefix(strings.TrimSuffix(s, "}"), "{")
}

//RemoveBracket removes '[' and ']' from any string
func RemoveBracket(s string) string {
	return strings.TrimPrefix(strings.TrimSuffix(s, "]"), "[")
}

//SplitPattern splits a pattern with the '/'
func SplitPattern(c string) []string {
	return strings.Split(c, "/")
}

//SplitPatternAndRemovePrefix splits a pattern with the '/'
func SplitPatternAndRemovePrefix(c string) []string {
	return strings.Split(strings.TrimPrefix(cleanPath(c), "/"), "/")
}

// CheckPriority is used to return the priority of a pattern. 0 for highest(when no parameters),1 for restricted parameters({id:[]}) and 2 for loose paramters. The first parameter catched is used for rating
func CheckPriority(patt string) int {
	sets := splitPattern(patt)

	for _, so := range sets {
		if morespecial.MatchString(so) {
			return 1
		}
		if special.MatchString(so) {
			return 2
		}
		continue
	}

	return 0
}

// cleanPath returns the canonical path for p, eliminating . and .. elements.
// Borrowed from the net/http package.
func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	np := path.Clean(p)
	// path.Clean removes trailing slash except for root;
	// put the trailing slash back if necessary.
	if p[len(p)-1] == '/' && np != "/" {
		np += "/"
	}
	return np
}

// HasKeyParam returns true/false if the special pattern {:[..]} exists in the string
func HasKeyParam(p string) bool {
	return morespecial.MatchString(p)
}

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
