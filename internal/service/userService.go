package service

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"tesla01/bisa_patungan/internal/model"
	repository2 "tesla01/bisa_patungan/internal/repository"
)

type UserService interface {
	RegisterUser(input model.RegisterUserInput) (model.User, error)
	LoginUser(input model.LoginInput) (model.User, error)
	IsEmailAvailable(input model.CheckEmailInput) (bool, error)
	SaveAvatar(id int, fileLocation string) (model.User, error)
	GetUserByID(id int) (model.User, error)
}

type UserServiceImpl struct {
	repository repository2.UserRepository
}

func NewUserService(repository repository2.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{repository}
}

func (s *UserServiceImpl) RegisterUser(input model.RegisterUserInput) (model.User, error) {
	user := model.User{}
	user.Name = input.Name
	user.Email = input.Email
	user.Occupation = input.Occupation

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}

	user.PasswordHash = string(passwordHash)
	user.Role = "user"

	newUser, err := s.repository.Save(user)
	if err != nil {
		return newUser, err
	}

	return newUser, nil

}

func (s *UserServiceImpl) LoginUser(input model.LoginInput) (model.User, error) {
	email := input.Email
	password := input.Password

	user, err := s.repository.FindByEmail(email)

	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))

	if err != nil {
		return user, err
	}

	return user, nil

}

func (s *UserServiceImpl) IsEmailAvailable(input model.CheckEmailInput) (bool, error) {
	email := input.Email

	user, err := s.repository.FindByEmail(email)

	if err != nil {
		return true, err
	}

	if user.ID == 0 {
		return true, nil
	}

	return false, nil
}

func (s *UserServiceImpl) SaveAvatar(id int, fileLocation string) (model.User, error) {
	user, err := s.repository.FindByID(id)

	if err != nil {
		return user, err
	}

	user.AvatarFileName = fileLocation

	updatedUser, err := s.repository.Update(user)

	if err != nil {
		return updatedUser, err
	}

	return updatedUser, nil

}

func (s *UserServiceImpl) GetUserByID(id int) (model.User, error) {
	user, err := s.repository.FindByID(id)
	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("user not found")
	}

	return user, nil
}
