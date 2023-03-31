package interfaces

import (
	"context"
	"errors"
	"testApplication/models"
)

var ErrNoRows = errors.New("no field found")

type UserRepo interface {
	List(ctx context.Context, offset int, limit int) ([]models.User, error)
	ById(ctx context.Context, id int) (models.User, error)
	ByEmail(ctx context.Context, email string) (models.User, error)
	CreateUser(ctx context.Context, newUser models.User) (models.User, error)
	UpdateUser(ctx context.Context, user models.User) (models.User, error)
	DeleteUser(ctx context.Context, id int) (models.User, error)

	UpdateRoles(ctx context.Context, user models.User, roles []models.Role) (models.User, error)

	GetAllUserGrants(ctx context.Context, user int) (grants map[string][]string, err error)
}
