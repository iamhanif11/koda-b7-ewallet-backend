package controller

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	// "github.com/gin-gonic/gin/binding"
	"github.com/iamhanif11/ewallet-backend/internal/dto"
	"github.com/iamhanif11/ewallet-backend/internal/service"
	"github.com/iamhanif11/ewallet-backend/pkg"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// Profile User
//
//	@Summary		Get current user Profile
//	@Description	retrieved detailed information of currently logged in
//	@Tags			user
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.Response[dto.UserProfileRes]
//	@Failure		500	{object}	dto.ErrorResponse
//	@Failure		401	{object}	dto.ErrorResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Router			/user/profile [get]
func (uc *UserController) GetProfile(ctx *gin.Context) {
	claims, ok := ctx.Get("user")

	if !ok {
		log.Println("claims gaada")
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Unauthorized Access",
			Success: false,
			Error:   "Error",
		})
		return

	}
	log.Println("claims ", claims)
	userClaims, ok := claims.(*pkg.Claims)
	if !ok {
		log.Println(userClaims)
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Internal Server Error",
			Success: false,
			Error:   "Failed to parse user claims",
		})
		return
	}

	log.Println("check: ", userClaims.Id)
	userProfile, err := uc.userService.GetProfile(ctx.Request.Context(), userClaims.Id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse{
			Message: "User not found",
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response[dto.UserProfileRes]{
		Data:    userProfile,
		Message: "Get Profile Success",
		Success: true,
	})
}

// Check PIN
//
//	@Summary		Check pin user
//	@Description	Verify if the provided pin is correct for the logged in user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			payload			body	dto.UserCheckPinReq	true	"PIN payload"
//	@Success		202	{object}	dto.Response[dto.UserCheckPinRes]
//	@Failure		500	{object}	dto.ErrorResponse
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		401	{object}	dto.ErrorResponse
//	@Router			/user/profile/pin/check [post]
func (uc *UserController) CheckPin(ctx *gin.Context) {
	claims, ok := ctx.Get("user")
	userClaims, ok := claims.(*pkg.Claims)
	log.Println("cek: ", userClaims)

	if !ok {
		return
	}

	var body dto.UserCheckPinReq
	log.Println("cek body", body)
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Bad request",
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	res, err := uc.userService.CheckPin(ctx.Request.Context(), userClaims.Id, body.Pin)

	if err != nil {
		log.Printf("[ERROR] CheckPin: %v", err)
		if errors.Is(err, service.ErrPin) {
			ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Message: "Invalid PIN",
				Success: false,
			})
			return
		}
		if errors.Is(err, service.ErrPin) {
			ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Message: "No PIN defined",
				Success: false,
			})
			return
		}
		log.Println("cek", res)
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Internal Server Error",
			Success: false,
		})
		log.Println("cek2: ", res)
		return

	}
	ctx.JSON(http.StatusAccepted, dto.Response[dto.UserCheckPinRes]{
		Message: "PIN Valid",
		Success: true,
		Data:    res,
	})

}

// Update Profile
//
//	@Summary		Update User Profile
//	@Description	Update detailed information of the currently logged in user
//	@Tags			user
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			fullname		formData	string	false	"Update Fullname"
//	@Param			phone			formData	string	false	"Update Phone"
//	@Param			picture			formData	file	false	"Update Profile Picture"
//	@Success		202	{object}	dto.Response[dto.UserUpdateProfilRes]
//	@Failure		500	{object}	dto.ErrorResponse
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		401	{object}	dto.ErrorResponse
//	@Failure		422	{object}	dto.ErrorResponse
//	@Router			/user/profile 	[patch]
func (uc *UserController) UpdateProfile(ctx *gin.Context) {
	claims, ok := ctx.Get("user")
	userClaims, ok := claims.(*pkg.Claims)
	if !ok {
		return
	}

	var body dto.UserUpdateProfileReq
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid Input Data",
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	var pictureURL *string

	if body.Picture != nil {
		const maxUploadSize = 1024 * 1024
		if body.Picture.Size > maxUploadSize {
			ctx.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
				Message: "File too large",
				Success: false,
			})
			return
		}

		ext := strings.ToLower(filepath.Ext(body.Picture.Filename))
		log.Println(ext)
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			ctx.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
				Message: "Invalid file format",
				Success: false,
			})
			return
		}

		filename := fmt.Sprintf("user_%d%s", time.Now().UnixNano(), ext)
		dst := filepath.Join("public", "img", "profiles", filename)

		if err := ctx.SaveUploadedFile(body.Picture, dst); err != nil {
			ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Message: "Failed to save image",
				Success: false,
			})
			return
		}

		generatedURL := "/img/profile/" + filename
		pictureURL = &generatedURL
	}

	res, err := uc.userService.UpdateProfile(ctx.Request.Context(), userClaims.Id, body, pictureURL)
	log.Println(res)
	if err != nil {
		log.Println("err: ", err)
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Upload Failed",
			Success: false,
		})
		return
	}

	ctx.JSON(http.StatusAccepted, dto.Response[dto.UserUpdateProfilRes]{
		Message: "Update Profile Succesfully",
		Success: true,
		Data:    res,
	})
}

// Update Password
//
//	@Summary		Update User Password
//	@Description	Change the current user password to new password
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			request			body	dto.UserUpdatePasswordReq	true	"Password Update"
//	@Success		202	{object}	dto.Response[any]
//	@Failure		500	{object}	dto.ErrorResponse
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		401	{object}	dto.ErrorResponse
//	@Router			/user/password 	[patch]
func (uc *UserController) UpdatePassword(ctx *gin.Context) {
	claims, ok := ctx.Get("user")
	userClaims, ok := claims.(*pkg.Claims)
	if !ok {
		return
	}

	var body dto.UserUpdatePasswordReq
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Bad Request",
			Success: false,
		})
		return
	}

	if err := uc.userService.UpdatePassword(ctx.Request.Context(), userClaims.Id, body); err != nil {
		if errors.Is(err, service.ErrInvalidPasswd) {
			ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Message: "Invalid Current Password",
				Success: false,
				Error:   err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Status Internal Server Error",
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response[any]{
		Message: "Update Password Successfully",
		Success: true,
	})
}

// Update PIN
//
//	@Summary		Update User PIN
//	@Description	Change or set new security PIN for the user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			request			body	dto.UserUpdatePinReq		true	"PIN Update"
//	@Success		200	{object}	dto.Response[any]
//	@Failure		500	{object}	dto.ErrorResponse
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		401	{object}	dto.ErrorResponse
//	@Router			/user/pin	 	[patch]
func (uc *UserController) UpdatePin(ctx *gin.Context) {
	claims, ok := ctx.Get("user")
	userClaims, ok := claims.(*pkg.Claims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Authentication Failed",
		})
		return
	}

	var body dto.UserUpdatePinReq
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Format Data Invalid",
			Error:   err.Error(),
			Success: false,
		})
		return

	}
	if err := uc.userService.UpdatePin(ctx.Request.Context(), userClaims.Id, body); err != nil {

		if errors.Is(err, service.ErrInvalidPin) {
			ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Message: "Current PIN wrong!",
				Success: false,
			})
			return
		}

		if errors.Is(err, service.ErrPin) {
			ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Message: "Please Input PIN",
				Success: false,
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Internal Server Error",
			Success: false,
		})
		return
	}
	ctx.JSON(http.StatusOK, dto.Response[any]{
		Message: "Update Pin Succesfully",
		Success: true,
	})

}

// Get Dashboard Info
//
//	@Summary		Get user Dashboard information
//	@Description	Retrieve data for the user dashboard
//	@Tags			user
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	dto.Response[dto.UserDashboardInformationRes]
//	@Failure		500	{object}	dto.ErrorResponse
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		401	{object}	dto.ErrorResponse
//	@Router			/user/wallet	[get]
func (uc *UserController) GetDashboardInformation(ctx *gin.Context) {
	claims, ok := ctx.Get("user")
	userClaims, ok := claims.(*pkg.Claims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Authentication Failed",
		})
		return
	}

	res, err := uc.userService.GetDashboardInformation(ctx.Request.Context(), userClaims.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Internal Server Error",
			Error:   err.Error(),
			Success: false,
		})
		return
	}
	ctx.JSON(http.StatusOK, dto.Response[dto.UserDashboardInformationRes]{
		Message: "Dashboard Information is Displayed",
		Data:    res,
		Success: true,
	})
}

// Get Transaction
//
//	@Summary		Get User Transaction Report
//	@Description	Retrieve transaction reports with duration
//	@Tags			user
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			duration		query 	string									false	"Duration of report, default is 7d" default(7d)
//	@Success		200	{object}	dto.Response[[]dto.UserTransactionReportRes]
//	@Failure		500	{object}	dto.ErrorResponse
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		401	{object}	dto.ErrorResponse
//	@Router			/user/reports	[get]
func (uc *UserController) GetTransactionReport(ctx *gin.Context) {
	duration := ctx.DefaultQuery("duration", "7d")
	if duration != "7d" {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Bad Request",
			Success: false,
		})
	}

	claims, ok := ctx.Get("user")
	userClaims, ok := claims.(*pkg.Claims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Authentication Failed",
			Success: false,
		})
		return
	}

	res, err := uc.userService.GetTransactionReport(ctx.Request.Context(), userClaims.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Internal Server Error",
			Error:   err.Error(),
			Success: false,
		})
		return
	}
	ctx.JSON(http.StatusOK, dto.Response[[]dto.UserTransactionReportRes]{
		Message: "Get Data Graph Succesfully",
		Data:    res,
		Success: true,
	})
}
