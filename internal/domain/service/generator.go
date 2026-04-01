package service

import (
	"crypto/sha256"
	"errors"
	"strings"

	"github.com/sthbryan/easypass/internal/domain/entity"
	"golang.org/x/crypto/argon2"
)

var (
	uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowercase = "abcdefghijklmnopqrstuvwxyz"
	numbers   = "0123456789"
	symbols   = "!@#$%^&*()_+-=[]{}|;:,.<>?"
)

type PasswordGenerator interface {
	Generate(password, masterPass string, config *entity.PasswordConfig) (string, error)
}

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) Generate(password, masterPass string, config *entity.PasswordConfig) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}
	if masterPass == "" {
		return "", errors.New("master password cannot be empty")
	}

	charset := g.buildCharset(config)
	if charset == "" {
		return "", errors.New("charset cannot be empty")
	}

	derived, err := g.deriveKey(password, masterPass, uint32(config.Length))
	if err != nil {
		return "", err
	}

	result := make([]byte, config.Length)
	for i, b := range derived {
		result[i] = charset[int(b)%len(charset)]
	}

	return string(result), nil
}

func (g *Generator) buildCharset(config *entity.PasswordConfig) string {
	var charset strings.Builder

	if config.UseUppercase {
		charset.WriteString(uppercase)
	}
	if config.UseLowercase {
		charset.WriteString(lowercase)
	}
	if config.UseNumbers {
		charset.WriteString(numbers)
	}
	if config.UseSymbols {
		charset.WriteString(symbols)
	}

	return charset.String()
}

func (g *Generator) deriveKey(password, salt string, length uint32) ([]byte, error) {

	saltSource := salt + password
	saltBytes := sha256.Sum256([]byte(saltSource))

	key := argon2.IDKey(
		[]byte(password),
		saltBytes[:],
		1,
		64*1024,
		2,
		length,
	)

	return key, nil
}
