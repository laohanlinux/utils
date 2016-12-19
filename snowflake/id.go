// Copyright © 2014 Terry Mao All rights reserved.
// This file is part of gosnowflake.

// gosnowflake is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// gosnowflake is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with gosnowflake.  If not, see <http://www.gnu.org/licenses/>.

// Reference: https://github.com/Terry-Mao/gosnowflake.git

package snowflake

import (
	"fmt"
	"sync"
	"time"

	log "github.com/laohanlinux/utils/gokitlog"
)

const (
	twepoch            = int64(1288834974657)
	workerIDBits       = uint(5)
	datacenterIDBits   = uint(5)
	maxWorkerID        = -1 ^ (-1 << workerIDBits)
	maxDatacenterID    = -1 ^ (-1 << datacenterIDBits)
	sequenceBits       = uint(12)
	workerIDShift      = sequenceBits
	datacenterIDShift  = sequenceBits + workerIDBits
	timestampLeftShift = sequenceBits + workerIDBits + datacenterIDBits
	sequenceMask       = -1 ^ (-1 << sequenceBits)
	maxNextIdsNum      = 100
)

type IDWorker struct {
	sequence      int64
	lastTimestamp int64
	workerID      int32
	twepoch       int64
	datacenterID  int32
	mutex         sync.Mutex
}

// NewIDWorker new a snowflake id generator object.
func NewIDWorker(workerID, datacenterID int32, twepoch int64) (*IDWorker, error) {
	idWorker := &IDWorker{}
	if workerID > maxWorkerID || workerID < 0 {
		log.Errorf("worker Id can't be greater than %d or less than 0", maxWorkerID)
		return nil, fmt.Errorf("worker Id: %d error", workerID)
	}
	if datacenterID > maxDatacenterID || datacenterID < 0 {
		log.Errorf("datacenter Id can't be greater than %d or less than 0", maxDatacenterID)
		return nil, fmt.Errorf("datacenter Id: %d error", datacenterID)
	}
	idWorker.workerID = workerID
	idWorker.datacenterID = datacenterID
	idWorker.lastTimestamp = -1
	idWorker.sequence = 0
	idWorker.twepoch = twepoch
	idWorker.mutex = sync.Mutex{}
	//log.Debugf("worker starting. timestamp left shift %d, datacenter id bits %d, worker id bits %d, sequence bits %d, workerid %d", timestampLeftShift, datacenterIDBits, workerIDBits, sequenceBits, workerID)
	return idWorker, nil
}

// timeGen generate a unix millisecond.
func timeGen() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// tilNextMillis spin wait till next millisecond.
func tilNextMillis(lastTimestamp int64) int64 {
	timestamp := timeGen()
	for timestamp <= lastTimestamp {
		timestamp = timeGen()
	}
	return timestamp
}

// NextID get a snowflake id.
func (id *IDWorker) NextID() (int64, error) {
	id.mutex.Lock()
	defer id.mutex.Unlock()
	timestamp := timeGen()
	if timestamp < id.lastTimestamp {
		log.Errorf("clock is moving backwards.  Rejecting requests until %d.", id.lastTimestamp)
		return 0, fmt.Errorf("Clock moved backwards.  Refusing to generate id for %d milliseconds", id.lastTimestamp-timestamp)
	}
	if id.lastTimestamp == timestamp {
		id.sequence = (id.sequence + 1) & sequenceMask
		if id.sequence == 0 {
			timestamp = tilNextMillis(id.lastTimestamp)
		}
	} else {
		id.sequence = 0
	}
	id.lastTimestamp = timestamp
	return ((timestamp - id.twepoch) << timestampLeftShift) | (int64(id.datacenterID) << datacenterIDShift) | (int64(id.workerID) << workerIDShift) | id.sequence, nil
}

// NextIds get snowflake ids.
func (id *IDWorker) NextIds(num int) ([]int64, error) {
	if num > maxNextIdsNum || num < 0 {
		log.Errorf("NextIds num can't be greater than %d or less than 0", maxNextIdsNum)
		return nil, fmt.Errorf("NextIds num: %d error", num)
	}
	ids := make([]int64, num)
	id.mutex.Lock()
	defer id.mutex.Unlock()
	for i := 0; i < num; i++ {
		timestamp := timeGen()
		if timestamp < id.lastTimestamp {
			log.Errorf("clock is moving backwards.  Rejecting requests until %d.", id.lastTimestamp)
			return nil, fmt.Errorf("Clock moved backwards.  Refusing to generate id for %d milliseconds", id.lastTimestamp-timestamp)
		}
		if id.lastTimestamp == timestamp {
			id.sequence = (id.sequence + 1) & sequenceMask
			if id.sequence == 0 {
				timestamp = tilNextMillis(id.lastTimestamp)
			}
		} else {
			id.sequence = 0
		}
		id.lastTimestamp = timestamp
		ids[i] = ((timestamp - id.twepoch) << timestampLeftShift) | (int64(id.datacenterID) << datacenterIDShift) | (int64(id.workerID << workerIDShift)) | id.sequence
	}
	return ids, nil
}
