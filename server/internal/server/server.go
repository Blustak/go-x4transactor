package server

import (
	"database/sql"
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
	s.Info("started listening", slog.Any("listener address", listener.Addr().String()))
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		s.handleConn(conn)

	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	s.Info("handling connection", slog.Any("remoteAddr", conn.RemoteAddr().String()))

	msg, err := protocol.ReadMessage(conn)
	if err != nil {
		writeErrorMessage("couldn't read message", conn)
		return
	}
	res, err := s.handleMessage(msg)
	if err != nil {
		writeErrorMessage(err.Error(), conn)
		return
	}
	n, err := res.WriteMessage(conn)
	if err != nil {
		s.State.Error("error writing message", slog.Any("error", err))
	}
	s.State.Info("wrote bytes to conn", slog.Int("bytes written", n), slog.Any("conn", conn))

}

func (s *Server) handleMessage(m *protocol.ProtocolMessage) (*protocol.ProtocolMessage, error) {
	var msg protocol.ProtoMessage
	switch m.ProtocolType {
	case protocol.ProtocolOK:
		// Just send back an OK message
		msg = &protocol.OKMessage{}
	default:
		// unrecognised, so send back an error
		return nil, fmt.Errorf("unrecognised protocol type: %v", m.ProtocolType)

	}
	return protocol.NewMessage(msg)
}

func writeErrorMessage(msg string, c net.Conn) (int, error) {
	errMsg := protocol.ErrorMessage{
		Msg: msg,
	}
	m, err := protocol.NewMessage(&errMsg)
	if err != nil {
		return 0, err
	}
	return m.WriteMessage(c)
}
