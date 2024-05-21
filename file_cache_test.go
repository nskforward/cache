package cache

import (
	"bytes"
	"io"
	"math/rand"
	"testing"
	"time"
)

func TestFileCache(t *testing.T) {
	c := NewFileCache("/Users/17847869/go/src/github.com/nskforward/cache/build")
	err := c.FindOrCreate("1").Write(bytes.NewReader([]byte("111")))
	if err != nil {
		t.Fatalf("cannot write key '1': %s", err)
	}
	entry := c.Find("1")
	if entry == nil {
		t.Fatalf("key '1' must exists")
	}
	if entry.Size() != 3 {
		t.Fatalf("key '1' size must be 3, actual: %d", entry.Size())
	}
	var buf bytes.Buffer
	err = entry.Read(&buf)
	if err != nil {
		t.Fatalf("cannot read key '1': %s", err)
	}
	if buf.String() != "111" {
		t.Fatalf("key '1' read value must be '111', actual: '%s'", buf.String())
	}
}

func BenchmarkFileCache(b *testing.B) {
	c := NewFileCache("/Users/17847869/go/src/github.com/nskforward/cache/build")
	keys := [][2]string{
		{"1", "111"},
		{"2", "222"},
		{"3", "333"},
		{"4", "444"},
		{"5", "555"},
		{"6", "666"},
		{"7", "777"},
		{"8", "888"},
		{"9", "999"},
	}
	for _, pair := range keys {
		entry := c.FindOrCreate(pair[0])
		entry.Write(bytes.NewReader([]byte(pair[1])))
		c.RemoveOnExpiration(pair[0], entry.SetTTL(2*time.Second))
	}

	for i := 0; i < b.N; i++ {
		pair := keys[rand.Intn(len(keys))]
		entry := c.Find(pair[0])
		if entry == nil {
			b.Fatalf("cannot find key '%s'", pair[0])
		}
		err := entry.Read(io.Discard)
		if err != nil {
			b.Fatal(err)
		}
		entry.SetTTL(2 * time.Second)
	}
}
