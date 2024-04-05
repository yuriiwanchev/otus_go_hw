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
	length int
	first  *ListItem
	last   *ListItem
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.first
}

func (l *list) Back() *ListItem {
	return l.last
}

func (l *list) PushFront(v interface{}) *ListItem {
	elementToPush := &ListItem{
		Value: v,
		Next:  l.first,
		Prev:  nil,
	}
	if l.length == 0 {
		l.first = elementToPush
		l.last = elementToPush
		l.length++
		return elementToPush
	}

	if l.length > 0 {
		l.first.Prev = elementToPush
		elementToPush.Next = l.first
		l.first = elementToPush
		l.length++
		return elementToPush
	}

	return l.first
}

func (l *list) PushBack(v interface{}) *ListItem {
	elementToPush := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  l.last,
	}
	if l.length == 0 {
		l.first = elementToPush
		l.last = elementToPush
		l.length++
		return elementToPush
	}

	if l.length > 0 {
		l.last.Next = elementToPush
		elementToPush.Prev = l.last
		l.last = elementToPush
		l.length++
		return elementToPush
	}

	return l.first
}

func (l *list) Remove(itemToRemove *ListItem) {
	switch {
	case itemToRemove.Prev == nil:
		l.first = l.first.Next
		l.first.Prev = nil
	case itemToRemove.Next == nil:
		l.last = l.last.Prev
		l.last.Next = nil
	default:
		itemToRemove.Prev.Next = itemToRemove.Next
		itemToRemove.Next.Prev = itemToRemove.Prev
	}

	l.length--
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}
