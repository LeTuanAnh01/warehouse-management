package postgres

import (
	"warehouse-api/internal/domain/entity"
	"warehouse-api/internal/domain/repository"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id uint) (*entity.User, error) {
	var user entity.User
	if err := r.db.Preload("Role").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(username string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Preload("Role").Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Preload("Role").Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindAll(limit, offset int) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64

	if err := r.db.Model(&entity.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Preload("Role").Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) Update(user *entity.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(user *entity.User) error {
	return r.db.Delete(user).Error
}

func (r *userRepository) UpdateLastLogin(id uint) error {
	return r.db.Model(&entity.User{}).Where("id = ?", id).
		Update("last_login", gorm.Expr("NOW()")).Error
}
