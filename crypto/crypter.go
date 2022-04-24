package crypto

import (
	"fmt"
	"strings"

	"github.com/Jrp0h/backpack/utils"
)

type Crypto interface {
	Encrypt(file string) error
	Decrypt(file string) error
}

func LoadFromConfig(data map[string]string) (Crypto, error) {
	if data == nil {
		utils.Log.Debug("crypto/crypto: Config data is nil")
		return nil, nil
	}

	t, exists := data["type"]

	if !exists {
		return nil, fmt.Errorf("crypto/crypto: Missing required field 'type'")
	}

	switch strings.ToLower(t) {
	case "aes":
		return loadAES(&data)
	default:
		return nil, fmt.Errorf("crypto/crypto: Unknown type '%s'", t)
	}
}
