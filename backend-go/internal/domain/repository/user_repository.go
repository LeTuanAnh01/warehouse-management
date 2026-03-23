package repository

import (
	"warehouse-api/internal/domain/entity"
)

type UserRepository interface {
	Create(user *entity.User) error
	FindByID(id uint) (*entity.User, error)
	FindByUsername(username string) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	Update(user *entity.User) error
	UpdateLastLogin(id uint) error
}