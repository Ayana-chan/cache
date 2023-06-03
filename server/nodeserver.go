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
	self      string // 自己的端口号，例如 ":8000"
	basePath  string
	mainCache cachestruct.Cache //缓存数据
}

// NewNodeServer 初始化一个NodeServer
func NewNodeServer(self string, cacheBytes int) *NodeServer {
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

// ServeHTTP 处理来自router的请求
func (p *NodeServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	//跨域
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("Access-Control-Allow-Methods", "*")            //允许所有方法

	//查询
	if r.Method == "GET" {

		//检查请求base地址是否正确
		if !strings.HasPrefix(r.URL.Path, p.basePath) {
			panic("HTTPPool serving unexpected path: " + r.URL.Path)
		}
		p.Log("%s %s", r.Method, r.URL.Path)

		// /<basepath>/<key>中取出key
		key := r.URL.Path[len(p.basePath):]

		if len(key) == 0 {

			//key为空，返回400
			http.Error(w, "bad request", http.StatusBadRequest)

			return
		}
		data, ok := p.mainCache.Get(key)
		if !ok {

			//未命中，返回404
			http.Error(w, "cache miss", 404)

			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(data.ByteSlice())
		return

		//设值
	} else if r.Method == "POST" {
		if !strings.HasPrefix(r.URL.Path, p.basePath) {
			panic("HTTPPool serving unexpected path: " + r.URL.Path)
		}
		p.Log("%value %value", r.Method, r.URL.Path)

		// 在/<basepath>/<key>中取出key
		key := r.URL.Path[len(p.basePath):]

		if len(key) == 0 {

			//key为空，返回400
			http.Error(w, "bad request", http.StatusBadRequest)

			return
		}
		value, _ := ioutil.ReadAll(r.Body)
		p.Log("SET " + key + ":" + string(value))
		p.mainCache.Add(key, linkedlist.Data{
			B: value,
		})
	} else if r.Method == "DELETE" {
		if !strings.HasPrefix(r.URL.Path, p.basePath) {
			panic("HTTPPool serving unexpected path: " + r.URL.Path)
		}
		p.Log("%value %value", r.Method, r.URL.Path)
		key := r.URL.Path[len(p.basePath):]
		if len(key) == 0 {

			//key为空，返回400
			http.Error(w, "bad request", http.StatusBadRequest)

			return
		}
		p.mainCache.Delete(key)
	}
}
