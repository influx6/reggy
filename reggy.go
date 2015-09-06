package reggy

import (
	"regexp"
	"strings"
)

/*
pattern: /name/{id:[/\d+/]}/log/{date:[/\w+\W+/]}
pattern: /name/:id
*/

// MapString defines a map of strings key and value
type MapString map[string]string

// MapGeneric defines a generic stringed key map
type MapGeneric map[string]interface{}

// Matchable defines an interface for matchers
type Matchable interface {
	validatePattern(n string)
}

//ClassicMatchMux provides a class array-path matcher
type ClassicMatchMux struct {
	Pattern string
	pix     ClassicList
	endless bool
}

// CreateClassic returns a new ClassicMatchMux
func CreateClassic(pattern string) *ClassicMatchMux {
	pm := ClassicPattern(stripLastSlash(pattern))
	return &ClassicMatchMux{pattern, pm, IsEndless(pattern)}
}

// Validate validates if a string matches the pattern and returns the parameter parts
func (m *ClassicMatchMux) Validate(f string) (bool, MapGeneric) {
	var state bool
	cleaned := strings.TrimSuffix(cleanPath(f), "/")
	src := splitPattern(cleaned)

	total := len(m.pix)
	srclen := len(src)

	if !m.endless && (total < srclen || total > srclen) {
		return false, nil
	}

	param := make(MapGeneric)

	for k, v := range m.pix {
		if k >= srclen {
			state = false
			break
		}

		if v.Validate(src[k]) {
			if v.param {
				param[v.original] = src[k]
			}
			state = true
			continue
		} else {
			state = false
			break
		}
	}

	return state, param
}

// ClassicList defines a list of matchers
type ClassicList []*ClassicMatcher

// ClassicMatcher defines a single piece of pattern to be matched against
type ClassicMatcher struct {
	*regexp.Regexp
	original string
	param    bool
}

//ClassicPattern returns list of ClassicMatcher
func ClassicPattern(pattern string) []*ClassicMatcher {
	sr := splitPattern(pattern)
	ms := make(ClassicList, len(sr))
	for k, val := range sr {
		ms[k] = GenerateClassicMatcher(val)
	}
	return ms
}

//GenerateClassicMatcher returns a *ClassicMatcher based on a pattern part
func GenerateClassicMatcher(val string) *ClassicMatcher {
	id, rx, b := YankSpecial(val)
	mrk := regexp.MustCompile(rx)

	return &ClassicMatcher{
		mrk,
		id,
		b,
	}
}

// Validate validates the value against the matcher
func (f *ClassicMatcher) Validate(i interface{}) bool {
	rs := i.(string)
	return f.MatchString(rs)
}
