package main

import (
	"fmt"
	"sync"
)

// Node LRU 节点
type Node struct {
	key   int
	value int
	prev  *Node
	next  *Node
}

// LRUCache 缓存结构体
type LRUCache struct {
	capacity int
	cache    map[int]*Node // 哈希表
	head     *Node         // 双向链表头（最近使用）
	tail     *Node         // 双向链表尾（最久未使用）
	mutex    sync.Mutex    // 并发安全
}

// NewLRUCache 新建 LRU 缓存
func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[int]*Node),
		head:     nil,
		tail:     nil,
	}
}

// 移动节点到头部（标记为最近使用）
func (l *LRUCache) moveToHead(node *Node) {
	if node == l.head {
		return
	}
	// 从链表中移除
	if node.prev != nil {
		node.prev.next = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	}
	if node == l.tail {
		l.tail = node.prev
	}
	// 移动到头部
	node.next = l.head
	node.prev = nil
	if l.head != nil {
		l.head.prev = node
	}
	l.head = node
	if l.tail == nil {
		l.tail = node
	}
}

// 添加新节点
func (l *LRUCache) addNode(node *Node) {
	node.prev = nil
	node.next = l.head
	if l.head != nil {
		l.head.prev = node
	}
	l.head = node
	if l.tail == nil {
		l.tail = node
	}
}

func (l *LRUCache) removeTail() *Node {
	if l.tail == nil {
		return nil
	}
	node := l.tail
	l.tail = node.prev
	if l.tail != nil {
		l.tail.next = nil
	} else {
		l.head = nil
	}
	return node
}

// Get 操作
func (l *LRUCache) Get(key int) int {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if node, ok := l.cache[key]; ok {
		l.moveToHead(node)
		return node.value
	}
	return -1
}

// Put 操作
func (l *LRUCache) Put(key int, value int) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if node, ok := l.cache[key]; ok {
		// 更新已有键
		node.value = value
		l.moveToHead(node)
		return
	}

	// 新增节点
	newNode := &Node{key: key, value: value}
	l.cache[key] = newNode
	l.addNode(newNode)

	// 超出容量，淘汰尾部
	if len(l.cache) > l.capacity {
		tail := l.removeTail()
		if tail != nil {
			delete(l.cache, tail.key)
		}
	}
}

func main() {
	cache := NewLRUCache(2)   // 容量为 2
	cache.Put(1, 1)           // 缓存: [1]
	cache.Put(2, 2)           // 缓存: [2, 1]
	fmt.Println(cache.Get(1)) // 输出: 1, 缓存: [1, 2]
	cache.Put(3, 3)           // 缓存满，淘汰 2, 缓存: [3, 1]
	fmt.Println(cache.Get(2)) // 输出: -1
	cache.Put(4, 4)           // 淘汰 1, 缓存: [4, 3]
	fmt.Println(cache.Get(1)) // 输出: -1
	fmt.Println(cache.Get(3)) // 输出: 3
}
