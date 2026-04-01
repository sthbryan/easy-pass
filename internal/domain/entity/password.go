package entity

import "errors"

var (
	ErrEmptyPassword   = errors.New("password cannot be empty")
	ErrEmptyMasterPass = errors.New("master password cannot be empty")
)

type Password struct {
	Name     string
	Secure   string
	Platform string
}

func NewPassword(name, secure, platform string) *Password {
	return &Password{
		Name:     name,
		Secure:   secure,
		Platform: platform,
	}
}
