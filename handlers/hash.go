package handlers

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"io/ioutil"

	"github.com/Jrp0h/backpack/config"
	"github.com/Jrp0h/backpack/utils"
)

type HashStatus int

const (
	HASH_NO_CHANGE HashStatus = iota
	HASH_CONTINUE
)

type hashResult struct {
	cfg *config.Config

	newHash []byte

	Result HashStatus
}

func HandleHash(cfg *config.Config, dataPath string, skipCompare bool) (*hashResult, error) {
	// Compute new hash
	data, err := ioutil.ReadFile(dataPath)
	if err != nil {
		return nil, err
	}

	h := sha512.New()
	_, err = h.Write(data)
	if err != nil {
		return nil, err
	}
	newHash := h.Sum(nil)

	// Check Prev Hash
	var prevHash []byte

	if !skipCompare {
		prev, err := ioutil.ReadFile(cfg.Hash)
		if err != nil {
			return nil, err
		}

		prevHash, err = hex.DecodeString(string(prev))
		if err != nil {
			return nil, err
		}
	}

	if !skipCompare && bytes.Equal(newHash, prevHash) {
		utils.Log.Info("Data hasn't changed. Skipping backup.")
		return &hashResult{
			cfg,
			newHash,
			HASH_NO_CHANGE,
		}, nil
	}

	utils.Log.Debug("handlers/hash: Hash return")
	return &hashResult{
		cfg,
		newHash,
		HASH_CONTINUE,
	}, nil

}

func (handler *hashResult) StoreHash(onlyIf *bool) error {
	// Store new hash
	if *onlyIf && utils.PathIsFile(handler.cfg.Hash) {
		err := ioutil.WriteFile(handler.cfg.Hash, []byte(hex.EncodeToString(handler.newHash)), 0644)
		if err != nil {
			utils.Log.Error("Failed to store hash. %s", err.Error())
			return err
		}
	}

	return nil
}
