package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/iamhanif11/ewallet-backend/internal/dto"
	"github.com/iamhanif11/ewallet-backend/internal/service"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// User Register
//
//	@Summary		Register a user
//	@Description	Create a new user account for E-Wallet
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body body dto.NewUser true "register payload"
//	@Success		201	{object}	dto.Response[dto.User]
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/auth/register [post]
func (a *AuthController) Register(ctx *gin.Context) {
	var body dto.NewUser
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Println("Error: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Registration failed",
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	res, err := a.authService.RegisterUser(ctx.Request.Context(), body)
	if err != nil {
		log.Println("Error: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Email is Registered",
			Success: false,
			Error:   "Internal Server Error",
		})
		return
	}
	ctx.JSON(http.StatusCreated, dto.Response[dto.User]{
		Data:    res,
		Message: "Register Success",
		Success: true,
	})
}

// User Login
//
//	@Summary		Login into a user
//	@Description	Login into user for E-Wallet
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body body dto.Login true "login payload"
//	@Success		200	{object}	dto.Response[dto.LoginResponse]
//	@Failure		500	{object}	dto.ErrorResponse
//	@Failure		401	{object}	dto.ErrorResponse
//	@Router			/auth/ [post]
func (ac *AuthController) Login(ctx *gin.Context) {
	var body dto.Login
	if err := ctx.ShouldBindBodyWith(&body, binding.JSON); err != nil {
		log.Println("Error: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Error",
			Success: false,
			Error:   "Internal Server Error",
		})
		return
	}

	token, _, err := ac.authService.LoginUser(ctx.Request.Context(), body)
	if err != nil {
		log.Println("Error: ", err.Error())
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Login Failed",
			Success: false,
			Error:   "Unauthorized Access",
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response[dto.LoginResponse]{
		Message: "Login Succesfully",
		Success: true,
		Data: dto.LoginResponse{
			Token: token,
		},
	},
	)
}
