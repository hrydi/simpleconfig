package simpleconfig

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type Option struct {
	Filename    string
	CustomPaths []string
	Validator   func(any) error
}

type OptionFunc func(*Option)

func WithFilename(name string) OptionFunc {
	return func(o *Option) {
		o.Filename = name
	}
}

func WithCustomPaths(paths ...string) OptionFunc {
	return func(o *Option) {
		o.CustomPaths = paths
	}
}

func WithValidator(fn func(any) error) OptionFunc {
	return func(o *Option) {
		o.Validator = fn
	}
}

type Result[T any] struct {
	Config     *T
	LoadedPath string
}

func getSearchPaths(opt Option) []string {
	fileName := opt.Filename
	if fileName == "" {
		fileName = "config.toml"
	}

	paths := make([]string, 0, len(opt.CustomPaths)+2)
	for _, p := range opt.CustomPaths {
		paths = append(paths, filepath.Join(p, fileName))
	}
	paths = append(paths,
		filepath.Join(".", "config", fileName),
		filepath.Join(".", fileName),
	)

	return paths
}

func Load[T any](opts ...OptionFunc) (*Result[T], error) {
	opt := Option{}
	for _, fn := range opts {
		fn(&opt)
	}

	searchPaths := getSearchPaths(opt)

	for _, path := range searchPaths {
		file, err := os.Open(path)
		if err != nil {
			continue
		}

		var cfg T
		err = func() error {
			defer file.Close()
			if err := toml.NewDecoder(file).Decode(&cfg); err != nil {
				return err
			}
			if opt.Validator != nil {
				if err := opt.Validator(&cfg); err != nil {
					return fmt.Errorf("validation failed: %w", err)
				}
			}
			return nil
		}()

		if err != nil {
			continue
		}

		return &Result[T]{Config: &cfg, LoadedPath: path}, nil
	}

	return nil, fmt.Errorf("couldn't find any config file in paths: %v", searchPaths)
}

func MustLoad[T any](opts ...OptionFunc) *Result[T] {
	res, err := Load[T](opts...)
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	return res
}
