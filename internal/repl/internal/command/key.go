package command

import (
	"fmt"
	"slices"
)

type CommandKey int

const (
    KeyHelp CommandKey = iota
)

var keyStringMap = map[CommandKey]string{
    KeyHelp: "help",
}


var keyAliases = map[CommandKey][]string{
    KeyHelp: {"h", "help"},
}

func (c CommandKey) String() string {
    return keyStringMap[c]
}

var commandUsageMap = map[CommandKey]string{
    KeyHelp: "print this help",
}

func (c CommandKey) Usage() string {
    return commandUsageMap[c]
}

func ParseString(s string) (CommandKey, error) {
    for k,v := range keyAliases {
        if slices.Contains(v,s) {
            return k, nil
        }
    }
    return 0, fmt.Errorf("%s not a recognised key",s)
}
