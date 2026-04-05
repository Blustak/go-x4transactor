package protocol

import "encoding/gob"

type CommandType int

const (
    CommandOK CommandType = iota
)

type StatusCode int

const (
    StatusOK StatusCode = iota
)

func Init () {
    var s StatusCode
    var c CommandType
    gob.Register(s)
    gob.Register(c)
}

type Message[P any] struct{
    Payload P
}

type Response[R any] struct{
    Status StatusCode
    Payload R
}

type Handshake struct{
    Command CommandType
}
