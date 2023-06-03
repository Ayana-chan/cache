package server

import (
	"bytes"
	"cache/consistenthash"
	"cache/lru/linkedlist"
	"cache/singleflight"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Router struct {
	self       string // 自己的端口号，例如 ":8000"
	basePath   string
	HashCircle *consistenthash.Map //哈希环
	loader     *singleflight.Group
}

// NewRouter 初始化一个Router
func NewRouter(self string, replicas int) *Router {
	return &Router{
		self:       self,
		basePath:   defaultBasePath,
		HashCircle: consistenthash.New(replicas, nil),
		loader:     &singleflight.Group{},
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
		value, statusCode := router.loader.Do(key, func() (interface{}, int) {
			return router.getRoute(key)
		})

		if statusCode == 404 || statusCode == 999 {

			//未命中，返回404
			http.Error(w, "cache miss", 404)

			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(value.(*linkedlist.Data).ByteSlice())

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
		resp, err := http.Post("http://"+nodeAddr+defaultBasePath+key, "raw", bytes.NewReader(value))
		defer resp.Body.Close()

		//发现节点宕机,删除该节点，并不断尝试设值直到成功将值设到正常的节点中，如果节点全部失效，关闭当前Router
		for err != nil && len(router.HashCircle.Keys) != 0 {
			router.Log("Node DOWN: " + nodeAddr)
			router.HashCircle.Delete(nodeAddr)
			nodeAddr = router.HashCircle.Get(key)
			resp, err = doPost("http://"+nodeAddr+defaultBasePath+key, "raw", bytes.NewReader(value))
		}

		//全部节点均失效
		if len(router.HashCircle.Keys) == 0 {
			router.Log("No node running")
			router.Log("router DOWN")
			http.Error(w, "router DOWN", http.StatusServiceUnavailable)
			os.Exit(-1)
		}

		router.Log("SET " + key + ":" + string(value))
	} else if r.Method == "DELETE" {
		router.Log("%s %s", r.Method, r.URL.Path)

		// /<basepath>/<key>中取出key
		key := r.URL.Path[len(router.basePath):]

		if len(key) == 0 {

			//key为空，返回400
			http.Error(w, "bad request", http.StatusBadRequest)

			return
		}
		nodeAddr := router.HashCircle.Get(key)
		req, _ := http.NewRequest("DELETE", "http://"+nodeAddr+defaultBasePath+key, nil)
		res, err := http.DefaultClient.Do(req)
		defer res.Body.Close()
		if err != nil {
			router.Log("Node down: " + nodeAddr)
			router.HashCircle.Delete(nodeAddr)
		}
	}
}

// 查询对应节点
func (router *Router) getRoute(key string) (*linkedlist.Data, int) {
	nodeAddr := router.HashCircle.Get(key)
	resp, err := http.Get("http://" + nodeAddr + defaultBasePath + key)

	//发现节点宕机，删除该节点
	if err != nil {
		router.Log("Node down: " + nodeAddr)
		router.HashCircle.Delete(nodeAddr)
		return nil, 999
	}

	defer resp.Body.Close()
	statusCode := resp.StatusCode
	if statusCode == 404 {
		return nil, 404
	}
	value, _ := ioutil.ReadAll(resp.Body)
	return &linkedlist.Data{B: value}, 200
}

// 为了关闭尝试节点是否有效的请求中的所有Body，创建此函数
func doPost(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	resp, err = http.Post(url, contentType, body)
	defer resp.Body.Close()
	return resp, err
}
