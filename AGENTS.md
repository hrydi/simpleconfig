# AGENTS.md

## Project Overview
Simple Go config package (`github.com/hrydi/simpleconfig`) that loads TOML configuration files using `github.com/pelletier/go-toml/v2`.

## Commands
- **Build**: `go build`
- **Test**: `go test ./...`
- **Lint**: `go vet ./...`
- **Format**: `go fmt ./...`
- **Tidy**: `go mod tidy`

## Package Structure
- `config.go` - Main package (`simpleconfig`) with `Load[T]`, `MustLoad[T]`, `Result[T]`, and functional options
- `config_test.go` - Tests
- Searches for config in: custom paths (with filename), `./config/config.toml`, `./config.toml`

## Key Conventions
- Module: `github.com/hrydi/simpleconfig`
- Package name: `simpleconfig` (matching module name)
- Go version: 1.24.13
- Uses generics for type-safe config loading
- Functional options pattern for configuration

## API
```go
// Load returns Result with Config and LoadedPath
res, err := simpleconfig.Load[MyConfig](
    simpleconfig.WithFilename("myconfig.toml"),
    simpleconfig.WithCustomPaths("/etc/myapp"),
    simpleconfig.WithValidator(func(cfg *MyConfig) error { ... }),
)

// MustLoad calls log.Fatal on error
res := simpleconfig.MustLoad[MyConfig](...)

// Result type
type Result[T any] struct {
    Config     *T
    LoadedPath string
}
```

## Search Priority
1. Custom paths (in order provided) + filename
2. `./config/<filename>`
3. `./<filename>`

## Usage Example
```go
res, err := simpleconfig.Load[MyConfig](
    simpleconfig.WithFilename("myconfig.toml"),
    simpleconfig.WithCustomPaths("/etc/myapp"),
    simpleconfig.WithValidator(func(cfg *MyConfig) error {
        if cfg.Server.Host == "" {
            return errors.New("server.host required")
        }
        return nil
    }),
)
fmt.Println("Loaded from:", res.LoadedPath)
```