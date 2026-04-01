package usecase

import (
	"fmt"

	"github.com/sthbryan/easypass/internal/application/dto"
	"github.com/sthbryan/easypass/internal/domain/repository"
)

type ShowUseCase struct {
	passwordStore repository.PasswordStore
}

func NewShowUseCase(ps repository.PasswordStore) *ShowUseCase {
	return &ShowUseCase{passwordStore: ps}
}

func (uc *ShowUseCase) Execute(input dto.ShowInput) (*dto.ShowOutput, error) {
	password, err := uc.passwordStore.Get(input.Name, input.MasterPass)
	if err != nil {
		return nil, fmt.Errorf("password not found: %w", err)
	}

	return &dto.ShowOutput{
		Password: password,
		Copied:   true,
	}, nil
}

func (uc *ShowUseCase) ListNames() ([]string, error) {
	return uc.passwordStore.List()
}
