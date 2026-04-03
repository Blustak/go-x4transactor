package command

import "github.com/Blustak/go-transActor/internal/config"

type Command struct{
    Key CommandKey
    callback func(c *config.Config) error
}
