/*
 * @Date: 2021-11-18 15:20:02
 * @LastEditTime: 2021-11-18 15:29:34
 * @FilePath: \gnet_server\test\linklist_test.go
 * @Description: 链表单元测试
 */
package main

import (
	"container/list"
	"fmt"
	"testing"

	"github.com/chatroom/lib"
)

func TestLinkList(t *testing.T) {
	var linkList = lib.NewLinkList()
	linkList.PushBack(1)
	if linkList.Size() != 1 {
		t.Errorf("test1: push err")
	}

	linkList.PushFront(2)
	if linkList.Size() != 2 {
		t.Errorf("test1: push err")
	}

	var val = linkList.PopBack()
	if val == nil {
		t.Errorf("test3: pop err")
	}
	if val.(int) != 1 {
		t.Errorf("test4: pop err")
	}

	val = linkList.PopFront()
	if val == nil {
		t.Errorf("test5: pop err")
	}
	if val.(int) != 2 {
		t.Errorf("test6: pop err")
	}

	linkList.PushBack(1)
	linkList.PushBack(2)
	linkList.PushBack(3)
	linkList.PushBack(4)
	linkList.PushBack(5)
	linkList.PushBack(6)

	linkList.RLockFunc(func(list *list.List) {
		var length = list.Len()

		if length > 0 {
			for e := list.Back(); e != nil; e = e.Prev() {
				var v = e.Value.(int)
				fmt.Println(v)
			}
		}
	})

	linkList.Clear()
	if linkList.Size() != 0 {
		t.Errorf("test7: clear err")
	}
}
