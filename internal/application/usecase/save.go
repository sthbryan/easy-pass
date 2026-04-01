package usecase

import (
	"fmt"

	"github.com/sthbryan/easypass/internal/application/dto"
	"github.com/sthbryan/easypass/internal/domain/entity"
	"github.com/sthbryan/easypass/internal/domain/repository"
	"github.com/sthbryan/easypass/internal/domain/service"
)

type SaveUseCase struct {
	generator    service.PasswordGenerator
	configRepo  repository.ConfigRepository
	passwordStore repository.PasswordStore
}

func NewSaveUseCase(
	g service.PasswordGenerator,
	cr repository.ConfigRepository,
	ps repository.PasswordStore,
) *SaveUseCase {
	return &SaveUseCase{
		generator:    g,
		configRepo:  cr,
		passwordStore: ps,
	}
}

func (uc *SaveUseCase) Execute(input dto.SaveInput) (*dto.SaveOutput, error) {
	config, err := uc.configRepo.Load()
	if err != nil {
		config = entity.DefaultConfig()
	}

	secure, err := uc.generator.Generate(input.Password, input.MasterPass, config)
	if err != nil {
		return nil, err
	}

	if err := uc.passwordStore.Save(input.Name, secure, input.MasterPass); err != nil {
		return nil, fmt.Errorf("failed to save: %w", err)
	}

	return &dto.SaveOutput{
		Success: true,
		Message: fmt.Sprintf("Password '%s' saved", input.Name),
	}, nil
}
