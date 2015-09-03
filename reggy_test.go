package reggy

import (
	"regexp"
	"testing"
)

func TestPicker(t *testing.T) {
	if HasPick(`id`) {
		t.Fatal(`/admin/id has no picker`)
	}
	if !HasPick(`:id`) {
		t.Fatal(`/admin/:id has picker`)
	}
}

func TestSpecialChecker(t *testing.T) {
	if !HasParam(`/admin/{id:[\d+]}`) {
		t.Fatal(`/admin/{id:[\d+]} is special`)
	}
	if HasParam(`/admin/id`) {
		t.Fatal(`/admin/id is not special`)
	}
	if !HasParam(`/admin/:id`) {
		t.Fatal(`/admin/:id is special`)
	}
	if HasKeyParam(`/admin/:id`) {
		t.Fatal(`/admin/:id is special`)
	}
}

func TestClassicPatternCreation(t *testing.T) {
	cpattern := `/name/{id:[\d+]}`
	r := ClassicPattern(cpattern)
	if r[0] == nil {
		t.Fatalf("invalid array %+s", r)
	}
}

func TestClassicMuxPicker(t *testing.T) {
	cpattern := `/name/:id`
	r := CreateClassic(cpattern)

	if r == nil {
		t.Fatalf("invalid array: %+s", r)
	}

	state, param := r.Validate(`/name/12`, false)

	if !state {
		t.Fatalf("incorrect pattern: %+s %t", param, state)
	}

}

func TestClassicMux(t *testing.T) {
	cpattern := `/name/{id:[\d+]}`
	r := CreateClassic(cpattern)

	if r == nil {
		t.Fatalf("invalid array: %+s", r)
	}

	state, param := r.Validate(`/name/12`, false)

	if !state {
		t.Fatalf("incorrect pattern: %+s %t", param, state)
	}

}

func TestFunctionalMux(t *testing.T) {
	npattern := `/name/id`
	numbOnly := regexp.MustCompile(`\d+`)
	validators := MapFunc{
		"id": func(i interface{}) bool {
			rs := i.(string)
			if numbOnly.MatchString(rs) {
				return true
			}
			return false
		},
	}

	r := CreateFunctional(npattern, validators)

	if r == nil {
		t.Fatalf("invalid array: %+s", r)
	}

	state, param := r.Validate(`/name/2`, false)

	if !state {
		t.Fatalf("incorrect pattern: %+s %t", param, state)
	}

}
