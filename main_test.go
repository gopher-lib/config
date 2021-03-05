package config

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

const simpleConfigStr = `
port: 8080
secret1: ${SOME_SECRET}
secret2: ${SECRET}
dollar: $$money
`
const simpleConfigEnvStr = `
SOME_SECRET=secret-value
`

const complexConfigStr = `
db:
  user: root
  password: ${DB_PASSWORD}
connectionString: "${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}?parseTime=true"
`
const complexConfigEnvStr = `
DB_USER=root
DB_PASSWORD=admin
DB_HOST=localhost
DB_PORT=3306
DB_NAME=core
`

type c struct {
	config, env string
}

type filenames struct {
	simple, complex c
}

var f *filenames

func TestMain(m *testing.M) {
	// Setup.
	simpleConfig, err := ioutil.TempFile("", "config.simple.*.yaml")
	simpleEnv, err := ioutil.TempFile("", ".env.simple")
	complexConfig, err := ioutil.TempFile("", "config.complex.*.yaml")
	complexEnv, err := ioutil.TempFile("", ".env.complex")
	if err != nil {
		log.Fatalf("failed to create temp files: %v", err)
	}
	_, err = simpleConfig.WriteString(simpleConfigStr)
	_, err = simpleEnv.WriteString(simpleConfigEnvStr)
	_, err = complexConfig.WriteString(complexConfigStr)
	_, err = complexEnv.WriteString(complexConfigEnvStr)
	if err != nil {
		log.Fatalf("failed to write temp files: %v", err)
	}

	// Initialize global variable.
	f = &filenames{c{simpleConfig.Name(), simpleEnv.Name()}, c{complexConfig.Name(), complexEnv.Name()}}

	// Run all tests.
	code := m.Run()

	// Tear down.
	os.Remove(simpleConfig.Name())
	os.Remove(simpleEnv.Name())
	os.Remove(complexConfig.Name())
	os.Remove(complexEnv.Name())

	os.Exit(code)
}
