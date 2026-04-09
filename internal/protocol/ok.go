package protocol

type OKMessage struct{}

func (m *OKMessage) GetProtoType() ProtoType {
	return ProtocolOK
}
