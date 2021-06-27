# [config](https://github.com/golang-config/config) Eloquent configuration for Golang apps.

[![Go Reference](https://pkg.go.dev/badge/github.com/gopher-lib/config.svg)](https://pkg.go.dev/github.com/gopher-lib/config)

## Features:

- Substitutes `$VARIABLE` and `${VARIABLE}` with variables found in a shell environment.
- Syntaxes for setting up default values and specifying mandatory variables:
  - ${VARIABLE:-default} evaluates to default if VARIABLE is unset or empty in the environment.
  - ${VARIABLE-default} evaluates to default only if VARIABLE is unset in the environment.
  - ${VARIABLE:?err} panics with an error message containing err if VARIABLE is unset or empty in the environment.
  - ${VARIABLE?err} panics with an error message containing err if VARIABLE is unset in the environment.

Examples:
