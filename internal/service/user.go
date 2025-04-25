package service

import (
	"context"
	"github.com/gekich/go-web/internal/db/user"
	"github.com/gekich/go-web/internal/dto/user"
	"github.com/gekich/go-web/internal/repository"
	string2 "github.com/gekich/go-web/internal/util/string"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) GetPublicUser(ctx context.Context, id int64) (dto.UserResponse, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
		Bio:   string2.FromNullString(u.Bio),
	}, nil
}

func (s *UserService) ListUsers(ctx context.Context) ([]dto.UserResponse, error) {
	users, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	var result []dto.UserResponse
	for _, u := range users {
		result = append(result, dto.UserResponse{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
			Bio:   string2.FromNullString(u.Bio),
		})
	}
	return result, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id int64, input dto.UpdateUserInput) error {
	return s.repo.Update(ctx, user.UpdateUserParams{
		ID:   id,
		Name: input.Name,
		Bio:  string2.ToNullString(input.Bio),
	})
}

func (s *UserService) DeleteUser(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
