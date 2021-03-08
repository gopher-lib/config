package config

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestLoadFile(t *testing.T) {
	viper.Reset()

	type db struct {
		User     string
		Password string
	}
	type config struct {
		Port             int
		AuthSecret       string
		Secret           string
		Dollar           string
		DB               db
		ConnectionString string
	}
	var conf config
	err := LoadFile(&conf, "./testdata/config.testing.yaml", "./testdata/.env.testing")
	if err != nil {
		t.Fatal(err)
	}
	expected := config{8080, "secret", "", "$dollar", db{"root", "admin"}, "root:admin@tcp(localhost:3306)/core?parseTime=true"}
	if !reflect.DeepEqual(conf, expected) {
		t.Errorf("not equal: %v != %v", conf, expected)
	}
}

func TestLoad(t *testing.T) {
	t.Run("complex config", func(t *testing.T) {
		viper.Reset()

		os.Setenv("DB_PASSWORD", "root")
		const confStr = `
port: 1234
db:
  user: postgres
  password: ${DB_PASSWORD}
`
		type DB struct {
			User     string
			Password string
		}
		type Config struct {
			Port int
			DB   DB
		}
		var conf Config
		err := Load(strings.NewReader(confStr), &conf, "yaml")
		if err != nil {
			t.Fatal(err)
		}
		expected := Config{1234, DB{"postgres", "root"}}
		if !reflect.DeepEqual(conf, expected) {
			t.Errorf("not equal: %v != %v", conf, expected)
		}
	})
}
