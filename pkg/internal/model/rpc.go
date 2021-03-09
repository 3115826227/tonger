package model

import (
	"encoding/json"
	"tonger/pkg/constant"
)

type RPCMessage struct {
	MessageType constant.RPCMessageType `json:"message_type"`
}

func (msg *RPCMessage) ToString() string {
	data, _ := json.Marshal(msg)
	return string(data)
}
