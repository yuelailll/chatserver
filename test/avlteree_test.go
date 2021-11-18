/*
 * @Date: 2021-11-18 14:12:50
 * @LastEditTime: 2021-11-18 14:41:51
 * @FilePath: \gnet_server\test\avlteree_test.go
 * @Description:平衡二叉树测试
 */
package main

import (
	"testing"

	"github.com/chatroom/lib"
)

func TestAVLTree(t *testing.T) {
	var tree = lib.NewAVLTree()

	tree.Set("a", 1)

	if !tree.Contains("a") {
		t.Errorf("key should exists:%s", "a")
	}

	if tree.Contains("b") {
		t.Errorf("key should not exists:%s", "b")
	}

	var value = tree.Get("a")
	if value == nil {
		t.Errorf("key %s value should not be nil", "a")
	}
	data, ok := value.(int)
	if !ok {
		t.Errorf("key %s value should not be int", "a")
	}
	if data != 1 {
		t.Errorf("key %s value should not be 1", "a")
	}

	value = tree.Get("b")
	if value != nil {
		t.Errorf("key %s value should be nil", "b")
	}

	tree.Remove("a")
	value = tree.Get("a")
	if value != nil {
		t.Errorf("remove err; key %s value should be nil", "a")
	}

	tree.Set("a", 1)
	tree.Set("b", 2)

	tree.IteratorAsc(func(k, v interface{}) bool {
		if k.(string) == "a" && v.(int) != 1 {
			t.Errorf("key %s value %d", k.(string), v.(int))
		}
		return true
	})
}
