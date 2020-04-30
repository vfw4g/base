package conv

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type Book struct {
	Name   string
	Author string
}

type Foo struct {
	Name  string
	Age   int
	Roles []string
	Books map[string]*Book
	*Book
	TFieldA string `fieldcopy:"field_a"`
}

type Bar struct {
	Roles []string
	Age   int
	Name  string `cp:"name"`
	Books map[string]*Book
	*Book
	TFielda string `fieldcopy:"field_a,omitempty"`
}

func TestStructFieldCopy(t *testing.T) {
	t.Run("sourceStruct", func(t *testing.T) {
		f := Foo{
			Name: "kook",
			Age:  18,
			Roles: []string{
				"admin",
			},
			Book: &Book{
				Name:   "b3",
				Author: "a3",
			},
			TFieldA: "fieldA",
		}
		b := Bar{}
		err := StructFieldCopy(f, &b)
		assert.NoError(t, err)
		assert.True(t, (f.Name == b.Name && f.Age == b.Age))
		assert.True(t, f.Book.Author == b.Book.Author)
		assert.True(t, f.TFieldA == b.TFielda)
		assert.ElementsMatch(t, f.Roles, b.Roles)
	})
	t.Run("sourcePtr", func(t *testing.T) {
		f := Foo{
			Name: "kook",
			Age:  18,
			Roles: []string{
				"admin",
			},
		}
		b := Bar{}
		err := StructFieldCopy(&f, &b)
		assert.NoError(t, err)
		assert.True(t, (f.Name == b.Name && f.Age == b.Age))
		assert.ElementsMatch(t, f.Roles, b.Roles)
	})
	t.Run("mapField", func(t *testing.T) {
		f := Foo{
			Name: "kook",
			Age:  18,
			Roles: []string{
				"admin",
			},
			Books: map[string]*Book{
				"b1": {Name: "b1", Author: "a1"},
				"b2": {Name: "b2", Author: "a2"},
			},
		}
		b := Bar{}
		err := StructFieldCopy(f, &b)
		assert.NoError(t, err)
		assert.True(t, b.Books["b2"].Author == "a2")
	})
	t.Run("structPtr", func(t *testing.T) {
		f := Foo{
			Name: "kook",
			Age:  18,
			Roles: []string{
				"admin",
			},
			Books: map[string]*Book{
				"b1": {Name: "b1", Author: "a1"},
				"b2": {Name: "b2", Author: "a2"},
			},
		}
		b := Bar{}
		err := StructFieldCopy(f, &b)
		assert.NoError(t, err)
		assert.True(t, b.Books["b2"].Author == "a2")
	})
	t.Run("omitempty", func(t *testing.T) {
		f := Foo{
			Name: "kook",
			Age:  18,
			Roles: []string{
				"admin",
			},
			Book: &Book{
				Name:   "b3",
				Author: "a3",
			},
			TFieldA: "", //zero omitempty
		}
		b := Bar{}
		err := StructFieldCopy(f, &b)
		assert.NoError(t, err)
		assert.True(t, "" == b.TFielda)
	})
}
