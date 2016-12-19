// // Copyright © 2014 Terry Mao All rights reserved.
// // This file is part of gosnowflake.

// // gosnowflake is free software: you can redistribute it and/or modify
// // it under the terms of the GNU General Public License as published by
// // the Free Software Foundation, either version 3 of the License, or
// // (at your option) any later version.

// // gosnowflake is distributed in the hope that it will be useful,
// // but WITHOUT ANY WARRANTY; without even the implied warranty of
// // MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// // GNU General Public License for more details.

// // You should have received a copy of the GNU General Public License
// // along with gosnowflake.  If not, see <http://www.gnu.org/licenses/>.

// // Reference: https://github.com/Terry-Mao/gosnowflake.git

package snowflake

import (
	"fmt"

	log "github.com/laohanlinux/utils/gokitlog"
)

const DefaultTwEpoch int64 = 1288834974657

// key = datacenterid + workerid
type Workers map[int32]*IDWorker

// NewWorkers new id workers instance.
func NewWorkers(twEpoch int64) (Workers, error) {
	if twEpoch == 0 {
		twEpoch = DefaultTwEpoch
	}
	idWorkers := make(map[int32]*IDWorker)
	maxID := -1 ^ (-1 << (workerIDBits + datacenterIDBits))
	for idx := 0; idx < maxID; idx++ {
		datacenter := idx >> workerIDBits
		worker := idx & maxWorkerID
		idWorker, err := NewIDWorker(int32(worker), int32(datacenter), twEpoch)
		if err != nil {
			log.Errorf("NewIdWorker(%d) error(%v)", worker, err)
			return nil, err
		}
		idWorkers[int32(idx)] = idWorker
	}
	return Workers(idWorkers), nil
}

// Get get a specified worker by workerId.
func (w Workers) Get(ID int32) (*IDWorker, error) {
	worker, ok := w[ID]
	if !ok {
		log.Warnf("id: %d not register", ID)
		return nil, fmt.Errorf("snowflake id: %d don't register in this service", ID)
	}
	return worker, nil
}
