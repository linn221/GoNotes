package main

import (
	"encoding/json"
	"errors"
	"sync"
	"time"
)

// AI generated one for quick demo purpose
type inMemoryCache struct {
	mu   sync.Mutex
	data map[string]any
	sets map[string]map[string]struct{}
}

func NewInMemoryCache() *inMemoryCache {
	return &inMemoryCache{
		data: make(map[string]any),
		sets: make(map[string]map[string]struct{}),
	}
}

// --- Implementation of CacheService ---

func (c *inMemoryCache) GetObject(key string, dest any) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	val, ok := c.data[key]
	if !ok {
		return false, nil
	}
	bytes, err := json.Marshal(val)
	if err != nil {
		return false, err
	}
	err = json.Unmarshal(bytes, dest)
	return err == nil, err
}

func (c *inMemoryCache) GetValue(key string) (string, bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	val, ok := c.data[key]
	if !ok {
		return "", false, nil
	}
	strVal, ok := val.(string)
	if !ok {
		return "", false, errors.New("value is not a string")
	}
	return strVal, true, nil
}

func (c *inMemoryCache) SetObject(key string, obj any, _ time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = obj
	return nil
}

func (c *inMemoryCache) SetValue(key string, value string, _ time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
	return nil
}

func (c *inMemoryCache) RemoveKey(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
	return nil
}

func (c *inMemoryCache) RemoveKeyWithCount(key string) (int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.data[key]; ok {
		delete(c.data, key)
		return 1, nil
	}
	return 0, nil
}

func (c *inMemoryCache) RemoveKeysWithCount(keys []string) (int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	var count int64
	for _, key := range keys {
		if _, ok := c.data[key]; ok {
			delete(c.data, key)
			count++
		}
	}
	return count, nil
}

func (c *inMemoryCache) RemoveKeys(keys ...string) error {
	_, err := c.RemoveKeysWithCount(keys)
	return err
}

// --- Set operations ---

func (c *inMemoryCache) AddSet(setKey string, member string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.sets[setKey]; !ok {
		c.sets[setKey] = make(map[string]struct{})
	}
	c.sets[setKey][member] = struct{}{}
	return nil
}

func (c *inMemoryCache) GetSetMembers(setKey string) ([]string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	set, ok := c.sets[setKey]
	if !ok {
		return nil, nil
	}
	members := make([]string, 0, len(set))
	for member := range set {
		members = append(members, member)
	}
	return members, nil
}

func (c *inMemoryCache) RemoveSetMember(setKey string, member string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if set, ok := c.sets[setKey]; ok {
		delete(set, member)
		if len(set) == 0 {
			delete(c.sets, setKey)
		}
	}
	return nil
}
