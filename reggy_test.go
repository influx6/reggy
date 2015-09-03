package reggy

import (
	"log"
	"testing"
)

func TestPriority(t *testing.T) {
	if CheckPriority(`/admin/id`) != 0 {
		t.Fatal(`/admin/id priority is not 2`)
	}
	if CheckPriority(`/admin/:id`) != 2 {
		t.Fatal(`/admin/:id priority is not 2`)
	}

	if CheckPriority(`/admin/{id:[\d+]}/:name`) != 1 {
		t.Fatal(`/admin/:id priority is not 1`)
	}
}

func TestEndless(t *testing.T) {
	if !IsEndless(`/admin/id/*`) {
		t.Fatal(`/admin/id/* is not endless`)
	}
	if !IsEndless(`/admin/id*`) {
		t.Fatal(`/admin/id* is not endless`)
	}
}

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

	state, param := r.Validate(`/name/12`)

	if !state {
		t.Fatalf("incorrect pattern: %+s %t", param, state)
	}

}

func TestClassicMux(t *testing.T) {
	cpattern := `/name/{id:[\d+]}/`
	r := CreateClassic(cpattern)

	if r == nil {
		t.Fatalf("invalid array: %+s", r)
	}

	state, param := r.Validate(`/name/12/d`)
	log.Printf("Match: %t %+s", state, param)

	if !state {
		t.Fatalf("incorrect pattern: %+s %t", param, state)
	}

}
