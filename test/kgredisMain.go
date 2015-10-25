package main

import (
    "fmt"
    "net/http"
    _ "net/http/pprof"
    "time"

    "math/rand"
    "runtime"
    "strconv"
    "sync"

    "github.com/laohanlinux/memCache/kgredis"
)

func doSetWork() {
    var g sync.WaitGroup
    for i := 0; i < 1024; i++ {
        g.Add(1)
        go func(key int) {
            defer g.Done()
            for {
                r := Get()

                p.SetString(strconv.Itoa(key), []byte(strconv.Itoa(key)), 10)
                runtime.Gosched()
            }
        }(i)
    }

    g.Wait()

}

func doGetWork() {
    var g sync.WaitGroup
    for i := 0; i < 102400; i++ {
        g.Add(1)
        go func() {
            defer g.Done()
            for {
                idx := rand.Intn(1024)
                value, err := p.GetString(strconv.Itoa(idx))
                if err != nil {
                    fmt.Println("get error:", err)
                } else {

                    fmt.Printf("key:%d, value:%v\n", idx, value)
                }

                runtime.Gosched()
            }
        }()
    }
    g.Wait()

}

var numbers = 10

func main() {

    runtime.GOMAXPROCS(runtime.NumCPU())
    //NewRedisPoolConn(address string, maxIdle, maxActive int, idleTimeout, defaultExpire int64)
    err := kgredis.NewRedisPoolChan("localhost:6379", 100)
    if err == nil {
        fmt.Println("redis pool is empty")
        return
    }
    go doSetWork()
    go doGetWork()

    http.ListenAndServe(":6060", nil)
    time.Sleep(time.Second * 30)
}
