package reggy

import (
	"regexp"
	"strings"
)

/*
pattern: /name/{id:[/\d+/]}/log/{date:[/\w+\W+/]}
*/

var (
	allSlashes = regexp.MustCompile(`/+`)
	paramd     = regexp.MustCompile(`^{[\w\W]+}$`)
)

type BoolFunc func(i interface{}) bool
type MapFunc map[string]BoolFunc
type MapString map[string]string
type FunctionalList []*FunctionalMatcher
type ClassicList []*ClassicMatcher
type MapGeneric map[string]interface{}

type Matchable interface {
	validatePattern(n string)
}

type ClassicMatcher struct {
	*regexp.Regexp
	Original string
	param    bool
}

func (f *ClassicMatcher) Validate(i interface{}) bool {
	rs := i.(string)
	return f.MatchString(rs)
}

type FunctionalMatcher struct {
	Fn       func(n interface{}) bool
	Original string
	param    bool
}

func (f *FunctionalMatcher) String() string {
	return f.Original
}

func (f *FunctionalMatcher) Validate(i interface{}) bool {
	return f.Fn(i)
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

//GenerateClassicMatcher returns a *ClassicMatcher based on a pattern part
func GenerateClassicMatcher(val string) *ClassicMatcher {
	if paramd.MatchString(val) {
		part := strings.Split(removeCurly(val), ":")
		mrk := regexp.MustCompile(removeBracket(part[1]))

		return &ClassicMatcher{
			mrk,
			part[0],
			true,
		}
	}

	return &ClassicMatcher{
		regexp.MustCompile(val),
		val,
		false,
	}
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

//ClassicPattern returns list of ClassicMatcher
func ClassicPattern(pattern string) []*ClassicMatcher {
	sr := splitPattern(pattern)
	ms := make(ClassicList, len(sr))
	for k, val := range sr {
		ms[k] = GenerateClassicMatcher(val)
	}
	return ms
}

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

type ClassicMatchMux struct {
	Pattern string
	Pix     ClassicList
}

type FunctionalMatchMux struct {
	Pattern string
	Pix     FunctionalList
}

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
				param[v.Original] = src[k]
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
				param[v.Original] = src[k]
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

func CreateClassic(pattern string) *ClassicMatchMux {
	pm := ClassicPattern(pattern)
	return &ClassicMatchMux{pattern, pm}
}

func CreateFunctional(pattern string, f MapFunc) *FunctionalMatchMux {
	pm := MappedPattern(pattern, f)
	return &FunctionalMatchMux{pattern, pm}
}
