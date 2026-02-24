package ws

import "encoding/json"

// Envelope: mensaje genérico SCAVO WS
// id: correlation id (client o server)
// type: req|res|evt
// action: nombre del comando (auth.login, system.ping, etc.)
// data: payload
// error: error (para res)

type MsgType string

const (
	MsgTypeReq MsgType = "req"
	MsgTypeRes MsgType = "res"
	MsgTypeEvt MsgType = "evt"
)

type Envelope struct {
	ID     string          `json:"id"`
	Type   MsgType         `json:"type"`
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data,omitempty"`
	Error  *ErrPayload     `json:"error,omitempty"`
}

type ErrPayload struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}
