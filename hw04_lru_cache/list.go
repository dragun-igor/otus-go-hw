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
	Value interface{} // Значение
	Next  *ListItem   // Указатель на следующий элемент
	Prev  *ListItem   // Указатель на предыдущий элемент
}

type list struct {
	length    int       // Длина списка
	firstItem *ListItem // Первый элемент
	lastItem  *ListItem // Последний элемент
}

func NewList() List {
	// Создаём новый список
	return new(list)
}

func (l list) Len() int {
	// Возвращаем длину
	return l.length
}

func (l list) Front() *ListItem {
	// Возвращаем указатель на первый элемент
	return l.firstItem
}

func (l list) Back() *ListItem {
	// Возвращаем указатель на последний элемент
	return l.lastItem
}

func (l *list) PushFront(v interface{}) *ListItem {
	// Добавление нового элемента в начало списка
	newItem := &ListItem{Value: v}
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
	// Добавление нового элемента в конец списка
	newItem := &ListItem{Value: v}
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
	// Уменьшаем длину списка на 1
	defer func() {
		l.length--
	}()
	// Если элемент был единственным, то убираем указатели на первый и последний элемент
	if i.Prev == nil && i.Next == nil {
		l.firstItem = nil
		l.lastItem = nil
		return
	}
	// Если элемент был первым
	if i.Prev == nil {
		l.firstItem = i.Next
		l.firstItem.Prev = nil
		return
	}
	// Если элемент был последним
	if i.Next == nil {
		l.lastItem = i.Prev
		l.lastItem.Next = nil
		return
	}
	// Если элемент был ни первый, ни последний
	if i.Prev != nil || i.Next != nil {
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
		return
	}
}

func (l *list) MoveToFront(i *ListItem) {
	// Перемещение элемента в начало очереди
	if l.length > 1 && i.Prev != nil {
		l.lastItem = i.Prev
		l.firstItem.Prev = i
		i.Prev.Next = i.Next
		i.Next = l.firstItem
		i.Prev = nil
		l.firstItem = i
	}
}
