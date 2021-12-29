package initchecker

import (
	"log"
	"net/http"
	"os"
)

// Checker 用来在router启动时检查其节点地址是否有效
type Checker struct {
	Addrs []string
}

func (c *Checker) Check() {
	for _, addr := range c.Addrs {
		_, err := http.Get("http://" + addr + "/cache/test")
		if err != nil {
			log.Println("[Router init check] [Error] " + "invalid address:" + addr)
			os.Exit(-1)
		} else {
			log.Println("[Router init check] " + "valid address:" + addr)
		}
	}
}
