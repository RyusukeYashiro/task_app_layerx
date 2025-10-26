package auth

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Clockは時刻を取得する関数型
type Clock func() time.Time

type JWTClaims struct {
	UID         int64 `json:"uid"`   
	TokenVersion int   `json:"tkn_ver"` 
	jwt.RegisteredClaims
}

// JWTServiceはJWTの生成と解析を行うインターフェース
type JWTService interface {
	GenerateToken(userID int64, tokenVersion int) (string, error)
	ParseToken(tokenString string) (*JWTClaims, error)
}

// jwtServiceはJWTServiceの実装
type jwtService struct {
	secret []byte
	issuer string
	ttl    time.Duration
	clock  Clock
}


func NewJWTService(secret, issuer string, ttl time.Duration, clock Clock) JWTService {
	if clock == nil {
		clock = time.Now
	}
	return &jwtService{
		secret: []byte(secret),
		issuer: issuer,
		ttl:    ttl,
		clock:  clock,
	}
}

// GenerateTokenはJWTを生成する
func (s *jwtService) GenerateToken(userID int64, tokenVersion int) (string, error) {
	now := s.clock()
	claims := &JWTClaims{
		UID:          userID,
		TokenVersion: tokenVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.ttl)),
			Subject:   strconv.FormatInt(userID, 10), 
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString(s.secret)
}

// ParseTokenはJWTを解析する
func (s *jwtService) ParseToken(tokenString string) (*JWTClaims, error) {
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}), 
		jwt.WithIssuer(s.issuer),                                    
		jwt.WithLeeway(30*time.Second),                              
	)

	token, err := parser.ParseWithClaims(tokenString, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token expired")
		}
		if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, errors.New("token not valid yet")
		}
		return nil, errors.New("token invalid")
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token invalid")
	}
	return claims, nil
}
