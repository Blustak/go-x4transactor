package protocol

type ErrorMessage struct {
	Msg string
}

func (err *ErrorMessage) GetProtoType() ProtoType {
	return ProtocolError
}
