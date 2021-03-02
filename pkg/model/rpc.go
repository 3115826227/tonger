package model

import "tonger/pkg/constant"

type RPCMessage struct {
	MessageType constant.RPCMessageType `json:"message_type"`
}
