# [config](https://github.com/golang-config/config): Eloquent configuration for Golang apps.

## Features:

- Substitutes `$VARIABLE` and `${VARIABLE}` with variables found in a shell environment.
- Syntaxes for setting up default values and specifying mandatory variables:
  - ${VARIABLE:-default} evaluates to default if VARIABLE is unset or empty in the environment.
  - ${VARIABLE-default} evaluates to default only if VARIABLE is unset in the environment.
  - ${VARIABLE:?err} exits with an error message containing err if VARIABLE is unset or empty in the environment.
  - ${VARIABLE?err} exits with an error message containing err if VARIABLE is unset in the environment.
- By default loads `.env`, but with optional agruments to `config.LoadFile` can load other files.

Examples:
