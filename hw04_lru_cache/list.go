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
	var result *ListItem

	if l.front != nil {
		result = &ListItem{Next: l.front, Value: v}
		l.front.Prev = result
	} else {
		result = &ListItem{Value: v}
	}

	l.front = result

	if l.back == nil {
		l.back = result
	}

	l.len++

	return result
}

func (l *list) PushBack(v interface{}) *ListItem {
	var result *ListItem

	if l.back != nil {
		result = &ListItem{Prev: l.back, Value: v}
		l.back.Next = result
	} else {
		result = &ListItem{Value: v}
	}

	l.back = result

	if l.front == nil {
		l.front = result
	}

	l.len++

	return result
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}

func (l *list) Remove(i *ListItem) {
	switch {
	case i.Prev != nil && i.Next != nil:
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	case i.Prev != nil && i.Next == nil:
		i.Prev.Next = nil
		l.back = i.Prev
	case i.Next != nil && i.Prev == nil:
		i.Next.Prev = nil
		l.front = i.Next
	}

	l.len--
}

func NewList() List {
	return new(list)
}
