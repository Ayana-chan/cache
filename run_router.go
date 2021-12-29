package main

import (
	"cache/initchecker"
	"cache/server"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func main() {
	len := len(os.Args)
	addrs := os.Args[1 : len-2]
	port := ":" + os.Args[len-2][1:]
	checker := initchecker.Checker{Addrs: addrs}
	checker.Check()
	replicas, _ := strconv.Atoi(os.Args[len-1][1:])
	router := server.NewRouter(port, replicas)
	router.HashCircle.Add(addrs...)
	fmt.Println("[Router " + port + "] Node is running")
	fmt.Println("[Router " + port + "] Port" + port)
	fmt.Println("[Router " + port + "] Replicas:" + strconv.Itoa(replicas))
	http.ListenAndServe(port, router)
}
