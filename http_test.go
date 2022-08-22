package geeCache

import (
	"fmt"
	"net/http"
	"testing"
)

func TestHttp(t *testing.T) {
	var db = map[string]string{
		"Tom":  "630",
		"Jack": "589",
		"Sam":  "567",
	}

	//生成groups
	NewGroup("scores", 2<<10, GetterFunc(func(key string) ([]byte, error) {
		if val, ok := db[key]; ok {
			return []byte(val), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}))
	pool := NewHttpPool()

	t.Fatal(http.ListenAndServe("localhost:999", pool))
}
