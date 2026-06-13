package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/iamhanif11/ewallet-backend/internal/dto"
	"github.com/iamhanif11/ewallet-backend/internal/repository"
	"github.com/iamhanif11/ewallet-backend/pkg"
)

var ErrEmailAlreadyExists = errors.New("Email is already registered, please use another email")

type AuthService struct {
	authRepository *repository.AuthRepository
	db             repository.DBTX
}

func NewAuthService(authRepository *repository.AuthRepository, db repository.DBTX) *AuthService {
	return &AuthService{
		authRepository: authRepository,
		db:             db,
	}
}

// register
func (as *AuthService) RegisterUser(ctx context.Context, user dto.NewUser) (dto.User, error) {
	existingUser, err := as.authRepository.GetUserByEmail(ctx, user.Email)

	if err == nil && existingUser.Email != "" {
		return dto.User{}, ErrEmailAlreadyExists
	}
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
		return "", dto.User{}, errors.New("Password not match")
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

func (as *AuthService) LogoutUser(ctx context.Context, token string, expiresAt time.Time) error {
	timeRemaining := time.Until(expiresAt)

	log.Printf("Token Expires At: %v\n", expiresAt)

	log.Printf("Sisa Waktu (TTL): %v\n", timeRemaining)
	if timeRemaining <= 0 {
		log.Println("WARNING: Token sudah expired, tidak disimpan ke Redis.")
		return nil
	}

	err := as.authRepository.BlacklistToken(ctx, token, timeRemaining)
	if err != nil {
		log.Printf("Failed blacklist token to redis: %v", err)
		return errors.New("Failed to logout")
	}

	log.Println("sukses")
	return nil
}

func (as *AuthService) CheckPinUser(ctx context.Context, email string) (bool, error) {
	return as.authRepository.CheckPinUserByEmail(ctx, email)
}

func (as *AuthService) VerifyEmail(ctx context.Context, req dto.VerifyEmailReq) error {
	exists, err := as.authRepository.CheckEmailExist(ctx, as.db, req.Email)

	if err != nil {
		return err
	}

	if !exists {
		return errors.New("email not found")
	}

	return nil
}

func (as *AuthService) ResetPassword(ctx context.Context, req dto.ResetPasswordReq) error {
	if req.NewPassword != req.ConfirmPassword {
		return errors.New("password confirmation does not match")
	}

	var hc pkg.HashConfig
	hc.UseRecommended()

	hashedPassword := hc.GenerateHash(req.NewPassword)

	err := as.authRepository.ResetPassword(ctx, as.db, req.Email, hashedPassword)

	if err != nil {
		return err
	}

	return nil

}
