package main

import (
	"cache/server"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func main() {
	addr := os.Args[1][1:]
	capacity := os.Args[2][1:]
	cap, _ := strconv.ParseInt(capacity, 10, 64)
	fmt.Println("Node is running")
	fmt.Println("Address:" + addr)
	fmt.Println("Capacity:" + capacity)
	http.ListenAndServe(addr, server.NewNodeServer("localhost:9999", cap))
}
