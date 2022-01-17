/*
* 1                                               42           52             64
* +-----------------------------------------------+------------+---------------+
* | timestamp(ms)                                 | workerID   | sequence      |
* +-----------------------------------------------+------------+---------------+
* | 0000000000 0000000000 0000000000 0000000000 0 | 0000000000 | 0000000000 00 |
* +-----------------------------------------------+------------+---------------+
*
* 1. 41位时间截(毫秒级)，注意这是时间截的差值（当前时间截 - 开始时间截)。可以使用约70年: (1L << 41) / (1000L * 60 * 60 * 24 * 365) = 69
* 2. 10位数据机器位，可以部署在1024个节点
* 3. 12位序列，毫秒内的计数，同一机器，同一时间截并发4096个序号
 */

package snowflake

import (
	"sync"
	"time"
)

const (
	twepoch        = int64(1483228800000)             //开始时间截 (2017-01-01)
	workerIdBits   = uint(10)                         //机器id所占的位数
	sequenceBits   = uint(12)                         //序列所占的位数
	workerIdMax    = int64(-1 ^ (-1 << workerIdBits)) //支持的最大机器id数量
	sequenceMask   = int64(-1 ^ (-1 << sequenceBits)) //
	workerIdShift  = sequenceBits                     //机器id左移位数
	timestampShift = sequenceBits + workerIdBits      //时间戳左移位数
)

// A Snowflake struct holds the basic information needed for a snowflake generator worker
type Snowflake struct {
	sync.Mutex
	timestamp int64
	workerID  int64
	sequence  int64
}

// NewSnowflake returns a new snowflake worker that can be used to generate snowflake IDs
func NewSnowflake(workerID int64) *Snowflake {

	if workerID < 0 || workerID > workerIdMax {
		workerID = 101
	}

	return &Snowflake{
		timestamp: 0,
		workerID:  workerID,
		sequence:  0,
	}
}

// Generate creates and returns a unique snowflake ID
func (s *Snowflake) Generate() int64 {
	s.Lock()
	now := time.Now().UnixNano() / 1000000

	if s.timestamp == now {
		s.sequence = (s.sequence + 1) & sequenceMask

		if s.sequence == 0 {
			for now <= s.timestamp {
				now = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		s.sequence = 0
	}

	s.timestamp = now

	r := (now-twepoch)<<timestampShift | (s.workerID << workerIdShift) | (s.sequence)

	s.Unlock()
	return r
}
