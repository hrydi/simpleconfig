# Usage

## Installation

```bash
go get github.com/hrydi/simpleconfig
```

## Basic Usage

```go
type Config struct {
	Server struct {
		Host string
		Port int
	}
	Database struct {
		DSN string
	}
}

res, err := simpleconfig.Load[Config]()
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Server: %s:%d\n", res.Config.Server.Host, res.Config.Server.Port)
```

## With Options

```go
res, err := simpleconfig.Load[Config](
	simpleconfig.WithFilename("app.toml"),
	simpleconfig.WithCustomPaths("/etc/myapp", "/opt/myapp/config"),
	simpleconfig.WithValidator(func(cfg *Config) error {
		if cfg.Server.Port == 0 {
			return fmt.Errorf("server.port must be > 0")
		}
		return nil
	}),
)
```

## MustLoad

```go
// Panics on error - useful for startup
res := simpleconfig.MustLoad[Config](
	simpleconfig.WithFilename("config.toml"),
)
```

## Search Priority

Files are checked in order:

1. Each custom path + filename: `/etc/myapp/app.toml`, `/opt/myapp/config/app.toml`
2. `./config/app.toml`
3. `./app.toml`

First match wins.

---

[← Back to index](index)