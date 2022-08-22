package geeCache

import (
	"net/http"
	"strings"
)

const default_base_url = "_geecache"

type HttpPool struct {
	basepath string
}

func NewHttpPool() *HttpPool {
	return &HttpPool{
		basepath: default_base_url,
	}
}

//ServeHTTP handle all http requests
func (p *HttpPool) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !strings.HasPrefix(req.URL.Path, p.basepath) {
		panic("HttpPool sering unexpected path: " + req.URL.Path)
	}

	//expected path : host/basepath/groupName/key
	params := strings.Split(req.URL.Path[len(p.basepath):], "/")
	if len(params) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	groupName := params[0]
	key := params[1]
	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "No such cache group: "+groupName, http.StatusNotFound)
	}
	val, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(val.ByteSlice())
}
