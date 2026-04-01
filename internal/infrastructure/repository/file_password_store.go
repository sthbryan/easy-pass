package repository

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type FilePasswordStore struct {
	filePath string
}

func NewFilePasswordStore() *FilePasswordStore {
	configDir := os.ExpandEnv(ConfigDir)
	filePath := filepath.Join(configDir, "passwords.enc")
	return &FilePasswordStore{filePath: filePath}
}

func (s *FilePasswordStore) Save(name, secure, masterPass string) error {
	if err := os.MkdirAll(os.ExpandEnv(ConfigDir), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	passwords, err := s.loadAll(masterPass)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	passwords[name] = secure
	return s.saveAll(passwords, masterPass)
}

func (s *FilePasswordStore) Get(name, masterPass string) (string, error) {
	passwords, err := s.loadAll(masterPass)
	if err != nil {
		return "", err
	}

	secure, ok := passwords[name]
	if !ok {
		return "", fmt.Errorf("password '%s' not found", name)
	}

	return secure, nil
}

func (s *FilePasswordStore) List() ([]string, error) {
	configDir := os.ExpandEnv(ConfigDir)
	metaPath := filepath.Join(configDir, "passwords.meta")

	if _, err := os.Stat(metaPath); os.IsNotExist(err) {
		return []string{}, nil
	}

	data, err := os.ReadFile(metaPath)
	if err != nil {
		return nil, err
	}

	var meta struct {
		Names []string `json:"names"`
	}
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, err
	}

	return meta.Names, nil
}

func (s *FilePasswordStore) Delete(name string) error {
	return nil
}

func (s *FilePasswordStore) loadAll(masterPass string) (map[string]string, error) {
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		return make(map[string]string), nil
	}

	encrypted, err := os.ReadFile(s.filePath)
	if err != nil {
		return nil, err
	}

	if len(encrypted) == 0 {
		return make(map[string]string), nil
	}

	plaintext, err := s.decrypt(encrypted, masterPass)
	if err != nil {
		return nil, fmt.Errorf("invalid master password")
	}

	var passwords map[string]string
	if err := json.Unmarshal(plaintext, &passwords); err != nil {
		return nil, err
	}

	return passwords, nil
}

func (s *FilePasswordStore) saveAll(passwords map[string]string, masterPass string) error {
	data, err := json.Marshal(passwords)
	if err != nil {
		return err
	}

	encrypted, err := s.encrypt(data, masterPass)
	if err != nil {
		return err
	}

	if err := os.WriteFile(s.filePath, encrypted, 0600); err != nil {
		return err
	}

	metaPath := filepath.Join(os.ExpandEnv(ConfigDir), "passwords.meta")
	names := make([]string, 0, len(passwords))
	for n := range passwords {
		names = append(names, n)
	}
	metaData, err := json.Marshal(map[string][]string{"names": names})
	if err != nil {
		return err
	}
	if err := os.WriteFile(metaPath, metaData, 0644); err != nil {
		return err
	}

	return nil
}

func (s *FilePasswordStore) encrypt(plaintext []byte, password string) ([]byte, error) {
	key := sha256.Sum256([]byte(password))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	return append(nonce, gcm.Seal(nil, nonce, plaintext, nil)...), nil
}

func (s *FilePasswordStore) decrypt(ciphertext []byte, password string) ([]byte, error) {
	key := sha256.Sum256([]byte(password))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
