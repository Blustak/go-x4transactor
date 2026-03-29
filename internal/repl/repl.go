package repl

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/Blustak/go-transActor/internal/config"
)

type replCommand struct {
	name     string
	desc     string
	args     []string
	callback func(c *config.Config) error
}

func (r *replCommand) Usage() string {
	return fmt.Sprintf("% 10s:%-s", r.name, r.desc)
}

type replCommandRegistry struct {
	reg map[string]replCommand
}

type ErrorCommandNotFound struct {
	key string
}

func (err *ErrorCommandNotFound) Error() string {
	return fmt.Sprintf("command %s not registered.", err.key)
}

func (r *replCommandRegistry) GetUsage(w io.Writer) error {
	for _, v := range r.reg {
		if _, err := w.Write([]byte(v.Usage() + "\n")); err != nil {
			return err
		}
	}
	return nil
}

func (r *replCommandRegistry) registerCommand(c *replCommand) {
	r.reg[c.name] = *c
}

func (r *replCommandRegistry) handleCommand(k string, args []string, cfg *config.Config) error {
	c, ok := r.reg[k]
	if !ok {
		return &ErrorCommandNotFound{key: k}
	}
	c.args = args
	return c.callback(cfg)
}

type Repler struct {
	registry *replCommandRegistry
	cfg      *config.Config
}

var defaultRegistry = func() *replCommandRegistry {
	r := replCommandRegistry{
		reg: make(map[string]replCommand),
	}

	r.registerCommand(&replCommand{
		name: "help",
		desc: "print this help",
		callback: func(c *config.Config) error {
			return r.GetUsage(c.Writer)
		},
	})

    r.registerCommand(&replCommand{
        name: "exit",
        desc: "quit the program",
        callback: func(_ *config.Config) error {
            return errors.New("exit")
        },
    })
	return &r

}()

func NewRepler(cfg *config.Config) *Repler {
	r := Repler{
		registry: defaultRegistry,
		cfg:      cfg,
	}

	return &r
}

func (r *Repler) Run() {
	if r.registry == nil {
		r.registry = defaultRegistry
	}
	bufw := bufio.NewWriter(r.cfg.Writer)
	bufScan := bufio.NewScanner(r.cfg.Reader)
	replLog := r.cfg.Log.With(slog.Any("repler", r))

	for {
		if err := func() error {
			fmt.Fprint(bufw, "> ")
			bufw.Flush() // Print the ticker flush, handle all then print once done
			defer bufw.Flush()

			bufScan.Scan()
			words := strings.Fields(bufScan.Text())
			if len(words) <= 0 {
				fmt.Fprintln(bufw, "no input given. type \"help\" for usage tips.")
			}
			key := words[0]
			args := words[1:]
			if handleErr := r.registry.handleCommand(key, args, r.cfg); handleErr != nil {
				if _, ok := errors.AsType[*ErrorCommandNotFound](handleErr); !ok {
					return handleErr
				} else {
					replLog.Warn("command not found", slog.String("command name", key))
					return nil
				}
			}
			fmt.Fprintln(bufw, "")
			return nil
		}(); err != nil {
			replLog.Warn("Got error", slog.String("error", err.Error()))
			break
		}
	}
}
