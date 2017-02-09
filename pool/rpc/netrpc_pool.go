package rpc

import (
	"github.com/laohanlinux/utils/netrpc"
)

type netrpcClients struct {
	idx     uint64
	sources []*netrpcClients
}

func (ncs *netrpcClients) getAliveConn() *netrpcClient {

}

type netrpcClient struct {
	alive bool
	*netrpc.Client
}
