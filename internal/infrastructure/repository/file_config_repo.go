package repository

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/sthbryan/easypass/internal/domain/entity"
)

type FileConfigRepository struct {
	viper *viper.Viper
}

func NewFileConfigRepository() *FileConfigRepository {
	v := viper.New()

	configDir := os.ExpandEnv("$HOME/.config/easypass")
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(configDir)
	v.AddConfigPath(".")

	v.SetDefault("password.length", 16)
	v.SetDefault("password.use_uppercase", true)
	v.SetDefault("password.use_lowercase", true)
	v.SetDefault("password.use_numbers", true)
	v.SetDefault("password.use_symbols", true)
	v.SetDefault("password.exclude_similar", false)
	v.SetDefault("password.custom_symbols", "!@#$%^&*()_+-=[]{}|;:,.<>?")
	v.SetDefault("password.min_symbols", 0)
	v.SetDefault("password.min_numbers", 0)
	v.SetDefault("password.min_uppercase", 0)
	v.SetDefault("password.algorithm", "argon2id")
	v.SetDefault("password.iterations", 1)

	return &FileConfigRepository{viper: v}
}

func (r *FileConfigRepository) Save(config *entity.PasswordConfig) error {
	r.viper.Set("password.length", config.Length)
	r.viper.Set("password.use_uppercase", config.UseUppercase)
	r.viper.Set("password.use_lowercase", config.UseLowercase)
	r.viper.Set("password.use_numbers", config.UseNumbers)
	r.viper.Set("password.use_symbols", config.UseSymbols)
	r.viper.Set("password.exclude_similar", config.ExcludeSimilar)
	r.viper.Set("password.custom_symbols", config.CustomSymbols)
	r.viper.Set("password.min_symbols", config.MinSymbols)
	r.viper.Set("password.min_numbers", config.MinNumbers)
	r.viper.Set("password.min_uppercase", config.MinUppercase)
	r.viper.Set("password.algorithm", config.Algorithm)
	r.viper.Set("password.iterations", config.Iterations)

	configDir := filepath.Dir(r.viper.ConfigFileUsed())
	if configDir != "" {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	configPath := filepath.Join(os.ExpandEnv("$HOME/.config/easypass"), "config.yaml")
	if err := r.viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

func (r *FileConfigRepository) Load() (*entity.PasswordConfig, error) {
	if err := r.viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return entity.DefaultConfig(), nil
		}
		return nil, err
	}

	config := &entity.PasswordConfig{
		Length:         r.viper.GetInt("password.length"),
		UseUppercase:   r.viper.GetBool("password.use_uppercase"),
		UseLowercase:   r.viper.GetBool("password.use_lowercase"),
		UseNumbers:     r.viper.GetBool("password.use_numbers"),
		UseSymbols:     r.viper.GetBool("password.use_symbols"),
		ExcludeSimilar: r.viper.GetBool("password.exclude_similar"),
		CustomSymbols:  r.viper.GetString("password.custom_symbols"),
		MinSymbols:     r.viper.GetInt("password.min_symbols"),
		MinNumbers:     r.viper.GetInt("password.min_numbers"),
		MinUppercase:   r.viper.GetInt("password.min_uppercase"),
		Algorithm:      r.viper.GetString("password.algorithm"),
		Iterations:     r.viper.GetInt("password.iterations"),
	}

	return config, nil
}
