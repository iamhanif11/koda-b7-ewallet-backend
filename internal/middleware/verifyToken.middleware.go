package middleware

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/iamhanif11/ewallet-backend/internal/dto"
	"github.com/iamhanif11/ewallet-backend/internal/repository"
	"github.com/iamhanif11/ewallet-backend/pkg"
)

type AuthMiddleware struct {
	authRepository *repository.AuthRepository
}

func NewAuthMiddleware(authRepository *repository.AuthRepository) *AuthMiddleware {
	return &AuthMiddleware{
		authRepository: authRepository,
	}
}

func (m *AuthMiddleware) VerifyToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bearerToken := ctx.GetHeader("Authorization")
		if bearerToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
				Message: "Unauthorized Access, Please Login",
				Success: false,
				Error:   "Unauthorized Access, Please Login",
			})
			return
		}
		splittedBearer := strings.Split(bearerToken, " ")
		if len(splittedBearer) != 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
				Message: "Unauthorized Access, Please Login",
				Success: false,
				Error:   "Invalid Token",
			})
			return
		}
		token := splittedBearer[1]
		claims := &pkg.Claims{}
		if err := claims.VerifyJWT(token); err != nil {
			log.Println("Error: ", err.Error())
			if errors.Is(err, jwt.ErrTokenInvalidIssuer) || errors.Is(err, jwt.ErrTokenExpired) {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
					Message: "Error",
					Success: false,
					Error:   err.Error(),
				})
				return
			}
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.Response{
				Message: "Error",
				Success: false,
				Error:   "Internal Server Error",
			})
			return

		}
		ctx.Set("user", claims)

		ctx.Next()
	}

}
