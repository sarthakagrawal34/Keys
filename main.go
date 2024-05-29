package main

import (
	"fmt"
	"time"
)

type CacheBlock struct {
	key       int
	frequency int
	timestamp time.Time
}

type LFUCache struct {
	heap     []CacheBlock
	keyValue map[int]int
	capacity int
	size     int
}

func swap(heap []CacheBlock, i, j int) {
	temp := heap[i]
	heap[i] = heap[j]
	heap[i] = temp
}

func parent(i int) int {
	return (i - 1) / 2
}

func left(i int) int {
	return 2*i + 1
}

func right(i int) int {
	return 2*i + 2
}

func heapify(heap []CacheBlock, i, n int) {
	l := left(i)
	r := right(i)

	minim := i
	if l < n && (heap[l].frequency < heap[minim].frequency || (heap[l].frequency == heap[minim].frequency && heap[l].timestamp.Before(heap[minim].timestamp))) {
		minim = l
	}
	if r < n && (heap[r].frequency < heap[minim].frequency || (heap[r].frequency == heap[minim].frequency && heap[r].timestamp.Before(heap[minim].timestamp))) {
		minim = r
	}

	if minim != i {
		swap(heap, minim, i)
		heapify(heap, minim, n)
	}
}

func increment(cache *LFUCache, i int) {
	cache.heap[i].frequency++
	cache.heap[i].timestamp = time.Now()
	heapify(cache.heap, i, cache.size)
}

func insert(cache *LFUCache, key int, value int) {
	if cache.size == cache.capacity {
		removedKey := cache.heap[0].key
		delete(cache.keyValue, removedKey)
		cache.heap[0] = cache.heap[cache.size-1]
		cache.size--
		cache.heap = cache.heap[:cache.size] // Resize the slice
		heapify(cache.heap, 0, cache.size)
	}
	cache.heap = append(cache.heap, CacheBlock{
		key:       key,
		frequency: 1,
		timestamp: time.Now(),
	})
	cache.keyValue[key] = value
	i := cache.size
	cache.size++

	for i > 0 && (cache.heap[parent(i)].frequency > cache.heap[i].frequency || (cache.heap[parent(i)].frequency == cache.heap[i].frequency && cache.heap[parent(i)].timestamp.After(cache.heap[i].timestamp))) {
		swap(cache.heap, i, parent(i))
		i = parent(i)
	}
}

func update(cache *LFUCache, key, value int) {
	if _, found := cache.keyValue[key]; !found {
		insert(cache, key, value)
	} else {
		cache.keyValue[key] = value
		for i := 0; i < cache.size; i++ {
			if cache.heap[i].key == key {
				increment(cache, i)
				break
			}
		}
	}
}

func get(cache *LFUCache, key int) (int, error) {
	if value, found := cache.keyValue[key]; !found {
		return 0, fmt.Errorf("Not found")
	} else {
		for i := 0; i < cache.size; i++ {
			if cache.heap[i].key == key {
				increment(cache, i)
				break
			}
		}
		return value, nil
	}
}

func main() {
	cache := LFUCache{
		heap:     make([]CacheBlock, 0, 4),
		keyValue: make(map[int]int),
		capacity: 4,
		size:     0,
	}
	update(&cache, 1, 1)
	fmt.Println(get(&cache, 1))
	update(&cache, 2, 2)
	update(&cache, 3, 3)
	update(&cache, 4, 4)
	update(&cache, 5, 5)
	fmt.Println(get(&cache, 2))
}
