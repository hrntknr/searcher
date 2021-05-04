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

	actual := tokenizer.Analyze([]string{"すもももももももものうち"})

	if diff := cmp.Diff(
		[][]string{{
			"スモモ",
			"モ",
			"モモ",
			"モ",
			"モモ",
			"ノ",
			"ウチ",
		}},
		actual,
	); diff != "" {
		t.Errorf(diff)
	}
}

func TestAnalyzeWhitespace(t *testing.T) {
	tokenizer, _ := newTokenizer()

	actual := tokenizer.Analyze([]string{" "})

	if diff := cmp.Diff(
		[][]string{{}},
		actual,
	); diff != "" {
		t.Errorf(diff)
	}
}

func TestAnalyzeSymbol(t *testing.T) {
	tokenizer, _ := newTokenizer()

	actual := tokenizer.Analyze([]string{"！？"})

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
