package pokecache

import (
	"sync"
	"time"
)

// get request로 받은 데이터들 저장할 cache 구조체
type Cache struct {
	entries map[string]cacheEntry
	mu      *sync.RWMutex
	// current string // 현재 페이지의 url(entries의 key) 저장
}

// 데이터 entry 한개 저장할 구조체
type cacheEntry struct {
	createdAt time.Time // 이 entry의 생성 시점 저장
	val       []byte
}

// 새 Cache 구조체를 반환하는 함수
func NewCache(interval time.Duration) Cache {
	newC := Cache{
		entries: make(map[string]cacheEntry), // map[string]cacheEntry{} 대신 make 사용해보기
		mu:      &sync.RWMutex{},
	}
	ticker := time.Tick(interval) // interval마다 time.Time 반환하는 read only 채널 (<-chan time.Time)
	// @@@ 해답은 time.NewTicker(interval) 사용 => Ticker구조체의 C field가 time.Tick(interval) 와 동일
	go newC.reapLoop(ticker, interval) // go루틴은 go에 걸린 함수가 return 하거나 main이 return 할때 killed

	return newC
}

// cache에 새 entry 입력하는 method
func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

// cache에서 entry를 찾아 반환하는 method
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.entries[key]
	// if ok {
	// 	return entry.val, true
	// } else {
	// 	return []byte{}, false
	// }
	// @@@ 해답 코드 반영: ok가 false일때(map에 해당 키가 없을 때) value 타입의 zero 값을 반환하므로 그대로 사용해도 된다

	return entry.val, ok
}

// 주어진 시간간격마다 오래된 entry들을 삭제하는 함수
func (c *Cache) reapLoop(ticker <-chan time.Time, interval time.Duration) {
	for timestamp := range ticker { // ticker 채널에 일정간격마다 time.Time이 들어올때 for 루프 한단계씩 진행되고 그 사이는 block
		c.mu.Lock()
		for key, entry := range c.entries {
			if timestamp.Sub(entry.createdAt) > interval {
				// @@@ 해답 코드 비교:
				// entry.createdAt.Before(timestamp.Add(-interval))
				// func (t time.Time) Before(u time.Time) bool
				// Before reports whether the time instant t is before u
				// if key == c.current { // 현재 페이지인 경우 삭제 제외
				// 	continue
				// }
				delete(c.entries, key)
			}
		}
		c.mu.Unlock()
	}
}

// // Cache의 current 필드 설정 메소드
// func (c *Cache) SetCurrent(url string) {
// 	c.current = url
// }

// // Cache의 current 필드 접근 메소드
// func (c *Cache) GetCurrent() string {
// 	return c.current
// }
