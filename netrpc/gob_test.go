package netrpc

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

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
		log.Fatal("err", err)
	}
	// decoding
	quo.Quo, quo.Rem = 0, 0
	if err := dec.Decode(quo); err != nil {
		log.Fatal("err", err)
	}
	log.Printf("[quo:%v]", *quo)
	return nil
}

func netrpcServer() {
	arith := new(Arith)
	l, e := net.Listen("tcp", "127.0.0.1:12345")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	defer l.Close()
	server := NewServer()
	server.Register(*arith)
	server.Register(&HealthCheck{})
	server.Accept(l)
}

func netrpcClient() {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:12345", time.Second)
	if err != nil {
		log.Fatalf("err:%v", err)
	}
	defer conn.Close()

	codec := NewGoClientCodec(conn)
	client := NewClientWithCodec(codec)
	// Asynchronous call
	quotient := new(Quotient)
	for i := 0; i < 10; i++ {
		args := &Args{i, 8}
		metaData := make(map[string]string)
		metaData["time"] = fmt.Sprintf("%v", time.Now().UnixNano())
		//divCall := client.Call("Arith.Divide", metaData, args, quotient)
		//replyCall := <-divCall.Done // will be equal to divCall
		// check errors, print, etc.
		err := client.Call("Arith.Divide", metaData, args, quotient)
		if err != nil {
			log.Fatalf("err:%v", err)
		}
		log.Printf("reply=>  %+v\n", quotient)
	}
}

func TestGobDec(t *testing.T) {
	go netrpcServer()

	time.Sleep(time.Millisecond * 500)
	netrpcClient()
}
