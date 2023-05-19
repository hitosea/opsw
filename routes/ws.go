package routes

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"opsw/utils"
	"opsw/vars"
	"sync"
)

var (
	wsClients  []vars.WsClientStruct
	wsMsg      = vars.WsMsgStruct{}
	wsMutex    = sync.Mutex{}
	wsUpgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

const (
	wsHeartbeat    = 0 // 心跳
	wsOnline       = 1 // 连接
	wsOffline      = 2 // 断开
	wsSendMsg      = 3 // 消息发送
	wsOnlineClient = 4 // 获取在线客户端
)

// AuthWs 启动 websocket
func (app *AppStruct) AuthWs() {
	conn, err := wsUpgrader.Upgrade(app.Context.Writer, app.Context.Request, nil)
	if err != nil {
		utils.GinResult(app.Context, http.StatusBadRequest, "连接失败", gin.H{"error": err.Error()})
		return
	}
	rid := utils.GenerateString(6)
	// 完成时关闭连接释放资源
	defer func(conn *websocket.Conn) {
		_ = conn.Close()
	}(conn)
	go func() {
		// 监听连接“完成”事件，其实也可以说丢失事件
		<-app.Context.Done()
		// 客户端离线
		app.wsOfflineClients(rid)
	}()
	// 添加客户端（上线）
	app.wsOnlineClients("user", rid, conn)
	sendByte, _ := json.Marshal(map[string]any{
		"type": wsOnline,
		"data": map[string]any{
			"rid": rid,
			"uid": app.UserInfo.Id,
		},
	})
	_ = conn.WriteMessage(websocket.TextMessage, sendByte)
	// 循环读取客户端发送的消息
	for {
		// 读取客户端发送过来的消息，如果没发就会一直阻塞住
		_, msg, err := conn.ReadMessage()
		if err != nil {
			app.wsOfflineClients(rid)
			break
		}
		err = json.Unmarshal(msg, &wsMsg)
		if err != nil {
			continue
		}
		if wsMsg.Data == nil {
			wsMsg.Data = make(map[string]any)
		}
		sendByte = nil
		if wsMsg.Type == wsHeartbeat {
			// 心跳
			sendByte, _ = json.Marshal(map[string]any{
				"type": wsHeartbeat,
			})
		} else if wsMsg.Type == wsSendMsg {
			// 消息发送
			toType, _ := wsMsg.Data.(map[string]any)["to_type"].(string) // 客户端类型
			toUid, _ := wsMsg.Data.(map[string]any)["to_uid"].(float64)  // 发送给谁
			msgData, _ := wsMsg.Data.(map[string]any)["msg_data"].(any)  // 消息内容
			if toUid == 0 || msgData == nil {
				continue
			}
			msgByte, _ := json.Marshal(map[string]any{
				"type": wsSendMsg,
				"data": msgData,
			})
			for _, v := range wsClients {
				if v.Type == toType && v.Uid == int32(toUid) {
					_ = v.Conn.WriteMessage(websocket.TextMessage, msgByte)
				}
			}
		} else if wsMsg.Type == wsOnlineClient {
			// 在线客户端
			var list []map[string]any
			for _, v := range wsClients {
				list = append(list, map[string]any{
					"rid": v.Rid,
					"uid": v.Uid,
				})
			}
			sendByte, _ = json.Marshal(map[string]any{
				"type": wsOnlineClient,
				"data": map[string]any{
					"count": len(list),
					"list":  list,
				},
			})
		}
		if sendByte != nil {
			_ = conn.WriteMessage(websocket.TextMessage, sendByte)
		}
	}
}

func (app *AppStruct) wsOnlineClients(type_, rid string, conn *websocket.Conn) {
	for _, v := range wsClients {
		if v.Rid == rid {
			return
		}
	}
	wsMutex.Lock()
	wsClients = append(wsClients, vars.WsClientStruct{
		Conn: conn,
		Type: type_,
		Uid:  app.UserInfo.Id,
		Rid:  rid,
	})
	wsMutex.Unlock()
	app.wsNotifyClients(rid, vars.WsMsgStruct{
		Type: wsOnline,
		Data: map[string]any{
			"rid": rid,
			"uid": app.UserInfo.Id,
		},
	})
}

func (app *AppStruct) wsOfflineClients(rid string) {
	for k, v := range wsClients {
		if v.Rid == rid {
			wsMutex.Lock()
			wsClients = append(wsClients[:k], wsClients[k+1:]...)
			_ = v.Conn.Close()
			wsMutex.Unlock()
			break
		}
	}
	app.wsNotifyClients(rid, vars.WsMsgStruct{
		Type: wsOffline,
		Data: map[string]any{
			"rid": rid,
			"uid": app.UserInfo.Id,
		},
	})
}

func (app *AppStruct) wsNotifyClients(rid string, msgData vars.WsMsgStruct) {
	sendByte, err := json.Marshal(msgData)
	if err != nil {
		return
	}
	for _, v := range wsClients {
		if v.Rid != rid {
			_ = v.Conn.WriteMessage(websocket.TextMessage, sendByte)
		}
	}
}
