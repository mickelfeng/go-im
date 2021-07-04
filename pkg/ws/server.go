/**
  @author:panliang
  @data:2021/6/18
  @note
**/
package ws

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"strconv"
)

//客户端管理
type ClientManager struct {
	Clients    map[*Client]bool //存放ws长链接
	Broadcast  chan []byte      //收集消息分发到client
	Register   chan *Client     //新创建的长链接
	Unregister chan *Client     //注销的长链接
}

//客户端
type Client struct {
	ID     string          //客户端id
	Socket *websocket.Conn //长链接
	Send   chan []byte     //需要发送的消息
}

//消息结构体
type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

//发送的消息
type Msg struct {
	Code   int    `json:"code,omitempty"`
	FromId int    `json:"from_id,omitempty"`
	Msg    string `json:"msg,omitempty"`
	ToId   int    `json:"to_id,omitempty"`
	Status int    `json:"status,omitempty"`
}

//离线和上线消息
type OnlineMsg struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

//定义的一些状态码

const (
	connOut = 5000 //断开链接
	connOk  = 1000 //连接成功
	SendOk  = 200  //发送成功
)

// 创建客户端管理器
var Manager = ClientManager{
	Broadcast:  make(chan []byte),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
	Clients:    make(map[*Client]bool),
}

//启动websocket
func (manager *ClientManager) Start() {

	for {
		select {
		//如果有新的链接通过channel把链接传递给conn
		case conn := <-manager.Register:
			//将客户端的链接设置为true
			manager.Clients[conn] = true
			//把返回链接成功的消息json格式化
			jsonMessage, _ := json.Marshal(&OnlineMsg{Code: connOk, Msg: "用户上线啦"})
			//调用客户端方法发送消息
			manager.Send(jsonMessage, conn)
		case conn := <-manager.Unregister:
			//判断连接的状态，如果是true,就关闭send，删除连接client的值
			if _, ok := manager.Clients[conn]; ok {
				close(conn.Send)
				delete(manager.Clients, conn)
				jsonMessage, _ := json.Marshal(&OnlineMsg{Code: connOut, Msg: "用户离线了"})
				manager.Send(jsonMessage, conn)
			}
			//消息消费
		case message := <-manager.Broadcast:
			data := EnMessage(message)
			fmt.Println(data.Content)
			msg := new(Msg)
			err := json.Unmarshal([]byte(data.Content), &msg)
			if err != nil {
				fmt.Println(err)
			}
			jsonMessage_from, _ := json.Marshal(&Msg{Code: SendOk, Msg: msg.Msg, FromId: msg.FromId, ToId: msg.ToId, Status: 0})
			for conn := range manager.Clients {
				id, _ := strconv.Atoi(conn.ID)
				if id == msg.ToId {
					go PutData(msg)
					conn.Send <- jsonMessage_from
				}
			}
		}
	}
}

// 发送消息到客户端
func (manager *ClientManager) Send(message []byte, ignore *Client) {
	for conn := range manager.Clients {
		if conn != ignore {
			conn.Send <- message
		}
	}
}

//消息投递
func (c *Client) Read() {

	//关闭客户端注册 关闭socket连接
	defer func() {
		Manager.Unregister <- c
		c.Socket.Close()
	}()
	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			Manager.Unregister <- c
			c.Socket.Close()
			break
		}
		if string(message) == "HeartBeat" {
			c.Socket.WriteMessage(websocket.TextMessage, []byte(`{"code":0,"data":"heartbeat ok"}`))
			continue
		}
		//将数据投递到对应的客户端
		jsonMessage, _ := json.Marshal(&Message{Sender: c.ID, Content: string(message)})
		//manager.broadcast <- jsonMessage
		Manager.Broadcast <- jsonMessage
	}
}

//从客户端消费消息
func (c *Client) Write() {
	//关闭socket连接
	defer func() {
		c.Socket.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}