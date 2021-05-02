package main

import (
	"fmt"
)

type service interface {
	regist(uri string, body string) error
	search(str string, offset, count uint) ([]searchResult, error)
}

func newService(
	sentenceSplitter sentenceSplitter,
	tokenizer tokenizer,
	charFilter []charFilter,
	wordFilter []wordFilter,
	db db,
) (service, error) {
	return &serviceImpl{
		sentenceSplitter: sentenceSplitter,
		tokenizer:        tokenizer,
		charFilter:       charFilter,
		wordFilter:       wordFilter,
		db:               db,
	}, nil
}

type serviceImpl struct {
	sentenceSplitter sentenceSplitter
	tokenizer        tokenizer
	charFilter       []charFilter
	wordFilter       []wordFilter
	db               db
}

func (s *serviceImpl) regist(uri string, body string) error {
	sentence, err := s.sentenceSplitter.split(body)
	if err != nil {
		return err
	}
	for _, f := range s.charFilter {
		sentence = f.filter(sentence)
	}
	sentenceTokens := s.tokenizer.analyze(sentence)
	for _, f := range s.wordFilter {
		sentenceTokens = f.filter(sentenceTokens)
	}
	positions := map[string][]position{}
	pos := 0
	for i, tokens := range sentenceTokens {
		for j, token := range tokens {
			if _, ok := positions[token]; !ok {
				positions[token] = []position{}
			}
			positions[token] = append(positions[token], position{
				SentencePosition: uint(j),
				PostingPosition:  uint(pos),
				Sentence:         sentence[i],
			})
			pos++
		}
	}
	errors := []error{}
	for token, ps := range positions {
		if err := s.db.store(token, ps, uri); err != nil {
			errors = append(errors, err)
		}
	}
	if len(errors) != 0 {
		msg := ""
		for i, err := range errors {
			msg += fmt.Sprintf("\n%d:\n%s", i+1, err.Error())
		}
		return fmt.Errorf("%d errors occurred during save index\n%s", len(errors), msg)
	}
	return nil
}

func (s *serviceImpl) search(body string, offset, count uint) ([]searchResult, error) {
	for _, f := range s.charFilter {
		body = f.filter([]string{body})[0]
	}
	tokens := s.tokenizer.analyze([]string{body})[0]
	result, err := s.db.search(tokens, offset, count)
	if err != nil {
		return nil, err
	}
	return result, nil
}
