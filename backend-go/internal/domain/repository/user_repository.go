package repository

import "warehouse-api/internal/domain/entity"

type UserRepository interface {
	FindByID(id uint) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	FindByUsername(username string) (*entity.User, error)
	FindAll(limit, offset int) ([]entity.User, int64, error)
	Create(user *entity.User) error
	Update(user *entity.User) error
	Delete(user *entity.User) error
	UpdateLastLogin(id uint) error
}
