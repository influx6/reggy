package reggy

import (
	"regexp"
	"testing"
)

var (
	numbOnly   = regexp.MustCompile(`\d+`)
	cpattern   = `/name/{id:[\d+]}`
	npattern   = `/name/id`
	validators = MapFunc{
		"id": func(i interface{}) bool {
			rs := i.(string)
			if numbOnly.MatchString(rs) {
				return true
			}
			return false
		},
	}
)

func TestClassicPatternCreation(t *testing.T) {
	r := ClassicPattern(cpattern)
	if r[0] == nil {
		t.Fatalf("invalid array", r)
	}
}

func TestClassicMux(t *testing.T) {
	r := CreateClassic(cpattern)

	if r == nil {
		t.Fatalf("invalid array", r)
	}

	state, param := r.Validate(`/name/12`, false)

	if !state {
		t.Fatalf("incorrect pattern", param, state)
	}

}

func TestFunctionalMux(t *testing.T) {
	r := CreateFunctional(npattern, validators)

	if r == nil {
		t.Fatalf("invalid array", r)
	}

	state, param := r.Validate(`/name/2`, false)

	if !state {
		t.Fatalf("incorrect pattern", param, state)
	}

}
