package snowflake

import (
	"time"

	log "github.com/laohanlinux/utils/gokitlog"
)

// Snowflake init setting
const (
	SnowService                = "SnowflakeNetRPC"
	SnowServiceTimestampNetRPC = "SnowflakeNetRPC.Timestamp"
	SnowServiceNextIDNetRPC    = "SnowflakeNetRPC.NextID"
	SnowServiceNextIDsNetRPC   = "SnowflakeNetRPC.NextIDs"
)

// NextIDsArgs args of rpc remote call
type NextIDsArgs struct {
	ID  int32
	Num int
}

// NewSnowflakeNetRPC return a new snowflakenetrpc point object
func NewSnowflakeNetRPC(workers Workers) *SnowflakeNetRPC {
	return &SnowflakeNetRPC{workers: workers}
}

// SnowflakeNetRPC object of rpc remote call
type SnowflakeNetRPC struct {
	workers Workers
}

// NextID returns next Snowflake id with special worker id
func (s *SnowflakeNetRPC) NextID(args *NextIDsArgs, id *int64) error {
	worker, err := s.workers.Get(args.ID)
	if err != nil {
		return err
	}
	if *id, err = worker.NextID(); err != nil {
		log.Error("err", err)
		return err
	}
	return nil
}

// NextIDs returns next mutiple Snowflake id with special worker id
func (s *SnowflakeNetRPC) NextIDs(args *NextIDsArgs, ids *[]int64) error {
	worker, err := s.workers.Get(args.ID)
	if err != nil {
		return err
	}

	if *ids, err = worker.NextIds(args.Num); err != nil {
		log.Error("worke.NextIDs", args.Num, "err", err)
		return err
	}
	return nil
}

// Timestamp returns Snowflake server time
func (s *SnowflakeNetRPC) Timestamp(ignore int, timestamp *int64) error {
	*timestamp = time.Now().Unix()
	return nil
}
