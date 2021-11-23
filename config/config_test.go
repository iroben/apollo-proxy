package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitConfig(t *testing.T) {
	InitConfig()
	assert.Greater(t, Config.App.Port, 0)
}
