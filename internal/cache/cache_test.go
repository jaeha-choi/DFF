package cache

import (
	"bytes"
	"encoding/gob"
	"github.com/jaeha-choi/DFF/internal/datatype"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	c := NewCache("version")

	c.GetPut(0, datatype.Default, Adc)
	c.GetPut(1, datatype.Aram, Top)
	c.GetPut(0, datatype.Urf, Support)

	if c.String() != "0\t1\t" {
		t.Error("Incorrect result for TestCache")
	}
}

func TestEncode(t *testing.T) {
	c := NewCache("version")

	c.GetPut(0, datatype.Default, Mid)
	c.GetPut(1, datatype.Aram, Jungle)
	c.GetPut(0, datatype.Urf, Adc)

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(&c); err != nil {
		t.Error(err)
	}
}

func TestEncodeDecode(t *testing.T) {
	// Encode
	c := NewCache("version")

	c.GetPut(0, datatype.Default, Adc)
	c.GetPut(1, datatype.Default, Top)
	c.GetPut(0, datatype.Urf, Support)
	cache, isCached := c.GetPut(3, datatype.Default, Mid)
	if isCached {
		t.Error("Incorrect result for TestCache")
	}
	cache.CreationTime = time.Now()
	cache.URL = "someURL"
	c.GetPut(2, datatype.Default, Jungle)

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
	_, isCached = decoded.GetPut(0, datatype.Default, Mid)
	if isCached {
		t.Error("Incorrect result for TestCache")
	}

	if decoded.Existing[3].Value.Default[Mid].URL != "someURL" {
		t.Error("Incorrect result for TestCache")
	}

	// Get operation affects node order in the linked list
	cache, isCached = decoded.GetPut(3, datatype.Default, Mid)
	if isCached && cache.URL != "someURL" {
		t.Error("Incorrect result for TestCache")
	}

	if decoded.String() != "3\t0\t2\t1\t" {
		t.Error("Incorrect result for TestCache")
	}
}
