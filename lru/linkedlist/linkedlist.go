package linkedlist

type Element struct {
	Key   string
	Value Data
	Pre   *Element
	Next  *Element
}

type LinkedList struct {
	Head *Element
	Tail *Element
}

func New() *LinkedList {
	head := Element{
		Key:   "Head",
		Value: Data{},
	}
	tail := Element{
		Key:   "Tail",
		Value: Data{},
	}
	head.Next = &tail
	tail.Pre = &head
	list := LinkedList{
		Head: &head,
		Tail: &tail,
	}
	return &list
}

func (list *LinkedList) Remove(element *Element) {
	next := element.Next
	pre := element.Pre
	next.Pre = pre
	pre.Next = next
}

func (list *LinkedList) AddToHead(element *Element) {
	head := list.Head
	headNext := head.Next
	head.Next = element
	headNext.Pre = element
	element.Next = headNext
	element.Pre = head
}
