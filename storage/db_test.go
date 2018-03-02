package storage

import (
	"testing"

	"github.com/laohanlinux/assert"
)

func TestDB(t *testing.T) {
	t.Run("test changes object", func(t *testing.T) {
		changes := NewChanges()
		assert.NotNil(t, changes)
		patch := NewPatch()
		assert.NotNil(t, patch)
		assert.Nil(t, patch.Changes("test"))
		patch.InsertChanges("test", changes)
		assert.NotNil(t, patch.Changes("test"))

		changes = patch.Changes("test")
		assert.NotNil(t, changes)

	})
}
