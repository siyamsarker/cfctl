package config

import (
	"fmt"

	"github.com/zalando/go-keyring"
)

const (
	keyringService = "cfctl"
)

// StoreCredential stores a credential in the system keyring
func StoreCredential(accountName, credential string) error {
	err := keyring.Set(keyringService, accountName, credential)
	if err != nil {
		return fmt.Errorf("store credential: %w", err)
	}
	return nil
}

// GetCredential retrieves a credential from the system keyring
func GetCredential(accountName string) (string, error) {
	credential, err := keyring.Get(keyringService, accountName)
	if err != nil {
		return "", fmt.Errorf("get credential: %w", err)
	}
	return credential, nil
}

// DeleteCredential deletes a credential from the system keyring
func DeleteCredential(accountName string) error {
	err := keyring.Delete(keyringService, accountName)
	if err != nil {
		return fmt.Errorf("delete credential: %w", err)
	}
	return nil
}
