# simpleconfig

Simple Go config package that loads TOML configuration files.

## Features

- Type-safe config loading with generics
- Functional options pattern
- Search priority: custom paths → `./config/config.toml` → `./config.toml`
- Optional validation hook
- Returns loaded file path

## Install

```bash
go get github.com/hrydi/simpleconfig
```

## Usage

```go
package main

import (
	"fmt"
	"github.com/hrydi/simpleconfig"
)

type Config struct {
	Server struct {
		Host string
		Port int
	}
	Database struct {
		DSN string
	}
}

func main() {
	res, err := simpleconfig.Load[Config](
		simpleconfig.WithFilename("myconfig.toml"),
		simpleconfig.WithCustomPaths("/etc/myapp"),
		simpleconfig.WithValidator(func(cfg *Config) error {
			if cfg.Server.Host == "" {
				return fmt.Errorf("server.host is required")
			}
			return nil
		}),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Loaded from:", res.LoadedPath)
	fmt.Printf("Server: %s:%d\n", res.Config.Server.Host, res.Config.Server.Port)
}
```

```toml
# myconfig.toml
[server]
host = "localhost"
port = 8080

[database]
dsn = "postgres://user:pass@localhost/db"
```

## API

### `Load[T](opts ...OptionFunc) (*Result[T], error)`

Loads config from the first matching file.

### `MustLoad[T](opts ...OptionFunc) *Result[T]`

Like `Load` but calls `log.Fatal` on error.

### `Result[T]`

```go
type Result[T any] struct {
	Config     *T
	LoadedPath string
}
```

### Options

- `WithFilename(name string)` - Config filename (default: "config.toml")
- `WithCustomPaths(paths ...string)` - Additional search directories
- `WithValidator(fn func(*T) error)` - Validation callback

## Search Priority

1. Custom paths (in order) + filename
2. `./config/<filename>`
3. `./<filename>`

## License

MIT