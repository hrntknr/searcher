package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMappingCharFilter(t *testing.T) {
	filter, _ := newMappingCharFilter(map[string]string{":)": "happy", ":(": "sad"})
	actual := filter.filter([]string{":), :("})

	diff := cmp.Diff(
		[]string{"happy, sad"},
		actual,
	)
	if diff != "" {
		t.Errorf(diff)
	}
}
