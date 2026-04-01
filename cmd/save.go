package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sthbryan/easypass/internal/application/dto"
	"github.com/sthbryan/easypass/internal/application/usecase"
	"github.com/sthbryan/easypass/internal/infrastructure/generator"
	"github.com/sthbryan/easypass/internal/infrastructure/repository"
)

var saveCmd = &cobra.Command{
	Use:   "save [name] [password] [master-password]",
	Short: "Save a derived password",
	Args:  cobra.ExactArgs(3),
	RunE:  runSave,
}

func runSave(cmd *cobra.Command, args []string) error {
	name := args[0]
	password := args[1]
	masterPass := args[2]

	gen := generator.NewKDFGenerator()
	configRepo := repository.NewFileConfigRepository()
	passwordStore := repository.NewFilePasswordStore()

	uc := usecase.NewSaveUseCase(gen, configRepo, passwordStore)

	result, err := uc.Execute(dto.SaveInput{
		Name:       name,
		Password:   password,
		MasterPass: masterPass,
	})
	if err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}

	fmt.Printf("✓ %s\n", result.Message)
	return nil
}

func init() {
	rootCmd.AddCommand(saveCmd)
}
