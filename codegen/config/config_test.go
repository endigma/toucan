package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/endigma/toucan/codegen/config"
)

func TestLoadConfig(t *testing.T) {
	t.Run("config does not exist", func(t *testing.T) {
		_, err := config.LoadConfig("doesnotexist.hcl")
		assert.Error(t, err)
	})
}

func TestReadConfig(t *testing.T) {
	t.Run("invalid config", func(t *testing.T) {
		_, err := config.ReadConfig([]byte("invalid"))
		assert.Error(t, err)
	})

	t.Run("malformed config", func(t *testing.T) {
		_, err := config.ReadConfig([]byte(`um the uhh the um`))
		assert.Error(t, err)
	})

	t.Run("extra keys", func(t *testing.T) {
		_, err := config.LoadConfig("testdata/unknownkeys.hcl")
		assert.Error(t, err)
	})
}
