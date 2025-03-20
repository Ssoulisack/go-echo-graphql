package resolver

import (
	"my-graphql-project/api/graph"
	"my-graphql-project/data/repositories"
	"my-graphql-project/data/services"

	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	UserSvc services.UserService
}

// Mutation returns graph.MutationResolver implementation.
func (r *Resolver) Mutation() graph.MutationResolver { return &mutationResolver{r} }

// Query returns graph.QueryResolver implementation.
func (r *Resolver) Query() graph.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

func NewResolver(userSvc services.UserService) *Resolver {
	return &Resolver{
		UserSvc: userSvc,
	}
}

func InitializeResolver(db *gorm.DB) *Resolver {
	userRepo := repositories.NewUserRepository(db)
	userSvc := services.NewUserService(userRepo)
	return NewResolver(userSvc)
}
