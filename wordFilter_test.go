package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLowercaseFilter(t *testing.T) {
	filter, _ := newLowercaseFilter()
	actual := filter.filter([][]string{{"AmaZON", "pRiME"}})

	diff := cmp.Diff(
		[][]string{{"amazon", "prime"}},
		actual,
	)
	if diff != "" {
		t.Errorf(diff)
	}
}

func TestStopWordFilter(t *testing.T) {
	filter, _ := newStopWordFilter([]string{"i", "a"})
	actual := filter.filter([][]string{{"i", "have", "a", "pen"}})

	diff := cmp.Diff(
		[][]string{{"have", "pen"}},
		actual,
	)
	if diff != "" {
		t.Errorf(diff)
	}
}

func TestStemmerFilter(t *testing.T) {
	filter, _ := newStemmerFilter()
	actual := filter.filter([][]string{{"it", "was", "raining"}})

	diff := cmp.Diff(
		[][]string{{"it", "was", "rain"}},
		actual,
	)
	if diff != "" {
		t.Errorf(diff)
	}
}
