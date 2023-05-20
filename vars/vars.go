package vars

var (
	Config    ConfStruct
	WsClients []WsClientStruct
)

const (
	WsHeartbeat    = 0 // 心跳
	WsOnline       = 1 // 连接
	WsOffline      = 2 // 断开
	WsSendMsg      = 3 // 消息发送
	WsOnlineClient = 4 // 获取在线客户端
)
