package rpc

import (
	"bytes"
	"encoding/gob"
	"errors"
	"net"
	"testing"
	"time"

	log "github.com/laohanlinux/utils/gokitlog"
	"github.com/laohanlinux/utils/netrpc"
	"golang.org/x/net/context"
)

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Arith int

func (t Arith) Multiply(ctx context.Context, args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t Arith) Divide(ctx context.Context, args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B

	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	dec := gob.NewDecoder(&network)
	// encoding
	if err := enc.Encode(quo); err != nil {
		log.Crit("err", err)
	}
	// decoding
	quo.Quo, quo.Rem = 0, 0
	if err := dec.Decode(quo); err != nil {
		log.Crit("err", err)
	}
	log.Infof("[quo:%v]", *quo)
	return nil
}

func netrpcServer() {
	var (
		address = []string{"127.0.0.1:12345", "127.0.0.1:12346", "127.0.0.1:12347"}
	)
	for _, addr := range address {
		arith := new(Arith)
		l, e := net.Listen("tcp", addr)
		if e != nil {
			log.Crit("listen error:", e)
		}
		defer l.Close()
		server := netrpc.NewServer()
		server.Register(*arith)
		server.Register(&netrpc.HealthCheck{})
		go server.Accept(l)
	}
	time.Sleep(time.Minute * 100000)
}

func netrpcClientWorker() {
	var (
		address = []string{"127.0.0.1:12345", "127.0.0.1:12346", "127.0.0.1:12347"}
		opts    = []NetRPCRingOpt{}
	)
	for _, addr := range address {
		opts = append(opts, NetRPCRingOpt{Addr: addr, NetWork: "tcp", PoolSize: 10, test: true})
	}
	r, err := NewNetRPCRing(opts)
	if err != nil {
		log.Crit("err", err)
	}
	// Asynchronous call
	quotient := new(Quotient)
	for i := 0; i < 1000; i++ {
		args := &Args{i, 8}
		// metaData := make(map[string]string)
		// metaData["time"] = fmt.Sprintf("%v", time.Now().UnixNano())
		err := r.Call("Arith.Divide", args, quotient)
		if err != nil {
			log.Error("err", err)
		}
		log.Infof("quotient:%v", *quotient)
		time.Sleep(time.Millisecond * 200)
	}
}

func TestStartClient(t *testing.T) {

	time.Sleep(time.Millisecond * 500)
	netrpcClientWorker()
}

func TestStartServer(t *testing.T) {
	netrpcServer()
}
