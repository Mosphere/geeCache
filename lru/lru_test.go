package lru

import (
	"testing"
)

type Str string

func (s Str) Len() int {
	return len(s)
}

func TestAdd(t *testing.T) {
	cache := New(0, nil)
	cases := []entry{{"k1", Str("val1")}, {"k2", Str("val2")}}
	for _, c := range cases {
		cache.Add(c.key, c.value)
	}
}
func TestGet(t *testing.T) {
	cache := New(0, nil)

	cache.Add("k1", Str("123"))
	if val, ok := cache.Get("k1"); !ok || string(val.(Str)) != "123" {
		t.Fatalf("cache hit k1=123 failed.")
	}
}

func TestRemove(t *testing.T) {

}
