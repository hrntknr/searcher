package main

import (
	"strings"

	"github.com/kljensen/snowball/english"
)

type WordFilter interface {
	Filter([][]string) [][]string
}

func newLowercaseFilter() (*lowercaseFilter, error) {
	return &lowercaseFilter{}, nil
}

type lowercaseFilter struct {
}

func (f *lowercaseFilter) Filter(tokens [][]string) [][]string {
	newTokens := make([][]string, len(tokens))
	for i, token := range tokens {
		newTokens[i] = make([]string, len(tokens[i]))
		for j, token := range token {
			newTokens[i][j] = strings.ToLower(token)
		}
	}
	return newTokens
}

func newStopWordFilter(stopWords []string) (*stopWordFilter, error) {
	filter := &stopWordFilter{
		stopWords: map[string]struct{}{},
	}
	for _, word := range stopWords {
		filter.stopWords[word] = struct{}{}
	}
	return filter, nil
}

type stopWordFilter struct {
	stopWords map[string]struct{}
}

func (f *stopWordFilter) Filter(tokens [][]string) [][]string {
	newTokens := make([][]string, len(tokens))
	for i, token := range tokens {
		newTokens[i] = []string{}
		for _, token := range token {
			if _, ok := f.stopWords[token]; ok {
				continue
			}
			newTokens[i] = append(newTokens[i], token)
		}
	}
	return newTokens
}

func newStemmerFilter() (*stemmerFilter, error) {
	return &stemmerFilter{}, nil
}

type stemmerFilter struct {
}

func (f *stemmerFilter) Filter(tokens [][]string) [][]string {
	newTokens := make([][]string, len(tokens))
	for i, token := range tokens {
		newTokens[i] = make([]string, len(tokens[i]))
		for j, token := range token {
			stemmed := english.Stem(token, false)
			newTokens[i][j] = stemmed
		}
	}
	return newTokens
}
