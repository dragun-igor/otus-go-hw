package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	length    int
	firstItem *ListItem
	lastItem  *ListItem
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.firstItem
}

func (l *list) Back() *ListItem {
	return l.lastItem
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := new(ListItem)
	newItem.Value = v
	if l.length == 0 {
		l.lastItem = newItem
	} else {
		newItem.Next = l.firstItem
		l.firstItem.Prev = newItem
	}
	l.firstItem = newItem
	l.length++
	return l.firstItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := new(ListItem)
	newItem.Value = v
	if l.length == 0 {
		l.firstItem = newItem
	} else {
		newItem.Prev = l.lastItem
		l.lastItem.Next = newItem
	}
	l.lastItem = newItem
	l.length++
	return l.lastItem
}

func (l *list) Remove(i *ListItem) {
	defer func() {
		l.length--
	}()
	if i.Prev == nil && i.Next == nil {
		l.firstItem = nil
		l.lastItem = nil
		return
	}
	if i.Prev == nil {
		l.firstItem = i.Next
		l.firstItem.Prev = nil
		return
	}
	if i.Next == nil {
		l.lastItem = i.Prev
		l.lastItem.Next = nil
		return
	}
	if i.Prev != nil || i.Next != nil {
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
		return
	}
}

func (l *list) MoveToFront(i *ListItem) {
	if l.length > 1 && i.Prev != nil {
		l.lastItem = i.Prev
		l.firstItem.Prev = i
		i.Prev.Next = i.Next
		i.Next = l.firstItem
		i.Prev = nil
		l.firstItem = i
	}
}
