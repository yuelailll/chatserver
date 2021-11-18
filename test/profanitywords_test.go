/*
 * @Date: 2021-11-18 15:58:20
 * @LastEditTime: 2021-11-18 15:59:36
 * @FilePath: \gnet_server\test\profanitywords_test.go
 * @Description: 敏感词测试
 */
package main

import (
	"testing"

	"github.com/chatroom/lib"
)

func TestIsMatch(t *testing.T) {
	sensitiveList := []string{"测试", "新的测试"}
	input := "这是一个新的测试"

	util := lib.NewDFAUtil(sensitiveList)
	if util.IsMatch(input) == false {
		t.Errorf("Expected true, but got false")
	}
}

func TestHandleWord(t *testing.T) {
	sensitiveList := []string{"测试", "新的测试"}
	input := "这是一个新的测试"

	util := lib.NewDFAUtil(sensitiveList)
	newInput := util.HandleWord(input, '*')
	expected := "这是一个****"
	if newInput != expected {
		t.Errorf("Expected %s, but got %s", expected, newInput)
	}
}
