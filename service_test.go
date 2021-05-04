package main

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/hrntknr/searcher/mock"
	"github.com/hrntknr/searcher/types"
	"gorm.io/gorm"
)

func TestServiceRegist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sentenceSplitter := mock.NewMockSentenceSplitter(ctrl)
	tokenizer := mock.NewMockTokenizer(ctrl)
	charFilter := mock.NewMockCharFilter(ctrl)
	wordFilter := mock.NewMockWordFilter(ctrl)
	db := mock.NewMockDB(ctrl)
	gomock.InOrder(
		sentenceSplitter.EXPECT().Split("これはペンです。これはりんごです。:)。").Return([]string{"これはペンです。", "これはりんごです。", ":)。"}, nil),
		charFilter.EXPECT().Filter([]string{"これはペンです。", "これはりんごです。", ":)。"}).Return([]string{"これはペンです。", "これはりんごです。", "happy。"}),
		tokenizer.EXPECT().Analyze([]string{"これはペンです。", "これはりんごです。", "happy。"}).Return([][]string{{"コレ", "ハ", "ペン", "デス", "。"}, {"コレ", "ハ", "リンゴ", "デス", "。"}, {"happy", "。"}}),
		wordFilter.EXPECT().Filter([][]string{{"コレ", "ハ", "ペン", "デス", "。"}, {"コレ", "ハ", "リンゴ", "デス", "。"}, {"happy", "。"}}).Return([][]string{{"コレ", "ペン", "デス"}, {"コレ", "リンゴ", "デス"}, {"happy"}}),
		db.EXPECT().DocumentFromUri("uri").Return(nil, nil),
		db.EXPECT().CreateDcoument(gomock.Any()).Return(&types.Document{
			Model: gorm.Model{
				ID: 1,
			},
			Uri:        "uri",
			TokenCount: 7,
			Time:       time.Now(),
		}, nil),
		db.EXPECT().DeleteSentenceFromDocumentID(uint(1)).Return(nil),
	)
	thisispen := &types.Sentence{
		Model: gorm.Model{
			ID: 1,
		},
		DocumentID: 1,
		Index:      0,
		Sentence:   "これはペンです。",
		TokenCount: 3,
	}
	thisisapple := &types.Sentence{
		Model: gorm.Model{
			ID: 2,
		},
		DocumentID: 1,
		Index:      1,
		Sentence:   "これはりんごです。",
		TokenCount: 3,
	}
	happy := &types.Sentence{
		Model: gorm.Model{
			ID: 3,
		},
		DocumentID: 1,
		Index:      2,
		Sentence:   "happy。",
		TokenCount: 3,
	}
	db.EXPECT().CreateSentence(&types.Sentence{
		DocumentID: 1,
		Index:      0,
		Sentence:   "これはペンです。",
		TokenCount: 3,
	}).Return(thisispen, nil)
	db.EXPECT().CreateSentence(&types.Sentence{
		DocumentID: 1,
		Index:      1,
		Sentence:   "これはりんごです。",
		TokenCount: 3,
	}).Return(thisisapple, nil)
	db.EXPECT().CreateSentence(&types.Sentence{
		DocumentID: 1,
		Index:      2,
		Sentence:   "happy。",
		TokenCount: 1,
	}).Return(happy, nil)

	db.EXPECT().TokenFromString("コレ").Return(&types.Token{
		Model: gorm.Model{
			ID: 1,
		},
	}, nil)
	db.EXPECT().TokenFromString("ペン").Return(&types.Token{
		Model: gorm.Model{
			ID: 2,
		},
	}, nil)
	db.EXPECT().TokenFromString("リンゴ").Return(&types.Token{
		Model: gorm.Model{
			ID: 3,
		},
	}, nil)
	db.EXPECT().TokenFromString("デス").Return(&types.Token{
		Model: gorm.Model{
			ID: 4,
		},
	}, nil)
	db.EXPECT().TokenFromString("happy").Return(&types.Token{
		Model: gorm.Model{
			ID: 5,
		},
	}, nil)

	db.EXPECT().CreatePosting(&types.Posting{
		TokenID:    1,
		DocumentID: 1,
		Sentences:  []*types.Sentence{thisispen, thisisapple},
	})
	db.EXPECT().CreatePosting(&types.Posting{
		TokenID:    2,
		DocumentID: 1,
		Sentences:  []*types.Sentence{thisispen},
	})
	db.EXPECT().CreatePosting(&types.Posting{
		TokenID:    3,
		DocumentID: 1,
		Sentences:  []*types.Sentence{thisisapple},
	})
	db.EXPECT().CreatePosting(&types.Posting{
		TokenID:    4,
		DocumentID: 1,
		Sentences:  []*types.Sentence{thisispen, thisisapple},
	})
	db.EXPECT().CreatePosting(&types.Posting{
		TokenID:    5,
		DocumentID: 1,
		Sentences:  []*types.Sentence{happy},
	})

	service, _ := newService(
		sentenceSplitter,
		tokenizer,
		[]CharFilter{charFilter},
		[]WordFilter{wordFilter},
		db,
	)

	err := service.Regist("uri", "これはペンです。これはりんごです。:)。")
	if err != nil {
		t.Error(err)
	}
}

func TestServiceSearch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sentenceSplitter := mock.NewMockSentenceSplitter(ctrl)
	tokenizer := mock.NewMockTokenizer(ctrl)
	charFilter := mock.NewMockCharFilter(ctrl)
	wordFilter := mock.NewMockWordFilter(ctrl)
	db := mock.NewMockDB(ctrl)

	gomock.InOrder(
		charFilter.EXPECT().Filter([]string{"これ ペン ペンギン"}).Return([]string{"これ ペン ペンギン"}),
		tokenizer.EXPECT().Analyze([]string{"これ ペン ペンギン"}).Return([][]string{{"コレ", "ペン", "ペンギン"}}),
		wordFilter.EXPECT().Filter([][]string{{"コレ", "ペン", "ペンギン"}}).Return([][]string{{"コレ", "ペン", "ペンギン"}}),
		db.EXPECT().CountDocument().Return(uint(100), nil),
	)
	db.EXPECT().TokenFromString("コレ").Return(&types.Token{
		Model: gorm.Model{
			ID: 3,
		},
	}, nil)
	db.EXPECT().TokenFromString("ペン").Return(&types.Token{
		Model: gorm.Model{
			ID: 4,
		},
	}, nil)
	db.EXPECT().TokenFromString("ペンギン").Return(nil, gorm.ErrRecordNotFound)

	db.EXPECT().PostingList(uint(3)).Return([]*types.Posting{{
		TokenID:    3,
		DocumentID: 5,
		Sentences: []*types.Sentence{{
			Model: gorm.Model{
				ID: 2,
			},
		}, {
			Model: gorm.Model{
				ID: 2,
			},
		}},
	}}, nil)
	db.EXPECT().PostingList(uint(4)).Return([]*types.Posting{{
		TokenID:    4,
		DocumentID: 5,
		Sentences: []*types.Sentence{{
			Model: gorm.Model{
				ID: 3,
			},
		}},
	}}, nil)
	db.EXPECT().CountTermInDocument(uint(5)).Return(uint(6), nil)
	db.EXPECT().SentenceMultiFromID(gomock.Len(2)).Return([]*types.Sentence{
		{
			Model: gorm.Model{
				ID: 2,
			},
			DocumentID: 5,
			Index:      0,
			Sentence:   "これだよ、これ。",
			TokenCount: 3,
		},
		{
			Model: gorm.Model{
				ID: 3,
			},
			DocumentID: 5,
			Index:      1,
			Sentence:   "ペンってすごい。",
			TokenCount: 3,
		},
	}, nil)
	db.EXPECT().DocumentFromID(uint(5)).Return(&types.Document{
		Model: gorm.Model{
			ID: 5,
		},
		Uri: "test",
	}, nil)

	service, _ := newService(
		sentenceSplitter,
		tokenizer,
		[]CharFilter{charFilter},
		[]WordFilter{wordFilter},
		db,
	)

	result, err := service.Search("これ ペン ペンギン", 0, 10)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(
		[]types.SearchResult{{
			Uri:       "test",
			Score:     0.6520038342380243,
			Sentences: []string{"これだよ、これ。", "ペンってすごい。"},
		}},
		result,
	); diff != "" {
		t.Errorf(diff)
	}
}
