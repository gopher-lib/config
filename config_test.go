package config

import (
	"reflect"
	"testing"
)

func TestLoad(t *testing.T) {
	t.Run("simple config with interpolation", func(t *testing.T) {
		type Config struct {
			Port    int
			Secret1 string
			Secret2 string
			Dollar  string
		}
		var conf Config
		err := Load(&conf, "testdata/config.testing.yaml", "testdata/.env.testing")
		if err != nil {
			t.Fatalf("error loading config: %v\n", err)
		}
		expected := Config{5432, "secret-value", "", "$money"}
		if !reflect.DeepEqual(conf, expected) {
			t.Errorf("not equal: %v != %v", conf, expected)
		}
	})

	t.Run("more complex config with interpolation", func(t *testing.T) {
		type DB struct {
			User     string
			Password string
		}
		type Config struct {
			Port             int
			DB               DB
			ConnectionString string
		}
		var conf Config
		err := Load(&conf, "testdata/config.testing.yaml", "testdata/.env.testing")
		if err != nil {
			t.Fatalf("error loading config: %v\n", err)
		}
		expected := Config{5432, DB{"admin", "root"}, "postgres:root@localhost:5432/core"}
		if !reflect.DeepEqual(conf, expected) {
			t.Errorf("not equal: %v != %v", conf, expected)
		}
	})
}
