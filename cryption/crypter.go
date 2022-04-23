package cryption

import (
	"fmt"
	"strings"

	"github.com/Jrp0h/backpack/utils"
)

type Crypter interface {
	Encrypt(file string) error
	Decrypt(file string) error
} 

func LoadFromConfig(data map[string]string) (Crypter, error) {
	if data == nil {
		utils.Log.Debug("cryption/crypter: Data is nil")
		return nil, nil
	}

	t, exists := data["type"]

	if !exists {
		return nil, fmt.Errorf("cryption/crypter: Missing required field 'type'")
	}

	switch strings.ToLower(t) {
	case "aes":
		return loadAES(&data)
	default:
		return nil, fmt.Errorf("cryption/crypter: Unknown type '%s'", t)
	}
}