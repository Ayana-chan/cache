package server

import (
	"cache/cachestruct"
	"cache/lru/linkedlist"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/cache/"

type NodeServer struct {
	// 例如 "https://example.net:8000"
	self      string
	basePath  string
	mainCache cachestruct.Cache //缓存数据
}

// NewNodeServer 初始化一个NodeServer
func NewNodeServer(self string, cacheBytes int64) *NodeServer {
	return &NodeServer{
		self:     self,
		basePath: defaultBasePath,
		mainCache: cachestruct.Cache{
			CacheBytes: cacheBytes,
		},
	}
}

func (p *NodeServer) Log(format string, v ...interface{}) {
	log.Printf("[NodeServer %s] %s", p.self, fmt.Sprintf(format, v...))
}

// ServeHTTP handle all http requests
func (p *NodeServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//查询
	if r.Method == "GET" {
		if !strings.HasPrefix(r.URL.Path, p.basePath) {
			panic("HTTPPool serving unexpected path: " + r.URL.Path)
		}
		p.Log("%s %s", r.Method, r.URL.Path)
		// /<basepath>/<key>中取出key
		key := r.URL.Path[len(p.basePath):]
		if len(key) == 0 {
			http.Error(w, "bad request", http.StatusBadRequest) //key为空，返回400
			return
		}
		data, ok := p.mainCache.Get(key)
		if !ok {
			http.Error(w, "cache miss", 404) //未命中，返回404
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(data.ByteSlice())
		return
	} else if r.Method == "POST" { //设值
		if !strings.HasPrefix(r.URL.Path, p.basePath) {
			panic("HTTPPool serving unexpected path: " + r.URL.Path)
		}
		p.Log("%value %value", r.Method, r.URL.Path)
		// /<basepath>/<key>中取出key
		key := r.URL.Path[len(p.basePath):]
		if len(key) == 0 {
			http.Error(w, "bad request", http.StatusBadRequest) //key为空，返回400
			return
		}
		value, _ := ioutil.ReadAll(r.Body)
		p.Log("SET " + key + ":" + string(value))
		p.mainCache.Add(key, linkedlist.Data{
			B: value,
		})
	}
}
