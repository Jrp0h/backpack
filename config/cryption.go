package config

import "github.com/Jrp0h/backpack/cryption"

type cryptionConfig struct {
	Enabled bool
	Crypter cryption.Crypter
}

func loadCryption(config *configFile) (cryptionConfig, error) {
	enc, err := cryption.LoadFromConfig(config.Encryption)

	if err != nil {
		return cryptionConfig{
			Enabled: false,
			Crypter: nil,
		}, err
	}

	return cryptionConfig{
		Enabled: enc != nil,
		Crypter: enc,
	}, nil
}
