package kgredis

import (
    "fmt"
    "github.com/fzzy/radix/redis"
    "time"
)

const (
    TIMEOUT_WAIT_REDISCLIENT = "waiting for redis client in timeout"
)

// define a []byte type memcache pool
type RedisPoolChan struct {
    redisClientChan chan *redis.Client
    addr            string
    blockSize       int
}

// A simple connection pool. It will create a small pool of initial connections,
// and if more connections are needed they will be created on demand. If a
// connection is returned and the pool is full it will be closed.

func NewRedisPoolChan(addr string, size int) (*RedisPoolChan, error) {

    rdpc := &RedisPoolChan{
        redisClientChan: make(chan *redis.Client, size),
        blockSize:       size,
        addr:            addr,
    }

    for i := 0; i < size; i++ {
        redisConn, _ := rdpc.newConn()
        rdpc.redisClientChan <- redisConn
    }

    redisConn, err := rdpc.newConn()
    if err == nil {
        return rdpc, redisConn.Cmd("PING").Err
    }
    return rdpc, err
}

func (rdpc *RedisPoolChan) Get(timeout int) (r *redis.Client, err error) {
    select {
    case r = <-rdpc.redisClientChan:
    case <-time.After(time.Second * time.Duration(timeout)):
        err = fmt.Errorf("%s", TIMEOUT_WAIT_REDISCLIENT)
    }
    return
}

func (rdpc *RedisPoolChan) Put(r *redis.Client, err error) {
    if err != nil {
        if r != nil {
            r.Close()
        }
        r, _ = rdpc.newConn()
    }
    rdpc.redisClientChan <- r
}

func (rdpc *RedisPoolChan) Ping() error {
    r, err := rdpc.Get(10)

    if err != nil {
        return err
    }

    return r.Cmd("PING").Err
}

func (rdpc *RedisPoolChan) newConn() (*redis.Client, error) {
    r, err := redis.Dial("tcp", rdpc.addr)
    if err != nil {
        r.ReadReply()
    }
    return r, err
}
