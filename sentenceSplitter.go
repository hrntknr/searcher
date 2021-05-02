package main

import (
	"bufio"
	"strings"

	"github.com/ikawaha/kagome/v2/filter"
)

type sentenceSplitter interface {
	split(body string) ([]string, error)
}

func newSentenceSplitter() (*sentenceSplitterImpl, error) {
	return &sentenceSplitterImpl{}, nil
}

type sentenceSplitterImpl struct {
}

func (sp *sentenceSplitterImpl) split(body string) ([]string, error) {
	scanner := bufio.NewScanner(strings.NewReader(body))
	scanner.Split(filter.ScanSentences)
	result := []string{}
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
