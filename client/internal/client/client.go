package client

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
	"strings"

	"github.com/Blustak/go-transActor/internal/protocol"
)

type Client struct {
	Hostname, Port string
	*slog.Logger
	In  *bufio.Scanner
	Out *bufio.Writer
}

func (c *Client) Run() error {
	for {
		fmt.Fprint(c.Out, "> ")
		if err := c.Out.Flush(); err != nil {
			return err
		}
		c.In.Scan()
		if err := c.In.Err(); err != nil {
			return err
		}
		if err := c.Step(); err != nil {
			return err
		}
		if err := c.Out.Flush(); err != nil {
			return err
		}
	}
}

func (c *Client) Step() error {
	args := strings.Fields(c.In.Text())

	if len(args) <= 0 {
		if _, err := fmt.Fprintln(c.Out, "Please enter something."); err != nil {
			return err
		}
		return nil
	}
	c.Debug("entered command", slog.Any("args", args))
	switch args[0] {
	case "q", "quit", "exit":
		return fmt.Errorf("Quit command")
	case "OK":
		return c.handleOk()
	default:
		if _, err := fmt.Fprintf(c.Out, "unrecognised command %s\n", args[0]); err != nil {
			return err
		}
		return nil
	}
}

func (c *Client) handleOk() error {
	req, err := protocol.NewMessage(&protocol.OKMessage{})
	if err != nil {
		return err
	}
	conn, err := net.Dial("tcp4", c.Hostname+":"+c.Port)
	if err != nil {
		c.Error("failed to connect to remote", slog.String("address", c.Hostname+":"+c.Port), slog.Any("error", err))
		return err
	}
	defer conn.Close()
	c.Info("connected to host", slog.Any("remote addr", conn.RemoteAddr().String()), slog.Any("local addr", conn.LocalAddr().String()))
	n, err := req.WriteMessage(conn)
	if err != nil {
		c.Error("error writing request", slog.Any("error", err), slog.Int("bytes written", n))
		return err
	}
	c.Debug("wrote message", slog.Any("request", req), slog.Int("bytes written", n))
	res, err := protocol.ReadMessage(conn)
	if err != nil {
		c.Error("error reading message", slog.Any("error", err))
		return err
	}
	switch res.ProtocolType {
	case protocol.ProtocolError:
		v := &protocol.ErrorMessage{}
		if err := res.UnmarshalMessage(v); err != nil {
			c.Error("error unmarshalling message", slog.Any("error", err))
			return err
		}
		fmt.Fprintf(c.Out, "Error in OK request: %s\n", v.Msg)
		return nil
	case protocol.ProtocolOK:
		v := &protocol.OKMessage{}
		if err := res.UnmarshalMessage(v); err != nil {
			c.Error("error unmarshalling message", slog.Any("error", err))
			return err
		}
		fmt.Fprintln(c.Out, "Server says OK!")
		return nil
	default:
		c.Error("unrecognised message type", slog.Any("protocolType", res.ProtocolType), slog.Any("message", res))
		return fmt.Errorf("unrecognised protocol type")
	}
}
