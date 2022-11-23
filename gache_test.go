package gache

import (
	"testing"
	"time"
)

func Test(t *testing.T) {
	c := NewCacheWithExpire(time.Second)
	c.Set("a", "b")

	v, ok := c.Get("a")

	if !ok {
		t.Fatal("!ok1")
	}
	var str string
	if str, ok = v.(string); !ok {
		t.Fatal("!ok2")
	}

	if str != "b" {
		t.Fatal(`str != "b"`)
	}

	time.Sleep(time.Second)

	t.Log(c.Get("a"))
}
