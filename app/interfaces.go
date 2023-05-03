package app

import (
	"context"
	"geekbrains/internal/models"

	"github.com/google/uuid"
)

type Repository interface {
	InitSchema(ctx context.Context) error
	GetUser(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetUsers(ctx context.Context) ([]models.User, error)
	GetUserArticles(ctx context.Context, userID uuid.UUID) ([]models.Article, error)

	AddUser(ctx context.Context, user models.User) (uuid.UUID, error)
	AddUserArticle(ctx context.Context, article models.Article) (uuid.UUID, error)
}
