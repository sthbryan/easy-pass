package main

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/sthbryan/easypass/internal/domain/entity"
	"github.com/sthbryan/easypass/internal/infrastructure/repository"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "View or modify configuration",
}

var showConfigCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	RunE:  runShowConfig,
}

var setLengthCmd = &cobra.Command{
	Use:   "length [n]",
	Short: "Set password length (8-128)",
	Args:  cobra.ExactArgs(1),
	RunE:  runSetLength,
}

var setUpperCmd = &cobra.Command{
	Use:   "uppercase [true|false]",
	Short: "Enable/disable uppercase letters",
	Args:  cobra.ExactArgs(1),
	RunE:  runSetUpper,
}

var setLowerCmd = &cobra.Command{
	Use:   "lowercase [true|false]",
	Short: "Enable/disable lowercase letters",
	Args:  cobra.ExactArgs(1),
	RunE:  runSetLower,
}

var setNumbersCmd = &cobra.Command{
	Use:   "numbers [true|false]",
	Short: "Enable/disable numbers",
	Args:  cobra.ExactArgs(1),
	RunE:  runSetNumbers,
}

var setSymbolsCmd = &cobra.Command{
	Use:   "symbols [true|false]",
	Short: "Enable/disable symbols",
	Args:  cobra.ExactArgs(1),
	RunE:  runSetSymbols,
}

var setExcludeSimilarCmd = &cobra.Command{
	Use:   "exclude-similar [true|false]",
	Short: "Exclude similar characters (0, O, l, 1, I)",
	Args:  cobra.ExactArgs(1),
	RunE:  runSetExcludeSimilar,
}

var setCustomSymbolsCmd = &cobra.Command{
	Use:   "custom-symbols [symbols]",
	Short: "Set custom symbols (e.g., \"!@#$%\")",
	Args:  cobra.ExactArgs(1),
	RunE:  runSetCustomSymbols,
}

var setMinSymbolsCmd = &cobra.Command{
	Use:   "min-symbols [n]",
	Short: "Minimum symbols required",
	Args:  cobra.ExactArgs(1),
	RunE:  runSetMinSymbols,
}

var setMinNumbersCmd = &cobra.Command{
	Use:   "min-numbers [n]",
	Short: "Minimum numbers required",
	Args:  cobra.ExactArgs(1),
	RunE:  runSetMinNumbers,
}

var setAlgorithmCmd = &cobra.Command{
	Use:   "algorithm [argon2id|pbkdf2|scrypt]",
	Short: "Set derivation algorithm",
	Args:  cobra.ExactArgs(1),
	RunE:  runSetAlgorithm,
}

func runShowConfig(cmd *cobra.Command, args []string) error {
	repo := repository.NewFileConfigRepository()
	cfg, err := repo.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Println("Current configuration:")
	fmt.Printf("  Length:           %d\n", cfg.Length)
	fmt.Printf("  Uppercase:        %v\n", cfg.UseUppercase)
	fmt.Printf("  Lowercase:        %v\n", cfg.UseLowercase)
	fmt.Printf("  Numbers:          %v\n", cfg.UseNumbers)
	fmt.Printf("  Symbols:          %v\n", cfg.UseSymbols)
	fmt.Printf("  Exclude similar:  %v\n", cfg.ExcludeSimilar)
	fmt.Printf("  Custom symbols:   %s\n", cfg.CustomSymbols)
	fmt.Printf("  Min symbols:      %d\n", cfg.MinSymbols)
	fmt.Printf("  Min numbers:      %d\n", cfg.MinNumbers)
	fmt.Printf("  Algorithm:        %s\n", cfg.Algorithm)

	return nil
}

func runSetLength(cmd *cobra.Command, args []string) error {
	n, err := strconv.Atoi(args[0])
	if err != nil || n < 8 || n > 128 {
		return fmt.Errorf("invalid length: %s (use 8-128)", args[0])
	}

	repo := repository.NewFileConfigRepository()
	cfg, err := repo.Load()
	if err != nil {
		cfg = entity.DefaultConfig()
	}
	cfg.Length = n

	if err := repo.Save(cfg); err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}

	fmt.Printf("Length set to %d\n", n)
	return nil
}

func runSetUpper(cmd *cobra.Command, args []string) error {
	return setBoolFlag("uppercase", args[0])
}

func runSetLower(cmd *cobra.Command, args []string) error {
	return setBoolFlag("lowercase", args[0])
}

func runSetNumbers(cmd *cobra.Command, args []string) error {
	return setBoolFlag("numbers", args[0])
}

func runSetSymbols(cmd *cobra.Command, args []string) error {
	return setBoolFlag("symbols", args[0])
}

func runSetExcludeSimilar(cmd *cobra.Command, args []string) error {
	return setBoolFlag("exclude_similar", args[0])
}

func runSetCustomSymbols(cmd *cobra.Command, args []string) error {
	repo := repository.NewFileConfigRepository()
	cfg, err := repo.Load()
	if err != nil {
		cfg = entity.DefaultConfig()
	}
	cfg.CustomSymbols = args[0]

	if err := repo.Save(cfg); err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}

	fmt.Printf("Custom symbols: %s\n", args[0])
	return nil
}

func runSetMinSymbols(cmd *cobra.Command, args []string) error {
	n, err := strconv.Atoi(args[0])
	if err != nil || n < 0 {
		return fmt.Errorf("invalid value: %s", args[0])
	}

	repo := repository.NewFileConfigRepository()
	cfg, err := repo.Load()
	if err != nil {
		cfg = entity.DefaultConfig()
	}
	cfg.MinSymbols = n

	if err := repo.Save(cfg); err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}

	fmt.Printf("Min symbols: %d\n", n)
	return nil
}

func runSetMinNumbers(cmd *cobra.Command, args []string) error {
	n, err := strconv.Atoi(args[0])
	if err != nil || n < 0 {
		return fmt.Errorf("invalid value: %s", args[0])
	}

	repo := repository.NewFileConfigRepository()
	cfg, err := repo.Load()
	if err != nil {
		cfg = entity.DefaultConfig()
	}
	cfg.MinNumbers = n

	if err := repo.Save(cfg); err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}

	fmt.Printf("Min numbers: %d\n", n)
	return nil
}

func runSetAlgorithm(cmd *cobra.Command, args []string) error {
	valid := map[string]bool{"argon2id": true, "pbkdf2": true, "scrypt": true}
	if !valid[args[0]] {
		return fmt.Errorf("invalid algorithm: %s (use argon2id, pbkdf2, or scrypt)", args[0])
	}

	repo := repository.NewFileConfigRepository()
	cfg, err := repo.Load()
	if err != nil {
		cfg = entity.DefaultConfig()
	}
	cfg.Algorithm = args[0]

	if err := repo.Save(cfg); err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}

	fmt.Printf("Algorithm: %s\n", args[0])
	return nil
}

func setBoolFlag(field, value string) error {
	v, err := strconv.ParseBool(value)
	if err != nil {
		return fmt.Errorf("invalid value: %s (use true or false)", value)
	}

	repo := repository.NewFileConfigRepository()
	cfg, err := repo.Load()
	if err != nil {
		cfg = entity.DefaultConfig()
	}

	switch field {
	case "uppercase":
		cfg.UseUppercase = v
	case "lowercase":
		cfg.UseLowercase = v
	case "numbers":
		cfg.UseNumbers = v
	case "symbols":
		cfg.UseSymbols = v
	case "exclude_similar":
		cfg.ExcludeSimilar = v
	}

	if err := repo.Save(cfg); err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}

	fmt.Printf("%s set to %v\n", field, v)
	return nil
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(showConfigCmd)
	configCmd.AddCommand(setLengthCmd)
	configCmd.AddCommand(setUpperCmd)
	configCmd.AddCommand(setLowerCmd)
	configCmd.AddCommand(setNumbersCmd)
	configCmd.AddCommand(setSymbolsCmd)
	configCmd.AddCommand(setExcludeSimilarCmd)
	configCmd.AddCommand(setCustomSymbolsCmd)
	configCmd.AddCommand(setMinSymbolsCmd)
	configCmd.AddCommand(setMinNumbersCmd)
	configCmd.AddCommand(setAlgorithmCmd)
}
