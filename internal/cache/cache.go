package cache

import (
	"encoding/gob"
	"github.com/jaeha-choi/DFF/internal/datatype"
	"os"
	"time"
)

// CAPACITY is the max allowed number of champions to hold
const CAPACITY int = 16

// EXPIRATION data expiration time in days
const EXPIRATION = 7

type GameMode int
type Position int

const (
	DEFAULT GameMode = iota
	ARAM
	URF
)

const (
	TOP Position = iota
	JUNGLE
	MID
	ADC
	SUPPORT
)

type Cache struct {
	Head     *Node
	Tail     *Node
	Capacity int
	Existing map[string]*Node
}

type Node struct {
	Key     string
	Next    *Node
	Prev    *Node
	URF     *CachedData
	ARAM    *CachedData
	Default []*CachedData
}

type CachedData struct {
	CreationTime time.Time
	Version      string
	Spells       datatype.Spells
	RunePages    datatype.RunePages
	ItemPages    []datatype.ItemPage
}

// NewCache create new cache
func NewCache() *Cache {
	head := &Node{}
	tail := &Node{}

	head.Next = tail
	tail.Prev = head

	return &Cache{
		Head:     head,
		Tail:     tail,
		Existing: make(map[string]*Node, CAPACITY),
		Capacity: CAPACITY,
	}
}

// RestoreCache restore saved cache
func RestoreCache(filename string) (cache *Cache, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	err = gob.NewDecoder(file).Decode(&cache)

	if cache != nil {
		for len(cache.Existing) >= CAPACITY {
			cache.delLast()
		}
	}

	return
}

// SaveCache save cache
func (c *Cache) SaveCache(filename string) (err error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	err = gob.NewEncoder(file).Encode(&c)
	return
}

// Get retrieves the cached data
func (c *Cache) Get(name string, mode GameMode, lane Position) (data *CachedData) {
	var node *Node
	var exist bool

	if node, exist = c.Existing[name]; !exist || node == nil {
		return nil
	}

	if mode == URF {
		data = node.URF
	} else if mode == ARAM {
		data = node.ARAM
	} else if mode == DEFAULT {
		data = node.Default[lane]
	}

	if data != nil {
		// If expiration date passed, remove data
		if t := time.Now().Sub(data.CreationTime); t >= time.Hour*24*EXPIRATION {
			if mode == URF {
				node.URF = nil
			} else if mode == ARAM {
				node.ARAM = nil
			} else if mode == DEFAULT {
				node.Default[lane] = nil
			}
			data = nil
		}
	}

	// Move the node to the front
	node.Prev.Next = node.Next
	node.Next.Prev = node.Prev

	node.Prev = c.Head
	node.Next = c.Head.Next

	c.Head.Next.Prev = node
	c.Head.Next = node

	return
}

// Put updates the cache
func (c *Cache) Put(name string, mode GameMode, lane Position, data *CachedData) {
	var node *Node
	var exist bool

	if len(c.Existing) >= c.Capacity {
		c.delLast()
	}

	if node, exist = c.Existing[name]; !exist || node == nil {
		node = &Node{
			Key:     name,
			URF:     nil,
			ARAM:    nil,
			Default: make([]*CachedData, 5),
			Next:    nil,
			Prev:    nil,
		}
		c.Existing[name] = node
	} else {
		// If already exist, remove from the linked list before adding to the front
		node.Prev.Next = node.Next
		node.Next.Prev = node.Prev
	}

	// Add/Move the node to the front
	node.Prev = c.Head
	node.Next = c.Head.Next

	c.Head.Next.Prev = node
	c.Head.Next = node

	if mode == URF {
		node.URF = data
	} else if mode == ARAM {
		node.ARAM = data
	} else if mode == DEFAULT {
		node.Default[lane] = data
	}
}

// delLast deletes the last node in the cache (excluding head/tail)
func (c *Cache) delLast() {
	if len(c.Existing) > 0 {
		delete(c.Existing, c.Tail.Prev.Key)
		c.Tail.Prev.Prev.Next = c.Tail
		c.Tail.Prev = c.Tail.Prev.Prev
	}
}
