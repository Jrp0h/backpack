package config

import (
	"github.com/Jrp0h/backpack/crypto"
)

type cryptoConfig struct {
	Enabled bool
	Crypter crypto.Crypto
}

func loadCrypto(config *configFile) (cryptoConfig, error) {
	enc, err := crypto.LoadFromConfig(config.Encryption)

	if err != nil {
		return cryptoConfig{
			Enabled: false,
			Crypter: nil,
		}, err
	}

	return cryptoConfig{
		Enabled: enc != nil,
		Crypter: enc,
	}, nil
}
