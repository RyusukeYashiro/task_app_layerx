package hash

import "golang.org/x/crypto/bcrypt"

// BcryptServiceはパスワードハッシュ化を行うインターフェース
type BcryptService interface {
	HashPassword(plain string) (string, error)
	VerifyPassword(hashed, plain string) bool
}

// bcryptServiceはBcryptServiceを実装
type bcryptService struct {
	cost int
}

// NewBcryptServiceで新しいBcryptServiceを作成する
func NewBcryptService(cost int) BcryptService {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = 12 
	}
	return &bcryptService{cost: cost}
}

// HashPasswordは平文パスワードをハッシュ化する
func (s *bcryptService) HashPassword(plain string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), s.cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// VerifyPasswordはハッシュ化されたパスワードと平文パスワードを比較する
func (s *bcryptService) VerifyPassword(hashed, plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	return err == nil
}
