[![GoDoc](https://godoc.org/github.com/golobby/container?status.svg)](https://godoc.org/github.com/golobby/container)
[![Build Status](https://travis-ci.org/golobby/container.svg?branch=master)](https://travis-ci.org/golobby/container)
[![Go Report Card](https://goreportcard.com/badge/github.com/golobby/container)](https://goreportcard.com/report/github.com/golobby/container)
[![Awesome](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/sindresorhus/awesome) 
[![Coverage Status](https://coveralls.io/repos/github/golobby/container/badge.svg?branch=master)](https://coveralls.io/github/golobby/container?branch=master)

# Container
A lightweight yet powerful IoC container for Go projects. It provides a simple, fluent and easy-to-use interface to make dependency injection in GoLang easier.
Atrico fork adds "lazy" creation of singletons and Container as an interface (not a global)

## Documentation

### Required Go Versions
It requires Go `v1.11` or newer versions.

### Installation
To install this package, run the following command in the root of your project.

```bash
go get github.com/atrico-go/container
```

### Introduction
GoLobby Container like any other IoC container is used to bind abstractions to their implementations.
Binding is a process of introducing an IoC container that which concrete (implementation) is appropriate for an abstraction. In this process, you also determine how it must be resolved, singleton or transient. 
In singleton binding, the container provides an instance once and returns it for each request. 
In transient binding, the container always returns a brand new instance for each request.
After the binding process, you can ask the IoC container to get the appropriate implementation of the abstraction that your code depends on. In this case, your code depends on abstractions, not implementations.
### Create container

```go
c := container.NewContainer()
```

### Binding

#### Singleton

Singleton binding using Container:

```go
c.Singleton(func() Abstraction {
  return Implementation
})
```

It takes a resolver function which its return type is the abstraction and the function body configures the related concrete (implementation) and returns it.

Example for a singleton binding:

```go
c.Singleton(func() Database {
  return &MySQL{}
})
```

#### Transient

Transient binding is also similar to singleton binding.

Example for a transient binding:

```go
c.Transient(func() Shape {
  return &Rectangle{}
})
```

### Resolving

Container resolves the dependencies with the method `make()`.

#### Using References

One way to get the appropriate implementation you need is to declare an instance of the abstraction type and pass its reference to Container this way:

```go
var a Abstraction
c.Make(&a)
// "a" will be implementation of the Abstraction
```

Example:

```go
var m Mailer
c.Make(&m)
m.Send("info@miladrahimi.com", "Hello Milad!")
```

#### Using Closures

Another way to resolve the dependencies is by using a function (receiver) that its arguments are the abstractions you 
need. Container will invoke the function and pass the related implementations for each abstraction.

```go
c.Make(func(a Abstraction) {
  // "a" will be implementation of the Abstraction
})
```

Example:

```go
c.Make(func(db Database) {
  // "db" will be the instance of MySQL
  db.Query("...")
})
```

You can also resolve multiple abstractions this way:

```go
c.Make(func(db Database, s Shape) {
  db.Query("...")
  s.Area()
})
```

#### Binding time

You can also resolve a dependency at the binding time in your resolver function like the following example.

```go
// Bind Config to JsonConfig
c.Singleton(func() Config {
    return &JsonConfig{...}
})

// Bind Database to MySQL
c.Singleton(func(c Config) Database {
    // "c" will be the instance of JsonConfig
    return &MySQL{
        Username: c.Get("DB_USERNAME"),
        Password: c.Get("DB_PASSWORD"),
    }
})
```

Notice: You can only resolve the dependencies in a binding resolver function that has already bound.

### Usage Tips

#### Performance
The package Container inevitably uses reflection in binding and resolving processes. 
If performance is a concern, you should use this package more carefully. 
Try to bind and resolve the dependencies out of the processes that are going to run many times 
(for example, on each request), put it where that run only once when you run your applications 
like main and init functions.

## License

GoLobby Container is released under the [MIT License](http://opensource.org/licenses/mit-license.php).
Atrico-go container follows this license
