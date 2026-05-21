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
