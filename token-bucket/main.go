// token bucket
package main

import (
	"fmt"
	"sync"
	"time"
)

type TokenBucket struct {
	ticker     *time.Ticker
	cap        int64
	avail      int64
	tokenMutex sync.Mutex
}

// 产生令牌
func (tb *TokenBucket) adjustTbDaemon() {
	tb.tokenMutex.Lock() // 临界、共享资源存取都需要加锁
	defer tb.tokenMutex.Unlock()

	for _ = range tb.ticker.C {
		if tb.avail < tb.cap {
			tb.avail++
		}
	}
}

func New(interval time.Duration, cap int64) *TokenBucket {
	tb := &TokenBucket{
		ticker: time.NewTicker(interval),
		cap:    cap,
	}

	go tb.adjustTbDaemon()
	return tb
}

// 消费令牌
func (tb *TokenBucket) TryTake(cnt int64) bool {
	tb.tokenMutex.Lock()
	defer tb.tokenMutex.Unlock()

	if cnt <= tb.avail {
		tb.avail -= cnt
		return true
	}
	return false
}

func main() {
	tb := New(time.Millisecond*200, 10)

	time.Sleep(time.Second)
	fmt.Println(tb.TryTake(5))
}
