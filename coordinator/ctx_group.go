package coordinator

import (
	"context"
	"sync"
)

func NewCtxGCoordinator(g *sync.WaitGroup) *CtxGCoordinator {
	ctx := context.Background()
	ctx, canFunc := context.WithCancel(ctx)
	return &CtxGCoordinator{Ctx: ctx, CanFunc: canFunc, Group: g}
}

// CtxGCoordinator must sure one goroutine one parent or one child CtxGCoordinator object
type CtxGCoordinator struct {
	Ctx     context.Context
	CanFunc context.CancelFunc
	Group   *sync.WaitGroup
}

func (cg *CtxGCoordinator) StartIncr() *CtxGCoordinator {
	ctx, canFunc := context.WithCancel(cg.Ctx)
	cgc := &CtxGCoordinator{
		Ctx:     ctx,
		CanFunc: canFunc,
		Group:   cg.Group,
	}
	cg.Group.Add(1)
	return cgc
}

func (cg *CtxGCoordinator) TimeoutSub() bool {
	if err := cg.Ctx.Err(); err != nil {
		cg.Group.Done()
		return true
	}
	return false
}

// Done is the ctx cancel channel.
// if client call the function, must call Sub too.
func (cg *CtxGCoordinator) Done() <-chan struct{} {
	return cg.Ctx.Done()
}

// Sub cg.Group one
func (cg *CtxGCoordinator) Sub() {
	cg.Group.Done()
}

// Stop all child goroutines and waitting theme exit
func (cg *CtxGCoordinator) Stop() {
	cg.CanFunc()
	cg.Group.Wait()
}
