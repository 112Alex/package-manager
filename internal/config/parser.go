package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var (
	ErrUnsupportedExt = errors.New("unsupported config file extension (use .json or .yaml/.yml)")
)

func loadFileBytes(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// LoadPacket reads packet config from json or yaml.
func LoadPacket(path string) (*PacketConfig, error) {
	b, err := loadFileBytes(path)
	if err != nil {
		return nil, err
	}
	var pc PacketConfig
	switch ext := filepath.Ext(path); ext {
	case ".json":
		err = json.Unmarshal(b, &pc)
	case ".yaml", ".yml":
		err = yaml.Unmarshal(b, &pc)
	default:
		return nil, ErrUnsupportedExt
	}
	if err != nil {
		return nil, err
	}
	return &pc, nil
}

// LoadPackages reads packages config.
func LoadPackages(path string) (*PackagesConfig, error) {
	b, err := loadFileBytes(path)
	if err != nil {
		return nil, err
	}
	var pc PackagesConfig
	switch ext := filepath.Ext(path); ext {
	case ".json":
		err = json.Unmarshal(b, &pc)
	case ".yaml", ".yml":
		err = yaml.Unmarshal(b, &pc)
	default:
		return nil, ErrUnsupportedExt
	}
	if err != nil {
		return nil, err
	}
	return &pc, nil
}
