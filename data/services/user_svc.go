package services

import (
	"my-graphql-project/data/repositories"
	"my-graphql-project/domain/entities"
	"my-graphql-project/domain/models"
)

type UserService interface {
	//Methods
	GetAllUsers() ([]models.User, error)
	CreateUser(req models.User) error
}

type userService struct {
	userRepo repositories.UserRepository
}

// GetAllUsers implements UserService.
func (u *userService) GetAllUsers() ([]models.User, error) {
	users, err := u.userRepo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	var res []models.User
	for _, user := range users {
		res = append(res, models.User{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		})
	}
	return res, nil
}

func (u *userService) CreateUser(req models.User) error {
	data := entities.User{
		Name:  req.Name,
		Email: req.Email,
	}
	if err := u.userRepo.CreateUser(data); err != nil {
		return err
	}
	return nil
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}
