package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/iamhanif11/ewallet-backend/internal/dto"
	"github.com/iamhanif11/ewallet-backend/internal/model"
	"github.com/iamhanif11/ewallet-backend/internal/repository"
	"github.com/iamhanif11/ewallet-backend/pkg"
	"github.com/redis/go-redis/v9"
)

type UserService struct {
	userRepository *repository.UserRepository
	rdb            *redis.Client
}

var ErrPin = errors.New("Please Input PIN")
var ErrInvalidPin = errors.New("Invalid")
var ErrInvalidPasswd = errors.New("Invalid Password")

func NewUserService(userRepository *repository.UserRepository, rdb *redis.Client) *UserService {
	return &UserService{
		userRepository: userRepository,
		rdb:            rdb,
	}
}

func (us *UserService) GetProfile(ctx context.Context, user_Id int) (dto.UserProfileRes, error) {
	cacheKey := fmt.Sprintf("user:%d:profile", user_Id)

	cachedData, err := us.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var profilRes dto.UserProfileRes

		if errUnmarshal := json.Unmarshal([]byte(cachedData), &profilRes); errUnmarshal == nil {
			return profilRes, nil
		}

	}
	profile, err := us.userRepository.GetProfileById(ctx, user_Id)
	if err != nil {
		return dto.UserProfileRes{}, err
	}

	res := dto.UserProfileRes{
		Fullname: profile.Fullname,
		Email:    profile.Email,
		Picture:  profile.Picture,
	}

	if resBytes, errMarshal := json.Marshal(res); errMarshal == nil {
		us.rdb.Set(ctx, cacheKey, resBytes, 24*time.Hour)
	}

	return res, nil
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

func (us *UserService) UpdateProfile(ctx context.Context, user_Id int, req dto.UserUpdateProfileReq, pictureURL *string) (dto.UserUpdateProfilRes, error) {
	user, err := us.userRepository.UpdateProfileById(ctx, user_Id, req.Fullname, req.Phone, pictureURL)
	if err != nil {
		return dto.UserUpdateProfilRes{}, err
	}

	cacheKey := fmt.Sprintf("user:%d:profile", user_Id)
	us.rdb.Del(ctx, cacheKey)

	var fullname, phone, picture string

	if user.Fullname != nil {
		fullname = *user.Fullname
	}

	if user.Phone != nil {
		phone = *user.Phone
	}

	if user.Picture != nil {
		picture = *user.Picture
	}
	return dto.UserUpdateProfilRes{
		Fullname: fullname,
		Phone:    phone,
		Picture:  picture,
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
	// cacheKey := fmt.Sprintf("user:%d:dashboard", userId)
	balanceKey := fmt.Sprintf("user:%d:balance", userId)
	incomeKey := fmt.Sprintf("user:%d:income", userId)
	expenseKey := fmt.Sprintf("user:%d:expense", userId)

	// cacheData, err := us.rdb.Get(ctx, cacheKey).Result()
	// if err == nil {
	// 	var dashboardres dto.UserDashboardInformationRes

	// 	if errUnmarshal := json.Unmarshal([]byte(cacheData), &dashboardres); errUnmarshal == nil {
	// 		return dashboardres, nil
	// 	}
	// }

	var res dto.UserDashboardInformationRes
	cacheMiss := false

	if balData, err := us.rdb.Get(ctx, balanceKey).Bytes(); err == nil {
		if errUnmarshal := json.Unmarshal(balData, &res.Balance); errUnmarshal != nil {
			cacheMiss = true
		}
	} else {
		cacheMiss = true
	}

	if incData, err := us.rdb.Get(ctx, incomeKey).Bytes(); err == nil {
		if errUnmarshal := json.Unmarshal(incData, &res.Income); errUnmarshal != nil {
			cacheMiss = true
		}
	} else {
		cacheMiss = true
	}

	if expData, err := us.rdb.Get(ctx, expenseKey).Bytes(); err == nil {
		if errUnmarshal := json.Unmarshal(expData, &res.Expense); errUnmarshal != nil {
			cacheMiss = true
		}
	} else {
		cacheMiss = true
	}

	if !cacheMiss {
		return res, nil
	}

	dashboard, err := us.userRepository.GetDashboardInformationById(ctx, userId)
	if err != nil {
		return dto.UserDashboardInformationRes{}, err
	}

	res.Balance = dashboard.Balance
	res.Income = dashboard.Income
	res.Expense = dashboard.Expense

	ttl := 15 * time.Minute

	if balByte, err := json.Marshal(res.Balance); err == nil {
		us.rdb.Set(ctx, balanceKey, balByte, ttl)
	}
	if incByte, err := json.Marshal(res.Income); err == nil {
		us.rdb.Set(ctx, incomeKey, incByte, ttl)
	}
	if expByte, err := json.Marshal(res.Expense); err == nil {
		us.rdb.Set(ctx, expenseKey, expByte, ttl)
	}

	return res, nil

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
