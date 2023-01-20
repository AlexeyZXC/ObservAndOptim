package store

import (
	"context"
	"elastic/e"
	"elastic/m"

	"github.com/mitchellh/mapstructure"
)

type iLogger interface {
	Log(format string, a ...any)
	Error(format string, a ...any)
}

type ArticleStore struct {
	E   e.E
	log iLogger
}

func NewArticleStore(log iLogger) (ArticleStore, error) {
	e, err := e.NewE("articles", log)
	if err != nil {
		return ArticleStore{}, err
	}
	return ArticleStore{E: e}, nil
}

func (s ArticleStore) Add(ctx context.Context, article m.Article) error {
	return s.E.Insert(ctx, article)
}

func (s ArticleStore) Search(ctx context.Context, query string) ([]m.Article, error) {
	result, err := s.E.Search(ctx, query)
	if err != nil {
		return nil, err
	}
	hits := result.Hits.Hits
	articles := []m.Article{}
	for _, hit := range hits {
		var article m.Article
		//map[string]interface{} -> struct
		err = mapstructure.Decode(hit.Source, &article)
		if err != nil {
			return nil, err
		}
		article.Id = hit.ID
		articles = append(articles, article)
	}
	return articles, nil
}

func (s ArticleStore) Get(ctx context.Context, id string) (m.Article, error) {
	result, err := s.E.Get(ctx, id)
	if err != nil {
		return m.Article{}, err
	}
	//l.L(result)
	s.log.Log("Article store: Get: result: %v", result)
	var article m.Article
	// err = mapstructure.Decode(result.Source, &article)
	// if err != nil {
	// 	return m.Article{}, err
	// }
	return article, nil
}
