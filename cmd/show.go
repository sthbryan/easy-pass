package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/sthbryan/easypass/internal/application/dto"
	"github.com/sthbryan/easypass/internal/application/usecase"
	"github.com/sthbryan/easypass/internal/infrastructure/clipboard"
	"github.com/sthbryan/easypass/internal/infrastructure/repository"
)

var showCmd = &cobra.Command{
	Use:   "show [name] [master-password]",
	Short: "Show saved password and copy to clipboard",
	Args:  cobra.ExactArgs(2),
	RunE:  runShow,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List saved passwords",
	RunE:  runList,
}

func runShow(cmd *cobra.Command, args []string) error {
	name := args[0]
	masterPass := args[1]

	passwordStore := repository.NewFilePasswordStore()
	uc := usecase.NewShowUseCase(passwordStore)

	result, err := uc.Execute(dto.ShowInput{
		Name:       name,
		MasterPass: masterPass,
	})
	if err != nil {
		return fmt.Errorf("failed to show: %w", err)
	}

	fmt.Printf("Password: %s\n", result.Password)

	clip := clipboard.NewSystemClipboard()
	if err := clip.Copy(result.Password); err != nil {
		fmt.Fprintf(os.Stderr, "⚠ Failed to copy to clipboard: %v\n", err)
	} else {
		fmt.Println("✓ Copied to clipboard")
	}

	return nil
}

func runList(cmd *cobra.Command, args []string) error {
	passwordStore := repository.NewFilePasswordStore()
	names, err := passwordStore.List()
	if err != nil {
		return fmt.Errorf("failed to list: %w", err)
	}

	if len(names) == 0 {
		fmt.Println("No saved passwords")
		return nil
	}

	fmt.Println("Saved passwords:")
	for _, name := range names {
		fmt.Printf("  - %s\n", name)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(listCmd)
}
