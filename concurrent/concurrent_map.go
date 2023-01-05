package concurrent

import (
	"encoding/json"
	"sync"
)

var ShardCount = 32

type Map[V any] []*MapShared[V]

type MapShared[V any] struct {
	items map[string]V
	sync.RWMutex
}

func New[V any]() Map[V] {
	m := make(Map[V], ShardCount)
	for i := 0; i < ShardCount; i++ {
		m[i] = &MapShared[V]{items: make(map[string]V)}
	}
	return m
}

func (m Map[V]) GetShard(key string) *MapShared[V] {
	return m[uint(fnv32(key))%uint(ShardCount)]
}

func (m Map[V]) MSet(data map[string]V) {
	for key, value := range data {
		shard := m.GetShard(key)
		shard.Lock()
		shard.items[key] = value
		shard.Unlock()
	}
}

func (m Map[V]) Set(key string, value V) {
	// Get map shard.
	shard := m.GetShard(key)
	shard.Lock()
	shard.items[key] = value
	shard.Unlock()
}

type UpsertCb[V any] func(exist bool, valueInMap V, newValue V) V

func (m Map[V]) Upsert(key string, value V, cb UpsertCb[V]) (res V) {
	shard := m.GetShard(key)
	shard.Lock()
	v, ok := shard.items[key]
	res = cb(ok, v, value)
	shard.items[key] = res
	shard.Unlock()
	return res
}

func (m Map[V]) SetIfAbsent(key string, value V) bool {
	shard := m.GetShard(key)
	shard.Lock()
	_, ok := shard.items[key]
	if !ok {
		shard.items[key] = value
	}
	shard.Unlock()
	return !ok
}

func (m Map[V]) Get(key string) (V, bool) {
	shard := m.GetShard(key)
	shard.RLock()
	val, ok := shard.items[key]
	shard.RUnlock()
	return val, ok
}

func (m Map[V]) Count() int {
	count := 0
	for i := 0; i < ShardCount; i++ {
		shard := m[i]
		shard.RLock()
		count += len(shard.items)
		shard.RUnlock()
	}
	return count
}

func (m Map[V]) Has(key string) bool {
	shard := m.GetShard(key)
	shard.RLock()
	_, ok := shard.items[key]
	shard.RUnlock()
	return ok
}

func (m Map[V]) Remove(key string) {
	shard := m.GetShard(key)
	shard.Lock()
	delete(shard.items, key)
	shard.Unlock()
}

type RemoveCb[V any] func(key string, v V, exists bool) bool

func (m Map[V]) RemoveCb(key string, cb RemoveCb[V]) bool {
	shard := m.GetShard(key)
	shard.Lock()
	v, ok := shard.items[key]
	remove := cb(key, v, ok)
	if remove && ok {
		delete(shard.items, key)
	}
	shard.Unlock()
	return remove
}

func (m Map[V]) Pop(key string) (v V, exists bool) {
	shard := m.GetShard(key)
	shard.Lock()
	v, exists = shard.items[key]
	delete(shard.items, key)
	shard.Unlock()
	return v, exists
}

func (m Map[V]) IsEmpty() bool {
	return m.Count() == 0
}

type Tuple[V any] struct {
	Key string
	Val V
}

func (m Map[V]) Iter() <-chan Tuple[V] {
	cans := snapshot(m)
	ch := make(chan Tuple[V])
	go fanIn(cans, ch)
	return ch
}

func (m Map[V]) IterBuffered() <-chan Tuple[V] {
	cans := snapshot(m)
	total := 0
	for _, c := range cans {
		total += cap(c)
	}
	ch := make(chan Tuple[V], total)
	go fanIn(cans, ch)
	return ch
}

func (m Map[V]) Clear() {
	for item := range m.IterBuffered() {
		m.Remove(item.Key)
	}
}

func snapshot[V any](m Map[V]) (cans []chan Tuple[V]) {
	if len(m) == 0 {
		panic(`ConcurrentMap is not initialized. Should run New() before usage.`)
	}
	cans = make([]chan Tuple[V], ShardCount)
	wg := sync.WaitGroup{}
	wg.Add(ShardCount)
	for index, shard := range m {
		go func(index int, shard *MapShared[V]) {
			shard.RLock()
			cans[index] = make(chan Tuple[V], len(shard.items))
			wg.Done()
			for key, val := range shard.items {
				cans[index] <- Tuple[V]{key, val}
			}
			shard.RUnlock()
			close(cans[index])
		}(index, shard)
	}
	wg.Wait()
	return cans
}

func fanIn[V any](cans []chan Tuple[V], out chan Tuple[V]) {
	wg := sync.WaitGroup{}
	wg.Add(len(cans))
	for _, ch := range cans {
		go func(ch chan Tuple[V]) {
			for t := range ch {
				out <- t
			}
			wg.Done()
		}(ch)
	}
	wg.Wait()
	close(out)
}

func (m Map[V]) Items() map[string]V {
	tmp := make(map[string]V)
	
	for item := range m.IterBuffered() {
		tmp[item.Key] = item.Val
	}
	
	return tmp
}

type IterCb[V any] func(key string, v V)

func (m Map[V]) IterCb(fn IterCb[V]) {
	for idx := range m {
		shard := (m)[idx]
		shard.RLock()
		for key, value := range shard.items {
			fn(key, value)
		}
		shard.RUnlock()
	}
}

func (m Map[V]) Keys() []string {
	count := m.Count()
	ch := make(chan string, count)
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(ShardCount)
		for _, shard := range m {
			go func(shard *MapShared[V]) {
				shard.RLock()
				for key := range shard.items {
					ch <- key
				}
				shard.RUnlock()
				wg.Done()
			}(shard)
		}
		wg.Wait()
		close(ch)
	}()
	
	keys := make([]string, 0, count)
	for k := range ch {
		keys = append(keys, k)
	}
	return keys
}

func (m Map[V]) MarshalJSON() ([]byte, error) {
	tmp := make(map[string]V)
	
	for item := range m.IterBuffered() {
		tmp[item.Key] = item.Val
	}
	return json.Marshal(tmp)
}

func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	keyLength := len(key)
	for i := 0; i < keyLength; i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}
