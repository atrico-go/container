package container_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/atrico-go/container"
)

type Shape interface {
	SetArea(int)
	GetArea() int
}

type Circle struct {
	a int
}

func (c *Circle) SetArea(a int) {
	c.a = a
}

func (c Circle) GetArea() int {
	return c.a
}

type Database interface {
	Connect() bool
}

type MySQL struct{}

func (m MySQL) Connect() bool {
	return true
}

func TestSingletonItShouldMakeAnInstanceOfTheAbstraction(t *testing.T) {
	area := 5

	c := container.NewContainer()
	c.Singleton(func() Shape {
		return &Circle{a: area}
	})

	c.Make(func(s Shape) {
		a := s.GetArea()
		assert.Equal(t, area, a)
	})
}

func TestSingletonItShouldMakeSameObjectEachMake(t *testing.T) {
	c := container.NewContainer()
	c.Singleton(func() Shape {
		return &Circle{a: 5}
	})

	area := 6

	c.Make(func(s1 Shape) {
		s1.SetArea(area)
	})

	c.Make(func(s2 Shape) {
		a := s2.GetArea()
		assert.Equal(t, a, area)
	})
}

func TestSingletonWithNonFunctionResolverItShouldPanic(t *testing.T) {
	value := "the resolver must be a function"
	c := container.NewContainer()
	assert.PanicsWithValue(t, value, func() {
		c.Singleton("STRING!")
	}, "Expected panic")
}

func TestSingletonItShouldResolveResolverArguments(t *testing.T) {
	area := 5
	c := container.NewContainer()
	c.Singleton(func() Shape {
		return &Circle{a: area}
	})

	c.Singleton(func(s Shape) Database {
		assert.Equal(t, s.GetArea(), area)
		return &MySQL{}
	})
}

func TestTransientItShouldMakeDifferentObjectsOnMake(t *testing.T) {
	area := 5

	c := container.NewContainer()
	c.Transient(func() Shape {
		return &Circle{a: area}
	})

	c.Make(func(s1 Shape) {
		s1.SetArea(6)
	})

	c.Make(func(s2 Shape) {
		a := s2.GetArea()
		assert.Equal(t, a, area)
	})
}

func TestTransientItShouldMakeAnInstanceOfTheAbstraction(t *testing.T) {
	area := 5

	c := container.NewContainer()
	c.Transient(func() Shape {
		return &Circle{a: area}
	})

	c.Make(func(s Shape) {
		a := s.GetArea()
		assert.Equal(t, a, area)
	})
}

func TestMakeWithSingleInputAndCallback(t *testing.T) {
	c := container.NewContainer()
	c.Singleton(func() Shape {
		return &Circle{a: 5}
	})

	c.Make(func(s Shape) {
		if _, ok := s.(*Circle); !ok {
			t.Error("Expected Circle")
		}
	})
}

func TestMakeWithMultipleInputsAndCallback(t *testing.T) {
	c := container.NewContainer()
	c.Singleton(func() Shape {
		return &Circle{a: 5}
	})

	c.Singleton(func() Database {
		return &MySQL{}
	})

	c.Make(func(s Shape, m Database) {
		if _, ok := s.(*Circle); !ok {
			t.Error("Expected Circle")
		}

		if _, ok := m.(*MySQL); !ok {
			t.Error("Expected MySQL")
		}
	})
}

func TestMakeWithSingleInputAndReference(t *testing.T) {
	c := container.NewContainer()
	c.Singleton(func() Shape {
		return &Circle{a: 5}
	})

	var s Shape

	c.Make(&s)

	if _, ok := s.(*Circle); !ok {
		t.Error("Expected Circle")
	}
}

func TestMakeWithMultipleInputsAndReference(t *testing.T) {
	c := container.NewContainer()
	c.Singleton(func() Shape {
		return &Circle{a: 5}
	})

	c.Singleton(func() Database {
		return &MySQL{}
	})

	var (
		s Shape
		d Database
	)

	c.Make(&s)
	c.Make(&d)

	if _, ok := s.(*Circle); !ok {
		t.Error("Expected Circle")
	}

	if _, ok := d.(*MySQL); !ok {
		t.Error("Expected MySQL")
	}
}

func TestMakeWithUnsupportedReceiver(t *testing.T) {
	value := "the receiver must be either a reference or a callback"
	c := container.NewContainer()
	assert.PanicsWithValue(t, value, func() {
		c.Make("STRING!")
	}, "Expected panic")
}

func TestMakeWithNonReference(t *testing.T) {
	value := "cannot detect type of the receiver, make sure your are passing reference of the object"
	c := container.NewContainer()
	assert.PanicsWithValue(t, value, func() {
		var s Shape
		c.Make(s)
	}, "Expected panic")
}

func TestMakeWithUnboundedAbstraction(t *testing.T) {
	value := "no concrete found for the abstraction container_test.Shape"
	c := container.NewContainer()
	assert.PanicsWithValue(t, value, func() {
		var s Shape
		c.Reset()
		c.Make(&s)
	}, "Expected panic")
}

func TestMakeWithCallbackThatHasAUnboundedAbstraction(t *testing.T) {
	value := "no concrete found for the abstraction: container_test.Database"
	c := container.NewContainer()
	assert.PanicsWithValue(t, value, func() {
		c.Reset()
		c.Singleton(func() Shape {
			return &Circle{}
		})
		c.Make(func(s Shape, d Database) {})
	}, "Expected panic")
}
