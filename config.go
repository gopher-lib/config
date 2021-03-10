package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func LoadFile(rawVal interface{}, filename string, envPath ...string) error {
	if len(envPath) > 0 && envPath[0] != "" {
		if err := godotenv.Load(envPath[0]); err != nil {
			return fmt.Errorf("config: failed to load env. file: %v", err)
		}
	} else {
		// Ignore error as we are loading default env. file.
		_ = godotenv.Load(".env")
	}
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("config: failed to open config file: %v", err)
	}
	return Load(f, rawVal, strings.TrimPrefix(filepath.Ext(filename), "."))
}

func Load(in io.Reader, rawVal interface{}, configType string) error {
	viper.SetConfigType(configType)
	err := viper.ReadConfig(in)
	if err != nil {
		return fmt.Errorf("config: failed to read in config: %w", err)
	}

	for _, key := range viper.AllKeys() {
		newKey := os.Expand(viper.GetString(key), mapping)
		viper.Set(key, newKey)
	}

	err = viper.Unmarshal(&rawVal)
	if err != nil {
		return fmt.Errorf("config: failed to unmarshal config: %w", err)
	}
	return nil
}

// mapping is second argument for os.Expand function.
func mapping(s string) string {
	if strings.HasPrefix(s, "$") {
		return s
	}
	return os.Getenv(s)
}
