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
	len   int
	front *ListItem
	back  *ListItem
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := &ListItem{
		Value: v,
	}
	if l.front == nil {
		l.front = newItem
		l.back = newItem
	} else {
		newItem.Next = l.front
		l.front.Prev = newItem
		l.front = newItem
	}
	l.len++

	return newItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := &ListItem{
		Value: v,
	}
	if l.back == nil {
		l.back = newItem
		l.front = newItem
	} else {
		newItem.Prev = l.back
		l.back.Next = newItem
		l.back = newItem
	}
	l.len++

	return newItem
}
func (l *list) Remove(i *ListItem) {
	defer func() {
		l.len--
	}()
	if i.Next == nil {
		i.Prev.Next = nil
		l.back = i.Prev
		return
	}
	switch i.Prev == nil {
	case true:
		i.Next.Prev = nil
		l.front = i.Next
	case false:
		i.Next.Prev = i.Prev
		i.Prev.Next = i.Next
	}
}

func (l *list) MoveToFront(i *ListItem) {
	if i.Prev == nil {
		return
	}

	switch i.Next == nil {
	case true:
		i.Prev.Next = nil
		l.back = i.Prev
		i.Next = l.front
		l.front.Prev = i
		i.Prev = nil
		l.front = i
	case false:
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
		i.Next = l.front
		l.front.Prev = i
		i.Prev = nil
		l.front = i
	}
}
