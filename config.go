package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/viper"
)

const (
	YAML = "yaml"
	JSON = "json"
	TOML = "toml"
)

var supportedFormats = map[string]bool{
	YAML: true,
	JSON: true,
	TOML: true,
}

// LoadFile reads configuration data from filename
// and unmarshals it into v.
//
// Internally it calls Load.
func LoadFile(v interface{}, filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("config: failed to open config file: %w", err)
	}
	return Load(f, v, strings.TrimPrefix(filepath.Ext(filename), "."))
}

// Load reads configuration data encoded in the format specified by configType from in
// and unmarshals it into v.
//
// Load returns an error when data is in wrong or unsupported format,
// or when it failed to unmarshal data into v.
func Load(in io.Reader, v interface{}, configType string) error {
	if !supportedFormats[configType] {
		return fmt.Errorf("config: %s - unsupported configuration format", configType)
	}

	viper.SetConfigType(configType)
	err := viper.ReadConfig(in)
	if err != nil {
		return fmt.Errorf("config: failed to read in config: %w", err)
	}

	for _, key := range viper.AllKeys() {
		newKey := os.Expand(viper.GetString(key), mapping)
		viper.Set(key, newKey)
	}

	err = viper.Unmarshal(&v)
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
			message := m[2:]
			if message == "" {
				message = key + " is empty or not set"
			}
			panic(message)
		}
		return val
	case strings.HasPrefix(m, "?"):
		val, found := os.LookupEnv(key)
		if found {
			return val
		}
		message := m[2:]
		if message == "" {
			message = key + " is not set"
		}
		panic(message)
	}
	return os.Getenv(key)
}
