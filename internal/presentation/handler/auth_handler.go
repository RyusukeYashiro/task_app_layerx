package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ryusuke/task_app_layerx/internal/presentation/middleware"
	authuc "github.com/ryusuke/task_app_layerx/internal/usecase/auth"
)

// AuthHandlerは認証関連のHTTPハンドラー
type AuthHandler struct {
	authUseCase *authuc.AuthUseCase
}

// NewAuthHandlerで新しいAuthHandlerを作成
func NewAuthHandler(authUseCase *authuc.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
	}
}

// Signupはユーザー登録を処理
// POST /auth/signup
func (h *AuthHandler) Signup(c echo.Context) error {
	var request SignupRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "invalid request body",
		})
	}

	usecaseRequest := authuc.SignupRequest{
		Email:    request.Email,
		Password: request.Password,
		Name:     request.Name,
	}

	resp, err := h.authUseCase.Signup(c.Request().Context(), usecaseRequest)
	if err != nil {
		return HandleError(c, err)
	}

	return c.JSON(http.StatusCreated, AuthResponse{
		Token: resp.Token,
		User: UserResponse{
			ID:    resp.User.ID,
			Email: resp.User.Email,
			Name:  resp.User.Name,
		},
	})
}

// Loginはログインを処理
// POST /auth/login
func (h *AuthHandler) Login(c echo.Context) error {
	var request LoginRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "invalid request body",
		})
	}

	usecaseRequest := authuc.LoginRequest{
		Email:    request.Email,
		Password: request.Password,
	}

	resp, err := h.authUseCase.Login(c.Request().Context(), usecaseRequest)
	if err != nil {
		return HandleError(c, err)
	}

	return c.JSON(http.StatusOK, AuthResponse{
		Token: resp.Token,
		User: UserResponse{
			ID:    resp.User.ID,
			Email: resp.User.Email,
			Name:  resp.User.Name,
		},
	})
}

// Logoutはログアウトを処理
// POST /auth/logout
func (h *AuthHandler) Logout(c echo.Context) error {
	// JWTミドルウェアでセットされたユーザーIDを取得
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    "UNAUTHORIZED",
			Message: "unauthorized",
		})
	}

	// UseCaseを呼び出し
	if err := h.authUseCase.Logout(c.Request().Context(), userID); err != nil {
		return HandleError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}
