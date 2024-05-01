package di

import (
	"errors"

	"go.uber.org/dig"
)

type (
	Constructor struct {
		Name string
		Fn   interface{}
	}
)

var container *dig.Container

func Provide(constructor interface{}, opts ...dig.ProvideOption) error {
	if container == nil {
		container = dig.New()
	}

	return container.Provide(constructor, opts...)
}

func Invoke(fn interface{}, opts ...dig.InvokeOption) error {
	if container == nil {
		return errors.New("no constructor provided")
	}
	return container.Invoke(fn, opts...)
}