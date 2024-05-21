package cache

import (
	"math/rand"
	"testing"
	"time"
)

func TestMemCache(t *testing.T) {
	c := NewMemCache()
	c.FindOrCreate("1").SetValue("111").SetTTL(200 * time.Millisecond)

	entry := c.Find("1")
	if entry == nil {
		t.Fatalf("key '1' must be existed")
	}
	if entry.GetValue() != "111" {
		t.Fatalf("key '1' must contain value '111', actual '%s'", entry.GetValue())
	}

	time.Sleep(300 * time.Millisecond)
	entry = c.Find("1")
	if entry == nil {
		t.Fatalf("key '1' must be deleted")
	}
	entry.SetValue("222").SetTTL(200 * time.Millisecond)

	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		entry = c.Find("1")
		if entry == nil {
			t.Fatalf("key '1' not prolonged")
		}
		if entry.GetValue() != "222" {
			t.Fatalf("prolonged key '1' must contain value '222', actual '%s'", entry.GetValue())
		}
		entry.SetTTL(100 * time.Millisecond)
	}

	time.Sleep(300 * time.Millisecond)
	entry = c.Find("1")
	if entry != nil {
		t.Fatalf("key '1' must be deleted")
	}
}

func BenchmarkMemCache(b *testing.B) {
	c := NewMemCache()
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
		entry := c.FindOrCreate(pair[0]).SetValue(pair[1])
		c.RemoveOnExpiration(pair[0], entry.SetTTL(100*time.Millisecond))
	}
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			pair := keys[rand.Intn(len(keys))]
			entry := c.Find(pair[0])
			if entry == nil {
				b.Fatalf("cannot find key '%s'", pair[0])
			}
			if entry.GetValue() != pair[1] {
				b.Fatalf("value must be '%s', actual '%s'", pair[1], entry.GetValue())
			}
			entry.SetTTL(100 * time.Millisecond)
		}
	})
}
