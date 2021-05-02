package main

import (
	"github.com/ikawaha/kagome-dict/ipa"
	kagome "github.com/ikawaha/kagome/v2/tokenizer"
)

type tokenizer interface {
	analyze(text []string) [][]string
}

func newTokenizer() (*tokenizerImpl, error) {
	kagomeTokenizer, err := kagome.New(ipa.Dict(), kagome.OmitBosEos())
	if err != nil {
		return nil, err
	}

	return &tokenizerImpl{
		kagome: kagomeTokenizer,
	}, nil
}

type tokenizerImpl struct {
	kagome *kagome.Tokenizer
}

func (t *tokenizerImpl) analyze(text []string) [][]string {
	result := make([][]string, len(text))
	for i, text := range text {
		tokens := t.kagome.Analyze(text, kagome.Search)
		res := []string{}
		for _, t := range tokens {
			features := t.Features()
			if features[1] == "空白" {
				continue
			}
			kana := t.Surface
			if len(features) >= 8 {
				kana = features[7]
			}
			res = append(res, kana)
		}
		result[i] = res
	}
	return result
}
