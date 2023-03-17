package interfaces

import (
	"context"
	"testApplication/models"
)

type ClientRepo interface {
	GetClients(ctx context.Context, offset int, limit int) []models.Client
	GetClientById(ctx context.Context, id int) (models.Client, error)
	CreateClient(context.Context, models.Client) (models.Client, error)
	UpdateClient(context.Context, models.Client) error
	DeleteClient(ctx context.Context, id int) error
}
