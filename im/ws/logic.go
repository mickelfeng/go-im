/**
  @author:panliang
  @data:2021/9/4
  @note
**/
package ws

import (
	"encoding/json"
	"fmt"
	messageModel "go_im/im/http/models/msg"
	"go_im/im/http/models/user"
	"go_im/pkg/helpler"
	"go_im/pkg/model"
	"strconv"
	"time"
)

//group message insert db
func PutGroupData(msg *Msg,is_read int,channel_type int) {
	channel_a := helpler.ProduceChannelGroupName(strconv.Itoa(msg.ToId))
	fid := uint64(msg.FromId)
	tid := uint64(msg.ToId)
	user := messageModel.ImMessage{FromId:fid,
		ToId: tid,
		Msg: msg.Msg,
		CreatedAt: time.Unix(time.Now().Unix(), 0,).Format("2006-01-02 15:04:05"),
		Channel:channel_a,IsRead: is_read,MsgType: msg.MsgType,ChannelType: channel_type}
	model.DB.Create(&user)
	return
}


//The private chat insert db
func PutData(msg *Msg,is_read int,channel_type int) {
	channel_a,_ := helpler.ProduceChannelName( strconv.Itoa(msg.FromId), strconv.Itoa(msg.ToId))
	fid := uint64(msg.FromId)
	tid := uint64(msg.ToId)
	user := messageModel.ImMessage{FromId:fid,
		ToId: tid,
		Msg: msg.Msg,
		CreatedAt: time.Unix(time.Now().Unix(), 0,).Format("2006-01-02 15:04:05"),
		Channel:channel_a,IsRead: is_read,MsgType: msg.MsgType,ChannelType: channel_type}
	model.DB.Create(&user)
	return
}


func PushUserOnlineNotification(conn *ImClient,id int64)  {
	var msgList []ImMessage
	list := model.DB.Where("to_id=? and is_read=?", id, 0).Find(&msgList)
	if list.Error != nil {
		fmt.Println(list.Error)
	}
	for key, _ := range msgList {
		data, _ := json.Marshal(&Msg{Code: SendOk, Msg: msgList[key].Msg,
			FromId: msgList[key].FromId, ToId: msgList[key].ToId,
			Status: 0, MsgType: msgList[key].MsgType,ChannelType: msgList[key].ChannelType})
		conn.Send <- data
	}
}

func PushUserOfflineNotification(manager *ImClientManager,conn *ImClient)  {
	if _,ok := manager.ImClientMap[conn.ID];ok {
		id, _ := strconv.ParseInt(conn.ID, 10, 64)
		user.SetUserStatus(uint64(id), 0)
		jsonMessage, _ := json.Marshal(&OnlineMsg{Code: connOut, Msg: "用户离线了" + conn.ID, ID: conn.ID,ChannelType: 3})
		manager.ImSend(jsonMessage, conn)
		conn.Socket.Close()
		close(conn.Send)
		delete(manager.ImClientMap, conn.ID)
	}
}