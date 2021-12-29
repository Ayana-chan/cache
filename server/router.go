package server

import (
	"cache/consistenthash"
	"cache/lru/linkedlist"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Router struct {
	self       string
	basePath   string
	HashCircle *consistenthash.Map
}

func NewRouter(self string, replicas int) *Router {
	return &Router{
		self:       self,
		basePath:   defaultBasePath,
		HashCircle: consistenthash.New(replicas, nil),
	}
}

func (router *Router) Log(format string, v ...interface{}) {
	log.Printf("[Router %s] %s", router.self, fmt.Sprintf(format, v...))
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, router.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	if r.Method == "GET" {
		router.Log("%s %s", r.Method, r.URL.Path)
		// /<basepath>/<key>中取出key
		key := r.URL.Path[len(router.basePath):]
		if len(key) == 0 {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		value, statusCode := router.getRoute(key)
		if statusCode == 404 {
			http.Error(w, "cache miss", 404) //未命中，返回404
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(value.ByteSlice())
	} else if r.Method == "POST" {
		router.Log("%s %s", r.Method, r.URL.Path)
		// /<basepath>/<key>中取出key
		key := r.URL.Path[len(router.basePath):]
		if len(key) == 0 {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		nodeAddr := router.HashCircle.Get(key)
		resp, _ := http.Post(nodeAddr+defaultBasePath+key, "raw", r.Body)
		defer resp.Body.Close()
		value, _ := ioutil.ReadAll(r.Body)
		router.Log("SET " + key + ":" + string(value))
	}
}

func (router *Router) getRoute(key string) (*linkedlist.Data, int) {
	nodeAddr := router.HashCircle.Get(key)
	resp, _ := http.Get(nodeAddr + defaultBasePath + key)
	defer resp.Body.Close()
	statusCode := resp.StatusCode
	if statusCode == 404 {
		return nil, 404
	}
	value, _ := ioutil.ReadAll(resp.Body)
	return &linkedlist.Data{B: value}, 200
}
