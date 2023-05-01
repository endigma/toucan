package schema_test

import (
	"os"
	"testing"

	"github.com/endigma/toucan/schema"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	t.Run("config does not exist", func(t *testing.T) {
		_, err := schema.LoadSchemaFile("doesnotexist.hcl")
		assert.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("config exists", func(r *testing.T) {
		_, err := schema.LoadSchemaFile("testdata/valid.hcl")
		assert.NoError(t, err)
	})

	t.Run("invalid model value", func(t *testing.T) {
		_, err := schema.LoadSchemaFile("testdata/invalid.hcl")
		assert.Error(t, err)
	})
}

func TestReadConfig(t *testing.T) {
	t.Parallel()

	t.Run("invalid config", func(t *testing.T) {
		_, err := schema.ReadSchemaFile([]byte("invalid"))
		assert.Error(t, err)
	})

	t.Run("malformed config", func(t *testing.T) {
		_, err := schema.ReadSchemaFile([]byte(`um the uhh the um`))
		assert.Error(t, err)
	})

	t.Run("extra keys", func(t *testing.T) {
		_, err := schema.ReadSchemaFile([]byte(`steve = "a man"`))
		assert.Error(t, err)
	})
}
