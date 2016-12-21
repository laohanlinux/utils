package snowflake

import (
	"testing"

	"github.com/laohanlinux/assert"
)

func TestSnowflake(t *testing.T) {
	workers, err := NewWorkers(0)
	assert.Nil(t, err)
	assert.Equal(t, len(workers), 1024)
	sn := NewSnowflakeNetRPC(workers)
	var id int64
	err = sn.NextID(&NextIDsArgs{ID: 2, Num: 1}, &id)
	assert.Nil(t, err)
	assert.Equal(t, true, id > 0)
}
