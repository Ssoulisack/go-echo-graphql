package repositories

import (
	"my-graphql-project/domain/entities"

	"gorm.io/gorm"
)

type UserRepository interface {
	//Methods
	GetAllUsers() ([]entities.User, error)
	CreateUser(req entities.User) error
}

type userRepository struct {
	db *gorm.DB
}

// CreateUser implements UserRepository.
func (u *userRepository) CreateUser(req entities.User) error {
	tx := u.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(&req).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
	
}

// GetAllUsers implements UserRepository.
func (u *userRepository) GetAllUsers() ([]entities.User, error) {
	var users []entities.User

	if err := u.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}
