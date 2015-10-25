package kgredis

import (
    "container/list"
    //	"errors"
    "fmt"
    "github.com/fzzy/radix/redis"
    "runtime/debug"
    "sync"
    "time"
)

var rdpc *RedisPoolChan

// define a []byte type memcache pool
type RedisPoolChan struct {
    send      chan *redis.Client
    recv      chan *redis.Client
    once      sync.Once
    addr      string
    blockSize uint64
}

// A simple connection pool. It will create a small pool of initial connections,
// and if more connections are needed they will be created on demand. If a
// connection is returned and the pool is full it will be closed.

func NewRedisPoolChan(addr string, size int) error {
    if rdpc != nil {
        return rdpc
    }

    rdpc = &RedisPoolChan{
        send:      make(chan chan *redis.Client),
        recv:      make(chan chan *redis.Client),
        blockSize: size,
        addr:      addr,
    }

    if rdpc != nil {
        nbc.once.Do(func() { go rdpc.doWork() })
    }

    redis, err := rdpc.newConn()
    if err != nil {
        return err
    }
    redis.Close()
    return nil
}

func Get() redis.Client {
    return <-rdpc.send
}

func Put(r redis.Client, err error) {
    if err != nil {
        if redis != nil {
            r.Close()
        }
        return
    }
    rdpc.recv <- r
}

func (nbc *NonBlockingChan) doWork() {
    defer func() {
        if err := recover(); err != nil {
            // do some something
            fmt.Println("MemCachePool unexpectedly terminated, err:", err)
            go nbc.doWork()
        }
        debug.FreeOSMemory()
    }()

    items := list.New()
    timeout := time.NewTimer(time.Minute)

    for {
        if items.Len() == 0 {
            for i := 0; i < rdpc.blockSize; i++ {
                conn, err := rdpc.newConn()
                if err != nil {
                    continue
                }
                items.PushFront(conn)
            }
            //  items.PushFront(nbc.makeBuffer())
        }
        select {
        case e := <-rdpc.recv:
            items.PushBack(e)
        case rdpc.send <- items.Remove(items.Front()).(redis.Client):
        case <-timeout.C:
            // do clear memcache work, because no one are use the memcache pool
            itemsLen := items.Len()

            for i := 0; i < itemsLen; i++ {
                conn := items.Front().Value.(redis.Client)
                conn.Close()
                items.Remove(items.Front())
            }
            debug.FreeOSMemory()
            timeout.Reset(time.Minute)
        }
    }
}

func (rdpc *RedisPoolChan) newConn() (redis.Client, error) {
    return redis.Dial("tcp", rpc.addr)
}
