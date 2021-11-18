/*
 * @Date: 2021-11-18 14:43:12
 * @LastEditTime: 2021-11-18 16:12:26
 * @FilePath: \gnet_server\lib\linklist.go
 * @Description: 双向链表
 */
package lib

import (
	"bytes"
	"container/list"
	"sync"
)

type LinkList struct {
	mu   sync.RWMutex
	list *list.List
}

func NewLinkList() *LinkList {
	return &LinkList{
		mu:   sync.RWMutex{},
		list: list.New(),
	}
}

type Element = list.Element

func (l *LinkList) PushFront(v interface{}) (e *Element) {
	l.mu.Lock()
	if l.list == nil {
		l.list = list.New()
	}
	e = l.list.PushFront(v)
	l.mu.Unlock()
	return
}

// PushBack inserts a new element <e> with value <v> at the back of list <l> and returns <e>.
func (l *LinkList) PushBack(v interface{}) (e *Element) {
	l.mu.Lock()
	if l.list == nil {
		l.list = list.New()
	}
	e = l.list.PushBack(v)
	l.mu.Unlock()
	return
}

// PushFronts inserts multiple new elements with values <values> at the front of list <l>.
func (l *LinkList) PushFronts(values []interface{}) {
	l.mu.Lock()
	if l.list == nil {
		l.list = list.New()
	}
	for _, v := range values {
		l.list.PushFront(v)
	}
	l.mu.Unlock()
}

// PushBacks inserts multiple new elements with values <values> at the back of list <l>.
func (l *LinkList) PushBacks(values []interface{}) {
	l.mu.Lock()
	if l.list == nil {
		l.list = list.New()
	}
	for _, v := range values {
		l.list.PushBack(v)
	}
	l.mu.Unlock()
}

// PopBack removes the element from back of <l> and returns the value of the element.
func (l *LinkList) PopBack() (value interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
		return
	}
	if e := l.list.Back(); e != nil {
		value = l.list.Remove(e)
	}
	return
}

// PopFront removes the element from front of <l> and returns the value of the element.
func (l *LinkList) PopFront() (value interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
		return
	}
	if e := l.list.Front(); e != nil {
		value = l.list.Remove(e)
	}
	return
}

// PopBacks removes <max> elements from back of <l>
// and returns values of the removed elements as slice.
func (l *LinkList) PopBacks(max int) (values []interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
		return
	}
	length := l.list.Len()
	if length > 0 {
		if max > 0 && max < length {
			length = max
		}
		values = make([]interface{}, length)
		for i := 0; i < length; i++ {
			values[i] = l.list.Remove(l.list.Back())
		}
	}
	return
}

// PopFronts removes <max> elements from front of <l>
// and returns values of the removed elements as slice.
func (l *LinkList) PopFronts(max int) (values []interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
		return
	}
	length := l.list.Len()
	if length > 0 {
		if max > 0 && max < length {
			length = max
		}
		values = make([]interface{}, length)
		for i := 0; i < length; i++ {
			values[i] = l.list.Remove(l.list.Front())
		}
	}
	return
}

// Len returns the number of elements of list <l>.
// The complexity is O(1).
func (l *LinkList) Len() (length int) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return
	}
	length = l.list.Len()
	return
}

// Size is alias of Len.
func (l *LinkList) Size() int {
	return l.Len()
}

// RemoveAll removes all elements from list <l>.
func (l *LinkList) RemoveAll() {
	l.mu.Lock()
	l.list = list.New()
	l.mu.Unlock()
}

// See RemoveAll().
func (l *LinkList) Clear() {
	l.RemoveAll()
}

// RLockFunc locks reading with given callback function <f> within RWMutex.RLock.
func (l *LinkList) RLockFunc(f func(list *list.List)) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list != nil {
		f(l.list)
	}
}

// LockFunc locks writing with given callback function <f> within RWMutex.Lock.
func (l *LinkList) LockFunc(f func(list *list.List)) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.list == nil {
		l.list = list.New()
	}
	f(l.list)
}

// Join joins list elements with a string <glue>.
func (l *LinkList) Join(glue string) string {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.list == nil {
		return ""
	}
	buffer := bytes.NewBuffer(nil)
	length := l.list.Len()
	if length > 0 {
		for i, e := 0, l.list.Front(); i < length; i, e = i+1, e.Next() {
			buffer.WriteString(String(e.Value))
			if i != length-1 {
				buffer.WriteString(glue)
			}
		}
	}
	return buffer.String()
}

// String returns current list as a string.
func (l *LinkList) String() string {
	return "[" + l.Join(",") + "]"
}
