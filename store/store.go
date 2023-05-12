package store

import (
	"context"
	// "elastic/l"
	// "geekbrains/models"
	"geekbrains/e"
	"geekbrains/internal/models"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// type dbItem interface {
// 	models.Article | models.User
// 	// models.HasIDInterface
// }

type pdbItem interface {
	*models.Article | *models.User
	models.HasIDInterface
}

type Store struct {
	EArticle e.E
	EUser    e.E
}

func NewStore() (*Store, error) {
	// func NewStore() Store {
	earticle, err := e.NewE("articles", "Title")
	if err != nil {
		return nil, err
		// return Store{}
	}

	eusers, err := e.NewE("users", "Name")
	if err != nil {
		return nil, err
		// return Store{}
	}

	return &Store{EArticle: earticle, EUser: eusers}, nil
	// return Store{EArticle: earticle, EUser: eusers}
}

func (s *Store) AddUserArticle(ctx context.Context, article models.Article) (uuid.UUID, error) {
	return s.EArticle.Insert(ctx, article)
}

func (s *Store) AddUser(ctx context.Context, user models.User) (uuid.UUID, error) {
	return s.EUser.Insert(ctx, user)
}

// func (s *Store) SearchArticles(ctx context.Context, query string) ([]models.Article, error) {
func (s *Store) SearchArticles(ctx context.Context, userID uuid.UUID) ([]models.Article, error) {
	result, err := s.EArticle.Search(ctx, "user_id", userID)
	if err != nil {
		return nil, err
	}

	array := []models.Article{}
	pArray, err := search[*models.Article](result)
	for _, v := range pArray {
		array = append(array, *v)
	}
	return array, err

	// hits := result.Hits.Hits
	// articles := []m.Article{}
	// for _, hit := range hits {
	// 	var article m.Article
	// 	//map[string]interface{} -> struct
	// 	err = mapstructure.Decode(hit.Source, &article)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	article.Id = hit.ID
	// 	articles = append(articles, article)
	// }
	// return articles, nil
}

func (s *Store) SearchUser(ctx context.Context, query string) ([]models.User, error) {
	result, err := s.EUser.Search(ctx, s.EUser.FullTextSearchColumn, query)
	if err != nil {
		return nil, err
	}

	array := []models.User{}
	pArray, err := search[*models.User](result)
	for _, v := range pArray {
		array = append(array, *v)
	}
	return array, err
}

func search[PDBI pdbItem](result e.SearchResponse) ([]PDBI, error) {
	hits := result.Hits.Hits
	dbis := []PDBI{}
	for _, hit := range hits {
		var pdbi PDBI

		err := mapstructure.Decode(hit.Source, &pdbi)

		if err != nil {
			return nil, err
		}
		dbis = append(dbis, pdbi)
	}
	return dbis, nil
}

func (s *Store) GetArticle(ctx context.Context, id string) (models.Article, error) {
	span, s_ctx := opentracing.StartSpanFromContext(ctx, "Article Store Get called")
	defer span.Finish()

	result, err := s.EArticle.Get(s_ctx, id)
	if err != nil {
		span.LogFields(log.Error(err))
		return models.Article{}, err
	}
	var article models.Article
	err = mapstructure.Decode(result, &article)
	if err != nil {
		return models.Article{}, err
	}
	return article, nil
}

// GetUser(ctx context.Context, id uuid.UUID) (*User, error)
// func (s Store) GetUser(ctx context.Context, id string) (models.User, error) {
func (s *Store) GetUser(ctx context.Context, id uuid.UUID) (*models.User, error) {
	span, s_ctx := opentracing.StartSpanFromContext(ctx, "User Store Get called")
	defer span.Finish()

	result, err := s.EUser.Get(s_ctx, id.String())
	if err != nil {
		span.LogFields(log.Error(err))
		return &models.User{}, err
	}

	if result == nil {
		return nil, nil
	}

	var user models.User
	err = mapstructure.Decode(result, &user)
	if err != nil {
		return &models.User{}, err
	}
	return &user, nil
}

// GetUsers(ctx context.Context) ([]User, error)
func (s *Store) GetUsers(ctx context.Context) ([]models.User, error) {

	span, s_ctx := opentracing.StartSpanFromContext(ctx, "GetUsers store called")
	defer span.Finish()

	return s.SearchUser(s_ctx, "*")
}

func (r *Store) InitSchema(ctx context.Context) error {
	return nil
}

func (a *Store) GetUserArticles(ctx context.Context, userID uuid.UUID) ([]models.Article, error) {
	span, s_ctx := opentracing.StartSpanFromContext(ctx, "GetUsers store called")
	defer span.Finish()

	return a.SearchArticles(s_ctx, userID)
}
