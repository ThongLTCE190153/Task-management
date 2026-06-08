package services

import (
	"trithong.com/task-golang/internal/entities"
	"trithong.com/task-golang/internal/repositories"
)

type UserService interface {
	GetByID(id int) (entities.User, error)
	GetByEmail(email string) (entities.User, error)
	Update(id int, user entities.User) (entities.User, error)
	Delete(id int) error
	GetAll() ([]entities.User, error) 
	UpdateRole(id int, role string) error
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetByID(id int) (entities.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *userService) GetByEmail(email string) (entities.User, error) {
	return s.userRepo.GetByEmail(email)
}

func (s *userService) Update(id int, user entities.User) (entities.User, error) {
	return s.userRepo.Update(id, user)
}

func (s *userService) Delete(id int) error {
	return s.userRepo.Delete(id)
}

func (s *userService) GetAll() ([]entities.User, error) {
	return s.userRepo.GetAll()
}

func (s *userService) UpdateRole(id int, role string) error {
	return s.userRepo.UpdateRole(id, role)
}
