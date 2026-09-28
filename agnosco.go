package agnosco

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrMissingScheme = errors.New("dsn must be in scheme:// form")
	ErrUnknownScheme = errors.New("no driver registered for dsn scheme")
)

type Driver[T any] struct {
	Name    string
	Schemes []string
	Open    func(dsn string) T
}

type Registry[T any] struct {
	drivers map[string]Driver[T]
}

func NewRegistry[T any]() *Registry[T] {
	return &Registry[T]{drivers: make(map[string]Driver[T])}
}

func (r *Registry[T]) Register(driver Driver[T]) {
	if driver.Open == nil {
		panic("agnosco: driver open must not be nil")
	}
	if len(driver.Schemes) == 0 {
		panic("agnosco: driver must have at least one scheme")
	}
	if r.drivers == nil {
		r.drivers = make(map[string]Driver[T])
	}
	for _, scheme := range driver.Schemes {
		if _, ok := r.drivers[scheme]; ok {
			panic("agnosco: driver already registered for scheme " + scheme)
		}
	}
	for _, scheme := range driver.Schemes {
		r.drivers[scheme] = driver
	}
}

func (r *Registry[T]) Open(dsn string) (T, error) {
	var zero T

	scheme, _, ok := strings.Cut(strings.ToLower(dsn), "://")
	if !ok {
		return zero, fmt.Errorf("%w: %q", ErrMissingScheme, dsn)
	}

	driver, ok := r.drivers[scheme]
	if !ok {
		return zero, fmt.Errorf("%w %q, dsn %q", ErrUnknownScheme, scheme, dsn)
	}

	return driver.Open(dsn), nil
}
