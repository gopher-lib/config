package config

import (
	"reflect"
	"strings"
	"testing"
)

func TestLoadFile(t *testing.T) {
	type Config struct {
		Port    int
		Secret1 string
		Secret2 string
		Dollar  string
	}
	var conf Config
	err := LoadFile(&conf, "./testdata/config.testing.yaml", "./testdata/.env.testing")
	if err != nil {
		t.Fatal(err)
	}
	expected := Config{8080, "secret-value", "", "$money"}
	if !reflect.DeepEqual(conf, expected) {
		t.Errorf("not equal: %v != %v", conf, expected)
	}
}

func TestLoad(t *testing.T) {
	t.Run("more complex config with interpolation", func(t *testing.T) {
		const confStr = `
		db:
  	user: root
  	password: ${DB_PASSWORD}
		connectionString: "${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}?parseTime=true"
		`
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
		err := Load(strings.NewReader(confStr), &conf)
		if err != nil {
			t.Fatalf("error loading config: %v\n", err)
		}
		expected := Config{8080, DB{"root", "admin"}, "root:admin@tcp(localhost:3306)/core?parseTime=true"}
		if !reflect.DeepEqual(conf, expected) {
			t.Errorf("not equal: %v != %v", conf, expected)
		}
	})
}
