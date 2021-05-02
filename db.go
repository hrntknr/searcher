package main

import (
	"math"
	"sort"

	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type db interface {
	store(token string, positions []position, Uri string) error
	search(tokens []string, offset, count uint) ([]searchResult, error)
}

func newDb(db *gorm.DB) (*dbImpl, error) {
	return &dbImpl{
		db: db,
	}, nil
}

type dbImpl struct {
	db *gorm.DB
}

func (db *dbImpl) store(token string, positions []position, uri string) error {
	if err := db.db.Transaction(func(tx *gorm.DB) error {
		var index invertedIndex
		err := tx.First(&index, "token = ?", token).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if err == gorm.ErrRecordNotFound {
			index = invertedIndex{
				Token:    token,
				Postings: []posting{},
			}
			if err := tx.Create(&index).Error; err != nil {
				return err
			}
		}
		var currentPosting posting
		err = tx.First(&currentPosting, "uri = ? AND inverted_index_id = ?", uri, index.ID).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		var postingID uint
		if err == nil {
			postingID = currentPosting.ID
			if err := tx.Where("posting_id = ?", currentPosting.ID).Delete(&sentence{}).Error; err != nil {
				return err
			}
			if err := tx.Where("posting_id = ?", currentPosting.ID).Delete(&tokenPosition{}).Error; err != nil {
				return err
			}
		} else {
			posting := posting{
				InvertedIndexID: index.ID,
				Uri:             uri,
				TokenCount:      uint(len(positions)),
				Sentences:       []sentence{},
				TokenPositions:  []tokenPosition{},
			}
			if err := tx.Create(&posting).Error; err != nil {
				return err
			}
			postingID = posting.ID
		}
		for _, pos := range positions {
			sentence := sentence{
				PostingID:      postingID,
				Body:           pos.Sentence,
				TokenPositions: []tokenPosition{},
			}
			if err := tx.Create(&sentence).Error; err != nil {
				return err
			}
			position := tokenPosition{
				SentenceID:       sentence.ID,
				SentencePosition: pos.SentencePosition,
				PostingID:        postingID,
				PostingPosition:  pos.PostingPosition,
			}
			if err := tx.Create(&position).Error; err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (db *dbImpl) search(tokens []string, offset, count uint) ([]searchResult, error) {
	var allPostingCount int64
	if err := db.db.Model(&posting{}).Count(&allPostingCount).Error; err != nil {
		return nil, err
	}
	invertedIndexes := make([]invertedIndex, len(tokens))
	var eg errgroup.Group
	for i, token := range tokens {
		i, token := i, token
		eg.Go(func() error {
			// GORMのjoinが使いにくすぎてゴリ押し
			if err := db.db.First(&invertedIndexes[i], "token = ?", token).Error; err != nil {
				return err
			}
			if err := db.db.Find(&invertedIndexes[i].Postings, "inverted_index_id = ?", invertedIndexes[i].ID).Error; err != nil {
				return err
			}
			for j, posting := range invertedIndexes[i].Postings {
				if err := db.db.Find(&invertedIndexes[i].Postings[j].Sentences, "posting_id = ?", posting.ID).Error; err != nil {
					return err
				}
				for k, sentence := range invertedIndexes[i].Postings[j].Sentences {
					if err := db.db.Find(&invertedIndexes[i].Postings[j].Sentences[k].TokenPositions, "sentence_id = ?", sentence.ID).Error; err != nil {
						return err
					}
				}
			}
			for j, posting := range invertedIndexes[i].Postings {
				if err := db.db.Find(&invertedIndexes[i].Postings[j].TokenPositions, "posting_id = ?", posting.ID).Error; err != nil {
					return err
				}
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	scores := map[string]float64{}
	postingCache := map[string]cache{}
	for _, tokenIndex := range invertedIndexes {
		idf := math.Log(float64(allPostingCount) / float64(len(tokenIndex.Postings)+1))
		for _, posting := range tokenIndex.Postings {
			tf := float64(len(posting.TokenPositions)) / float64(posting.TokenCount)
			if _, ok := scores[posting.Uri]; !ok {
				scores[posting.Uri] = 1
			}
			if _, ok := postingCache[posting.Uri]; !ok {
				postingCache[posting.Uri] = cache{
					uri:       posting.Uri,
					sentences: map[uint]string{},
				}
			}
			cache := postingCache[posting.Uri]
			for _, sentence := range posting.Sentences {
				for _, pos := range sentence.TokenPositions {
					cache.sentences[pos.PostingPosition-pos.SentencePosition] = sentence.Body
				}
			}
			postingCache[posting.Uri] = cache
			scores[posting.Uri] *= (tf * idf)
		}
	}

	resultOrder := []string{}
	for key := range scores {
		resultOrder = append(resultOrder, key)
	}
	sort.Slice(resultOrder, func(i, j int) bool {
		return scores[resultOrder[i]] > scores[resultOrder[j]]
	})

	result := []searchResult{}
	cursor := int(offset)
	for len(result) < int(count) {
		if len(resultOrder) <= cursor {
			return result, nil
		}

		sentenceIndexes := []uint{}
		for key := range postingCache[resultOrder[cursor]].sentences {
			sentenceIndexes = append(sentenceIndexes, key)
		}
		sort.Slice(sentenceIndexes, func(i, j int) bool {
			return sentenceIndexes[i] < sentenceIndexes[j]
		})
		sentences := make([]string, len(sentenceIndexes))
		for i, key := range sentenceIndexes {
			sentences[i] = postingCache[resultOrder[cursor]].sentences[key]
		}
		result = append(result, searchResult{
			Uri:       postingCache[resultOrder[cursor]].uri,
			Score:     scores[resultOrder[cursor]],
			Sentences: sentences,
		})
		cursor++
	}

	return result, nil
}

type cache struct {
	uri       string
	sentences map[uint]string
}

type searchResult struct {
	Uri       string
	Score     float64
	Sentences []string
}

type invertedIndex struct {
	gorm.Model
	Token    string `gorm:"type:varchar(191);uniqueIndex"`
	Postings []posting
}

type posting struct {
	gorm.Model
	InvertedIndexID uint
	Uri             string `gorm:"type:varchar(191);index"`
	TokenCount      uint
	Sentences       []sentence
	TokenPositions  []tokenPosition
}

type sentence struct {
	gorm.Model
	PostingID      uint
	Body           string
	TokenPositions []tokenPosition
}

type tokenPosition struct {
	gorm.Model
	SentenceID       uint
	SentencePosition uint
	PostingID        uint
	PostingPosition  uint
}

type position struct {
	SentencePosition uint
	PostingPosition  uint
	Sentence         string
}
