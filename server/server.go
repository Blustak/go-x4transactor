package server

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	"log/slog"
	"net"
	"sync"

	"github.com/Blustak/go-transActor/internal/config"
	"github.com/Blustak/go-transActor/internal/database"
	"github.com/Blustak/go-transActor/internal/protocol"
	_ "github.com/mattn/go-sqlite3"
)

type Server struct {
	*config.State
	Hostname, Port, protocol string
}

func (s *Server) URL() string {
	return fmt.Sprintf("%s://%s:%s", s.protocol, s.Hostname, s.Port)
}

func NewServer(hostName, port, dbPath string, log *slog.Logger) (s *Server, err error) {
	s = &Server{
		protocol: "tcp4",
		State: &config.State{
			DB:      nil,
			RWMutex: sync.RWMutex{},
		},
	}
	dbConn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return
	}

	s.State.DB = database.New(dbConn)
	s.Hostname = hostName
	s.Port = port
	s.Logger = log

	if s.Logger != nil {
		s.Logger = log.With(slog.Group("server", slog.String("URL", s.URL()), slog.String("DBPath", dbPath)))
	}

	return
}

func (s *Server) ListenAndServe() error {
	listener, err := net.Listen("tcp4", s.Hostname+":"+s.Port)
	if err != nil {
		return err
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
        handleConn(conn,s.State)

	}
}

func handleConn(conn net.Conn, cfg *config.State) {
    defer conn.Close()

    var hs protocol.Handshake
    enc := gob.NewEncoder(conn)
    dec := gob.NewDecoder(conn)
    if err := dec.Decode(&hs); err != nil {
        cfg.Error("handshake failed",slog.Any("error", err))
        return 
    }
    switch hs.Command{
        case protocol.CommandOK:
            if err := serveConn(dec,enc,handleOk); err != nil {
                cfg.Error("error serving OK request", slog.String("error", err.Error()))
                return
            }

    }
    
}

func serveConn[P,R any](dec *gob.Decoder, enc *gob.Encoder, process func(protocol.Message[P]) (protocol.Response[R], error)) error {
    var msg protocol.Message[P]
    if err := dec.Decode(&msg); err != nil {return err}
    res, err := process(msg)
    if err != nil {
        return err
    }
    if err := enc.Encode(res); err != nil {
        return err
    }
    return nil
}

func handleOk(msg protocol.Message[protocol.OkRequest]) (protocol.Response[protocol.OkResponse], error) {
    return protocol.Response[protocol.OkResponse]{
        Status: protocol.StatusOK,
        Payload: protocol.OkResponse{},
    }, nil
}
