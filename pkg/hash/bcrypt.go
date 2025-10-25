package hash

import "golang.org/x/crypto/bcrypt"

var cost = 12

func HashPassword(plain string) (string, error) {
	bytes , err := bcrypt.GenerateFromPassword([]byte(plain), cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func VerifyPassword(hashed , plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	return err == nil
}

