package service

import (
	"context"
	"go-template/gin_sqlc_setup/services/auth/config"
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
	Profile(ctx context.Context, email string) (db.User, error)
}

type authService struct {
	userRepository repository.UserRepository
	config         *config.Config
}

func NewAuthService(userRepository repository.UserRepository, config *config.Config) AuthService {
	return &authService{
		userRepository: userRepository,
		config:         config,
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

	exp_time, err := time.ParseDuration(s.config.JWTExpiration)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"name":  user.Name,
		"id":    user.ID,
		"exp":   jwt.TimeFunc().Add(exp_time).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *authService) Profile(ctx context.Context, email string) (db.User, error) {
	user, err := s.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return db.User{}, err
	}

	return user, nil
}
