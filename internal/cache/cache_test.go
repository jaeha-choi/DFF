package cache

import (
	"bytes"
	"encoding/gob"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	c := NewCache()

	c.Put("A", DEFAULT, ADC, nil)
	c.Put("B", ARAM, TOP, nil)
	c.Put("A", URF, SUPPORT, nil)

	if c.String() != "A\tB\t" {
		t.Error("Incorrect result for TestCache")
	}
}

func TestEncode(t *testing.T) {
	c := NewCache()

	c.Put("A", DEFAULT, MID, &CachedData{})
	c.Put("B", DEFAULT, JUNGLE, &CachedData{})
	c.Put("A", URF, ADC, &CachedData{})

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(&c); err != nil {
		t.Error(err)
	}
}

func TestEncodeDecode(t *testing.T) {
	// Encode
	c := NewCache()

	c.Put("A", DEFAULT, ADC, &CachedData{})
	c.Put("B", DEFAULT, TOP, &CachedData{})
	c.Put("A", URF, SUPPORT, &CachedData{})
	c.Put("D", DEFAULT, MID, &CachedData{
		CreationTime: time.Now(),
		Version:      "VersionD",
	})
	c.Put("C", DEFAULT, JUNGLE, &CachedData{})

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(&c); err != nil {
		t.Error(err)
	}

	// Decode
	var decoded Cache
	if err := gob.NewDecoder(&buf).Decode(&decoded); err != nil {
		t.Error(err)
	}

	// Get operation affects node order in the linked list
	// Because CreationTime is not set, the following line should return <nil>
	if decoded.Get("A", DEFAULT, MID) != nil {
		t.Error("Incorrect result for TestCache")
	}

	if decoded.Existing["D"].Value.Default[MID].Version != "VersionD" {
		t.Error("Incorrect result for TestCache")
	}

	// Get operation affects node order in the linked list
	if decoded.Get("D", DEFAULT, MID).Version != "VersionD" {
		t.Error("Incorrect result for TestCache")
	}

	if decoded.String() != "D\tA\tC\tB\t" {
		t.Error("Incorrect result for TestCache")
	}
}
