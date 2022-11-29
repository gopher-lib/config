package config

import (
	"os"
	"strings"
	"testing"
)

func TestLoadFile(t *testing.T) {
	os.Setenv("VAR_1", "value_1")
	os.Setenv("VAR_3", "")
	defer os.Clearenv()

	var cfg struct {
		Key string
	}
	err := LoadFile(".github/testdata/config.yaml", &cfg)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Key != "value_1" {
		t.Errorf("cfg.Key = %v, want %v", cfg.Key, "value_1")
	}
}

func TestLoad(t *testing.T) {
	os.Setenv("VAR_1", "value_1")
	os.Setenv("CONFIG_OVERWRITE_key2", "")
	os.Setenv("CONFIG_OVERWRITE_s.innerkey", "3")
	defer os.Clearenv()

	var cfg struct {
		Key, Key2 string
		S         struct {
			InnerKey uint8
		}
	}
	in := strings.NewReader(
		`{"key": "${VAR_1}", "key2": "${VAR_1}", "s": {"innerKey": 2}}`,
	)
	err := Load(in, JSON, &cfg)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Key != "value_1" {
		t.Errorf("cfg.Key = %v, want %v", cfg.Key, "value_1")
	}
	if cfg.Key2 != "" {
		t.Errorf("cfg.Key2 = %v, want %v", cfg.Key2, "")
	}
	if cfg.S.InnerKey != 3 {
		t.Errorf("S.InnerKey = %v, want %v", cfg.S.InnerKey, 3)
	}
}

func Test_mapping(t *testing.T) {
	os.Setenv("VAR_1", "value_1")
	os.Setenv("VAR_3", "")
	defer os.Clearenv()
	tests := []struct {
		name   string
		argStr string
		want   string
	}{
		{"#1", "${VAR_1}", "value_1"},
		{"#2", "key:${VAR_1}", "key:value_1"},
		{"#3", "prefix_${VAR_1}_postfix", "prefix_value_1_postfix"},
		{"#4", "$$15/h", "$15/h"},
		{"#5", "${VAR_2-default}", "default"},
		{"#6", "${VAR_3-default}", ""},
		{"#7", "${VAR_3:-default}", "default"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := os.Expand(tt.argStr, mapping); got != tt.want {
				t.Errorf("os.Expand() = %v, want %v", got, tt.want)
			}
		})
	}
}
