package password

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	// bcrypt cost factor
	cost = 12
)

// Hash 对密码进行 bcrypt 哈希
func Hash(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(bytes), nil
}

// Verify 验证密码是否匹配
func Verify(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
