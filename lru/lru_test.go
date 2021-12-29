package lru

import (
	"cache/lru/linkedlist"
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	data := linkedlist.Data{B: []byte("hello")}
	fmt.Println(data)
	c := New(3, func(key string, val linkedlist.Data) {
		fmt.Println("delete " + key)
	})
	head := c.list.Head
	c.Add("key1", linkedlist.Data{B: []byte("val1")})
	printList(head, *c)
	c.Add("key2", linkedlist.Data{B: []byte("val2")})
	printList(head, *c)
	c.Add("key3", linkedlist.Data{B: []byte("val3")})
	printList(head, *c)
	c.Add("key4", linkedlist.Data{B: []byte("val4")})
	printList(head, *c)
	c.Add("key5", linkedlist.Data{B: []byte("val5")})
	printList(head, *c)
	if value, ok := c.Get("key4"); ok {
		fmt.Println(value)
	} else {
		fmt.Println("null")
	}
	printList(head, *c)
	if value, ok := c.Get("key1"); ok {
		fmt.Println(value)
	} else {
		fmt.Println("null")
	}
}

func printList(node *linkedlist.Element, c Cache) {
	fmt.Print("capacity=")
	fmt.Print(c.GetCapacity())
	fmt.Print("  ")
	for node != nil {
		fmt.Print(node.Key + "->")
		node = node.Next
	}
	fmt.Println()
}
