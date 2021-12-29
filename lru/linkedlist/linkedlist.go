package linkedlist

// Element 队列元素，包含Data用来存储数据
type Element struct {
	Key   string
	Value Data
	Pre   *Element
	Next  *Element
}

// LinkedList 一个双向队列
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

// Remove 删除一个元素
func (list *LinkedList) Remove(element *Element) {
	next := element.Next
	pre := element.Pre
	next.Pre = pre
	pre.Next = next
}

// AddToHead 将一个元素加到队列头部
func (list *LinkedList) AddToHead(element *Element) {
	head := list.Head
	headNext := head.Next
	head.Next = element
	headNext.Pre = element
	element.Next = headNext
	element.Pre = head
}
