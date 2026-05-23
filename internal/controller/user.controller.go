package controller

import (
	"errors"
	"log"
	"net/http"

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

func (uc *UserController) GetProfile(ctx *gin.Context) {
	claims, ok := ctx.Get("user")

	if !ok {
		log.Println("claims gaada")
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Message: "Unauthorized",
			Success: false,
			Error:   "Error",
		})
		return

	}
	log.Println("claims ", claims)
	userClaims, ok := claims.(*pkg.Claims)
	if !ok {
		log.Println(userClaims)
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Message: "Internal Server Error",
			Success: false,
			Error:   "Failed to parse user claims",
		})
		return
	}

	log.Println("check: ", userClaims.Id)
	userProfile, err := uc.userService.GetProfile(ctx.Request.Context(), userClaims.Id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.Response{
			Message: "User not found",
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Data:    userProfile,
		Message: "Get Profile Success",
		Success: true,
	})
}

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
		ctx.JSON(http.StatusBadRequest, dto.Response{
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
			ctx.JSON(http.StatusUnauthorized, dto.Response{
				Message: "Invalid PIN",
				Success: false,
			})
			return
		}
		if errors.Is(err, service.ErrPin) {
			ctx.JSON(http.StatusUnauthorized, dto.Response{
				Message: "No PIN defined",
				Success: false,
			})
			return
		}
		log.Println("cek", res)
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Message: "Internal Server Error",
			Success: false,
		})
		log.Println("cek2: ", res)
		return

	}
	ctx.JSON(http.StatusAccepted, dto.Response{
		Message: "PIN Valid",
		Success: true,
		Data:    res,
	})

}

func (uc *UserController) UpdateProfile(ctx *gin.Context) {
	claims, ok := ctx.Get("user")
	userClaims, ok := claims.(*pkg.Claims)
	if !ok {
		return
	}

	var body dto.UserUpdateProfileReq
	if err := ctx.ShouldBindBodyWith(&body, binding.JSON); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Message: "Invalid",
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	res, err := uc.userService.UpdateProfile(ctx.Request.Context(), userClaims.Id, body)
	log.Println(res)
	if err != nil {
		log.Println("err: ", err)
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Message: "Internal Server Error",
			Success: false,
		})
		return
	}

	ctx.JSON(http.StatusAccepted, dto.Response{
		Message: "Update Profile Succesfully",
		Success: true,
		Data:    res,
	})
}

func (uc *UserController) UpdatePassword(ctx *gin.Context) {
	claims, ok := ctx.Get("user")
	userClaims, ok := claims.(*pkg.Claims)
	if !ok {
		return
	}

	var body dto.UserUpdatePasswordReq
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Message: "Bad Request",
			Success: false,
		})
		return
	}

	if err := uc.userService.UpdatePassword(ctx.Request.Context(), userClaims.Id, body); err != nil {
		if errors.Is(err, service.ErrInvalidPasswd) {
			ctx.JSON(http.StatusUnauthorized, dto.Response{
				Message: "Invalid Current Password",
				Success: false,
				Error:   err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Message: "Status Internal Server Error",
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message: "Update Password Successfully",
		Success: true,
	})
}

func (uc *UserController) UpdatePin(ctx *gin.Context) {
	claims, ok := ctx.Get("user")
	userClaims, ok := claims.(*pkg.Claims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Message: "Authentication Failed",
		})
		return
	}

	var body dto.UserUpdatePinReq
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Message: "Format Data Invalid",
			Error:   err.Error(),
			Success: false,
		})
		return

	}
	if err := uc.userService.UpdatePin(ctx.Request.Context(), userClaims.Id, body); err != nil {

		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Message: "Internal Server Error",
			Success: false,
		})
		return
	}
	ctx.JSON(http.StatusOK, dto.Response{
		Message: "Update Pin Succesfully",
		Success: true,
	})

}

func (uc *UserController) GetDashboardInformation(ctx *gin.Context) {
	claims, ok := ctx.Get("user")
	userClaims, ok := claims.(*pkg.Claims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Message: "Authentication Failed",
		})
		return
	}

	res, err := uc.userService.GetDashboardInformation(ctx.Request.Context(), userClaims.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Message: "Internal Server Error",
			Error:   err.Error(),
			Success: false,
		})
		return
	}
	ctx.JSON(http.StatusOK, dto.Response{
		Message: "Dashboard Information is Displayed",
		Data:    res,
		Success: true,
	})
}

func (uc *UserController) GetTransactionReport(ctx *gin.Context) {
	duration := ctx.DefaultQuery("duration", "7d")
	if duration != "7d" {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Message: "Bad Request",
			Success: false,
		})
	}

	claims, ok := ctx.Get("user")
	userClaims, ok := claims.(*pkg.Claims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Message: "Authentication Failed",
			Success: false,
		})
		return
	}

	res, err := uc.userService.GetTransactionReport(ctx.Request.Context(), userClaims.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Message: "Internal Server Error",
			Error:   err.Error(),
			Success: false,
		})
		return
	}
	ctx.JSON(http.StatusOK, dto.Response{
		Message: "Get Data Graph Succesfully",
		Data:    res,
		Success: true,
	})
}
