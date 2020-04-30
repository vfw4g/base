package configs

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestInitConfigLocation(t *testing.T) {
	root, _ := filepath.Abs(".")
	c, err := InitConfigLocation(filepath.Join(root, "example.toml"))
	if err != nil {
		t.Errorf("%+v\n", err)
	}
	testVaule := c.GetString("TEST_KEY")
	assert.Equal(t, "test", testVaule)
}

func TestInitConfig(t *testing.T) {
	paths := []string{
		".",
	}
	c, err := InitConfig("example", paths...)
	if err != nil {
		t.Errorf("%+v\n", err)
	}
	testVaule := c.GetString("TEST_KEY")
	assert.Equal(t, "test", testVaule)
}

func TestInitConfigInputErrorConfigName(t *testing.T) {
	paths := []string{
		".",
	}
	noExistConfigName := "example"
	_, err := InitConfig(noExistConfigName, paths...)
	assert.Errorf(t, err, "%+v\n", err)
}

func TestInitDefault(t *testing.T) {
	c, err := InitDefault()
	if err != nil {
		t.Errorf("%+v\n", err)
	}
	testVaule := c.GetString("TEST_KEY")
	assert.Equal(t, "test", testVaule)
}
