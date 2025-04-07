package service

import (
	"context"
	"go-template/gin_sqlc_setup/services/auth/db"
	"go-template/gin_sqlc_setup/services/auth/repository"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (string, error)
	Register(ctx context.Context, email, password, name string) (int32, error)
}

type authService struct {
	userRepository repository.UserRepository
}

func NewAuthService(userRepository repository.UserRepository) AuthService {
	return &authService{
		userRepository: userRepository,
	}
}
func (s *authService) Register(ctx context.Context, email, password, name string) (int32, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return 0, err
	}

	var nameValue pgtype.Text
	if err := nameValue.Scan(name); err != nil {
		return 0, err
	}

	user, err := s.userRepository.CreateUser(ctx, db.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
		Name:         nameValue,
	})
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepository.GetUserByEmail(ctx, email)

	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"name":  user.Name,
		"id":    user.ID,
		"exp":   jwt.TimeFunc().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte("JWT_SECRET_KEY"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
