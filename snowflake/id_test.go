package snowflake

import (
	"fmt"
	"testing"
	"time"

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

func BenchmarkIDWorker(b *testing.B) {

	idWorker, _ := NewIDWorker(1, 1, 1288834974657)

	var n uint64
	for i := 0; i < b.N; i++ {
		if j, err := idWorker.NextID(); err != nil {
			fmt.Println(err)
		} else {
			n++
			fmt.Println(j, time.Now().Unix(), idWorker.lastTimestamp)
		}
	}

	fmt.Println("n", n)
}
