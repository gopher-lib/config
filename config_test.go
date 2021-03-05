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
		err := Load(&conf, f.simple.config, f.simple.env)
		if err != nil {
			t.Fatalf("error loading config: %v\n", err)
		}
		expected := Config{8080, "secret-value", "", "$money"}
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
		err := Load(&conf, f.complex.config, f.complex.env)
		if err != nil {
			t.Fatalf("error loading config: %v\n", err)
		}
		expected := Config{8080, DB{"root", "admin"}, "root:admin@tcp(localhost:3306)/core?parseTime=true"}
		if !reflect.DeepEqual(conf, expected) {
			t.Errorf("not equal: %v != %v", conf, expected)
		}
	})
}
