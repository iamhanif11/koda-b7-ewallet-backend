package service

import (
	"context"
	"errors"
	"time"

	"github.com/iamhanif11/ewallet-backend/internal/dto"
	"github.com/iamhanif11/ewallet-backend/internal/model"
	"github.com/iamhanif11/ewallet-backend/internal/repository"
	"github.com/iamhanif11/ewallet-backend/pkg"
)

type UserService struct {
	userRepository *repository.UserRepository
}

var ErrPin = errors.New("Please Input PIN")
var ErrInvalidPin = errors.New("Invalid")
var ErrInvalidPasswd = errors.New("Invalid Password")

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (us *UserService) GetProfile(ctx context.Context, user_Id int) (dto.UserProfileRes, error) {
	profile, err := us.userRepository.GetProfileById(ctx, user_Id)
	if err != nil {
		return dto.UserProfileRes{}, err
	}

	return dto.UserProfileRes{
		Fullname: profile.Fullname,
		Email:    profile.Email,
		Picture:  profile.Picture,
	}, nil
}

func (us *UserService) CheckPin(ctx context.Context, user_Id int, pin string) (dto.UserCheckPinRes, error) {
	user, err := us.userRepository.GetPinById(ctx, user_Id)
	if err != nil {
		return dto.UserCheckPinRes{}, err
	}
	if user.Pin == nil {
		return dto.UserCheckPinRes{}, ErrPin
	}
	return dto.UserCheckPinRes{IsValid: true}, nil
}

func (us *UserService) UpdateProfile(ctx context.Context, user_Id int, req dto.UserUpdateProfileReq) (dto.UserUpdateProfilRes, error) {
	user, err := us.userRepository.UpdateProfileById(ctx, user_Id, req.Fullname, req.Phone, req.Picture)
	if err != nil {
		return dto.UserUpdateProfilRes{}, err
	}
	return dto.UserUpdateProfilRes{
		Fullname: *user.Fullname,
		Email:    user.Email,
		Phone:    *user.Phone,
		Picture:  *user.Picture,
	}, nil
}

func (us *UserService) UpdatePassword(ctx context.Context, userId int, req dto.UserUpdatePasswordReq) error {
	user, err := us.userRepository.GetPasswordById(ctx, userId)
	if err != nil {
		return err
	}

	var hash pkg.HashConfig
	if err := hash.Compare(req.CurrentPassword, user.Password); err != nil {
		return errors.New("Invalid Password")
	}

	hash.UseRecommended()
	hashedPassword := hash.GenerateHash(req.NewPassword)

	return us.userRepository.UpdatePasswordById(ctx, userId, hashedPassword)
}

func (us *UserService) UpdatePin(ctx context.Context, userId int, req dto.UserUpdatePinReq) error {
	return us.userRepository.UpdatedPinById(ctx, userId, req.Pin)
}

func (us *UserService) GetDashboardInformation(ctx context.Context, userId int) (dto.UserDashboardInformationRes, error) {
	dashboard, err := us.userRepository.GetDashboardInformationById(ctx, userId)
	if err != nil {
		return dto.UserDashboardInformationRes{}, err
	}

	return dto.UserDashboardInformationRes{
		Balance: dashboard.Balance,
		Income:  dashboard.Income,
		Expense: dashboard.Expense,
	}, nil
}

func (us *UserService) GetTransactionReport(ctx context.Context, userId int) ([]dto.UserTransactionReportRes, error) {
	endDate := time.Now().Truncate(24 * time.Hour)
	startDate := endDate.AddDate(0, 0, -6)

	reports, err := us.userRepository.GetTransactionReportById(ctx, userId, startDate, endDate)
	if err != nil {
		return nil, err
	}

	reportMap := make(map[string]model.UserTransactionReport, len(reports))
	for _, report := range reports {
		reportMap[report.Date.Format(time.DateOnly)] = report
	}

	res := make([]dto.UserTransactionReportRes, 0, 7)
	for i := 0; i < 7; i++ {
		currentDate := startDate.AddDate(0, 0, i)
		dateStr := currentDate.Format(time.DateOnly)

		data, found := reportMap[dateStr]

		var income, expense int
		if found {
			income = data.Income
			expense = data.Expense
		}

		res = append(res, dto.UserTransactionReportRes{
			Date:    dateStr,
			Day:     currentDate.Format("Mon"),
			Income:  income,
			Expense: expense,
		})
	}
	return res, nil
}
