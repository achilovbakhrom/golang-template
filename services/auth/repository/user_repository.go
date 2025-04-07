package repository

import (
	"context"
	"go-template/gin_sqlc_setup/services/auth/db"
)

type UserRepository interface {
	CreateUser(context context.Context, user db.User) (db.User, error)
	GetUserByID(context context.Context, id int64) (db.User, error)
	GetUserByEmail(context context.Context, email string) (db.User, error)
}

type userRepository struct {
	queries *db.Queries
}

func NewUserRepository(queries *db.Queries) UserRepository {
	return &userRepository{
		queries: queries,
	}
}

func (r *userRepository) CreateUser(context context.Context, user db.User) (db.User, error) {

	createdUser, err := r.queries.CreateUser(context, db.CreateUserParams{
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Name:         user.Name,
	})
	if err != nil {
		return db.User{}, err
	}
	return createdUser, nil
}
func (r *userRepository) GetUserByID(context context.Context, id int64) (db.User, error) {
	user, err := r.queries.GetUserByID(context, int32(id))
	if err != nil {
		return db.User{}, err
	}
	return user, nil
}
func (r *userRepository) GetUserByEmail(context context.Context, email string) (db.User, error) {
	user, err := r.queries.GetUserByEmail(context, email)
	if err != nil {
		return db.User{}, err
	}
	return user, nil
}
