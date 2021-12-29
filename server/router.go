package server

import (
	"bytes"
	"cache/consistenthash"
	"cache/lru/linkedlist"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Router struct {
	self       string // 自己的端口号，例如 ":8000"
	basePath   string
	HashCircle *consistenthash.Map //哈希环
}

// NewRouter 初始化一个Router
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

	//检查请求base地址是否正确
	if !strings.HasPrefix(r.URL.Path, router.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}

	//查询
	if r.Method == "GET" {
		router.Log("%s %s", r.Method, r.URL.Path)

		// /<basepath>/<key>中取出key
		key := r.URL.Path[len(router.basePath):]

		if len(key) == 0 {

			//key为空，返回400
			http.Error(w, "bad request", http.StatusBadRequest)

			return
		}
		value, statusCode := router.getRoute(key)
		if statusCode == 404 {

			//未命中，返回404
			http.Error(w, "cache miss", 404)

			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(value.ByteSlice())

		//设值
	} else if r.Method == "POST" {
		router.Log("%s %s", r.Method, r.URL.Path)

		// /<basepath>/<key>中取出key
		key := r.URL.Path[len(router.basePath):]

		if len(key) == 0 {

			//key为空，返回400
			http.Error(w, "bad request", http.StatusBadRequest)

			return
		}
		nodeAddr := router.HashCircle.Get(key)
		value, _ := ioutil.ReadAll(r.Body)

		//复制value的reader
		reader := bytes.NewReader(value)

		resp, _ := http.Post("http://"+nodeAddr+defaultBasePath+key, "raw", reader)
		defer resp.Body.Close()
		router.Log("SET " + key + ":" + string(value))
	}
}

//查询对应节点
func (router *Router) getRoute(key string) (*linkedlist.Data, int) {
	nodeAddr := router.HashCircle.Get(key)
	resp, _ := http.Get("http://" + nodeAddr + defaultBasePath + key)
	defer resp.Body.Close()
	statusCode := resp.StatusCode
	if statusCode == 404 {
		return nil, 404
	}
	value, _ := ioutil.ReadAll(resp.Body)
	return &linkedlist.Data{B: value}, 200
}
