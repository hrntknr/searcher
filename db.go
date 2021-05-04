package main

import (
	"github.com/hrntknr/searcher/types"
	"gorm.io/gorm"
)

type DB interface {
	// 全ドキュメント数
	CountDocument() (uint, error)
	// ドキュメントの中の単語数
	CountTermInDocument(documentID uint) (uint, error)

	// URIからドキュメントに
	DocumentFromUri(uri string) (*types.Document, error)
	// IDからドキュメントを取得
	DocumentFromID(id uint) (*types.Document, error)
	// ドキュメントを作成
	CreateDcoument(document *types.Document) (*types.Document, error)

	// トークン文字列からトークンに
	TokenFromString(token string) (*types.Token, error)
	// IDからトークンを取得
	TokenFromID(id uint) (*types.Token, error)
	// トークンを作成
	CreateToken(token *types.Token) (*types.Token, error)

	// ポスティングリストを取得、センテンスのアソシエーションを結合
	PostingList(tokenID uint) ([]*types.Posting, error)
	// ポスティングを作成
	CreatePosting(posting *types.Posting) (*types.Posting, error)

	// 複数IDからセンテンスを同時取得、ソートはID順
	SentenceMultiFromID(ids []uint) ([]*types.Sentence, error)
	// センテンスを作成
	CreateSentence(sentence *types.Sentence) (*types.Sentence, error)
	// 指定したドキュメントのセンテンスを一括削除（更新用）、ついでにポスティング、アソシエーションも消す
	DeleteSentenceFromDocumentID(documentID uint) error
}

func newDb(db *gorm.DB) (*dbImpl, error) {
	return &dbImpl{
		db: db,
	}, nil
}

type dbImpl struct {
	db *gorm.DB
}

func (db *dbImpl) CountDocument() (uint, error) {
	var count int64
	if err := db.db.Model(&types.Document{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return uint(count), nil
}

func (db *dbImpl) CountTermInDocument(documentID uint) (uint, error) {
	var document types.Document
	if err := db.db.Model(&types.Document{}).Where("id = ?", documentID).First(&document).Error; err != nil {
		return 0, err
	}
	return uint(document.TokenCount), nil
}

func (db *dbImpl) DocumentFromUri(uri string) (*types.Document, error) {
	var document types.Document
	err := db.db.Model(&types.Document{}).Where("uri = ?", uri).First(&document).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &document, nil
}

func (db *dbImpl) DocumentFromID(id uint) (*types.Document, error) {
	var document types.Document
	err := db.db.Model(&types.Document{}).Where("id = ?", id).First(&document).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &document, nil
}

func (db *dbImpl) CreateDcoument(document *types.Document) (*types.Document, error) {
	if err := db.db.Model(&types.Document{}).Create(document).Error; err != nil {
		return nil, err
	}
	return document, nil
}

func (db *dbImpl) TokenFromString(token string) (*types.Token, error) {
	var tkn types.Token
	err := db.db.Model(&types.Token{}).Where("token = ?", token).First(&tkn).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &tkn, nil
}

func (db *dbImpl) TokenFromID(id uint) (*types.Token, error) {
	var tkn types.Token
	err := db.db.Model(&types.Token{}).Where("id = ?", id).First(&tkn).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &tkn, nil
}

func (db *dbImpl) CreateToken(token *types.Token) (*types.Token, error) {
	if err := db.db.Model(&types.Token{}).Create(token).Error; err != nil {
		return nil, err
	}
	return token, nil
}

func (db *dbImpl) PostingList(tokenID uint) ([]*types.Posting, error) {
	lsit := []*types.Posting{}
	if err := db.db.Model(&types.Posting{}).Where("token_id = ?", tokenID).Preload("Sentences").Find(&lsit).Error; err != nil {
		return nil, err
	}
	return lsit, nil
}

func (db *dbImpl) CreatePosting(posting *types.Posting) (*types.Posting, error) {
	if err := db.db.Model(&types.Posting{}).Create(posting).Error; err != nil {
		return nil, err
	}
	return posting, nil
}

func (db *dbImpl) SentenceMultiFromID(ids []uint) ([]*types.Sentence, error) {
	sentences := []*types.Sentence{}
	if err := db.db.Model(&types.Sentence{}).Where(ids).Find(&sentences).Error; err != nil {
		return nil, err
	}
	return sentences, nil
}

func (db *dbImpl) CreateSentence(sentence *types.Sentence) (*types.Sentence, error) {
	if err := db.db.Model(&types.Sentence{}).Create(sentence).Error; err != nil {
		return nil, err
	}
	return sentence, nil
}

func (db *dbImpl) DeleteSentenceFromDocumentID(documentID uint) error {
	if err := db.db.Transaction(func(tx *gorm.DB) error {
		sentences := []*types.Sentence{}
		if err := tx.Model(&types.Sentence{}).Where("document_id = ?", documentID).Preload("Postings").Find(&sentences).Error; err != nil {
			return err
		}
		if len(sentences) > 0 {
			if err := tx.Model(&types.Sentence{}).Select("Postings").Delete(&sentences).Error; err != nil {
				return err
			}
			postings := []int{}
			for _, sentence := range sentences {
				for _, posting := range sentence.Postings {
					postings = append(postings, int(posting.ID))
				}
			}
			if len(postings) > 0 {
				if err := tx.Model(&types.Posting{}).Delete(&types.Posting{}, postings).Error; err != nil {
					return err
				}
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
