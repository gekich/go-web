package service

import (
	"context"
	"github.com/gekich/go-web/internal/db/user"
	"github.com/gekich/go-web/internal/util/message"
	util "github.com/gekich/go-web/internal/util/string"
	"github.com/golang-jwt/jwt/v4"
	"time"

	"github.com/gekich/go-web/internal/dto/auth"
	"github.com/gekich/go-web/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo        repository.UserRepository
	jwtSecret   string
	tokenExpiry time.Duration
}

func NewAuthService(repo repository.UserRepository, jwtSecret string, tokenExpiry time.Duration) *AuthService {
	return &AuthService{repo: repo, jwtSecret: jwtSecret, tokenExpiry: tokenExpiry}
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": foundUser.ID,
		"exp":     time.Now().Add(s.tokenExpiry).Unix(),
	})

	signedToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
