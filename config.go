package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Load loads configuration from provided file, interpolates
// environement variables and then unmarshals it into provided struct.
func Load(rawVal interface{}, filename string, envPath ...string) error {
	if len(envPath) > 0 && envPath[0] != "" {
		if err := godotenv.Load(envPath[0]); err != nil {
			return fmt.Errorf("failed to load env. file: %v", err)
		}
	} else {
		// Ignore error as we are loading default env. file.
		_ = godotenv.Load(".env")
	}

	viper.SetConfigName(strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filepath.Base(filename))))
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filepath.Dir(filename))
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		return fmt.Errorf("failed to read in config: %w", err)
	}

	for _, key := range viper.AllKeys() {
		viper.Set(key, os.Expand(viper.GetString(key), mapping))
	}

	err = viper.Unmarshal(&rawVal)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
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
