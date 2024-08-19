package hash

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const maxPasswordLength = 72

type BcryptHashService struct{}

func (b *BcryptHashService) HashPassword(password string) (string, error) {
	logger := logrus.New()
	logger.SetOutput(logrus.StandardLogger().Out)
	logger.SetLevel(logrus.InfoLevel)
	logger.Info("Hashing password")

	if len(password) > maxPasswordLength {
		return "", errors.New("password length exceeds 72 bytes")
	}

	logger.Info("Generating hashed password")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	logger.Info("Generated hashed password")
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

func (b *BcryptHashService) ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
