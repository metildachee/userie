package interfaces

import (
	"github.com/metildachee/userie/models"
)

type UserDao interface {
	create(u models.User) error

	Create(u models.User) error
	BatchCreate(u []models.User) error
	Update(u models.User) error
	UpdateName(id, name string) error
	Delete(id string) error
	GetById(id string) (models.User, error)
	GetAll(limit int) ([]models.User, error)
}
