package config

import "github.com/Jrp0h/backuper/cryption"

type cryptionConfig struct {
	Enable bool
	Crypter cryption.Crypter
}

func loadCryption(config *configFile) (cryptionConfig, error) {
	enc, err := cryption.LoadFromConfig(config.Encryption)

	if err != nil {
		return cryptionConfig{
			Enable: false,
			Crypter: nil,
		}, err
	}

	return cryptionConfig{
		Enable: enc != nil,
		Crypter: enc,
	}, nil
}