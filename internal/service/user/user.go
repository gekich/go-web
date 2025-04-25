package service

import (
	"context"
	"github.com/gekich/go-web/internal/db/user"
	"github.com/gekich/go-web/internal/dto/user"
	repository "github.com/gekich/go-web/internal/repository/user"
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
		Bio:   u.Bio,
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
			Bio:   u.Bio,
		})
	}
	return result, nil
}

func (s *UserService) CreateUser(ctx context.Context, input dto.CreateUserInput) (dto.UserResponse, error) {
	u, err := s.repo.Create(ctx, user.CreateUserParams{
		Name:  input.Name,
		Email: input.Email,
		Bio:   input.Bio,
	})
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
		Bio:   u.Bio,
	}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id int64, input dto.UpdateUserInput) error {
	return s.repo.Update(ctx, user.UpdateUserParams{
		ID:   id,
		Name: input.Name,
		Bio:  input.Bio,
	})
}

func (s *UserService) DeleteUser(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
