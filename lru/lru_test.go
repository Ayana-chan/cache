package lru

import (
	"cache/lru/linkedlist"
	"fmt"
	"testing"
)

type String string

func (s String) Length() int {
	return len(s)
}

func TestGet(t *testing.T) {
	c := New(3, func(key string, val linkedlist.Value) {
		fmt.Println("delet " + key)
	})
	head := c.list.Head
	c.Add("key1", String("val1"))
	printList(head, *c)
	c.Add("key2", String("val2"))
	printList(head, *c)
	c.Add("key3", String("val3"))
	printList(head, *c)
	c.Add("key4", String("val4"))
	printList(head, *c)
	c.Add("key5", String("val5"))
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
