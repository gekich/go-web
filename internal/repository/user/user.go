package repository

import (
	"context"

	db "github.com/gekich/go-web/internal/db/user"
)

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (db.User, error)
	List(ctx context.Context) ([]db.User, error)
	Create(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	Update(ctx context.Context, arg db.UpdateUserParams) error
	Delete(ctx context.Context, id int64) error
}

type userRepository struct {
	q *db.Queries
}

func NewUserRepository(q *db.Queries) UserRepository {
	return &userRepository{q: q}
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (db.User, error) {
	return r.q.GetUser(ctx, id)
}

func (r *userRepository) List(ctx context.Context) ([]db.User, error) {
	return r.q.ListUsers(ctx)
}

func (r *userRepository) Create(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	return r.q.CreateUser(ctx, arg)
}

func (r *userRepository) Update(ctx context.Context, arg db.UpdateUserParams) error {
	return r.q.UpdateUser(ctx, arg)
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	return r.q.DeleteUser(ctx, id)
}
