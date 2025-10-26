package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/ryusuke/task_app_layerx/internal/domain"
	"github.com/ryusuke/task_app_layerx/pkg/auth"
)

// JWTMiddlewareはJWT認証を行うミドルウェア
func JWTMiddleware(jwtService auth.JWTService, userRepo domain.UserRepository, txManager domain.TxManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Authorizationヘッダーからトークンを取得
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"code":    "MISSING_TOKEN",
					"message": "missing authorization header",
				})
			}

			// Bearerプレフィックスを除去
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"code":    "INVALID_TOKEN_FORMAT",
					"message": "authorization header must be 'Bearer {token}'",
				})
			}

			// トークンを検証
			claims, err := jwtService.ParseToken(tokenString)
			if err != nil {
				if strings.Contains(err.Error(), "expired") {
					return c.JSON(http.StatusUnauthorized, map[string]string{
						"code":    "TOKEN_EXPIRED",
						"message": "token has expired, please login again",
					})
				}
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"code":    "INVALID_TOKEN",
					"message": "invalid token",
				})
			}

			// ユーザーを取得
			executor := txManager.AsExecutor()
			user, err := userRepo.FindByID(c.Request().Context(), executor, claims.UID)
			if err != nil {
				if errors.Is(err, domain.ErrUserNotFound) {
					return c.JSON(http.StatusUnauthorized, map[string]string{
						"code":    "USER_NOT_FOUND",
						"message": "user not found",
					})
				}
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"code":    "INTERNAL_ERROR",
					"message": "internal server error",
				})
			}

			// TokenVersionチェック
			if user.TokenVersion != claims.TokenVersion {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"code":    "TOKEN_INVALIDATED",
					"message": "token has been invalidated, please login again",
				})
			}

			// コンテキストにユーザーIDを保存
			c.Set("userID", user.ID)

			return next(c)
		}
	}
}

// GetUserIDはコンテキストからユーザーIDを取得
func GetUserID(c echo.Context) int64 {
	userID, ok := c.Get("userID").(int64)
	if !ok {
		return 0
	}
	return userID
}
