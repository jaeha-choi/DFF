package cache

import (
	"bytes"
	"encoding/gob"
	"github.com/jaeha-choi/DFF/internal/datatype"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	c := NewCache()

	c.GetPut("A", datatype.DEFAULT, ADC)
	c.GetPut("B", datatype.ARAM, TOP)
	c.GetPut("A", datatype.URF, SUPPORT)

	if c.String() != "A\tB\t" {
		t.Error("Incorrect result for TestCache")
	}
}

func TestEncode(t *testing.T) {
	c := NewCache()

	c.GetPut("A", datatype.DEFAULT, MID)
	c.GetPut("B", datatype.ARAM, JUNGLE)
	c.GetPut("A", datatype.URF, ADC)

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(&c); err != nil {
		t.Error(err)
	}
}

func TestEncodeDecode(t *testing.T) {
	// Encode
	c := NewCache()

	c.GetPut("A", datatype.DEFAULT, ADC)
	c.GetPut("B", datatype.DEFAULT, TOP)
	c.GetPut("A", datatype.URF, SUPPORT)
	cache, isCached := c.GetPut("D", datatype.DEFAULT, MID)
	if isCached {
		t.Error("Incorrect result for TestCache")
	}
	cache.CreationTime = time.Now()
	cache.Version = "VersionD"
	c.GetPut("C", datatype.DEFAULT, JUNGLE)

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
	// Because CreationTime is not set, the following line should return false for isCached
	_, isCached = decoded.GetPut("A", datatype.DEFAULT, MID)
	if isCached {
		t.Error("Incorrect result for TestCache")
	}

	if decoded.Existing["D"].Value.Default[MID].Version != "VersionD" {
		t.Error("Incorrect result for TestCache")
	}

	// Get operation affects node order in the linked list
	cache, isCached = decoded.GetPut("D", datatype.DEFAULT, MID)
	if isCached && cache.Version != "VersionD" {
		t.Error("Incorrect result for TestCache")
	}

	if decoded.String() != "D\tA\tC\tB\t" {
		t.Error("Incorrect result for TestCache")
	}
}
