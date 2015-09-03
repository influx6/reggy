package reggy

import "regexp"

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

// BoolFunc defines a function that returns a bool
type BoolFunc func(i interface{}) bool

// MapFunc defines a map of boolean returning functions
type MapFunc map[string]BoolFunc

// FunctionalList defines a list of FunctionalMatchers
type FunctionalList []*FunctionalMatcher

// FunctionalMatcher defines a piece of a functional match
type FunctionalMatcher struct {
	Fn       func(n interface{}) bool
	original string
	param    bool
}

//MappedPattern returns a list of FunctionalMatcher
func MappedPattern(pattern string, f MapFunc) []*FunctionalMatcher {
	src := splitPattern(pattern)
	ms := make(FunctionalList, len(src))
	for k, val := range src {
		if fn, ok := f[val]; ok {
			if ok {
				ms[k] = GenerateFunctionalMatcher(val, fn)
			} else {
				ms[k] = GenerateFunctionalMatcher(val, nil)
			}
		} else {
			ms[k] = GenerateFunctionalMatcher(val, nil)
		}
	}

	return ms
}

//GenerateFunctionalMatcher returns a FunctionalMatcher
func GenerateFunctionalMatcher(val string, fn func(data interface{}) bool) *FunctionalMatcher {
	if fn == nil {
		return &FunctionalMatcher{
			func(i interface{}) bool {
				return i == val
			},
			val,
			false,
		}
	}

	return &FunctionalMatcher{
		fn,
		val,
		true,
	}
}

// String returns the original value
func (f *FunctionalMatcher) String() string {
	return f.original
}

// Validate checks the value against the function
func (f *FunctionalMatcher) Validate(i interface{}) bool {
	return f.Fn(i)
}

// FunctionalMatchMux provides a map like validator matcher
type FunctionalMatchMux struct {
	Pattern string
	Pix     FunctionalList
}

// CreateFunctional returns a new FunctionalMatchMux
func CreateFunctional(pattern string, f MapFunc) *FunctionalMatchMux {
	pm := MappedPattern(pattern, f)
	return &FunctionalMatchMux{pattern, pm}
}

// Validate validates if a string matches the pattern and returns the parameter parts
func (m *FunctionalMatchMux) Validate(f string, strictlen bool) (bool, MapGeneric) {
	var state bool
	src := splitPattern(f)
	param := make(MapGeneric)

	if !!strictlen {
		if len(src) != len(m.Pix) {
			state = false
			return state, param
		}
	}

	for k, v := range m.Pix {
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

//ClassicMatchMux provides a class array-path matcher
type ClassicMatchMux struct {
	Pattern string
	Pix     ClassicList
}

// CreateClassic returns a new ClassicMatchMux
func CreateClassic(pattern string) *ClassicMatchMux {
	pm := ClassicPattern(pattern)
	return &ClassicMatchMux{pattern, pm}
}

// Validate validates if a string matches the pattern and returns the parameter parts
func (m *ClassicMatchMux) Validate(f string, strictlen bool) (bool, MapGeneric) {
	var state bool
	src := splitPattern(f)
	param := make(MapGeneric)

	if !!strictlen {
		if len(src) != len(m.Pix) {
			state = false
			return state, param
		}
	}

	for k, v := range m.Pix {
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
