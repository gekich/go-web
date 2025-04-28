package service

import (
	"context"
	"github.com/gekich/go-web/internal/auth"
	"github.com/gekich/go-web/internal/db/user"
	"github.com/gekich/go-web/internal/dto/auth"
	"github.com/gekich/go-web/internal/repository"
	"github.com/gekich/go-web/internal/util/message"
	util "github.com/gekich/go-web/internal/util/string"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

type AuthService struct {
	repo       repository.UserRepository
	jwtManager *auth.JWTManager
}

func NewAuthService(repo repository.UserRepository, jwtManager *auth.JWTManager) *AuthService {
	return &AuthService{repo: repo, jwtManager: jwtManager}
}

func (s *AuthService) Register(ctx context.Context, input dto.RegisterUserInput) error {
	exists, err := s.repo.ExistsUserByEmail(ctx, input.Email)

	if err != nil {
		return err
	}
	if exists {
		return errors.ErrEmailInUse
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = s.repo.Create(ctx, user.CreateUserParams{
		Name:     input.Name,
		Email:    input.Email,
		Bio:      util.ToNullString(input.Bio),
		Password: string(hashedPassword),
	})
	return err
}

func (s *AuthService) Login(ctx context.Context, input dto.LoginInput) (string, error) {
	foundUser, err := s.repo.GetByEmail(ctx, input.Email)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(input.Password)); err != nil {
		return "", errors.ErrInvalidCredentials
	}

	signedToken, err := s.jwtManager.Generate(strconv.FormatInt(foundUser.ID, 10))

	if err != nil {
		return "", err
	}

	return signedToken, nil
}
