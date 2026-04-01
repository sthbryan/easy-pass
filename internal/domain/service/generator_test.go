package service

import (
	"testing"

	"github.com/sthbryan/easypass/internal/domain/entity"
)

func TestGenerator_Generate(t *testing.T) {
	g := NewGenerator()
	config := entity.DefaultConfig()

	tests := []struct {
		name       string
		password   string
		masterPass string
		wantErr    bool
	}{
		{
			name:       "empty password",
			password:   "",
			masterPass: "master123",
			wantErr:    true,
		},
		{
			name:       "empty master",
			password:   "password",
			masterPass: "",
			wantErr:    true,
		},
		{
			name:       "valid values",
			password:   "limonada",
			masterPass: "miMaster",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := g.Generate(tt.password, tt.masterPass, config)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) != config.Length {
				t.Errorf("Generate() length = %v, want %v", len(got), config.Length)
			}
		})
	}
}

func TestGenerator_SameInputSameOutput(t *testing.T) {
	g := NewGenerator()
	config := entity.DefaultConfig()

	result1, _ := g.Generate("limonada", "master123", config)
	result2, _ := g.Generate("limonada", "master123", config)

	if result1 != result2 {
		t.Errorf("Same inputs should produce same outputs: %s != %s", result1, result2)
	}
}

func TestGenerator_DifferentInputDifferentOutput(t *testing.T) {
	g := NewGenerator()
	config := entity.DefaultConfig()

	result1, _ := g.Generate("limonada", "master123", config)
	result2, _ := g.Generate("limonada", "master456", config)

	if result1 == result2 {
		t.Error("Different inputs should produce different outputs")
	}
}
