# Examples

## Multiple Config Files

```go
// Load base config, then override with environment-specific
base, _ := simpleconfig.Load[Config](WithFilename("config.toml"))
prod, _ := simpleconfig.Load[Config](WithFilename("config.prod.toml"))

// Merge manually or use prod values
```

## Environment-based Loading

```go
func LoadConfig() (*Config, error) {
	env := os.Getenv("APP_ENV")
	filename := "config.toml"
	
	switch env {
	case "production", "prod":
		filename = "config.prod.toml"
	case "staging", "stage":
		filename = "config.staging.toml"
	case "development", "dev":
		filename = "config.dev.toml"
	}
	
	return simpleconfig.Load[Config](
		WithFilename(filename),
		WithCustomPaths("/etc/myapp"),
	)
}
```

## Nested Configuration

```go
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Logging  LoggingConfig
}

type ServerConfig struct {
	Host string
	Port int
	TLS  TLSConfig
}

type TLSConfig struct {
	Enabled  bool
	CertFile string
	KeyFile  string
}
```

```toml
[server]
host = "0.0.0.0"
port = 8080

[server.tls]
enabled = true
cert_file = "/etc/ssl/cert.pem"
key_file = "/etc/ssl/key.pem"

[database]
dsn = "postgres://user:pass@localhost/db"

[redis]
addr = "localhost:6379"
password = ""
db = 0

[logging]
level = "info"
format = "json"
```

## Kubernetes ConfigMap

```go
// In Kubernetes, mount ConfigMap at /etc/config
res := simpleconfig.MustLoad[Config](
	simpleconfig.WithCustomPaths("/etc/config"),
)
```

## Docker Compose

```yaml
# docker-compose.yml
services:
  app:
    volumes:
      - ./config:/etc/myapp
```

```go
res := simpleconfig.MustLoad[Config](
	simpleconfig.WithCustomPaths("/etc/myapp"),
)
```

## Default Values

```go
func LoadWithDefaults() *Config {
	res := simpleconfig.MustLoad[Config]()
	
	// Apply defaults for missing values
	if res.Config.Server.Port == 0 {
		res.Config.Server.Port = 8080
	}
	if res.Config.Server.Host == "" {
		res.Config.Server.Host = "localhost"
	}
	
	return res.Config
}
```

---

[← Back to index](index)