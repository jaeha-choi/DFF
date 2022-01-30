package cache

import (
	"bytes"
	"encoding/gob"
	"github.com/jaeha-choi/DFF/internal/datatype"
	"os"
	"time"
)

// CAPACITY is the max allowed number of champions to hold
const CAPACITY int = 16

// EXPIRATION data expiration time in days
const EXPIRATION = 7

type Position int

const (
	NONE Position = iota - 1
	TOP
	JUNGLE
	MID
	ADC
	SUPPORT
)

type Cache struct {
	Head     *Node
	Tail     *Node
	Capacity int
	Size     int
	Existing map[string]*Node
}

type Node struct {
	Next  *Node
	Prev  *Node
	Value *NodeValue
}

type NodeValue struct {
	Key                string
	URF                CachedData
	ARAM               CachedData
	Default            []CachedData
	AvailablePositions []Position
	DefaultPosition    Position
}

type CachedData struct {
	CreationTime     time.Time
	Version          string
	PositionPickRate string
	URL              string

	Spells    datatype.Spells
	RunePages []datatype.DFFRunePage
	ItemPages datatype.ItemPage
}

// NewCache create new cache
func NewCache() *Cache {
	head := &Node{}
	tail := &Node{}

	head.Next = tail
	tail.Prev = head

	return &Cache{
		Size:     0,
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

func (c *Cache) GetPutNode(name string) (node *Node, isCached bool) {
	var exist bool

	if node, exist = c.Existing[name]; !exist || node == nil {
		node = &Node{
			Value: &NodeValue{
				Key:     name,
				URF:     CachedData{},
				ARAM:    CachedData{},
				Default: make([]CachedData, 5),
			},
		}
		c.Existing[name] = node
		c.Size++
	} else {
		// If already exist, remove from the linked list before adding to the front
		node.Prev.Next = node.Next
		node.Next.Prev = node.Prev
		isCached = true
	}

	// Move the node to the front
	node.Prev = c.Head
	node.Next = c.Head.Next

	c.Head.Next.Prev = node
	c.Head.Next = node

	return
}

func (c *Cache) GetPut(name string, mode datatype.GameMode, position Position) (data *CachedData, isCached bool) {
	var node *Node
	node, isCached = c.GetPutNode(name)

	if mode == datatype.URF {
		data = &node.Value.URF
	} else if mode == datatype.ARAM {
		data = &node.Value.ARAM
	} else if mode == datatype.DEFAULT {
		if position == NONE {
			position = node.Value.DefaultPosition
		}
		data = &node.Value.Default[position]
	}

	if data != nil {
		// TODO: Compare game version vs cache version to perform version check
		// If expiration date passed, remove data
		if t := time.Now().Sub(data.CreationTime); t >= time.Hour*24*EXPIRATION {
			if mode == datatype.URF {
				node.Value.URF = CachedData{}
				data = &node.Value.URF
			} else if mode == datatype.ARAM {
				node.Value.ARAM = CachedData{}
				data = &node.Value.ARAM
			} else if mode == datatype.DEFAULT {
				node.Value.Default[position] = CachedData{}
				data = &node.Value.Default[position]
			}
			isCached = false
		}
		if data.RunePages == nil {
			isCached = false
		}
	}
	//fmt.Println("Using cached data: ", isCached)

	return
}

// delLast deletes the last node in the cache (excluding head/tail)
func (c *Cache) delLast() {
	if len(c.Existing) > 0 {
		delete(c.Existing, c.Tail.Prev.Value.Key)
		c.Tail.Prev.Prev.Next = c.Tail
		c.Tail.Prev = c.Tail.Prev.Prev
		c.Size--
	}
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (c Cache) MarshalBinary() (data []byte, err error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	if err = encoder.Encode(&c.Capacity); err != nil {
		return nil, err
	}

	if err = encoder.Encode(&c.Size); err != nil {
		return nil, err
	}

	curr := c.Head.Next
	for i := 0; i < c.Size; i++ {
		if err = encoder.Encode(&curr.Value); err != nil {
			return nil, err
		}
		curr = curr.Next
	}

	return buf.Bytes(), err
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (c *Cache) UnmarshalBinary(data []byte) (err error) {
	decoder := gob.NewDecoder(bytes.NewReader(data))

	if err = decoder.Decode(&c.Capacity); err != nil {
		return err
	}

	if err = decoder.Decode(&c.Size); err != nil {
		return err
	}

	c.Head = &Node{}
	c.Tail = &Node{}

	c.Existing = make(map[string]*Node, CAPACITY)

	curr := c.Head
	prev := curr
	for i := 0; i < c.Size; i++ {
		curr.Next = &Node{}
		curr = curr.Next
		curr.Prev = prev
		prev = curr
		if err = decoder.Decode(&curr.Value); err != nil {
			return err
		}
		c.Existing[curr.Value.Key] = curr
	}

	curr.Next = c.Tail
	c.Tail.Prev = curr

	return err
}

// String implements the fmt.Stringer interface.
func (c *Cache) String() (str string) {
	curr := c.Head.Next
	for i := 0; i < c.Size; i++ {
		str += curr.Value.Key + "\t"
		curr = curr.Next
	}

	return
}
