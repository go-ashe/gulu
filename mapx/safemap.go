package mapx

import "sync"

// SafeMap 是一个并发安全的 map，支持泛型键和值
type SafeMap[K comparable, V any] struct {
	mu sync.RWMutex
	m  map[K]V
}

// NewSafeMap 创建一个新的 SafeMap 实例
func NewSafeMap[K comparable, V any]() *SafeMap[K, V] {
	return &SafeMap[K, V]{
		m: make(map[K]V),
	}
}

// Set 设置 map 中的键值对
func (sm *SafeMap[K, V]) Set(key K, value V) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.m[key] = value
}

// Get 获取 map 中的值，返回值和是否存在标志
func (sm *SafeMap[K, V]) Get(key K) (V, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	value, exists := sm.m[key]
	return value, exists
}

// Delete 删除 map 中的键值对
func (sm *SafeMap[K, V]) Delete(key K) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.m, key)
}

// Range 遍历 map 中的所有键值对，类似于 sync.Map 的 Range 方法
func (sm *SafeMap[K, V]) Range(f func(K, V) bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	for key, value := range sm.m {
		if !f(key, value) {
			break
		}
	}
}
