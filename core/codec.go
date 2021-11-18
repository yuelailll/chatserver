/*
 * @Date: 2021-11-12 17:58:02
 * @LastEditTime: 2021-11-18 14:35:11
 * @FilePath: \gnet_server\core\codec.go
 * @Description: 自定义编解码器
 */
package core

import (
	"container/list"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chatroom/codec"
	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pool/goroutine"
)

type CustomCodecServer struct {
	*gnet.EventServer
	addr       string
	multicore  bool
	async      bool
	codec      gnet.ICodec
	workerPool *goroutine.Pool
}

func (cs *CustomCodecServer) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	log.Printf("custom codec server is listening on %s (multi-cores: %t, loops: %d)\n",
		srv.Addr.String(), srv.Multicore, srv.NumEventLoop)
	return
}

func (cs *CustomCodecServer) React(frame []byte, conn gnet.Conn) (out []byte, action gnet.Action) {
	fmt.Println("into react: length of framePayload is ", len(frame))
	fmt.Println("into react: framePayload", string(frame))
	// packet decode
	var pkt codec.Packet
	pkt, err := codec.Decode(frame)
	if err != nil {
		return
	}

	switch pkt.(type) {
	case *codec.Submit:
		fmt.Println("into submit")
		var submitPkt = pkt.(*codec.Submit)
		var _id = submitPkt.ID
		var _payload = submitPkt.Payload

		fmt.Println(string(_payload))

		var data = make(map[string]interface{})
		if err := json.Unmarshal(_payload, &data); err != nil {
			fmt.Println(err.Error())
			return nil, gnet.Close
		}

		var _event_id = data["event_id"]
		var _message = data["data"].(map[string]interface{})

		if _, ok := _message["user_name"]; !ok {
			return nil, gnet.Close
		}

		if _event_id != "login" {
			if !Users.Contains(_message["user_name"].(string)) {
				return nil, gnet.Close
			}
		}

		var submitAckPkt = &codec.SubmitAck{
			ID:     _id,
			Result: 0,
		}

		if _event_id == "login" {
			var user User
			if ok := Users.Contains(_message["user_name"].(string)); !ok {
				user = User{
					Addr:      conn.RemoteAddr().String(),
					Name:      _message["user_name"].(string),
					Conn:      conn,
					LoginTime: time.Now().Unix(),
				}
				Users.Set(_message["user_name"].(string), user)
				AddrUsers.Set(conn.RemoteAddr().String(), user)
			} else {
				submitAckPkt.Result = 2
				submitAckPkt.Payload = []byte(`{"message":"user already exists"}`)
				_submitAckPkt, _ := codec.Encode(submitAckPkt)
				out = _submitAckPkt
				return
			}

			UserInRoom.Set(_message["user_name"].(string), "room0")
			ChatRooms.Get("room0").(*ChatRoom).Users.Set(_message["user_name"].(string), user)

			cs.workerPool.Submit(func() {
				var chatRoom = ChatRooms.Get("room0").(*ChatRoom)
				var messages = chatRoom.Messages

				messages.RLockFunc(func(_messages *list.List) {
					var length = _messages.Len()

					var res []map[string]interface{}
					var dataLength int

					if length > 50 {
						res = make([]map[string]interface{}, 0, 50)
						dataLength = 50
					} else {
						res = make([]map[string]interface{}, 0, length)
						dataLength = length
					}

					if length > 0 {
						for i, e := 0, _messages.Back(); i < dataLength; i, e = i+1, e.Prev() {
							var message = e.Value.(Message)
							var _tmp = map[string]interface{}{
								"from":      message.From,
								"timestamp": message.Timestamp,
								"message":   message.Payload,
							}
							res = append(res, _tmp)
						}
					}

					var payload = map[string]interface{}{
						"event_id": "message",
						"data":     res,
					}

					_payload, _ := json.Marshal(payload)
					submitAckPkt.Payload = _payload
					_submitAckPkt, _ := codec.Encode(submitAckPkt)
					conn.AsyncWrite(_submitAckPkt)
				})
			})
		} else if _event_id == "join" {
			var roomName = _message["room_name"].(string)
			if ok := ChatRooms.Contains(roomName); !ok {
				submitAckPkt.Result = 1
				_res, _ := codec.Encode(submitAckPkt)
				out = _res
				return
			}

			var userName = _message["user_name"].(string)
			var user = Users.Get(userName).(User)

			var nowRoom = UserInRoom.Get(userName).(string)

			UserInRoom.Set(userName, roomName)
			ChatRooms.Get(nowRoom).(*ChatRoom).Users.Remove(userName)
			ChatRooms.Get(roomName).(*ChatRoom).Users.Set(userName, user)

			var payload = map[string]interface{}{
				"event_id": "join",
				"data": map[string]string{
					"room_name": roomName,
					"result":    "success",
				},
			}

			_payload, _ := json.Marshal(payload)
			submitAckPkt.Payload = _payload

			_res, _ := codec.Encode(submitAckPkt)
			out = _res
		} else if _event_id == "message" {
			cs.workerPool.Submit(func() {
				var chatRoom = ChatRooms.Get(UserInRoom.Get(_message["user_name"].(string)).(string)).(*ChatRoom)

				if strings.HasPrefix(_message["message"].(string), "/stats") {
					var _tmp = strings.Split(_message["message"].(string), " ")
					var searchUser = _tmp[1]
					if !Users.Contains(searchUser) {
						submitAckPkt.Result = 1
						submitAckPkt.Payload = []byte("User not found")
						_res, _ := codec.Encode(submitAckPkt)
						conn.AsyncWrite(_res)
						return
					}
					var user = Users.Get(searchUser).(User)
					var userInRoom = UserInRoom.Get(user.Name).(string)
					var loginTime = user.LoginTime
					var onlineTime = time.Now().Unix() - loginTime

					var res = map[string]interface{}{
						"login_time":  loginTime,
						"online_time": onlineTime,
						"room_name":   userInRoom,
					}

					var payload = map[string]interface{}{
						"event_id": "stats",
						"data":     res,
					}

					_payload, _ := json.Marshal(payload)
					submitAckPkt.Payload = _payload

					_res, _ := codec.Encode(submitAckPkt)
					conn.AsyncWrite(_res)
					return
				}

				if strings.HasPrefix(_message["message"].(string), "/popular") {
					var _tmp = strings.Split(_message["message"].(string), " ")
					var roomName = _tmp[1]
					if ok := ChatRooms.Contains(roomName); !ok {
						submitAckPkt.Result = 1
						submitAckPkt.Payload = []byte("User not found")
						_res, _ := codec.Encode(submitAckPkt)
						conn.AsyncWrite(_res)
						return
					}

					var popular = GetPopular(roomName)

					var payload = map[string]interface{}{
						"event_id": "popular",
						"data":     popular,
					}

					_payload, _ := json.Marshal(payload)
					submitAckPkt.Payload = _payload
					_res, _ := codec.Encode(submitAckPkt)
					conn.AsyncWrite(_res)

					return
				}

				var _text = _message["message"].(string)
				var text = ProfanityWordsList.HandleWord(_text, '*')

				chatRoom.Mux.Lock()
				times, ok := chatRoom.WordsAni[text]
				if !ok {
					chatRoom.WordsAni[text] = 1
				} else {
					chatRoom.WordsAni[text] = times + 1
				}
				chatRoom.Mux.Unlock()

				var message = Message{
					Room:      string(chatRoom.Name),
					From:      _message["user_name"].(string),
					Timestamp: time.Now().Unix(),
					Payload:   text,
				}

				chatRoom.Messages.PushBack(message)

				var res = make([]map[string]interface{}, 0, 1)
				res = append(res, map[string]interface{}{
					"from":      message.From,
					"timestamp": message.Timestamp,
					"message":   message.Payload,
				})

				var payload = map[string]interface{}{
					"event_id": "message",
					"data":     res,
				}

				_payload, _ := json.Marshal(payload)
				submitAckPkt.Payload = _payload

				_res, _ := codec.Encode(submitAckPkt)

				var users = chatRoom.Users
				users.IteratorAsc(func(k interface{}, v interface{}) bool {
					var _userName = k.(string)
					if _userName != _message["user_name"] {
						var user = v.(User)
						var conn = user.Conn
						conn.AsyncWrite(_res)
					}
					return true
				})
			})
		}

		return
	default:
		return nil, gnet.Close
	}
}

func (cs *CustomCodecServer) OnOpened(conn gnet.Conn) (out []byte, action gnet.Action) {
	// fmt.Printf("connect addrss: %s", conn.RemoteAddr().String())
	// var submitAckPkt = &codec.SubmitAck{
	// 	ID:      "000000001",
	// 	Result:  1,
	// 	Payload: []byte(`{"name":"gnet_server","version":"1.0.0"}`),
	// }
	// _res, _ := codec.Encode(submitAckPkt)
	// out = _res
	return
}

func (cs *CustomCodecServer) OnClosed(conn gnet.Conn, err error) (action gnet.Action) {
	var addr = conn.RemoteAddr().String()

	if AddrUsers.Contains(addr) {
		var user = AddrUsers.Get(addr).(User)
		// 删除地址用户映射
		AddrUsers.Remove(addr)

		var userName = user.Name
		var userInRoom = UserInRoom.Get(userName).(string)
		// 删除用户聊天室映射
		UserInRoom.Remove(userName)

		chatRoom := ChatRooms.Get(userInRoom).(*ChatRoom)
		// 删除聊天室用户
		chatRoom.Users.Remove(userName)

		// 删除用户
		Users.Remove(userName)
	}
	return
}

func NewCustomCodecServer(addr string, multicore, async bool) {
	var err error

	var _codec = codec.Frame{}

	var cs = &CustomCodecServer{addr: addr, multicore: multicore, async: async, codec: _codec, workerPool: goroutine.Default()}

	err = gnet.Serve(cs, addr, gnet.WithMulticore(multicore), gnet.WithTCPKeepAlive(time.Minute*5), gnet.WithCodec(_codec))

	if err != nil {
		panic(err)
	}
}
