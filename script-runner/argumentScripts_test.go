package scriptrunner

import (
	"errors"
	"slices"
	"testing"
)

func TestOrderedArgs(t *testing.T) {
	s := argumentScript{fixedArgs: []string{"-n", "200"}, params: []string{"FIRST", "SECOND"}}

	// map order is random, so run it a few times to be sure order is stable
	for i := 0; i < 20; i++ {
		got, err := s.orderedArgs(map[string]string{"SECOND": "b", "FIRST": "a", "EXTRA": "ignored"})
		if err != nil {
			t.Fatal(err)
		}
		if want := []string{"-n", "200", "a", "b"}; !slices.Equal(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	}

	if _, err := s.orderedArgs(map[string]string{"FIRST": "a"}); !errors.Is(err, ErrorMissingArgument) {
		t.Fatalf("expected missing argument, got %v", err)
	}
}
