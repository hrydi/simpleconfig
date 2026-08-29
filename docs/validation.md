# Validation

## Validator Function

The `WithValidator` option accepts a function that receives the parsed config and returns an error if validation fails.

```go
simpleconfig.WithValidator(func(cfg *Config) error {
	if cfg.Server.Host == "" {
		return fmt.Errorf("server.host is required")
	}
	if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
		return fmt.Errorf("server.port must be 1-65535")
	}
	if cfg.Database.DSN == "" {
		return fmt.Errorf("database.dsn is required")
	}
	return nil
})
```

## Error Handling

Validation errors wrap the original error:

```go
res, err := simpleconfig.Load[Config](WithValidator(validator))
if err != nil {
	// err contains "validation failed: server.host is required"
	var validationErr *simpleconfig.ValidationError
	if errors.As(err, &validationErr) {
		// handle validation-specific error
	}
}
```

## Common Patterns

### Required Fields

```go
WithValidator(func(cfg *Config) error {
	if cfg.APIKey == "" {
		return errors.New("api_key required")
	}
	return nil
})
```

### Cross-field Validation

```go
WithValidator(func(cfg *Config) error {
	if cfg.TLS.Enabled && cfg.TLS.CertFile == "" {
		return errors.New("tls.cert_file required when tls.enabled=true")
	}
	return nil
})
```

### Conditional Validation

```go
WithValidator(func(cfg *Config) error {
	if cfg.Environment == "production" {
		if cfg.Debug {
			return errors.New("debug must be false in production")
		}
		if cfg.Database.DSN == "" {
			return errors.New("database.dsn required in production")
		}
	}
	return nil
})
```

---

[← Back to index](index)