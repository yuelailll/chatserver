/*
 * @Date: 2021-11-15 18:10:36
 * @LastEditTime: 2021-11-18 15:19:06
 * @FilePath: \gnet_server\core\init.go
 * @Description: 核心数据结构定义
 */
package core

import (
	"bufio"
	"fmt"
	"os"
	"sync"

	"github.com/chatroom/lib"
	"github.com/panjf2000/gnet"
)

// 用户结构
type User struct {
	Addr      string
	Name      string
	Conn      gnet.Conn
	LoginTime int64
}

// 用户列表
var Users *lib.AVLTree

// 连接地址对应用户
var AddrUsers *lib.AVLTree

// 聊天室结构
type ChatRoom struct {
	Mux      *sync.RWMutex
	Name     string
	Users    *lib.AVLTree
	Messages *lib.LinkList
	WordsAni map[string]int
}

// 聊天室列表
var ChatRooms *lib.AVLTree

var UserInRoom *lib.AVLTree

// 消息结构
type Message struct {
	Room      string
	From      string
	Timestamp int64
	Payload   string
}

var ProfanityWordsList *lib.DFAUtil

func init() {
	ChatRooms = lib.NewAVLTree()

	for i := 0; i < 10; i++ {
		ChatRooms.Set(fmt.Sprintf("room%d", i), &ChatRoom{
			Mux:      &sync.RWMutex{},
			Name:     fmt.Sprintf("room%d", i),
			Users:    lib.NewAVLTree(),
			Messages: lib.NewLinkList(),
			WordsAni: map[string]int{},
		})
	}

	Users = lib.NewAVLTree()
	AddrUsers = lib.NewAVLTree()

	UserInRoom = lib.NewAVLTree()

	initProfanityWords()
}

func initProfanityWords() {
	var profanityWords = make([]string, 0)

	var path = lib.GetAppPath() + "/static/profanity-words.txt"
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	var scanner = bufio.NewScanner(f)
	for scanner.Scan() {
		profanityWords = append(profanityWords, scanner.Text())
	}

	ProfanityWordsList = lib.NewDFAUtil(profanityWords)
}

func GetPopular(roomName string) string {
	var room = ChatRooms.Get(roomName).(*ChatRoom)
	room.Mux.RLock()
	defer room.Mux.RUnlock()

	var wordsAni = room.WordsAni

	var maxHead = lib.NewMaxHeap()

	for k, v := range wordsAni {
		maxHead.Enqueue(&lib.Node{
			Word:  string(k),
			Times: v,
		})
	}

	return maxHead.Dequeue().Word
}
