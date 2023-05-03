package app

import (
	"context"
	"fmt"
	"geekbrains/internal/models"
	"sync"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type cachedRepository struct {
	repository   Repository
	cache        *cache.Cache
	getUsersList []string
	// M            *sync.Mutex
	M sync.Mutex
}

func (r *cachedRepository) InitSchema(ctx context.Context) error {
	return r.repository.InitSchema(ctx)
}

func (r *cachedRepository) GetUser(ctx context.Context, id uuid.UUID) (*models.User, error) {
	key := fmt.Sprintf("user:%s", id)
	var user models.User
	err := r.cache.Get(ctx, key, &user)
	switch err {
	case nil:
		return &user, nil
	//Нет записи в кэш
	case cache.ErrCacheMiss:
		dbUser, dbErr := r.repository.GetUser(ctx, id)
		if dbErr != nil {
			return nil, dbErr
		}
		err = r.cache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   key,
			Value: dbUser,
			// TTL:   time.Hour,
			TTL: 2 * time.Second,
		})
		if err != nil {
			return nil, err
		}
		return dbUser, nil
	}
	return nil, err
}

// TODO: кэшировать запрос
func (r *cachedRepository) GetUsers(ctx context.Context) ([]models.User, error) {

	query := ""
	key := fmt.Sprintf("getUsers:%s", query)
	var users []models.User
	err := r.cache.Get(ctx, key, &users)
	switch err {
	case nil:
		return users, nil
	//Нет записи в кэш
	case cache.ErrCacheMiss:
		dbUsers, dbErr := r.repository.GetUsers(ctx)
		if dbErr != nil {
			return nil, dbErr
		}
		err = r.cache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   key,
			Value: dbUsers,
			TTL:   2 * time.Second,
		})
		if err != nil {
			return nil, err
		}
		r.M.Lock()
		r.getUsersList = append(r.getUsersList, key)
		r.M.Unlock()
		return dbUsers, nil
	}
	return nil, err
}

func (r *cachedRepository) GetUserArticles(ctx context.Context, userID uuid.UUID) ([]models.Article, error) {
	key := fmt.Sprintf("user_articles:%s", userID)
	var articles []models.Article
	err := r.cache.Get(ctx, key, &articles)
	switch err {
	case nil:
		return articles, nil
	case cache.ErrCacheMiss:
		dbArticles, dbErr := r.repository.GetUserArticles(ctx, userID)
		if dbErr != nil {
			return nil, dbErr
		}
		err = r.cache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   key,
			Value: dbArticles,
			TTL:   2 * time.Second,
		})
		if err != nil {
			return nil, err
		}
		return dbArticles, nil
	}
	return nil, err
}

func (r *cachedRepository) AddUser(ctx context.Context, user models.User) (uuid.UUID, error) {
	// r.M.Lock()
	// getUsersList := []string{}
	// copy(getUsersList, r.getUsersList)
	// r.M.Unlock()
	// for _, key := range getUsersList {
	// 	r.cache.Delete(ctx, key)
	// }

	var dbErr error

	user.ID, dbErr = r.repository.AddUser(ctx, user)
	if dbErr != nil {
		return uuid.UUID{}, dbErr
	}

	key := fmt.Sprintf("user:%s", user.ID)
	err := r.cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: user,
		TTL:   2 * time.Second,
	})

	if err != nil {
		return uuid.UUID{}, err
	}

	return user.ID, nil
}

func (r *cachedRepository) AddUserArticle(ctx context.Context, article models.Article) (uuid.UUID, error) {
	// todo
	return uuid.UUID{}, nil
}

func NewCachedRepository(repository Repository) Repository {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	rCache := cache.New(&cache.Options{
		Redis:      rdb,
		LocalCache: cache.NewTinyLFU(1000, time.Second),
	})
	return &cachedRepository{
		repository:   repository,
		cache:        rCache,
		getUsersList: []string{},
	}
}
