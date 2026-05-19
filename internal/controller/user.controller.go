package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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
