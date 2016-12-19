package snowflake

import (
	"fmt"
	"testing"

	"github.com/laohanlinux/assert"
)

func TestIDWorker(t *testing.T) {
	idWorkers, err := NewIDWorker(1, 1, 1288834974657)
	assert.Nil(t, err)
	assert.NotNil(t, idWorkers)
	ids, err := idWorkers.NextIds(200)
	assert.NotNil(t, err)
	for idx, v := range ids {
		fmt.Println("idx:", idx, "uid:", v)
	}
}
