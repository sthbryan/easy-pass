package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sthbryan/easypass/internal/application/dto"
	"github.com/sthbryan/easypass/internal/application/usecase"
	"github.com/sthbryan/easypass/internal/infrastructure/generator"
	"github.com/sthbryan/easypass/internal/infrastructure/repository"
)

var generateCmd = &cobra.Command{
	Use:   "generate [password] [master-password]",
	Short: "Generate a secure password",
	Args:  cobra.ExactArgs(2),
	RunE:  runGenerate,
}

func runGenerate(cmd *cobra.Command, args []string) error {
	password := args[0]
	masterPass := args[1]

	gen := generator.NewKDFGenerator()
	configRepo := repository.NewFileConfigRepository()
	uc := usecase.NewGenerateUseCase(gen, configRepo)

	result, err := uc.Execute(dto.GenerateInput{
		Password:   password,
		MasterPass: masterPass,
	})
	if err != nil {
		return fmt.Errorf("failed to generate: %w", err)
	}

	fmt.Println(result.SecurePassword)
	return nil
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().BoolP("copy", "c", false, "Copy to clipboard")
}
