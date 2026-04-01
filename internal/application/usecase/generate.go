package usecase

import (
	"github.com/sthbryan/easypass/internal/application/dto"
	"github.com/sthbryan/easypass/internal/domain/entity"
	"github.com/sthbryan/easypass/internal/domain/repository"
	"github.com/sthbryan/easypass/internal/domain/service"
)

type GenerateUseCase struct {
	generator   service.PasswordGenerator
	configRepo  repository.ConfigRepository
}

func NewGenerateUseCase(g service.PasswordGenerator, cr repository.ConfigRepository) *GenerateUseCase {
	return &GenerateUseCase{
		generator:  g,
		configRepo: cr,
	}
}

func (uc *GenerateUseCase) Execute(input dto.GenerateInput) (*dto.GenerateOutput, error) {
	config, err := uc.configRepo.Load()
	if err != nil {
		config = entity.DefaultConfig()
	}

	secure, err := uc.generator.Generate(input.Password, input.MasterPass, config)
	if err != nil {
		return nil, err
	}

	return &dto.GenerateOutput{
		SecurePassword: secure,
		Length:         len(secure),
	}, nil
}
