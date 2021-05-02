package main

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
)

func init() {
	// ignore message
	gin.SetMode("release")
}

func TestAnalyze(t *testing.T) {
	tokenizer, _ := newTokenizer()

	actual := tokenizer.analyze([]string{"すもももももももものうち"})

	if diff := cmp.Diff(
		[][]string{{
			"すもも",
			"も",
			"もも",
			"も",
			"もも",
			"の",
			"うち",
		}},
		actual,
	); diff != "" {
		t.Errorf(diff)
	}
}

func TestAnalyzeWhitespace(t *testing.T) {
	tokenizer, _ := newTokenizer()

	actual := tokenizer.analyze([]string{" "})

	if diff := cmp.Diff(
		[][]string{{}},
		actual,
	); diff != "" {
		t.Errorf(diff)
	}
}

func TestAnalyzeSymbol(t *testing.T) {
	tokenizer, _ := newTokenizer()

	actual := tokenizer.analyze([]string{"！？"})

	if diff := cmp.Diff(
		[][]string{{
			"！",
			"？",
		}},
		actual,
	); diff != "" {
		t.Errorf(diff)
	}
}
