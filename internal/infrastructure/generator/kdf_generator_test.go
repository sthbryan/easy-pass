package generator

import (
	"testing"

	"github.com/sthbryan/easypass/internal/domain/entity"
)

func TestKDFGenerator_Generate(t *testing.T) {
	g := NewKDFGenerator()
	config := entity.DefaultConfig()

	result, err := g.Generate("password", "master123", config)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if len(result) != config.Length {
		t.Errorf("Generate() length = %v, want %v", len(result), config.Length)
	}

	charset := g.buildCharset(config)
	for _, c := range result {
		found := false
		for _, ch := range charset {
			if c == ch {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Generate() contains invalid char: %c", c)
		}
	}
}

func TestKDFGenerator_Deterministic(t *testing.T) {
	g := NewKDFGenerator()
	config := entity.DefaultConfig()

	r1, _ := g.Generate("test", "salt", config)
	r2, _ := g.Generate("test", "salt", config)

	if r1 != r2 {
		t.Error("Should be deterministic")
	}
}

func TestKDFGenerator_DifferentOutputs(t *testing.T) {
	g := NewKDFGenerator()
	config := entity.DefaultConfig()

	p1, _ := g.Generate("pass1", "master", config)
	p2, _ := g.Generate("pass2", "master", config)

	if p1 == p2 {
		t.Error("Different passwords should produce different outputs")
	}
}

func TestKDFGenerator_ExcludeSimilar(t *testing.T) {
	g := NewKDFGenerator()
	config := &entity.PasswordConfig{
		Length:         16,
		UseUppercase:   true,
		UseLowercase:   true,
		UseNumbers:     true,
		UseSymbols:     false,
		ExcludeSimilar: true,
	}

	result, err := g.Generate("test", "master", config)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	similar := "0O1lI"
	for _, c := range result {
		for _, s := range similar {
			if c == s {
				t.Errorf("ExcludeSimilar: found '%c' in result: %s", c, result)
			}
		}
	}
}
