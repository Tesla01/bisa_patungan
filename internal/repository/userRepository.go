package repository

import (
	"gorm.io/gorm"
	"tesla01/bisa_patungan/internal/model"
)

type UserRepository interface {
	Save(user model.User) (model.User, error)
	FindByEmail(email string) (model.User, error)
	FindByID(id int) (model.User, error)
	Update(user model.User) (model.User, error)
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db}
}

func (r *UserRepositoryImpl) Save(user model.User) (model.User, error) {
	err := r.db.Create(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *UserRepositoryImpl) FindByEmail(email string) (model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).Find(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *UserRepositoryImpl) FindByID(id int) (model.User, error) {
	var user model.User
	err := r.db.Where("id = ?", id).Find(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *UserRepositoryImpl) Update(user model.User) (model.User, error) {
	err := r.db.Save(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}
