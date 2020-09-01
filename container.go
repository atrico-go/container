// Package container provides an IoC container for Go projects.
// It provides simple, fluent and easy-to-use interface to make dependency injection in GoLang easier.
package container

import (
	"reflect"
)

// Container interface
type Container interface {
	// Singleton will bind an abstraction to a concrete for further singleton resolves.
	// It takes a resolver function which returns the concrete and its return type matches the abstraction (interface).
	// The resolver function can have arguments of abstraction that have bound already in Container.
	Singleton(resolver interface{})
	// Transient will bind an abstraction to a concrete for further transient resolves.
	// It takes a resolver function which returns the concrete and its return type matches the abstraction (interface).
	// The resolver function can have arguments of abstraction that have bound already in Container.
	Transient(resolver interface{})
	// Reset will reset the container and remove all the bindings.
	Reset()
	// Make will resolve the dependency and return a appropriate concrete of the given abstraction.
	// It can take an abstraction (interface reference) and fill it with the related implementation.
	// It also can takes a function (receiver) with one or more arguments of the abstractions (interfaces) that need to be
	// resolved, Container will invoke the receiver function and pass the related implementations.
	Make(receiver interface{})
}

// Create default container
func NewContainer() Container {
	return &container{}
}

// Default implementation of Container
// container is the IoC container that will keep all of the bindings.
type container map[reflect.Type]binding

func (c *container) Singleton(resolver interface{}) {
	c.bind(resolver, true)
}

func (c *container) Transient(resolver interface{}) {
	c.bind(resolver, false)
}

// Reset will reset the container and remove all the bindings.
func (c *container) Reset() {
	*c = map[reflect.Type]binding{}
}

// Make will resolve the dependency and return a appropriate concrete of the given abstraction.
// It can take an abstraction (interface reference) and fill it with the related implementation.
// It also can takes a function (receiver) with one or more arguments of the abstractions (interfaces) that need to be
// resolved, Container will invoke the receiver function and pass the related implementations.
func (c *container) Make(receiver interface{}) {
	receiverTypeOf := reflect.TypeOf(receiver)
	if receiverTypeOf == nil {
		panic("cannot detect type of the receiver, make sure your are passing reference of the object")
	}

	if receiverTypeOf.Kind() == reflect.Ptr {
		abstraction := receiverTypeOf.Elem()

		if instance := c.resolve(abstraction); instance != nil {
			reflect.ValueOf(receiver).Elem().Set(reflect.ValueOf(instance))
			return
		}

		panic("no concrete found for the abstraction " + abstraction.String())
	}

	if receiverTypeOf.Kind() == reflect.Func {
		arguments := c.arguments(receiver)
		reflect.ValueOf(receiver).Call(arguments)
		return
	}

	panic("the receiver must be either a reference or a callback")
}

// invoke will call the given function and return its returned value.
// It only works for functions that return a single value.
func (c *container) invoke(function interface{}) interface{} {
	return reflect.ValueOf(function).Call(c.arguments(function))[0].Interface()
}

// bind will map an abstraction to a concrete.
func (c *container) bind(resolver interface{}, singleton bool) {
	resolverTypeOf := reflect.TypeOf(resolver)
	if resolverTypeOf.Kind() != reflect.Func {
		panic("the resolver must be a function")
	}

	for i := 0; i < resolverTypeOf.NumOut(); i++ {
		(*c)[resolverTypeOf.Out(i)] = binding{
			resolver: resolver,
			instance: nil,
			singleton: singleton,
		}
	}
}

// arguments will return resolved arguments of the given function.
func (c *container) arguments(function interface{}) []reflect.Value {
	functionTypeOf := reflect.TypeOf(function)
	argumentsCount := functionTypeOf.NumIn()
	arguments := make([]reflect.Value, argumentsCount)

	for i := 0; i < argumentsCount; i++ {
		abstraction := functionTypeOf.In(i)
		instance := c.resolve(abstraction)
		if  instance == nil {
			panic("no concrete found for the abstraction: " + abstraction.String())
		}
		arguments[i] = reflect.ValueOf(instance)
	}

	return arguments
}

// resolve will return the concrete of related abstraction.
func (c *container) resolve(abstraction reflect.Type) interface{} {
	if b, ok := (*c)[abstraction]; ok {
		// Return singleton if already resolved
		if b.instance != nil {
			return b.instance
		}
		instance := c.invoke(b.resolver)
		if b.singleton {
			b.instance = instance
			(*c)[abstraction] = b
		}
		return instance
	}
	return nil
}

// binding keeps a binding resolver and instance (for singleton bindings).
type binding struct {
	resolver  interface{} // resolver function
	instance  interface{} // instance stored for singleton bindings (on first resolve)
	singleton bool
}

