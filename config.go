package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
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
		// Ignore error because we are loading default env. file.
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

var pattern = regexp.MustCompile("(?i)[_a-z][_a-z0-9]*")

// mapping is second argument for os.Expand function.
func mapping(s string) string {
	loc := pattern.FindStringIndex(s)
	// if no match then silently ignore it.
	if loc == nil {
		return s
	}
	key := s[loc[0]:loc[1]]
	switch m := s[loc[1]:]; {
	case strings.HasPrefix(m, ":-"):
		val := os.Getenv(key)
		if val == "" {
			return m[2:]
		}
		return val
	case strings.HasPrefix(m, "-"):
		val, found := os.LookupEnv(key)
		if found {
			return val
		}
		return m[1:]
	case strings.HasPrefix(m, ":?"):
		val := os.Getenv(key)
		if val == "" {
			panic(m[2:])
		}
		return val
	case strings.HasPrefix(m, "?"):
		val, found := os.LookupEnv(key)
		if found {
			return val
		}
		panic(m[1:])
	}
	return os.Getenv(key)
}
