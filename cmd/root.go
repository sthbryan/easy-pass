package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/sthbryan/easypass/internal/infrastructure/clipboard"
	"github.com/spf13/viper"
	"github.com/sthbryan/easypass/internal/application/dto"
	"github.com/sthbryan/easypass/internal/application/usecase"
	"github.com/sthbryan/easypass/internal/infrastructure/generator"
	"github.com/sthbryan/easypass/internal/infrastructure/repository"
)

var rootCmd = &cobra.Command{
	Use:   "ep [password] [master-password]",
	Short: "EasyPass - Generate secure passwords from easy phrases",
	Long: `EasyPass generates secure passwords from easy-to-remember phrases.
Uses Argon2id to derive secure keys.

Example:
  ep "mypassword" "myMasterSecret"`,
	Args: cobra.RangeArgs(0, 2),
	RunE: runRoot,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose mode")
	rootCmd.PersistentFlags().BoolP("copy", "c", false, "Copy to clipboard")
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/easypass")
	viper.AddConfigPath(".")

	viper.SetDefault("password.length", 16)
	viper.SetDefault("password.use_uppercase", true)
	viper.SetDefault("password.use_lowercase", true)
	viper.SetDefault("password.use_numbers", true)
	viper.SetDefault("password.use_symbols", true)
	viper.SetDefault("password.exclude_similar", false)
	viper.SetDefault("password.custom_symbols", "!@#$%^&*()_+-=[]{}|;:,.<>?")
	viper.SetDefault("password.min_symbols", 0)
	viper.SetDefault("password.min_numbers", 0)
	viper.SetDefault("password.min_uppercase", 0)
	viper.SetDefault("password.algorithm", "argon2id")
	viper.SetDefault("password.iterations", 1)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintf(os.Stderr, "Error reading config: %v\n", err)
		}
	}
}

func runRoot(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help()
	}

	if len(args) == 1 && isSubcommand(args[0]) {
		return nil
	}

	if len(args) < 2 {
		return fmt.Errorf("requires password and master-password")
	}

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

	clip := clipboard.NewSystemClipboard()

	if copyFlag, _ := cmd.Flags().GetBool("copy"); copyFlag {
		if err := clip.Copy(result.SecurePassword); err != nil {
			return fmt.Errorf("failed to copy: %w", err)
		}
		fmt.Println("✓ Copied to clipboard")
	} else {
		fmt.Println(result.SecurePassword)
	}

	return nil
}

func isSubcommand(cmd string) bool {
	subcommands := map[string]bool{
		"config": true,
		"save":   true,
		"show":   true,
		"list":   true,
		"help":   true,
	}
	return subcommands[cmd]
}
