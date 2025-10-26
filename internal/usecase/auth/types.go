package auth

// SignupRequestはユーザー登録のリクエスト
type SignupRequest struct {
	Email    string
	Password string
	Name     string
}

// LoginRequestはログインのリクエスト
type LoginRequest struct {
	Email    string
	Password string
}

// AuthResponseは認証成功時のレスポンス
type AuthResponse struct {
	Token string
	User  UserResponse
}

// UserResponseはユーザー情報のレスポンス
type UserResponse struct {
	ID    int64
	Email string
	Name  string
}
