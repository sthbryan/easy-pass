package repository

import "github.com/sthbryan/easypass/internal/domain/entity"

type ConfigRepository interface {
	Save(config *entity.PasswordConfig) error
	Load() (*entity.PasswordConfig, error)
}

type PasswordStore interface {
	Save(name, secure, masterPass string) error
	Get(name, masterPass string) (string, error)
	List() ([]string, error)
	Delete(name string) error
}
