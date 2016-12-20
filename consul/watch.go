package consul

import (
	"fmt"
	"time"

	"github.com/hashicorp/consul/api"
)

const (
	DefaultWaitTime = 3600
)

// WatchKV watch the kv until it's ModifyIndex state change, the change condition are
// 1. kv change (deleted, updated, created)
// 2. session is removed
// Notic: the kv may be updated after call the function
func WatchKV(client *api.Client, kv *api.KVPair, WaitTime int) error {
	if WaitTime == 0 {
		WaitTime = DefaultWaitTime
	}
	for {
		tmpKV, _, err := client.KV().Get(kv.Key,
			&api.QueryOptions{
				WaitIndex: kv.ModifyIndex,
				WaitTime:  time.Second * time.Duration(WaitTime)})
		if err != nil {
			return err
		}
		if tmpKV == nil {
			kv = nil
			return fmt.Errorf("the key was deleted")
		}
		if tmpKV.ModifyIndex > kv.ModifyIndex {
			kv.ModifyIndex = tmpKV.ModifyIndex
			return nil
		}
		if tmpKV.Session == "" {
			return nil
		}
	}
	return nil
}

func UpdateKVModifyIndex(client *api.Client, kv *api.KVPair) error {
	tmpKV, _, err := client.KV().Get(kv.Key, &api.QueryOptions{WaitIndex: kv.ModifyIndex, WaitTime: time.Second})
	if err != nil {
		return err
	}
	if tmpKV == nil {
		return fmt.Errorf("the key was deleted")
	}
	kv.ModifyIndex = tmpKV.ModifyIndex
	return nil
}
