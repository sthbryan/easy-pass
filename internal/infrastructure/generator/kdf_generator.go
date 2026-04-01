package generator

import (
	"crypto/sha256"
	"strings"

	"github.com/sthbryan/easypass/internal/domain/entity"
	"golang.org/x/crypto/argon2"
)

type KDFGenerator struct{}

func NewKDFGenerator() *KDFGenerator {
	return &KDFGenerator{}
}

func (g *KDFGenerator) Generate(password, masterPass string, config *entity.PasswordConfig) (string, error) {
	charset := g.buildCharset(config)

	saltSource := masterPass + password + g.configString(config)
	salt := sha256.Sum256([]byte(saltSource))

	derived := argon2.IDKey(
		[]byte(password),
		salt[:],
		1,
		64*1024,
		2,
		uint32(config.Length),
	)

	result := make([]byte, config.Length)
	for i, b := range derived {
		result[i] = charset[int(b)%len(charset)]
	}

	result = g.ensureCharset(result, config)

	return string(result), nil
}

func (g *KDFGenerator) buildCharset(config *entity.PasswordConfig) string {
	var charset strings.Builder

	if config.UseUppercase {
		upper := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		if config.ExcludeSimilar {
			upper = "ABCDEFGHJKMNPQRSTUVWXYZ"
		}
		charset.WriteString(upper)
	}
	if config.UseLowercase {
		lower := "abcdefghijklmnopqrstuvwxyz"
		if config.ExcludeSimilar {
			lower = "abcdefghjkmnpqrstuvwxyz"
		}
		charset.WriteString(lower)
	}
	if config.UseNumbers {
		nums := "0123456789"
		if config.ExcludeSimilar {
			nums = "23456789"
		}
		charset.WriteString(nums)
	}
	if config.UseSymbols {

		if config.CustomSymbols != "" {
			charset.WriteString(config.CustomSymbols)
		} else {
			charset.WriteString("!@#$%^&*()_+-=[]{}|;:,.<>?")
		}
	}

	return charset.String()
}

func (g *KDFGenerator) configString(config *entity.PasswordConfig) string {
	result := ""
	if config.UseUppercase {
		result += "U"
	}
	if config.UseLowercase {
		result += "L"
	}
	if config.UseNumbers {
		result += "N"
	}
	if config.UseSymbols {
		result += "S"
	}
	return result + string(rune(config.Length))
}

func (g *KDFGenerator) ensureCharset(result []byte, config *entity.PasswordConfig) []byte {

	upperChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowerChars := "abcdefghijklmnopqrstuvwxyz"
	numberChars := "0123456789"
	symbolChars := config.CustomSymbols
	if symbolChars == "" {
		symbolChars = "!@#$%^&*()_+-=[]{}|;:,.<>?"
	}

	if config.ExcludeSimilar {
		upperChars = "ABCDEFGHJKMNPQRSTUVWXYZ"
		lowerChars = "abcdefghjkmnpqrstuvwxyz"
		numberChars = "23456789"
	}

	pos := 0

	if config.UseUppercase && pos < len(result) {
		result[pos] = upperChars[int(result[pos])%len(upperChars)]
		pos++
	}

	if config.UseLowercase && pos < len(result) {
		result[pos] = lowerChars[int(result[pos+1])%len(lowerChars)]
		pos++
	}

	if config.UseNumbers && pos < len(result) {
		result[pos] = numberChars[int(result[pos+2])%len(numberChars)]
		pos++
	}

	if config.UseSymbols && pos < len(result) {
		result[pos] = symbolChars[int(result[pos+3])%len(symbolChars)]
	}

	return result
}
