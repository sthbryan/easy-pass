package entity

import "errors"

type PasswordConfig struct {
	Length        int  `json:"length"`
	UseUppercase  bool `json:"use_uppercase"`
	UseLowercase  bool `json:"use_lowercase"`
	UseNumbers    bool `json:"use_numbers"`
	UseSymbols    bool `json:"use_symbols"`

	ExcludeSimilar bool   `json:"exclude_similar"`
	CustomSymbols  string `json:"custom_symbols"`
	MinSymbols     int    `json:"min_symbols"`
	MinNumbers     int    `json:"min_numbers"`
	MinUppercase   int    `json:"min_uppercase"`

	Algorithm  string `json:"algorithm"`
	Iterations int    `json:"iterations"`
}

func DefaultConfig() *PasswordConfig {
	return &PasswordConfig{
		Length:        16,
		UseUppercase:  true,
		UseLowercase:  true,
		UseNumbers:    true,
		UseSymbols:    true,
		ExcludeSimilar: false,
		CustomSymbols: "!@#$%^&*()_+-=[]{}|;:,.<>?",
		MinSymbols:    0,
		MinNumbers:    0,
		MinUppercase:  0,
		Algorithm:     "argon2id",
		Iterations:    1,
	}
}

func (c *PasswordConfig) Validate() error {
	if c.Length < 8 || c.Length > 128 {
		return errors.New("length must be between 8 and 128")
	}
	if !c.UseUppercase && !c.UseLowercase && !c.UseNumbers && !c.UseSymbols {
		return errors.New("must use at least one character type")
	}
	if c.MinSymbols > c.Length {
		return errors.New("min symbols cannot be greater than length")
	}
	return nil
}
