package protocol

import (
	"bytes"
	"encoding/gob"
	"net"
)

type ProtoType int

const (
	ProtocolOK ProtoType = iota
	ProtocolError
)

type ProtoMessage interface {
	GetProtoType() ProtoType
}

type ProtocolMessage struct {
	ProtocolType ProtoType
	Body         []byte
}

func Init() {
	gob.Register(OKMessage{})
	gob.Register(ErrorMessage{})
}

func NewMessage(v ProtoMessage) (msg *ProtocolMessage, err error) {
	m := ProtocolMessage{
		Body: make([]byte, 0),
	}
	m.ProtocolType = v.GetProtoType()
	if m.Body, err = Encode(v); err != nil {
		return
	}
	msg = &m
	return
}

func (m *ProtocolMessage) UnmarshalMessage(v ProtoMessage) error {
	bufr := bytes.NewReader(m.Body)
	dec := gob.NewDecoder(bufr)
	return dec.Decode(v)
}

func Decode(data []byte, v any) error {
	r := bytes.NewReader(data)
	dec := gob.NewDecoder(r)
	return dec.Decode(v)
}

func Encode(v any) ([]byte, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m *ProtocolMessage) ReadMessage(c net.Conn) error {
	var buf []byte
	empty := ProtocolMessage{}
	m = &empty
	bufr := bytes.NewReader(buf)
	if _, err := c.Read(buf); err != nil {
		return err
	}
	dec := gob.NewDecoder(bufr)
	if err := dec.Decode(m); err != nil {
		return err
	}
	return nil
}

func (m *ProtocolMessage) WriteMessage(c net.Conn) (int, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(m); err != nil {
		return 0, err
	}
	return c.Write(buf.Bytes())
}

func ReadMessage(c net.Conn) (*ProtocolMessage, error) {
	msg := &ProtocolMessage{
		Body: make([]byte, 0),
	}
	if err := msg.ReadMessage(c); err != nil {
		return nil, err
	}
	return msg, nil
}
