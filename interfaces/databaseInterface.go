package interfaces

import "testApplication/models"

type IDatabase interface {
	GetClients() []models.Client
	GetClientById(id int) (models.Client, error)
	CreateClient() (int, error)
	UpdateClient() (int, error)
	DeleteClient() (int, error)
}
