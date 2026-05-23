package service

import (
	"context"
	"errors"
	"log"

	"github.com/iamhanif11/ewallet-backend/internal/dto"
	"github.com/iamhanif11/ewallet-backend/internal/repository"
	"github.com/iamhanif11/ewallet-backend/pkg"
)

type AuthService struct {
	authRepository *repository.AuthRepository
}

func NewAuthService(authRepository *repository.AuthRepository) *AuthService {
	return &AuthService{
		authRepository: authRepository,
	}
}

// register
func (as *AuthService) RegisterUser(ctx context.Context, user dto.NewUser) (dto.User, error) {
	var hc pkg.HashConfig
	hc.UseRecommended()
	hashPwd := hc.GenerateHash(user.Password)
	newUser, err := as.authRepository.AddUser(ctx, user.Email, hashPwd)
	if err != nil {
		return dto.User{}, err
	}
	return dto.User{
		Id:        newUser.Id,
		Email:     newUser.Email,
		CreatedAt: newUser.CreatedAt,
	}, nil
}

// login
func (as *AuthService) LoginUser(ctx context.Context, user dto.Login) (string, dto.User, error) {
	log.Println(user)
	login, err := as.authRepository.GetUserByEmail(ctx, user.Email)
	if err != nil {
		log.Printf("Database Error: %v", err)
		return "", dto.User{}, errors.New("Email or Password Invalid 1")
	}

	var hash pkg.HashConfig
	if err := hash.Compare(user.Password, login.Password); err != nil {
		return "", dto.User{}, errors.New("Email or Password Invalid")
	}

	claims := pkg.NewClaims(login.Id, user.Email)
	token, err := claims.GenerateJWT()
	if err != nil {
		return "", dto.User{}, err
	}

	log.Println(token)

	return token, dto.User{
		Email: user.Email,
	}, nil
}
