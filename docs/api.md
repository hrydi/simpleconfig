# API Reference

## Types

### Result[T]

```go
type Result[T any] struct {
	Config     *T
	LoadedPath string
}
```

Returned by `Load` and `MustLoad`. Contains the parsed config and the path of the loaded file.

---

## Functions

### Load

```go
func Load[T any](opts ...OptionFunc) (*Result[T], error)
```

Loads config from the first matching file. Returns `(*Result[T], nil)` on success, or `(nil, error)` if no file found or decode/validation fails.

### MustLoad

```go
func MustLoad[T any](opts ...OptionFunc) *Result[T]
```

Same as `Load` but calls `log.Fatal` on error. Use for application startup where config is required.

---

## Option Functions

### WithFilename

```go
func WithFilename(name string) OptionFunc
```

Sets the config filename. Default: `"config.toml"`.

### WithCustomPaths

```go
func WithCustomPaths(paths ...string) OptionFunc
```

Adds custom search directories. Each path is joined with the filename.

```go
WithCustomPaths("/etc/myapp", "/opt/myapp/config")
// Searches: /etc/myapp/config.toml, /opt/myapp/config/config.toml, ...
```

### WithValidator

```go
func WithValidator(fn func(any) error) OptionFunc
```

Adds a validation function. Called after successful TOML decode. Receives `*T` (pointer to config struct).

```go
WithValidator(func(cfg *Config) error {
	if cfg.Server.Host == "" {
		return errors.New("server.host required")
	}
	return nil
})
```

---

## Type Parameters

`T` must be a struct type (or pointer to struct) with exported fields that match TOML keys.

---

[← Back to index](index)