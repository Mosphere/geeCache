package geeCache

import (
	"fmt"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

// A GetterFunc implements Getter with a function
type GetterFunc func(key string) ([]byte, error)

// Get implements Getter interface function
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

//函数类型实现某一个接口，称之为接口型函数，
//方便使用者调用时既能够传入函数作为参数，
//又能够传入实现了该接口的结构体作为参数

//定义一个函数类型F,
//并且实现接口A的方法，在这个方法中调用自己。
//这是将函数（参数返回值和F一致）转化为接口A的常用方法

type Group struct {
	name   string
	cache  cache
	getter Getter
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		return nil
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:   name,
		cache:  cache{cacheBytes: cacheBytes},
		getter: getter,
	}
	groups[name] = g
	return g
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if val, ok := g.cache.Get(key); ok {
		return val, nil
	}

	//缓存不存在
	return g.load(key)
}

func (g *Group) load(key string) (ByteView, error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, nil
	}

	val := ByteView{b: bytes}
	g.populateCache(key, val)
	return val, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.cache.Add(key, value)
}

func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}
