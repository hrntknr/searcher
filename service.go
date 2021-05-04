package main

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/hrntknr/searcher/types"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type Service interface {
	Regist(uri string, body string) error
	Search(str string, offset, count uint) ([]types.SearchResult, error)
}

func newService(
	sentenceSplitter SentenceSplitter,
	tokenizer Tokenizer,
	charFilter []CharFilter,
	wordFilter []WordFilter,
	db DB,
) (Service, error) {
	return &serviceImpl{
		sentenceSplitter: sentenceSplitter,
		tokenizer:        tokenizer,
		charFilter:       charFilter,
		wordFilter:       wordFilter,
		db:               db,
	}, nil
}

type serviceImpl struct {
	sentenceSplitter SentenceSplitter
	tokenizer        Tokenizer
	charFilter       []CharFilter
	wordFilter       []WordFilter
	db               DB
}

func (s *serviceImpl) Regist(uri string, body string) error {
	// ドキュメントを文章ごとの配列に分割
	sentences, err := s.sentenceSplitter.Split(body)
	if err != nil {
		return err
	}
	// 前処理
	for _, f := range s.charFilter {
		sentences = f.Filter(sentences)
	}
	// トークン化
	sentencesTokens := s.tokenizer.Analyze(sentences)
	// 後処理
	for _, f := range s.wordFilter {
		sentencesTokens = f.Filter(sentencesTokens)
	}

	// tokenCount
	tokenCount := 0
	for _, sentenceToken := range sentencesTokens {
		tokenCount += len(sentenceToken)
	}

	// ドキュメントIDを作成、取得
	document, err := s.db.DocumentFromUri(uri)
	if err != nil {
		return err
	}
	if document == nil {
		_document, err := s.db.CreateDcoument(&types.Document{
			Uri:        uri,
			TokenCount: uint(tokenCount),
			Time:       time.Now(),
		})
		if err != nil {
			return err
		}
		document = _document
	}
	// アップデート用に既存の文章を削除
	if err := s.db.DeleteSentenceFromDocumentID(document.ID); err != nil {
		return err
	}

	// 文章を追加
	dbSentences := make([]*types.Sentence, len(sentences))
	egAddSentence := errgroup.Group{}
	for i, sentence := range sentences {
		i, sentence := i, sentence
		egAddSentence.Go(func() error {
			sentence, err := s.db.CreateSentence(&types.Sentence{
				DocumentID: document.ID,
				Index:      uint(i),
				Sentence:   sentence,
				TokenCount: uint(len(sentencesTokens[i])),
			})
			if err != nil {
				return err
			}
			dbSentences[i] = sentence
			return nil
		})
	}
	if err := egAddSentence.Wait(); err != nil {
		return err
	}

	// トークンをユニークキーにポスティングリストを作成
	positionList := map[string][]positionCache{}
	pos := 0
	for i, tokens := range sentencesTokens {
		for j, token := range tokens {
			if _, ok := positionList[token]; !ok {
				positionList[token] = []positionCache{}
			}
			positionList[token] = append(positionList[token], positionCache{
				SentencePosition: uint(j),
				PostingPosition:  uint(pos),
				Sentence:         dbSentences[i],
			})
			pos++
		}
	}

	// ポスティングリストを追加
	egAddPostingList := errgroup.Group{}
	for tokenStr, positions := range positionList {
		tokenStr, positions := tokenStr, positions
		egAddPostingList.Go(func() error {
			token, err := s.db.TokenFromString(tokenStr)
			if err != nil {
				return err
			}
			if token == nil {
				_token, err := s.db.CreateToken(&types.Token{
					Token: tokenStr,
				})
				if err != nil {
					return err
				}
				token = _token
			}
			sentences := make([]*types.Sentence, len(positions))
			for i, position := range positions {
				sentences[i] = position.Sentence
			}

			if _, err := s.db.CreatePosting(&types.Posting{
				DocumentID: document.ID,
				TokenID:    token.ID,
				Sentences:  sentences,
			}); err != nil {
				return err
			}
			return nil
		})
	}
	if err := egAddPostingList.Wait(); err != nil {
		return err
	}

	return nil
}

func (s *serviceImpl) Search(body string, offset, count uint) ([]types.SearchResult, error) {
	b := []string{body}
	// 前処理
	for _, f := range s.charFilter {
		b = f.Filter(b)
	}
	// トークン化
	bt := s.tokenizer.Analyze(b)
	// 後処理
	for _, f := range s.wordFilter {
		bt = f.Filter(bt)
	}
	tokens := bt[0]
	if len(tokens) == 0 {
		return nil, fmt.Errorf("invalid input")
	}

	// ドキュメントの件数取得
	allCount, err := s.db.CountDocument()
	if err != nil {
		return nil, err
	}

	// トークンを検索。ない場合はスキップ
	dbTokens := []*types.Token{}
	dbTokensLock := sync.Mutex{}
	egListToken := errgroup.Group{}
	for _, token := range tokens {
		token := token
		egListToken.Go(func() error {
			t, err := s.db.TokenFromString(token)
			if err != nil && err != gorm.ErrRecordNotFound {
				return err
			}
			if t != nil {
				dbTokensLock.Lock()
				dbTokens = append(dbTokens, t)
				dbTokensLock.Unlock()
			}
			return nil
		})
	}
	if err := egListToken.Wait(); err != nil {
		return nil, err
	}
	// tokenがなければ絶望的、１つでも存在すればそれで検索する。（サジェスト的な）
	if len(dbTokens) == 0 {
		return []types.SearchResult{}, nil
	}

	// トークンごとにポスティングテーブルを取得
	postingLists := map[uint][]*types.Posting{}
	postingListsLock := sync.Mutex{}
	egGetPostingLists := errgroup.Group{}
	for _, token := range dbTokens {
		token := token
		egGetPostingLists.Go(func() error {

			// トークンごとにポスティングテーブルを取得
			postingList, err := s.db.PostingList(token.ID)
			if err != nil {
				return err
			}
			postingListsLock.Lock()
			postingLists[token.ID] = postingList
			postingListsLock.Unlock()
			return nil
		})
	}
	if err := egGetPostingLists.Wait(); err != nil {
		return nil, err
	}

	// ポスティングリストのANDを取る
	documentList := []uint{}
	postingCountByPage := map[uint]uint{}
	for _, postingList := range postingLists {
		for _, posting := range postingList {
			if _, ok := postingCountByPage[posting.DocumentID]; !ok {
				postingCountByPage[posting.DocumentID] = 0
			}
			postingCountByPage[posting.DocumentID] += 1
			if postingCountByPage[posting.DocumentID] == uint(len(dbTokens)) {
				documentList = append(documentList, posting.DocumentID)
			}
		}
	}

	// 各ページのスコアを計算
	scores := map[uint]float64{}
	for _, documentID := range documentList {
		termCount, err := s.db.CountTermInDocument(documentID)
		if err != nil {
			return nil, err
		}
		for _, token := range dbTokens {
			idf := math.Log(float64(allCount) / float64(len(postingLists[token.ID])+1))
			tf := float64(len(postingLists[token.ID])) / float64(termCount)
			scores[documentID] = tf * idf
		}
	}

	// スコアによって並べ替え
	sort.Slice(documentList, func(i, j int) bool {
		return scores[documentList[i]] > scores[documentList[j]]
	})

	// 検索対象範囲を絞る
	result := []types.SearchResult{}
	cursor := int(offset)
	for len(result) < int(count) {
		if len(documentList) <= cursor {
			return result, nil
		}
		// 検索結果を追加、
		documentID := documentList[cursor]
		// 重複削除
		sentenceMap := map[uint]struct{}{}
		for _, token := range dbTokens {
			for _, list := range postingLists[token.ID] {
				for _, sentence := range list.Sentences {
					sentenceMap[sentence.ID] = struct{}{}
				}
			}
		}
		// DBから文章をひっぱってくる
		sentenceIDs := []uint{}
		for sentenceID := range sentenceMap {
			sentenceIDs = append(sentenceIDs, sentenceID)
		}
		sentences, err := s.db.SentenceMultiFromID(sentenceIDs)
		if err != nil {
			return nil, err
		}
		sentenceStrs := make([]string, len(sentences))
		for i, sentence := range sentences {
			sentenceStrs[i] = sentence.Sentence
		}

		document, err := s.db.DocumentFromID(documentID)
		if err != nil {
			return nil, err
		}

		result = append(result, types.SearchResult{
			Uri:       document.Uri,
			Score:     scores[documentID],
			Sentences: sentenceStrs,
		})
		cursor++
	}

	return result, nil
}

type positionCache struct {
	SentencePosition uint
	PostingPosition  uint
	Sentence         *types.Sentence
}
